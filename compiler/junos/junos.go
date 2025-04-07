package junos

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"

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

	if in.RoutingPolicy != nil {
		err := r.CreateRoutingPolicyConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile routing-policy config: %w", err)
		}
		err = r.CreatePrefixListConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile prefix-list config: %w", err)
		}
		err = r.CreateCommunitiesConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile communities config: %w", err)
		}

	}

	if in.ACL != nil {
		err := r.CreateACLConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile ACL config: %w", err)
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
			neighbor := group.GetOrCreateNeighbor(bgpNeighbor.Config.NeighborAddress)
			peerASN := strconv.Itoa(bgpNeighbor.Config.PeerAs)
			neighbor.PeerAs = &peerASN
			description := bgpNeighbor.Config.Description
			if description != "" {
				neighbor.Description = &description
			}
		}

	}

	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing instances config: %s", err)
	}
	return nil
}

func (r Juniper) CreateRoutingOptsConfig(in model.Target) error {
	if len(in.NetworkInstances.NetworkInstance) < 1 {
		return nil
	}
	ri := r.root.Configuration.GetOrCreateRoutingInstances()

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

	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing options config: %w", err)
	}
	return nil
}

func (r Juniper) CreateRoutingPolicyConfig(in model.Target) error {
	if in.RoutingPolicy == nil {
		return nil
	}

	if len(in.RoutingPolicy.PolicyDefinitions.PolicyDefinition) < 1 {
		return nil
	}

	rp := r.root.Configuration.GetOrCreatePolicyOptions()

	for name, definition := range in.RoutingPolicy.PolicyDefinitions.PolicyDefinition {
		policy := rp.GetOrCreatePolicyStatement(name)

		for term, statement := range definition.Statements.Statement {
			t := policy.GetOrCreateTerm(term)
			if statement.Config.Seq != 0 {
				// Does JunOS ignore the sequence number?
			}
			if len(statement.Conditions.Config.CallPolicies) > 0 {
				from := t.GetOrCreateFrom()
				from.Policy = statement.Conditions.Config.CallPolicies
			}
			if statement.Conditions.Config.CommunitySet != "" {
				from := t.GetOrCreateFrom()
				from.Community = append(from.Community, statement.Conditions.Config.CommunitySet)
			}
			if statement.Conditions.BGPConditions.Config.CommunitySet != "" {
				from := t.GetOrCreateFrom()
				from.Community = append(from.Community, statement.Conditions.BGPConditions.Config.CommunitySet)
			}
			if statement.Conditions.BGPConditions.Config.ExtCommunitySet != "" {
				from := t.GetOrCreateFrom()
				from.Community = append(from.Community, statement.Conditions.BGPConditions.Config.ExtCommunitySet)
			}
			// From route-filter? -> rfc1918?

			then := t.GetOrCreateThen()
			switch statement.Actions.Config.PolicyResult {
			case "ACCEPT_ROUTE":
				then.Accept = true
			case "NEXT_ENTRY":
				then.Next = junos.JunosConfRoot_Configuration_PolicyOptions_PolicyStatement_Term_Then_Next_policy
			case "REJECT_ROUTE":
				then.Reject = true
			}

			if statement.Actions.BGPActions.Config.SetMed != 0 {
				then.GetOrCreateMetric().Metric = junos.UnionUint32(statement.Actions.BGPActions.Config.SetMed)
			}
			// Check this in the model
			//
			// if statement.Actions.BGPActions.SetExtCommunity.Reference.Config.ExtCommunitySetRef != "" {
			// 	then.GetOrCreateCommunity(
			// 		junos.JunosConfRoot_Configuration_PolicyOptions_PolicyStatement_Term_Then_Community_ChoiceIdent_add,
			// 		"add",
			// 		statement.Actions.BGPActions.SetExtCommunity.Reference.Config.ExtCommunitySetRef)
			// }

		}

	}

	err := rp.Validate()
	if err != nil {
		return fmt.Errorf("invalid policy definition inputs: %w", err)
	}
	return nil
}

func (r Juniper) CreatePrefixListConfig(in model.Target) error {
	if in.RoutingPolicy == nil {
		return nil
	}

	if len(in.RoutingPolicy.DefinedSets.PrefixSets.PrefixSet) < 1 {
		return nil
	}

	rp := r.root.Configuration.GetOrCreatePolicyOptions()

	for name, prefixset := range in.RoutingPolicy.DefinedSets.PrefixSets.PrefixSet {
		policy := rp.GetOrCreatePolicyStatement(name)

		for _, statement := range prefixset.Prefixes.Prefix {
			t := policy.GetOrCreateTerm("accept")
			from := t.GetOrCreateFrom()
			orlonger := junos.JunosConfRoot_Configuration_PolicyOptions_PolicyStatement_From_RouteFilter_ChoiceIdent_orlonger
			from.GetOrCreateRouteFilter(statement.IPPrefix, orlonger, "")
		}

		policy.GetOrCreateThen().Accept = true
	}

	err := rp.Validate()
	if err != nil {
		return fmt.Errorf("invalid prefix set inputs: %w", err)
	}
	return nil
}

func (r Juniper) CreateCommunitiesConfig(in model.Target) error {
	if in.RoutingPolicy == nil {
		return nil
	}

	if len(in.RoutingPolicy.DefinedSets.BGPDefinedSets.CommunitySets.CommunitySet) < 1 {
		return nil
	}

	rp := r.root.Configuration.GetOrCreatePolicyOptions()

	for name, commset := range in.RoutingPolicy.DefinedSets.BGPDefinedSets.CommunitySets.CommunitySet {
		c := rp.GetOrCreateCommunity(name)
		c.Members = append(c.Members, commset.Config.CommunityMember...)
	}

	err := rp.Validate()
	if err != nil {
		return fmt.Errorf("invalid prefix set inputs: %w", err)
	}
	return nil
}

func (r Juniper) CreateACLConfig(in model.Target) error {
	if in.ACL == nil {
		return nil
	}

	if len(in.ACL.ACLSets.ACLSet) < 1 {
		return nil
	}
	fw := r.root.Configuration.GetOrCreateFirewall()

	for name, aclset := range in.ACL.ACLSets.ACLSet {

		if len(aclset.ACLEntries.ACLEntry) < 1 {
			continue
		}
		names := strings.Split(name, "-")
		postype := names[0]

		filter := fw.GetOrCreateFilter(name)
		for _, acl := range aclset.ACLEntries.ACLEntry {
			term := postype + "-" + acl.Config.Description
			t := filter.GetOrCreateTerm(term)
			from := t.GetOrCreateFrom()
			then := t.GetOrCreateThen()

			for _, addr := range acl.IPv4.Config.SourceAddresses {
				from.GetOrCreateSourceAddress(addr)
			}
			for _, addr := range acl.IPv4.Config.DestinationAddresses {
				from.GetOrCreateDestinationAddress(addr)
			}

			if acl.Transport.Config.DestinationPort != "" {
				ports := strings.Split(acl.Transport.Config.DestinationPort, "..")
				switch len(ports) > 1 {
				case true:
					from.DestinationPort = []string{ports[0] + "-" + ports[1]}
				default:
					from.DestinationPort = []string{acl.Transport.Config.DestinationPort}
				}
			}

			if acl.IPv4.Config.Protocol != "" {
				switch acl.IPv4.Config.Protocol {
				case "IP_UDP":
					from.Protocol = []string{"udp"}
				case "IP_TCP":
					from.Protocol = []string{"tcp"}
				case "IP_ICMP":
					icmmpMap := map[string]string{
						"ECHO_REPLY":      "echo-reply",
						"DST_UNREACHABLE": "unreachable",
						"ECHO":            "echo-request",
						"TIME_EXCEEDED":   "time-exceeded",
					}
					from.Protocol = []string{"icmp"}
					for _, msg := range acl.IPv4.ICMPv4.Config.Types {
						from.IcmpType = append(from.IcmpType, icmmpMap[msg])
					}
				}
			}

			switch acl.Actions.Config.ForwardingAction {
			case "ACCEPT":
				then.Accept = true
			case "DROP":
				then.GetOrCreateDiscard()
			}

			then.Count = &term

			switch acl.Actions.Config.TargetGroup {
			// match-dscp forwarding-class ef loss-priority low code-points ef
			case "ef":
				fc := "ef"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class be loss-priority high code-points be
			case "be":
				fc := "be"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_medium_low
			// match-dscp forwarding-class q1 loss-priority low code-points cs2
			case "cs2":
				fc := "q1"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class q1 loss-priority high code-points cs1
			// match-dscp forwarding-class q1 loss-priority high code-points af11
			case "af11", "cs1":
				fc := "q1"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_medium_low
			// match-dscp forwarding-class q3 loss-priority high code-points cs3
			// match-dscp forwarding-class q3 loss-priority high code-points af31
			case "cs3", "af31":
				fc := "q3"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_medium_low
			// match-dscp forwarding-class q3 loss-priority low code-points cs5
			case "cs5":
				fc := "q3"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class q4 loss-priority low code-points cs4
			case "cs4":
				fc := "q4"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class q4 loss-priority high code-points af41
			// match-dscp forwarding-class q4 loss-priority high code-points af42
			case "af41", "af42":
				fc := "q4"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_medium_low
			// match-dscp forwarding-class sc loss-priority low code-points cs6
			case "cs6":
				fc := "sc"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class nc loss-priority low code-points cs7
			case "cs7":
				fc := "nc"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class q2 loss-priority low code-points af21
			case "af21":
				fc := "q2"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_low
			// match-dscp forwarding-class q2 loss-priority high code-points af22
			case "af22":
				fc := "q2"
				then.ForwardingClass = &fc
				then.LossPriority = junos.JunosConfRoot_Configuration_Firewall_Family_Inet_Filter_Term_Then_LossPriority_medium_low
			default:
			}

		}

	}

	for name, iface := range in.ACL.Interfaces.Interface {
		names := strings.Split(name, ".")
		var unit string
		if len(names) > 1 {
			unit = names[1]
		}

		intf := r.root.Configuration.GetOrCreateInterfaces().GetOrCreateInterface(names[0])

		filter := intf.GetOrCreateUnit(unit).GetOrCreateFamily().GetOrCreateInet().GetOrCreateFilter()
		for _, out := range iface.EgressACLSets.EgressACLSet {
			filter.GetOrCreateOutput().FilterName = &out.SetName
		}
		for _, in := range iface.IngressACLSets.IngressACLSet {
			filter.GetOrCreateInput().FilterName = &in.SetName
		}
	}

	err := fw.Validate()
	if err != nil {
		return fmt.Errorf("invalid prefix set inputs: %w", err)
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
