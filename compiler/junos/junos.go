package junos

import (
	"fmt"
	"net/netip"
	"strconv"

	"github.com/nleiva/yang-data-structures/junos"
	"github.com/nleiva/yang-config-gen/model"
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
	if len(in.Interfaces) > 0 {
		err := r.CreateIntConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile interface config: %w", err)
		}
	}

	// Compile BGP config
	if len(in.BGPSessions) > 0 {
		err := r.CreateRoutingInstancesConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile routing instance config: %w", err)
		}

		err = r.CreateRoutingOptsConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile routing options config: %w", err)
		}
	}

	return nil
}

func (r Juniper) CreateIntConfig(in model.Target) error {
	// Create interface root config
	intfs := r.root.Configuration.GetOrCreateInterfaces()

	for _, iface := range in.Interfaces {
		if iface.Name == "" {
			return fmt.Errorf("can't create interface without a name")
		}
		intf := intfs.GetOrCreateInterface(iface.Name)

		if iface.Unit == "" {
			if iface.Description != "" {
				intf.Description = &iface.Description
			}
			continue
		}
		unit, err := intf.NewUnit(iface.Unit)
		if err != nil {
			return fmt.Errorf("can't compile interface unit: %w", err)
		}

		if iface.Description != "" {
			unit.Description = &iface.Description
		}
		if iface.Address.IsValid() {
			unit.GetOrCreateFamily().GetOrCreateInet().GetOrCreateAddress(iface.Address.String())
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

	for _, bgpSession := range in.BGPSessions {
		bgp := ri.GetOrCreateInstance(bgpSession.VRF).GetOrCreateProtocols().GetOrCreateBgp()
		group := bgp.GetOrCreateGroup(bgpSession.Group)
		neighbor := group.GetOrCreateNeighbor(bgpSession.Neighbor.String())

		localIP, err := netip.ParsePrefix(bgpSession.LocalAddress.String())
		if err != nil {
			return err
		}

		localIPStr := localIP.Addr().String()

		neighbor.LocalAddress = &localIPStr
		PeerAS := strconv.FormatInt(bgpSession.PeerAS, 10)
		neighbor.PeerAs = &PeerAS

		if bgpSession.Status == "offline" {
			neighbor.GetOrCreateShutdown()
		}
	}

	err := ri.Validate()
	if err != nil {
		return fmt.Errorf("invalid routing instances config: %s", err)
	}
	return nil
}

func (r Juniper) CreateRoutingOptsConfig(in model.Target) error {
	as := r.root.Configuration.GetOrCreateRoutingOptions().GetOrCreateAutonomousSystem()
	ASN := strconv.FormatInt(in.ASN, 10)
	as.AsNumber = &ASN

	err := as.Validate()
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