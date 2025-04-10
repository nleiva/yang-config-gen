# YANG config generator

It translates network config parameters (as close to OpenConfig as possible) into vendor-specific YANG syntax that are network-device compatible. 

![Alt text](images/_yang-conf-gen_images_excalidraw.svg)

## Generate to router

Compile the code with `make build` or download an executable from [releases](https://github.com/nleiva/yang-config-gen/releases).

### Juniper

Translating OC'ish input to JunOS YANG (JSON).

#### Interfaces

```bash
$ ./bin/confgen_mac model/testdata/interface.json
{
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
}
```

#### BGP

```bash
$ ./bin/confgen_mac model/testdata/bgp.json 
{
  "configuration": {
    "routing-instances": {
      "instance": [
        {
          "name": "ofce",
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
                      "name": "10.99.99.1"
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
}
```

#### Policy options

```bash
$ ./bin/confgen_mac model/testdata/routingpolicy.json 
{
  "configuration": {
    "policy-options": {
      "community": [
        {
          "members": [
            "64496:509"
          ],
          "name": "r2d2"
        },
        {
          "members": [
            "64496:529"
          ],
          "name": "r2d2-b"
        },
        {
          "members": [
            "origin:100.0.0.0:1"
          ],
          "name": "wedge"
        },
        {
          "members": [
            "64496:519"
          ],
          "name": "yoda"
        },
        {
          "members": [
            "64496:539"
          ],
          "name": "yoda-b"
        }
      ],
      "policy-statement": [
        {
          "name": "amidala",
          "term": [
            {
              "from": {
                "route-filter": [
                  {
                    "address": "172.2.0.0/16",
                    "choice-ident": "orlonger",
                    "choice-value": ""
                  }
                ]
              },
              "name": "accept"
            }
          ],
          ....
```


## Parse JSON config from router

Save the following output to variable `routerLo0` like in the [junos_test.go](compiler/junos/junos_test.go) test file. I had to make the unit name a string instead of a number (*"got float64 type for field name, expect string"* error). 


```json
root@JunOS# show interfaces lo0 | display json 
{
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
                    "name" : 0, 
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
}
```

Then unmarshal it like this:

```go
  load := &junos.Junos{}
  if err := junos.Unmarshal([]byte(routerLo0), load); err != nil {
    t.Errorf("Can't unmarshal JSON: %v", err)
  }
```