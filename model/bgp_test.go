package model

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBGP(t *testing.T) {
	tt := []struct {
		name      string
		file      string
		netwInst  string
		length    int
		protocol  string
		ASN       int
		neighbor  string
		peergroup string
		localaddr string
		err       string
	}{
		{
			name:      "BGP 1",
			file:      "testdata/bgp.json",
			netwInst:  "ofce",
			length:    1,
			protocol:  "BGP",
			ASN:       65001,
			neighbor:  "10.99.99.1",
			peergroup: "customers",
			localaddr: "192.168.0.1",
		},
		{
			name:      "BGP 2",
			file:      "testdata/bgp.json",
			netwInst:  "ofce",
			length:    1,
			protocol:  "BGP",
			ASN:       65001,
			neighbor:  "10.99.99.1",
			peergroup: "customers",
			localaddr: "192.168.0.1",
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

			assert.Equal(t, tc.length, len(target.NetworkInstances.NetworkInstance))

			for name, ni := range target.NetworkInstances.NetworkInstance {
				assert.Equal(t, tc.netwInst, name)
				protocol, ok := ni.Protocols.Protocol[tc.protocol]
				if !ok {
					t.Errorf("can't unmarshal bgp config from %s", tc.file)
				}
				assert.Equal(t, tc.protocol, protocol.Name)

				assert.Equal(t, tc.ASN, protocol.BGP.Global.Config.As)

				neighbor, ok := protocol.BGP.Neighbors.Neighbor[tc.neighbor]
				if !ok {
					t.Errorf("can't unmarshal bgp neighbor config from %s", tc.file)
				}
				assert.Equal(t, tc.neighbor, neighbor.NeighborAddress)

				peergroup, ok := protocol.BGP.PeerGroups.PeerGroup[tc.peergroup]
				if !ok {
					t.Errorf("can't unmarshal bgp peergroup config from %s", tc.file)
				}
				assert.Equal(t, tc.peergroup, peergroup.PeerGroupName)
				assert.Equal(t, tc.localaddr, peergroup.Transport.Config.LocalAddress)

			}
		})
	}
}
