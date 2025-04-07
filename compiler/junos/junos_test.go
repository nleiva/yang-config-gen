package junos

import (
	"strings"
	"testing"

	"github.com/nleiva/yang-config-gen/model"
	"github.com/nleiva/yang-data-structures/junos"
	"github.com/stretchr/testify/assert"
)

func TestCreateIntConfigText(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "config-lo0.0",
			config: model.Target{
				Interfaces: model.Interfaces{
					Interface: map[string]model.Interface{
						"lo0": {
							Subinterfaces: model.SubInterfaces{
								SubInterface: map[string]model.SubInterface{
									"0": {
										Config: model.SubInterfaceConfig{
											Description: "Test-Description",
										},
									},
								},
							},
						},
					},
				},
			},
			expected: ygotLo00,
		},
		{
			name: "config-lo0",
			config: model.Target{
				Interfaces: model.Interfaces{
					Interface: map[string]model.Interface{
						"lo0": {
							Config: model.InterfaceConfig{
								Description: "Test-Description",
							},
						},
					},
				},
			},
			expected: ygotLo0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			j := NewCompiler()
			err := j.CreateIntConfig(tc.config)
			if err != nil {
				t.Errorf("can't create interface config: %s", err.Error())
			}

			jsonConfig, err := j.EmitConfig()
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

func TestNetworkInstance(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "config-lo0.0",
			config: model.Target{
				NetworkInstances: model.NetworkInstances{},
			},
			expected: "{}",
		},
		{
			name: "config-lo0",
			config: model.Target{
				NetworkInstances: model.NetworkInstances{},
			},
			expected: "{}",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			j := NewCompiler()
			err := j.CreateRoutingInstancesConfig(tc.config)
			if err != nil {
				t.Errorf("can't create routing instance config: %s", err.Error())
			}

			jsonConfig, err := j.EmitConfig()
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

func TestCreateBGPConfig(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "config-bgp-1",
			config: model.Target{
				NetworkInstances: model.NetworkInstances{
					NetworkInstance: map[string]model.NetworkInstance{
						"test": {
							// Interfaces: model.NtwInstInterfaces{
							// 	Interface: map[string]model.NtwInstInterface{},
							// },
							Protocols: model.NtwInstProtocols{
								Protocol: map[string]model.NtwInstProtocol{
									"BGP": {
										BGP: model.BGP{
											Global: model.Global{
												Config: model.GlobalConfig{},
											},
											Neighbors: model.Neighbors{
												Neighbor: map[string]model.Neighbor{
													"10.99.99.1": {
														NeighborAddress: "10.99.99.1",
														Config: model.NeighborConfig{
															NeighborAddress: "10.99.99.1",
															PeerAs:          65002,
															Enabled:         true,
															PeerGroup:       "customers",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: bgp1,
		},
		{
			name: "config-bgp-2",
			config: model.Target{
				NetworkInstances: model.NetworkInstances{},
			},
			expected: "{}",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			j := NewCompiler()
			err := j.CreateBGPConfig(tc.config)
			if err != nil {
				t.Errorf("can't create BGP config: %s", err.Error())
			}

			jsonConfig, err := j.EmitConfig()
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

func TestCreatePolicyConfig(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "config-policy-1",
			config: model.Target{
				RoutingPolicy: &model.RoutingPolicy{
					DefinedSets: model.DefinedSets{
						PrefixSets: model.PrefixSets{
							PrefixSet: map[string]model.PrefixSet{
								"test": {
									Name:   "test",
									Config: model.PrefixSetConfig{},
									Prefixes: model.Prefixes{
										Prefix: map[string]model.Prefix{
											"192.68.28.0/22..32": {
												IPPrefix:        "192.68.28.0/22",
												MasklengthRange: "22..32",
												Config: model.PrefixConfig{
													IPPrefix:        "192.68.28.0/22",
													MasklengthRange: "22..32",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: policy1,
		},
		{
			name: "config-policy-2",
			config: model.Target{
				RoutingPolicy: &model.RoutingPolicy{},
			},
			expected: "{}",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			j := NewCompiler()
			err := j.CreatePrefixListConfig(tc.config)
			if err != nil {
				t.Errorf("can't create Policy config: %s", err.Error())
			}

			jsonConfig, err := j.EmitConfig()
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

func TestCreateACLConfig(t *testing.T) {
	tt := []struct {
		name     string
		config   model.Target
		expected string
		err      string
	}{
		{
			name: "config-acl-1",
			config: model.Target{
				ACL: &model.ACL{
					ACLSets: model.ACLSets{
						ACLSet: map[string]model.ACLSet{
							"test": {
								ACLEntries: model.ACLEntries{
									ACLEntry: map[string]model.ACLEntry{
										"myACL": {
											SequenceID: 10,
											IPv4: model.ACLIPv4{
												Config: model.ACLIPv4Config{
													SourceAddresses:      []string{"10.0.0.0/8"},
													Protocol:             "IP-TCP",
													DestinationAddresses: []string{"100.64.0.0/24"},
													DSCP:                 5,
												},
											},
											Actions: model.ACLActions{
												Config: model.ACLActionsConfig{
													ForwardingAction: "ACCEPT",
													TargetGroup:      "cs2",
												},
											},
											Transport: model.ACLTransport{
												Config: model.ACLTransportConfig{
													DestinationPort: "443",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: acl1,
		},
		{
			name: "config-acl-2",
			config: model.Target{
				ACL: &model.ACL{},
			},
			expected: "{}",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			j := NewCompiler()
			err := j.CreateACLConfig(tc.config)
			if err != nil {
				t.Errorf("can't create ACL config: %s", err.Error())
			}

			jsonConfig, err := j.EmitConfig()
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

func TestReadJSONFromRouter(t *testing.T) {
	tt := []struct {
		name     string
		intf     string
		config   string
		expected string
		err      string
	}{
		{
			name:     "config-lo0",
			intf:     "lo0",
			config:   routerLo0,
			expected: "Loopback-CLI",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()
			load := &junos.Junos{}
			if err := junos.Unmarshal([]byte(routerLo0), load); err != nil {
				t.Errorf("Can't unmarshal JSON: %v", err)
			}

			assert.Equal(t, tc.expected, *load.Configuration.Interfaces.Interface[tc.intf].Description)

		})
	}
}

func cleanString(in string) (out string) {
	// Remove tabs and spaces
	out = strings.Replace(in, " ", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	return
}

const ygotLo0 = `{
  "configuration": {
    "interfaces": {
      "interface": [
        {
          "description": "Test-Description",
          "name": "lo0"
        }
      ]
    }
  }
}`

const ygotLo00 = `{
  "configuration": {
    "interfaces": {
      "interface": [
        {
          "name": "lo0",
          "unit": [
            {
              "description": "Test-Description",
              "name": "0"
            }
          ]
        }
      ]
    }
  }
}`

// OUTPUT from device. I made the unit name a string instead of a number.
// 1. got float64 type for field name, expect string
const routerLo0 = `{
    "configuration" : {
        "@" : {
            "junos:changed-seconds" : "1743092661", 
            "junos:changed-localtime" : "2025-03-27 16:24:21 UTC"
        }, 
        "interfaces" : {
            "interface" : [
            {
                "name" : "lo0", 
                "description" : "Loopback-CLI", 
                "unit" : [
                {
                    "name" : "0", 
                    "description" : "Test-Description", 
                    "family" : {
                        "inet" : {
                            "address" : [
                            {
                                "name" : "2.2.2.2/32"
                            }
                            ]           
                        }, 
                        "iso" : {
                            "address" : [
                            {
                                "name" : "39.752f.0100.0014.0000.9000.0020.0102.4308.2198.00"
                            }
                            ]
                        }
                    }
                }
                ]
            }
            ]
        }
    }
}`

const bgp1 = `{
  "configuration": {
    "routing-instances": {
      "instance": [
        {
          "name": "test",
          "protocols": {
            "bgp": {
              "group": [
                {
                  "name": "customers",
                  "neighbor": [
                    {
                      "name": "10.99.99.1",
                      "peer-as": "65002"
                    }
                  ]
                }
              ]
            }
          }
        }
      ]
    }
  }
}`

const policy1 = `{
  "configuration": {
    "policy-options": {
      "policy-statement": [
        {
          "name": "test",
          "term": [
            {
              "from": {
                "route-filter": [
                  {
                    "address": "192.68.28.0/22",
                    "choice-ident": "orlonger",
                    "choice-value": ""
                  }
                ]
              },
              "name": "accept"
            }
          ],
          "then": {
            "accept": [
              null
            ]
          }
        }
      ]
    }
  }
}`

const acl1 = `{
  "configuration": {
    "firewall": {
      "filter": [
        {
          "name": "test",
          "term": [
            {
              "from": {
                "destination-address": [
                  {
                    "name": "100.64.0.0/24"
                  }
                ],
                "destination-port": [
                  "443"
                ],
                "source-address": [
                  {
                    "name": "10.0.0.0/8"
                  }
                ]
              },
              "name": "test-",
              "then": {
                "accept": [
                  null
                ],
                "count": "test-",
                "forwarding-class": "q1",
                "loss-priority": "low"
              }
            }
          ]
        }
      ]
    }
  }
}`
