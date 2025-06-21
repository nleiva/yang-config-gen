package eos

import (
	"fmt"

	"github.com/nleiva/yang-config-gen/model"
	"github.com/nleiva/yang-data-structures/arista"
	"github.com/openconfig/ygot/ygot"
)

type Arista struct {
	SW        string
	root      *arista.Eos
	PrefixMod bool
}

func NewCompiler() Arista {
	a := new(arista.Eos)

	return Arista{
		SW:   "Arista",
		root: a,
	}

}

func (r Arista) GetPlatform() string {
	return r.SW
}

func (r Arista) CreateFullConfig(in model.Target) error {
	err := r.CreateInterfaceConfig(in)
	if err != nil {
		return fmt.Errorf("can't create Interface config: %w", err)
	}

	err = r.CreateNTPConfig(in)
	if err != nil {
		return fmt.Errorf("can't create NTP config: %w", err)
	}

	return nil
}

func (r Arista) CreateInterfaceConfig(in model.Target) error {
	for name, intf := range in.Interfaces.Interface {

		intfs := r.root.GetOrCreateInterface(name)

		data := intf.Subinterfaces.SubInterface["0"]
		ips := make([]model.AddressConfig, 0)
		for _, a := range data.IPv4.Addresses.Address {
			ips = append(ips, a.Config)
		}
		ip := ips[0].IP
		mask := ips[0].PrefixLength

		sub := intfs.GetOrCreateSubinterface(0)
		addr := sub.GetOrCreateIpv4().GetOrCreateAddress(ip)
		addr.Ip = ygot.String(ip)
		addr.PrefixLength = ygot.Uint8(uint8(mask))

	}
	return nil
}

func (r Arista) CreateNTPConfig(in model.Target) error {
	if in.System.NTP == nil {
		return nil
	}

	ntpSrv := in.System.NTP.Server

	if len(ntpSrv) < 1 {
		return nil
	}

	ntp := r.root.GetOrCreateSystem().GetOrCreateNtp()
	for _, ip := range ntpSrv {
		ntp.GetOrCreateServer(ip)
	}
	return nil
}

func (r Arista) EmitConfig(section string) (config string, err error) {
	var gs ygot.GoStruct
	switch section {
	case "system":
		gs = r.root.System
	default:
		gs = r.root
	}

	json, err := ygot.EmitJSON(gs, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: r.PrefixMod,
		},
	})

	if err != nil {
		return "", fmt.Errorf("can't emit JSON config: %w", err)
	}

	return json, nil
}
