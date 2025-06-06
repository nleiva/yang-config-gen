package model

type Interfaces struct {
	Interface map[string]Interface `json:"interface"`
}

type Interface struct {
	Name          string          `json:"name"`
	Config        InterfaceConfig `json:"config"`
	Subinterfaces SubInterfaces   `json:"subinterfaces"`
	Ethernet      Ethernet        `json:"ethernet"`
}

type InterfaceConfig struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
	MTU         int    `json:"mtu"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
	Bandwidth   int    `json:"bandwidth"`
}

type SubInterfaces struct {
	SubInterface map[string]SubInterface `json:"subinterface"`
}

type SubInterface struct {
	Config SubInterfaceConfig `json:"config"`
	Index  int                `json:"index"`
	IPv4   IPv4               `json:"ipv4"`
}

type SubInterfaceConfig struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
}

type IPv4 struct {
	Addresses Addresses  `json:"addresses"`
	Config    IPv4Config `json:"config"`
}

type Addresses struct {
	Address map[string]Address `json:"address"`
}

type Address struct {
	IP     string        `json:"ip"`
	Config AddressConfig `json:"config"`
}

type AddressConfig struct {
	IP           string `json:"ip"`
	PrefixLength int    `json:"prefix-length"`
}

type IPv4Config struct {
	Enabled bool `json:"enabled"`
}

type Ethernet struct {
	Config       EthernetConfig `json:"config"`
	SwitchedVLAN SwitchedVLAN   `json:"switched-vlan"`
}

type EthernetConfig struct {
	DuplexMode    string `json:"duplex-mode"`
	PortSpeed     int    `json:"port-speed"`
	AutoNegotiate bool   `json:"auto-negotiate"`
	AggregateID   int    `json:"aggregate-id"`
	Encapsulation string `json:"encapsulation"`
}

type SwitchedVLAN struct {
	Config SwitchedVLANConfig `json:"config"`
}

type SwitchedVLANConfig struct {
	InterfaceMode string `json:"interface-mode"`
	NativeVlan    int    `json:"native-vlan"`
	TrunkVlans    []int  `json:"trunk-vlans"`
}
