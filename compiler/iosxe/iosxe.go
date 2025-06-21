package iosxe

import (
	"fmt"

	"github.com/nleiva/yang-config-gen/compiler/common"
	"github.com/nleiva/yang-config-gen/model"
	"github.com/nleiva/yang-data-structures/cisco"
	"github.com/openconfig/ygot/ygot"
)

type Cisco struct {
	SW        string
	root      *cisco.Iosxe
	PrefixMod bool
}

func NewCompiler() Cisco {
	c := new(cisco.Iosxe)
	c.GetOrCreateNative()

	return Cisco{
		SW:   "Cisco",
		root: c,
	}

}

func (r Cisco) GetPlatform() string {
	return r.SW
}

func (r Cisco) CreateHostnameConfig(in model.Target) error {
	if in.Hostname != "" {
		r.root.GetOrCreateNative().Hostname = ygot.String(in.Hostname)
		return nil
	}
	return nil
}

func (r Cisco) CreateFullConfig(in model.Target) error {
	err := r.CreateHostnameConfig(in)
	if err != nil {
		return fmt.Errorf("can't create Hostname config: %w", err)
	}

	err = r.CreateNewInterfaceConfig(in)
	if err != nil {
		return fmt.Errorf("can't create Interface config: %w", err)
	}

	err = r.CreateNewNTPConfig(in)
	if err != nil {
		return fmt.Errorf("can't create NTP config: %w", err)
	}

	return nil
}

func (r Cisco) CreateNewInterfaceConfig(in model.Target) error {
	intfs := r.root.Native.GetOrCreateInterface()

	for name, intf := range in.Interfaces.Interface {
		_, port := common.SplitOnFirstNumber(name)

		data := intf.Subinterfaces.SubInterface["0"]
		ips := make([]model.AddressConfig, 0)
		for _, a := range data.IPv4.Addresses.Address {
			ips = append(ips, a.Config)
		}
		ip := ips[0].IP
		mask := ips[0].PrefixLength

		if intf.Config.Description == "LOOPBACK" {
			iface := intfs.GetOrCreateLoopback(0)
			addr := iface.GetOrCreateIp().GetOrCreateAddress().GetOrCreatePrimary()
			addr.Address = ygot.String(ip)
			addr.Mask = ygot.String(common.IntMasks[mask])
			continue
		}

		switch in.Platform {
		case "Cisco 9500":
			iface := intfs.GetOrCreateTwentyFiveGigE(port)
			addr := iface.GetOrCreateIp().GetOrCreateAddress().GetOrCreatePrimary()
			addr.Address = ygot.String(ip)
			addr.Mask = ygot.String(common.IntMasks[mask])
		case "Cisco 8000v":
			iface := intfs.GetOrCreateGigabitEthernet(port)
			addr := iface.GetOrCreateIp().GetOrCreateAddress().GetOrCreatePrimary()
			addr.Address = ygot.String(ip)
			addr.Mask = ygot.String(common.IntMasks[mask])
		default:
		}
	}
	return nil
}

func (r Cisco) CreateNewNTPConfig(in model.Target) error {
	if in.System.NTP == nil {
		return nil
	}

	ntpSrv := in.System.NTP.Server

	if len(ntpSrv) < 1 {
		return nil
	}

	ntp := r.root.GetOrCreateNative().GetOrCreateNtp()
	for _, ip := range ntpSrv {
		ntp.GetOrCreateServer().GetOrCreateServerList(ip)
	}
	return nil
}

func (r Cisco) EmitConfig() (config string, err error) {
	json, err := ygot.EmitJSON(r.root, &ygot.EmitJSONConfig{
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
