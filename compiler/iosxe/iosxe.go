package iosxe

import (
	"fmt"

	"github.com/nleiva/yang-config-gen/model"
	"github.com/nleiva/yang-data-structures/iosxe"
	"github.com/openconfig/ygot/ygot"
)

type Cisco struct {
	SW   string
	root *iosxe.Iosxe
}

func NewCompiler() Cisco {
	c := new(iosxe.Iosxe)
	c.GetOrCreateNative()

	return Cisco{
		SW:   "Cisco",
		root: c,
	}

}
func (r Cisco) GetPlatform() string {
	return r.SW
}

func (r Cisco) CompileConfig(in model.Target) error {
	// Compile interfaces config
	if in.Hostname != "" {
		err := r.CreateHostnameConfig(in)
		if err != nil {
			return fmt.Errorf("can't compile interface config: %w", err)
		}
	}
	return nil
}

func (r Cisco) CreateHostnameConfig(in model.Target) error {
	if in.Hostname == "" {
		return nil
	}

	r.root.Native.Hostname = ygot.String(in.Hostname)
	return nil
}

func (r Cisco) EmitConfig() (config string, err error) {
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
