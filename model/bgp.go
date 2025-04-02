package model

type BGP struct {
	Global     Global     `json:"global"`
	PeerGroups PeerGroups `json:"peer-groups"`
	Neighbors  Neighbors  `json:"neighbors"`
}

type Global struct {
	Config           GlobalConfig     `json:"config"`
	GracefulRestart  GracefulRestart  `json:"graceful-restart"`
	PolicyDefault    PolicyDefault    `json:"policy_default"`
	UseMultiplePaths UseMultiplePaths `json:"use-multiple-paths"`
}

type GlobalConfig struct {
	As              int  `json:"as"`
	PeerIPTrack     bool `json:"peer-ip-track"`
	RapidWithdrawal bool `json:"rapid-withdrawal"`
}

type GracefulRestart struct {
	Config EnableConfig `json:"config"`
}

type PolicyDefault struct {
	EBGP EBGPPolicyDefault `json:"ebgp"`
}

type EBGPPolicyDefault struct {
	Config DefaultBGPPolicyConfig `json:"config"`
}

type DefaultBGPPolicyConfig struct {
	ExportDefaultReject bool `json:"export-default-reject"`
	ImportDefaultReject bool `json:"import-default-reject"`
}

type UseMultiplePaths struct {
	Config EnableConfig        `json:"config"`
	EBGP   UseMultiplePathsBGP `json:"ebgp"`
	IBGP   UseMultiplePathsBGP `json:"ibgp"`
}

type UseMultiplePathsBGP struct {
	Config MultiplePathsConfig `json:"config"`
}

type MultiplePathsConfig struct {
	AllowMultipleAs bool `json:"allow-multiple-as"`
	MaximumPaths    int  `json:"maximum-paths"`
}

type PeerGroups struct {
	PeerGroup map[string]PeerGroup `json:"peer-group"`
}

type PeerGroup struct {
	AfiSafis         AfiSafis         `json:"afi-safis"`
	ApplyPolicy      ApplyPolicy      `json:"apply-policy"`
	ASPathOptions    ASPathOptions    `json:"as-path-options"`
	Config           PeerGroupConfig  `json:"config"`
	PeerGroupName    string           `json:"peer-group-name"`
	Timers           Timers           `json:"timers"`
	Transport        Transport        `json:"transport"`
	UseMultiplePaths UseMultiplePaths `json:"use-multiple-paths"`
}

type ApplyPolicy struct {
	Config BGPPolicyConfig `json:"config"`
}

type ASPathOptions struct {
	Config ASPathOptionsConfig `json:"config"`
}

type Timers struct {
	Config TimersConfig `json:"config"`
}

type Transport struct {
	Config TransportConfig `json:"config"`
}

type BGPPolicyConfig struct {
	ExportPolicy []string `json:"export-policy"`
	ImportPolicy []string `json:"import-policy"`
}

type TimersConfig struct {
	HoldTime          int `json:"hold-time"`
	KeepaliveInterval int `json:"keepalive-interval"`
}

type TransportConfig struct {
	LocalAddress string `json:"local-address"`
}

type ASPathOptionsConfig struct {
	ReplacePeerAs bool `json:"replace-peer-as"`
}

type PeerGroupConfig struct {
	Description      string `json:"description"`
	Enabled          bool   `json:"enabled"`
	LocalAs          int    `json:"local-as"`
	NextHopSelf      bool   `json:"next-hop-self"`
	NextHopUnchanged bool   `json:"next-hop-unchanged"`
	PeerAs           int    `json:"peer-as"`
	PeerGroupName    string `json:"peer-group-name"`
	PeerType         string `json:"peer-type"`
	RemovePrivateAs  string `json:"remove-private-as"`
	SplitHorizon     bool   `json:"split-horizon"`
}

type AfiSafis struct {
	AfiSafi map[string]AfiSafi `json:"afi-safi"`
}

type AfiSafi struct {
	AfiSafiName string        `json:"afi-safi-name"`
	Config      AfiSafiConfig `json:"config"`
	IPv4Unicast IPv4Unicast   `json:"IPV4_UNICAST"`
}

type AfiSafiConfig struct {
	AfiSafiName string `json:"afi-safi-name"`
	Enabled     bool   `json:"enabled"`
}

type IPv4Unicast struct {
	PrefixLimit IPv4UnicastPrefixLimit `json:"prefix-limit"`
}

type IPv4UnicastPrefixLimit struct {
	Config PrefixLimitConfig `json:"config"`
}

type PrefixLimitConfig struct {
	MaxPrefixes         int  `json:"max-prefixes"`
	PreventTeardown     bool `json:"prevent-teardown"`
	WarningThresholdPct int  `json:"warning-threshold-pct"`
}

type Neighbors struct {
	Neighbor map[string]Neighbor `json:"neighbor"`
}

type Neighbor struct {
	NeighborAddress string         `json:"neighbor-address"`
	Config          NeighborConfig `json:"config"`
	EnableBFD       EnableBFD      `json:"enable-bfd"`
}

type EnableBFD struct {
	Config EnableConfig `json:"config"`
}

type EnableConfig struct {
	Enabled bool `json:"enabled"`
}

type NeighborConfig struct {
	NeighborAddress string `json:"neighbor-address"`
	PeerAs          int    `json:"peer-as"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	PeerGroup       string `json:"peer-group"`
}
