package model

import (
	"net/netip"
)

type Target struct {
	Hostname    string
	Platform    string
	Vendor      string
	Site        string
	Role        string
	Model       string
	ASN         int64
	ID          int32
	Interfaces  []Interface
	BGPSessions []BGPSession
}

type Interface struct {
	Name        string
	Unit        string
	Display     string
	Description string
	Device      string
	Address     netip.Prefix
}

type BGPSession struct {
	LocalAddress netip.Prefix
	LocalAs      int64
	VRF          string
	Group        string
	Neighbor     netip.Addr
	PeerAS       int64
	Status       string
	Interface    Interface
}