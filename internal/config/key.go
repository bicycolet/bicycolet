package config

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	Base64Prefix = "base64:"
)

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	defaultValue string
	hidden       bool
	deprecated   string
	validator    func(string) error
}

// key defines the type of the value of a particular config key, along with
// other knobs such as default, validator, etc.
type key struct {
	// If the key is not set in a Map, use this value instead.
	// Defaults go via the getters and setters (only change)
	// So should return the value desired not the map state
	Default    string
	Hidden     bool   // Hide this key when dumping the object.
	Deprecated string // Optional message to set if this config value is deprecated.

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

// Tells if the given value can be assigned to this particular Value instance.
func (v *key) validate(value string) error {
	validator := v.Validator
	if validator == nil {
		// Dummy validator
		validator = func(string) error { return nil }
	}

	// Handle unsetting
	if value == "" {
		return validator(v.Default)
	}

	if v.Deprecated != "" && value != v.Default {
		return errors.Errorf("deprecated: %s", v.Deprecated)
	}

	// Run external validation function
	return validator(value)
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

// newBaseOptions a options instance with default values.
func newBaseOptions() *options {
	return &options{
		hidden:     false,
		deprecated: "",
	}
}

// newOptionsWithDefaultValidator a options instance with default validator.
func newOptionsWithDefaultValidator(validator func(string) error) *options {
	opts := newBaseOptions()
	opts.validator = validator
	return opts
}

// NewBoolKey create a boolean key
func NewBoolKey(options ...Option) key {
	opts := newOptionsWithDefaultValidator(validateBool)
	for _, option := range options {
		option(opts)
	}

	return key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     normalizeBool,
		Getter:     boolGetter,
	}
}

// NewStringKey create a boolean key
func NewStringKey(options ...Option) key {
	opts := newBaseOptions()
	for _, option := range options {
		option(opts)
	}

	return key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
	}
}

// NewInt64Key create a boolean key
func NewInt64Key(options ...Option) key {
	opts := newOptionsWithDefaultValidator(validateInt64)
	for _, option := range options {
		option(opts)
	}

	return key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
		Getter:     int64Getter,
	}
}

// NewDurationKey create a boolean key
func NewDurationKey(options ...Option) key {
	opts := newOptionsWithDefaultValidator(validateDuration)
	for _, option := range options {
		option(opts)
	}

	return key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     nil,
		Getter:     durationGetter,
	}
}

// NewByteArrayKey create a boolean key
func NewByteArrayKey(options ...Option) key {
	opts := newBaseOptions()
	for _, option := range options {
		option(opts)
	}

	return key{
		Default:    opts.defaultValue,
		Hidden:     opts.hidden,
		Deprecated: opts.deprecated,
		Validator:  opts.validator,
		Setter:     Base64Encode,
		Getter:     Base64Decode,
	}
}

// Validation

// validateBool
func validateBool(value string) error {
	if value == "" {
		return nil
	}
	if !contains(strings.ToLower(value), booleans) {
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

// Setters

// Normalize a boolean value, converting it to the string "true" or "false".
func normalizeBool(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	if contains(strings.ToLower(value), truthy) {
		return "true", nil
	}
	return "false", nil
}

// Base64Encode encodes a string to base64
func Base64Encode(value string) (string, error) {
	if strings.HasPrefix(value, Base64Prefix) {
		return value, nil
	}
	b64Value := base64.StdEncoding.EncodeToString([]byte(value))
	return fmt.Sprintf("%s%s", Base64Prefix, b64Value), nil
}

// Getters

// Base64Decode encodes a string to base64
func Base64Decode(b64Value string) (interface{}, error) {
	if !strings.HasPrefix(b64Value, Base64Prefix) {
		return "", errors.New("can't decode non encoded string")
	}
	value, err := base64.StdEncoding.DecodeString(b64Value[len(Base64Prefix):])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func boolGetter(v string) (interface{}, error) {
	return contains(strings.ToLower(v), truthy), nil
}

func int64Getter(v string) (interface{}, error) {
	return strconv.ParseInt(v, 10, 64)
}

func durationGetter(v string) (interface{}, error) {
	return time.ParseDuration(v)
}
