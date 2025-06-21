package model

type Target struct {
	Hostname         string
	Platform         string
	Vendor           string
	Site             string
	Role             string
	Model            string
	ASN              int64
	ID               int32
	Interfaces       Interfaces       `json:"interfaces"`
	System           System           `json:"system"`
	NetworkInstances NetworkInstances `json:"network-instances"`
	RoutingPolicy    *RoutingPolicy   `json:"routing-policy"`
	ACL              *ACL             `json:"acl"`
}
