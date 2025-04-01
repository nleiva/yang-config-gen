# Generating and reading YANG

## Generate to router

### Juniper

Translating OC'ish input to JunOS YANG.

#### Interfaces

```bash
$ go run main.go ../model/testdata/interface.json
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
$ go run main.go ../model/testdata/bgp.json 
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

## Parse JSON config from router

I removed YANG annotations from the device output to give to `ygot` Unmarshal method. I also made the unit name a string instead of a number.
1. it doesnt't like "@" : {}
2. got float64 type for field name, expect string

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