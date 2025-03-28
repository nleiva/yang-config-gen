# Generating and reading YANG

```bash
⇨  make test
go test ./... -v
?       github.com/nleiva/yang-config-gen/model [no test files]
=== RUN   TestCreateIntConfigText
=== RUN   TestCreateIntConfigText/config-lo0.0
=== PAUSE TestCreateIntConfigText/config-lo0.0
=== RUN   TestCreateIntConfigText/config-lo0
=== PAUSE TestCreateIntConfigText/config-lo0
=== CONT  TestCreateIntConfigText/config-lo0.0
=== CONT  TestCreateIntConfigText/config-lo0
--- PASS: TestCreateIntConfigText (0.00s)
    --- PASS: TestCreateIntConfigText/config-lo0.0 (0.08s)
    --- PASS: TestCreateIntConfigText/config-lo0 (0.08s)
=== RUN   TestReadJSONFromRouter
=== RUN   TestReadJSONFromRouter/config-lo0
--- PASS: TestReadJSONFromRouter (0.00s)
    --- PASS: TestReadJSONFromRouter/config-lo0 (0.00s)
PASS
```

## Juniper

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