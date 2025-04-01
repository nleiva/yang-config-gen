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
	Config GracefulRestartConfig `json:"config"`
}

type GracefulRestartConfig struct {
	Enabled bool `json:"enabled"`
}

type PolicyDefault struct {
	Ebgp struct {
		Config DefaultBGPPolicyConfig `json:"config"`
	} `json:"ebgp"`
}

type DefaultBGPPolicyConfig struct {
	ExportDefaultReject bool `json:"export-default-reject"`
	ImportDefaultReject bool `json:"import-default-reject"`
}

type UseMultiplePaths struct {
	Config struct {
		Enabled bool `json:"enabled"`
	} `json:"config"`
	Ebgp struct {
		Config MultiplePathsConfig `json:"config"`
	} `json:"ebgp"`
	Ibgp struct {
		Config MultiplePathsConfig `json:"config"`
	} `json:"ibgp"`
}

type MultiplePathsConfig struct {
	AllowMultipleAs bool `json:"allow-multiple-as"`
	MaximumPaths    int  `json:"maximum-paths"`
}

type PeerGroups struct {
	PeerGroup map[string]PeerGroup `json:"peer-group"`
}

type PeerGroup struct {
	AfiSafis    AfiSafis `json:"afi-safis"`
	ApplyPolicy struct {
		Config BGPPolicyConfig `json:"config"`
	} `json:"apply-policy"`
	ASPathOptions struct {
		Config ASPathOptionsConfig `json:"config"`
	} `json:"as-path-options"`
	Config        PeerGroupConfig `json:"config"`
	PeerGroupName string          `json:"peer-group-name"`
	Timers        struct {
		Config TimersConfig `json:"config"`
	} `json:"timers"`
	Transport struct {
		Config TransportConfig `json:"config"`
	} `json:"transport"`
	UseMultiplePaths UseMultiplePaths `json:"use-multiple-paths"`
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
	PrefixLimit struct {
		Config PrefixLimitConfig `json:"config"`
	} `json:"prefix-limit"`
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
	EnableBFD       struct {
		Config EnableBFDConfig `json:"config"`
	} `json:"enable-bfd"`
}

type EnableBFDConfig struct {
	Enabled bool `json:"enabled"`
}

type NeighborConfig struct {
	NeighborAddress string `json:"neighbor-address"`
	PeerAs          int    `json:"peer-as"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	PeerGroup       string `json:"peer-group"`
}
