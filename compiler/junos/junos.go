package junos

import (
	"fmt"
	"net/netip"
	"strconv"

	"github.com/nleiva/yang-config-gen/model"
	"github.com/nleiva/yang-data-structures/junos"
	"github.com/openconfig/ygot/ygot"
)

type Juniper struct {
	SW   string
	root *junos.Junos
}

func NewCompiler() Juniper {
	j := new(junos.Junos)
	j.GetOrCreateConfiguration()

	return Juniper{
		SW:   "Junos",
		root: j,
	}

}
func (r Juniper) GetPlatform() string {
	return r.SW
}

func (r Juniper) CompileConfig(in model.Target) error {
	// Compile interfaces config
	if len(in.Interfaces.Interface) > 0 {
		err := r.CreateIntConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile interface config: %w", err)
		}
	}

	// Compile NetworkInstances config
	ni := in.NetworkInstances.NetworkInstance
	if len(ni) > 0 {
		err := r.CreateRoutingInstancesConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile routing instance config: %w", err)
		}

		err = r.CreateBGPConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile BGP config: %w", err)
		}

		err = r.CreateRoutingOptsConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile routing options config: %w", err)
		}
	}

	return nil
}

func (r Juniper) CreateIntConfig(in model.Target) error {
	if len(in.Interfaces.Interface) < 1 {
		return nil
	}

	// Create interface root config
	intfs := r.root.Configuration.GetOrCreateInterfaces()

	for name, iface := range in.Interfaces.Interface {
		if name == "" {
			return fmt.Errorf("can't create interface without a name")
		}
		intf := intfs.GetOrCreateInterface(name)

		cfg := iface.Config
		// if !cfg.Enabled {
		// 	intf.Disable = true
		// }
		if cfg.MTU != 0 {
			intf.Mtu = junos.UnionUint32(iface.Config.MTU)
		}

		eth := iface.Ethernet
		if eth.SwitchedVLAN.Config.NativeVlan != 0 {
			intf.NativeVlanId = junos.UnionUint32(eth.SwitchedVLAN.Config.NativeVlan)
		}

		ethCfg := eth.Config
		switch ethCfg.PortSpeed {
		case 1000:
			intf.Speed = junos.JunosConfRoot_Configuration_Interfaces_Interface_Speed_1g
		case 10000:
			intf.Speed = junos.JunosConfRoot_Configuration_Interfaces_Interface_Speed_10g
		case 100000:
			intf.Speed = junos.JunosConfRoot_Configuration_Interfaces_Interface_Speed_100g
		default:
			intf.Speed = junos.JunosConfRoot_Configuration_Interfaces_Interface_Speed_UNSET
		}

		if ethCfg.DuplexMode != "" {
			switch ethCfg.DuplexMode {
			case "FULL":
				intf.LinkMode = junos.JunosConfRoot_Configuration_Interfaces_Interface_LinkMode_full_duplex
			default:
				intf.LinkMode = junos.JunosConfRoot_Configuration_Interfaces_Interface_LinkMode_automatic
			}
		}

		if ethCfg.Encapsulation != "" {
			switch ethCfg.Encapsulation {
			case "dot1q":
				intf.Encapsulation = junos.JunosConfRoot_Configuration_Interfaces_Interface_Encapsulation_vlan_ccc
			default:
				intf.Encapsulation = junos.JunosConfRoot_Configuration_Interfaces_Interface_Encapsulation_ethernet
			}
		}

		if len(iface.Subinterfaces.SubInterface) < 1 {
			if iface.Config.Description != "" {
				intf.Description = &iface.Config.Description
			}
			continue
		}

		for idx, subintf := range iface.Subinterfaces.SubInterface {

			unit, err := intf.NewUnit(idx)
			if err != nil {
				return fmt.Errorf("can't compile interface unit: %w", err)
			}

			if len(eth.SwitchedVLAN.Config.TrunkVlans) != 0 {
				for _, vlan := range eth.SwitchedVLAN.Config.TrunkVlans {
					unit.VlanIdList = append(unit.VlanIdList, strconv.Itoa(vlan))
				}
			}

			if subintf.Config.Description != "" {
				unit.Description = &subintf.Config.Description
			}
			if len(subintf.IPv4.Addresses.Address) > 0 {
				for _, address := range subintf.IPv4.Addresses.Address {
					addr, err := netip.ParseAddr(address.Config.IP)
					if err != nil {
						return fmt.Errorf("invalid IP address %v: %w", address.Config.IP, err)
					}
					prefix := netip.PrefixFrom(addr, address.Config.PrefixLength)
					unit.GetOrCreateFamily().GetOrCreateInet().GetOrCreateAddress(prefix.String())
				}
			}
		}
	}

	err := intfs.Validate()
	if err != nil {
		return fmt.Errorf("invalid interface config: %w", err)
	}

	return nil
}

func (r Juniper) CreateRoutingInstancesConfig(in model.Target) error {
	// Create routing instance root config
	ri := r.root.Configuration.GetOrCreateRoutingInstances()

	if len(in.NetworkInstances.NetworkInstance) > 0 {
		for name, ni := range in.NetworkInstances.NetworkInstance {
			instance := ri.GetOrCreateInstance(name)

			if len(ni.Interfaces.Interface) > 0 {
				for iname, _ := range ni.Interfaces.Interface {
					_ = instance.GetOrCreateInterface(iname)
				}
			}
		}
	}
	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing instances config: %s", err)
	}
	return nil
}

func (r Juniper) CreateBGPConfig(in model.Target) error {
	if len(in.NetworkInstances.NetworkInstance) < 1 {
		return nil
	}

	ri := r.root.Configuration.GetOrCreateRoutingInstances()

	for name, ni := range in.NetworkInstances.NetworkInstance {
		if len(ni.Protocols.Protocol) < 1 {
			return nil
		}
		protocol, ok := ni.Protocols.Protocol["BGP"]
		if !ok {
			continue
		}

		for groupname, bgpGroup := range protocol.BGP.PeerGroups.PeerGroup {
			instance := ri.GetOrCreateInstance(name)

			bgp := instance.GetOrCreateProtocols().GetOrCreateBgp()
			group := bgp.GetOrCreateGroup(groupname)
			group.Export = bgpGroup.ApplyPolicy.Config.ExportPolicy
			group.Import = bgpGroup.ApplyPolicy.Config.ImportPolicy

		}

		for _, bgpNeighbor := range protocol.BGP.Neighbors.Neighbor {
			instance := ri.GetOrCreateInstance(name)

			bgp := instance.GetOrCreateProtocols().GetOrCreateBgp()

			group := bgp.GetOrCreateGroup(bgpNeighbor.Config.PeerGroup)
			group.GetOrCreateNeighbor(bgpNeighbor.Config.NeighborAddress)
		}

	}

	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing instances config: %s", err)
	}
	return nil
}

func (r Juniper) CreateRoutingOptsConfig(in model.Target) error {
	ri := r.root.Configuration.GetOrCreateRoutingInstances()

	if len(in.NetworkInstances.NetworkInstance) > 0 {
		for name, netInst := range in.NetworkInstances.NetworkInstance {
			instance := ri.GetOrCreateInstance(name)

			protocol, ok := netInst.Protocols.Protocol["BGP"]
			if !ok {
				continue
			}

			ASN := protocol.BGP.Global.Config.As

			if ASN != 0 {
				a := strconv.Itoa(protocol.BGP.Global.Config.As)
				instance.GetOrCreateRoutingOptions().GetOrCreateAutonomousSystem().AsNumber = &a
			}

			if len(netInst.Interfaces.Interface) > 0 {
				for iname, _ := range netInst.Interfaces.Interface {
					_ = instance.GetOrCreateInterface(iname)
				}
			}
		}
	}

	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing options config: %w", err)
	}
	return nil
}

func (r Juniper) EmitConfig() (config string, err error) {
	json, err := ygot.EmitJSON(r.root, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		},
	})

	if err != nil {
		return "", fmt.Errorf("can't emit JSON config: %w", err)
	}

	return json, nil
}
