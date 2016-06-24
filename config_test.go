package main

import (
	"reflect"
	"testing"
)

var configTests = []struct {
	env    []string
	config Config
}{
	{
		env: []string{
			"SERVICES=fake:fake",
			"FAKE_DESCRIPTION=a fake service",
			"FAKE_CONFIG=true",
			"FAKE_USER_MAP=foo:bar",
		},
		config: Config{
			Services: map[string]ServiceConfig{
				"fake": ServiceConfig{
					Type:        "fake",
					Description: "a fake service",
					Config:      "true",
					UserMap: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
	},
}

func TestParseConfig(t *testing.T) {
	for _, tt := range configTests {
		config := ParseConfig(tt.env)
		if !reflect.DeepEqual(tt.config, config) {
			t.Errorf("expected ParseConfig(%+v) to return %+v, but returned %+v", tt.env, tt.config, config)
		}
	}
}

var getEnvTests = []struct {
	env   []string
	key   string
	value string
}{
	{
		[]string{"FOOBAR=foo", "FOO=bar"},
		"FOOBAR",
		"foo",
	},
	{
		[]string{"FOOBAR=foo", "FOO=bar"},
		"FOO",
		"bar",
	},
	{
		[]string{"FOOBAR=foo", "FOO=bar"},
		"doesn't exist",
		"",
	},
}

func TestGetEnv(t *testing.T) {
	for _, tt := range getEnvTests {
		actualValue := getEnv(tt.env, tt.key)
		if actualValue != tt.value {
			t.Errorf("getEnv(%v, %q) = %q, but expected %q", tt.env, tt.key, actualValue, tt.value)
		}
	}
}
