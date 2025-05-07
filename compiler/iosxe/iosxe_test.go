package iosxe

import (
	"strings"
	"testing"

	"github.com/nleiva/yang-config-gen/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateHostnameConfig(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "hostname:test",
			config: model.Target{
				Hostname: "test",
			},
			expected: host1,
		},
		{
			name: "hostname:my-device",
			config: model.Target{
				Hostname: "my-device",
			},
			expected: host2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := NewCompiler()
			err := c.CreateHostnameConfig(tc.config)
			if err != nil {
				t.Errorf("can't create interface config: %s", err.Error())
			}

			jsonConfig, err := c.EmitConfig()
			if err != nil {
				t.Errorf("can't emit config: %s", err.Error())
			}

			switch tc.err {
			case "":
				{
					expected := cleanString(tc.expected)
					result := cleanString(jsonConfig)
					assert.Equal(t, expected, result)
				}
			default:
				assert.ErrorContains(t, err, tc.err)
			}

		})
	}
}

func cleanString(in string) (out string) {
	// Remove tabs and spaces
	out = strings.Replace(in, " ", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	return
}

const host1 = `{
  "native": {
    "hostname": "test"
  }
}`

const host2 = `{
  "native": {
    "hostname": "my-device"
  }
}`
