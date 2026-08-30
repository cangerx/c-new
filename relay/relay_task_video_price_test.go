package relay

import (
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
	"github.com/stretchr/testify/require"
)

func TestApplyEstimatedTaskBillingRatiosSkipsFixedVideoPrice(t *testing.T) {
	info := &relaycommon.RelayInfo{
		PriceData: types.PriceData{
			Quota:         100,
			UsePrice:      true,
			UseVideoPrice: true,
		},
	}

	applyEstimatedTaskBillingRatios(info, "video-model", map[string]float64{
		"seconds": 10,
		"size":    1.5,
	})

	require.Equal(t, 100, info.PriceData.Quota)
	require.Nil(t, info.PriceData.OtherRatios())
}

func TestApplyEstimatedTaskBillingRatiosPreservesExistingFixedPriceBehavior(t *testing.T) {
	info := &relaycommon.RelayInfo{
		PriceData: types.PriceData{
			Quota:    100,
			UsePrice: true,
		},
	}

	applyEstimatedTaskBillingRatios(info, "video-model", map[string]float64{
		"seconds": 10,
		"size":    1.5,
	})

	require.Equal(t, 1500, info.PriceData.Quota)
	require.Equal(t, map[string]float64{"seconds": 10, "size": 1.5}, info.PriceData.OtherRatios())
}
