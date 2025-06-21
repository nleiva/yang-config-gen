package model

type System struct {
	AAA         *AAA     `json:"aaa,omitzero"`
	DNS         *DNS     `json:"dns,omitzero"`
	DomainName  string   `json:"domain-name,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
	Logging     *Logging `json:"logging,omitzero"`
	LoginBanner string   `json:"login-banner,omitempty"`
	MotdBanner  string   `json:"motd-banner,omitempty"`
	NTP         *NTP     `json:"ntp,omitzero"`
	SSHServer   *SSH     `json:"ssh-server,omitzero"`
}

type AAA struct {
	// Accounting     *AAA_Accounting     `json:"accounting"`
	// Authentication *AAA_Authentication `json:"authentication"`
	// Authorization  *AAA_Authorization  `json:"authorization"`
	ServerGroup map[string]string `json:"server-groups"`
}

type DNS struct {
	HostEntry map[string]string `json:"host-entries"`
	Search    []string          `json:"search"`
	Server    map[string]string `json:"servers"`
}
type Logging struct {
	Console      bool              `json:"console"`
	RemoteServer map[string]string `json:"remote-servers"`
}

type NTP struct {
	AuthMismatch  uint64            `json:"auth-mismatch,omitempty"`
	EnableNtpAuth bool              `json:"enable-ntp-auth,omitempty"`
	Enabled       bool              `json:"enabled,omitempty"`
	NTPKey        map[uint16]string `json:"ntp-keys,omitempty"`
	Server        map[string]string `json:"servers,omitempty"`
}

type SSH struct {
	Enable          bool   `json:"enable"`
	ProtocolVersion int    `json:"protocol-version"`
	RateLimit       uint16 `json:"rate-limit"`
	SessionLimit    uint16 `json:"session-limit"`
	Timeout         uint16 `json:"timeout"`
}
