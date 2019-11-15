package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/bicycolet/bicycolet/client"
	"github.com/bicycolet/bicycolet/client/info"
	pkgclient "github.com/bicycolet/bicycolet/pkg/client"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/flagset"
	yaml "gopkg.in/yaml.v2"
)

type baseCmd struct {
	ui      clui.UI
	flagset *flagset.FlagSet

	debug  bool
	format string
}

func (c *baseCmd) init() {
	c.flagset.BoolVar(&c.debug, "debug", false, "debug logging")
	c.flagset.StringVar(&c.format, "format", "yaml", "format to output the information json|yaml|tabular")
}

// UI returns a UI for interaction.
func (c *baseCmd) UI() clui.UI {
	return c.ui
}

// FlagSet returns the FlagSet associated with the command. All the flags are
// parsed before running the command.
func (c *baseCmd) FlagSet() *flagset.FlagSet {
	return c.flagset
}

func (c *baseCmd) Output(value interface{}) error {
	if !contains([]string{"json", "yaml", "tabular"}, c.format) {
		return errors.Errorf("invalid format type (expected: json|yaml|tabular) got: %s", c.format)
	}

	var result string
	switch c.format {
	case "yaml":
		bytes, err := yaml.Marshal(value)
		if err != nil {
			return errors.WithStack(err)
		}
		result = string(bytes)
	case "json":
		bytes, err := json.MarshalIndent(value, "", "\t")
		if err != nil {
			return errors.WithStack(err)
		}
		result = string(bytes)
	case "tabular":
		bytes, err := constructTabularOutput(value)
		if err != nil {
			return errors.WithStack(err)
		}
		result = string(bytes)
	default:
		result = fmt.Sprintf("%+v", value)
	}
	c.ui.Output(result)
	return nil
}

func constructTabularOutput(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 2, 2, 3, ' ', 0)

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	switch t.Kind() {
	case reflect.Map:
		for _, idx := range v.MapKeys() {
			out, err := constructTabularOutput(struct {
				Key   interface{}
				Value interface{}
			}{
				Key:   idx.Interface(),
				Value: v.MapIndex(idx).Interface(),
			})
			if err != nil {
				return nil, errors.WithStack(err)
			}
			fmt.Fprintln(buf, string(out))
		}
	case reflect.Struct:
		var headers []string
		var values []string
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tab := field.Tag.Get("tab")
			if tab == "" {
				tab = field.Name
			}
			headers = append(headers, tab)
			values = append(values, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
		fmt.Fprintln(w, strings.Join(headers, "\t"))
		fmt.Fprintln(w, strings.Join(values, "\t"))
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			out, err := constructTabularOutput(struct {
				Index int
				Value interface{}
			}{
				Index: i,
				Value: v.Index(i).Interface(),
			})
			if err != nil {
				return nil, errors.WithStack(err)
			}
			fmt.Fprintln(buf, string(out))
		}
	}
	return nil, errors.Errorf("unexpected type %s", t.Kind().String())
}

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

func getClient(address string, logger log.Logger) (*client.Client, error) {
	return client.New(
		address,
		client.WithLogger(log.WithPrefix(logger, "component", "client")),
	)
}

func getInfoClient(address string, logger log.Logger) (*info.Info, error) {
	client, err := getClient(address, logger)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return info.New(clientInfoShim{client: client}), nil
}

type clientInfoShim struct {
	client *client.Client
}

func (s clientInfoShim) Get(url string, fn func(info.Response, info.Metadata) error) error {
	return s.client.Get(url, func(res *pkgclient.Response, meta client.Metadata) error {
		return fn(clientResponseShim{
			response: res,
		}, meta)
	})
}

type clientResponseShim struct {
	response *pkgclient.Response
}

func (s clientResponseShim) Metadata() json.RawMessage {
	return s.response.Metadata
}
