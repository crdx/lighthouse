package adapterR_test

import (
	"testing"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/adapterR"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestFindByMACAddress(t *testing.T) {
	testCases := []struct {
		inputMACAddress string
		expectedAdapter *m.Adapter // Assuming m.Adapter is your model
		expectedFound   bool
	}{
		{"AA:AA:AA:AA:AA:AA", &m.Adapter{MACAddress: "AA:AA:AA:AA:AA:AA"}, true},
		{"ZZ:ZZ:ZZ:ZZ:ZZ:ZZ", nil, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputMACAddress, func(t *testing.T) {
			actualAdapter, actualFound := adapterR.FindByMACAddress(testCase.inputMACAddress)

			if testCase.expectedFound {
				assert.Equal(t, testCase.expectedAdapter.MACAddress, actualAdapter.MACAddress)
			}
			assert.Equal(t, testCase.expectedFound, actualFound)
		})
	}
}
