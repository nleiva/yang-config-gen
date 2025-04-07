package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(t *testing.T) {
	tt := []struct {
		name     string
		file     string
		expected string
		err      string
	}{
		{
			name:     "generate-interface",
			file:     "../model/testdata/interface.json",
			expected: intf,
		},
		{
			name:     "generate-bgp",
			file:     "../model/testdata/bgp.json",
			expected: bgp,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := run(tc.file)
			if err != nil {
				t.Errorf("Can't read file %s: %v", tc.file, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

const intf = `{
  "configuration": {
    "interfaces": {
      "interface": [
        {
          "description": "description",
          "encapsulation": "vlan-ccc",
          "link-mode": "full-duplex",
          "mtu": 1500,
          "name": "eth-0/1/1",
          "native-vlan-id": 22,
          "speed": "100g"
        },
        {
          "name": "ge-2/2/0",
          "unit": [
            {
              "family": {
                "inet": {
                  "address": [
                    {
                      "name": "10.10.10.1/24"
                    }
                  ]
                }
              },
              "name": "10"
            }
          ]
        }
      ]
    },
    "routing-instances": {
      "instance": [
        {
          "interface": [
            {
              "name": "ge-2/2/0.10"
            }
          ],
          "name": "red"
        }
      ]
    }
  }
}`

const bgp = `{
  "configuration": {
    "routing-instances": {
      "instance": [
        {
          "name": "office",
          "protocols": {
            "bgp": {
              "group": [
                {
                  "export": [
                    "ce-export-policy"
                  ],
                  "import": [
                    "ce-import-policy"
                  ],
                  "name": "customers",
                  "neighbor": [
                    {
                      "description": "cust_a",
                      "name": "10.99.99.1",
                      "peer-as": "65002"
                    }
                  ]
                }
              ]
            }
          },
          "routing-options": {
            "autonomous-system": {
              "as-number": "65001"
            }
          }
        }
      ]
    }
  }
}`
