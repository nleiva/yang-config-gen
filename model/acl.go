package model

type ACL struct {
	ACLSets    ACLSets       `json:"acl-sets"`
	Interfaces ACLInterfaces `json:"interfaces"`
}

type ACLSets struct {
	ACLSet map[string]ACLSet `json:"acl-set"`
}

type ACLSet struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Config     ACLSetConfig `json:"config"`
	ACLEntries ACLEntries   `json:"acl-entries"`
}

type ACLSetConfig struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	DefaultAction string `json:"default-action"`
	Direction     string `json:"direction"`
	ReplaceConfig bool   `json:"replace-config"`
}

type ACLEntries struct {
	ACLEntry map[string]ACLEntry `json:"acl-entry"`
}

type ACLEntry struct {
	SequenceID int            `json:"sequence-id"`
	Config     ACLEntryConfig `json:"config"`
	IPv4       ACLIPv4        `json:"ipv4"`
	Actions    ACLActions     `json:"actions"`
	Transport  ACLTransport   `json:"transport"`
}

type ACLActions struct {
	Config ACLActionsConfig `json:"config"`
}

type ACLIPv4 struct {
	Config ACLIPv4Config `json:"config"`
	ICMPv4 ICMPv4        `json:"icmpv4"`
}

type ICMPv4 struct {
	Config ICMPv4Config `json:"config"`
}

type ICMPv4Config struct {
	Types []string "types"
}

type ACLIPv4Config struct {
	SourceAddresses      []string `json:"source-addresses"`
	DestinationAddresses []string `json:"destination-addresses"`
	DSCP                 int      `json:"dscp"`
	Protocol             string   `json:"protocol"`
}

type ACLActionsConfig struct {
	ForwardingAction string `json:"forwarding-action"`
	TargetGroup      string `json:"target-group"`
}

type ACLTransport struct {
	Config ACLTransportConfig `json:"config"`
}

type ACLTransportConfig struct {
	SourcePort      string `json:"source-port"`
	DetailMode      string `json:"detail-mode"`
	BuiltinDetail   string `json:"builtin-detail"`
	DestinationPort string `json:"destination-port"`
}

type ACLEntryConfig struct {
	SequenceID  int    `json:"sequence-id"`
	Description string `json:"description"`
}

type ACLInterfaces struct {
	Interface map[string]ACLInterface `json:"interface"`
}

type ACLInterface struct {
	ID             string             `json:"id"`
	Config         ACLInterfaceConfig `json:"config"`
	IngressACLSets IngressACLSets     `json:"ingress-acl-sets"`
	EgressACLSets  EgressACLSets      `json:"egress-acl-sets"`
}

type ACLInterfaceConfig struct {
	ID string `json:"id"`
}

type IngressACLSets struct {
	IngressACLSet map[string]IngressACLSet `json:"ingress-acl-set"`
}

type IngressACLSet struct {
	SetName string              `json:"set-name"`
	Type    string              `json:"type"`
	Config  IngressACLSetConfig `json:"config"`
}

type IngressACLSetConfig struct {
	SetName string `json:"set-name"`
	Type    string `json:"type"`
}

type EgressACLSets struct {
	EgressACLSet map[string]EgressACLSet `json:"egress-acl-set"`
}

type EgressACLSet struct {
	SetName string             `json:"set-name"`
	Type    string             `json:"type"`
	Config  EgressACLSetConfig `json:"config"`
}

type EgressACLSetConfig struct {
	SetName string `json:"set-name"`
	Type    string `json:"type"`
}
