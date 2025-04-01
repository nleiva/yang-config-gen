package model

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInterface(t *testing.T) {
	tt := []struct {
		name       string
		file       string
		intf       string
		length     int
		L2         bool
		mtu        int
		encap      string
		trunkVLANs []int
		L3         bool
		ip         string
		mask       int
		subint     string

		err string
	}{
		{
			name:       "Interface 1",
			file:       "testdata/interface.json",
			intf:       "eth-0/1/1",
			length:     2,
			L2:         true,
			mtu:        1500,
			encap:      "dot1q",
			trunkVLANs: []int{10, 20, 40},
		},
		{
			name:   "Interface 2",
			file:   "testdata/interface.json",
			intf:   "ge-2/2/0",
			length: 2,
			L3:     true,
			subint: "10",
			ip:     "10.10.10.1",
			mask:   24,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Read source data from file
			file, err := os.Open(tc.file)
			if err != nil {
				t.Errorf("can't open testdata file %s: %s", tc.file, err.Error())
			}
			defer file.Close()

			target := new(Target)

			err = ReadData(file, target)
			if err != nil {
				t.Errorf("can't unmarshal data from %s: %s", tc.file, err.Error())
			}

			assert.Equal(t, tc.length, len(target.Interfaces.Interface))
			assert.Equal(t, tc.intf, target.Interfaces.Interface[tc.intf].Name)

			if tc.L2 {
				iface := target.Interfaces.Interface[tc.intf]
				assert.Equal(t, tc.mtu, iface.Config.MTU)
				assert.Equal(t, tc.trunkVLANs, iface.Ethernet.SwitchedVLAN.Config.TrunkVlans)
				assert.Equal(t, tc.encap, iface.Ethernet.Config.Encapsulation)
			}

			if tc.L3 {
				ipv4 := target.Interfaces.Interface[tc.intf].Subinterfaces.SubInterface[tc.subint].IPv4
				assert.Equal(t, tc.ip, ipv4.Addresses.Address[tc.ip].IP)
				assert.Equal(t, tc.mask, ipv4.Addresses.Address[tc.ip].Config.PrefixLength)
			}
		})
	}
}
