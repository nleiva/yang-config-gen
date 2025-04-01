package model

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRoutingPolicy(t *testing.T) {
	tt := []struct {
		name string
		file string

		definitions int
		prefixsets  int
		commsets    int
		err         string
	}{
		{
			name:        "Routing Policy 1",
			file:        "testdata/routingpolicy.json",
			definitions: 4,
			prefixsets:  13,
			commsets:    5,
		},
		{
			name:        "Routing Policy 2",
			file:        "testdata/routingpolicy.json",
			definitions: 4,
			prefixsets:  13,
			commsets:    5,
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

			assert.Equal(t, tc.prefixsets, len(target.RoutingPolicy.DefinedSets.PrefixSets.PrefixSet))
			assert.Equal(t, tc.commsets, len(target.RoutingPolicy.DefinedSets.BGPDefinedSets.CommunitySets.CommunitySet))
			assert.Equal(t, tc.definitions, len(target.RoutingPolicy.PolicyDefinitions.PolicyDefinition))

		})
	}
}
