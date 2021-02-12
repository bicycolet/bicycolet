package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Key defines the type of the value of a particular config key, along with
// other knobs such as default, validator, etc.
type Key struct {
	// If the key is not set in a Map, use this value instead.
	// Defaults go via the getters and setters (only change)
	// So should return the value desired not the map state
	Default string
	// Hide this key when dumping the object.
	Hidden bool
	// Optional message to set if this config value is deprecated.
	Deprecated string

	// Optional function used to validate the values. It's called by Map
	// all the times the value associated with this key is going to be
	// changed.
	Validator func(string) error

	// Optional function to manipulate a value before it's actually saved
	// in a Map. Called for both load and save.
	Setter func(string) (string, error)

	// Optional function to manipulate a value before it's returned from the
	// Map
	Getter func(string) (interface{}, error)
}

// Validate tells if the given value can be assigned to this particular Value
// instance.
func (v *Key) Validate(value string) error {
	// Handle unsetting
	if value == "" {
		if v.Validator == nil {
			return nil
		}
		return v.Validator(v.Default)
	}

	if v.Deprecated != "" && value != v.Default {
		return errors.Errorf("deprecated: %s", v.Deprecated)
	}

	// Run external validation function
	if v.Validator == nil {
		return nil
	}
	return v.Validator(value)
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	defaultValue string
	hidden       bool
	deprecated   string
	validator    func(string) error
}

// WithDefault sets the default option
func WithDefault(value string) Option {
	return func(options *options) {
		options.defaultValue = value
	}
}

// WithHidden sets the hidden option
func WithHidden(hidden bool) Option {
	return func(options *options) {
		options.hidden = hidden
	}
}

// WithDeprecated sets the deprecated option
func WithDeprecated(deprecated string) Option {
	return func(options *options) {
		options.deprecated = deprecated
	}
}

// WithValidator sets the deprecated option
func WithValidator(validator func(string) error) Option {
	return func(options *options) {
		options.validator = validator
	}
}

// newOptions a options instance with default values.
func newOptions() *options {
	return &options{
		hidden:     false,
		deprecated: "",
		validator: func(string) error {
			return nil
		},
	}
}

// newOptionsWithValidator a options instance with default validator.
func newOptionsWithValidator(validator func(string) error) *options {
	opts := newOptions()
	opts.validator = validator
	return opts
}

// NewBoolKey create a boolean key with default validators.
func NewBoolKey(options ...Option) Key {
	opts := newOptionsWithValidator(validateBool)
	for _, option := range options {
		option(opts)
	}

	return Key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     normalizeBool,
		Getter:     boolGetter,
	}
}

// NewStringKey create a boolean key
func NewStringKey(options ...Option) Key {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	return Key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
	}
}

// NewInt64Key create a boolean key
func NewInt64Key(options ...Option) Key {
	opts := newOptionsWithValidator(validateInt64)
	for _, option := range options {
		option(opts)
	}

	return Key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
		Getter:     int64Getter,
	}
}

// NewDurationKey create a boolean key
func NewDurationKey(options ...Option) Key {
	opts := newOptionsWithValidator(validateDuration)
	for _, option := range options {
		option(opts)
	}

	return Key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
		Getter:     durationGetter,
	}
}

var (
	truthy = []string{
		"true", "1", "yes", "on",
	}
	falsy = []string{
		"false", "0", "no", "off",
	}
	boolean = append(truthy, falsy...)
)

func contains(list []string, key string) bool {
	for _, entry := range list {
		if entry == key {
			return true
		}
	}
	return false
}

// Normalize a boolean value, converting it to the string "true" or "false".
func normalizeBool(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	if contains(truthy, strings.ToLower(value)) {
		return "true", nil
	}
	return "false", nil
}

func validateBool(value string) error {
	if value == "" {
		return nil
	}
	if !contains(boolean, strings.ToLower(value)) {
		return errors.Errorf("invalid boolean")
	}
	return nil
}

// validateInt64
func validateInt64(value string) error {
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return errors.Errorf("invalid integer")
	}
	return nil
}

// validateDuration
func validateDuration(value string) error {
	_, err := time.ParseDuration(value)
	if err != nil {
		return errors.Errorf("invalid duration")
	}
	return nil
}

func boolGetter(v string) (interface{}, error) {
	return contains(truthy, strings.ToLower(v)), nil
}

func int64Getter(v string) (interface{}, error) {
	return strconv.ParseInt(v, 10, 64)
}

func durationGetter(v string) (interface{}, error) {
	return time.ParseDuration(v)
}
