package model

type NetworkInstances struct {
	NetworkInstance map[string]NetworkInstance `json:"network-instance"`
}

type NetworkInstance struct {
	Interfaces NtwInstInterfaces `json:"interfaces"`
	Protocols  NtwInstProtocols  `json:"protocols"`
}

type NtwInstInterfaces struct {
	Interface map[string]NtwInstInterface `json:"interface"`
}

type NtwInstInterface struct {
	Config NtwInstInterfaceConfig `json:"config"`
	ID     string                 `json:"id"`
}

type NtwInstInterfaceConfig struct {
	ID           string `json:"id"`
	Interface    string `json:"interface"`
	Subinterface int    `json:"subinterface"`
}

type NtwInstProtocols struct {
	Protocol map[string]NtwInstProtocol `json:"protocol"`
}

type NtwInstProtocol struct {
	ID     string                `json:"identifier"`
	Name   string                `json:"name"`
	Config NtwInstProtocolConfig `json:"config"`
	BGP    BGP                   `json:"BGP"`
	OSPFv2 OSPFv2                `json:"OSPF2"`
	Static STATIC                `json:"STATIC"`
}

type NtwInstProtocolConfig struct {
	ID      string `json:"identifier"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
