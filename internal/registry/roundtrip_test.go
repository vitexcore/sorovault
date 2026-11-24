package registry_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soroworks/sorovault/internal/model"
	"github.com/soroworks/sorovault/internal/registry"
)

// TestBackfillSummaryJSONRoundTrip checks that a BackfillSummary value
// survives a JSON encode → decode cycle unchanged.
func TestBackfillSummaryJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := registry.BackfillSummary{
		Checked: 100,
		Found:   85,
		Unknown: 10,
		Failed:  5,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded registry.BackfillSummary
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original, decoded)
}

// TestWasmResultJSONRoundTrip checks that a WasmResult value survives a
// JSON encode → decode cycle unchanged.
func TestWasmResultJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{
		ABIVersion: model.ABIVersion,
		Functions: []model.Function{
			{Name: "has_voted", Inputs: []model.Param{{Name: "nullifier", Type: model.Type{Kind: model.KindBytesN, Display: "BytesN<32>", N: 32}}}},
		},
		Types:  model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}},
		Events: []model.Event{},
	}

	original := registry.WasmResult{
		Network:   "testnet",
		WasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		Interface: iface,
		DecodedAt: ts,
		Created:   true,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded registry.WasmResult
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.Network, decoded.Network)
	assert.Equal(t, original.WasmHash, decoded.WasmHash)
	assert.Equal(t, original.DecodedAt.UTC(), decoded.DecodedAt.UTC())
	assert.Equal(t, original.Created, decoded.Created)
	require.NotNil(t, decoded.Interface)
	assert.Equal(t, iface.ABIVersion, decoded.Interface.ABIVersion)
	assert.Equal(t, iface.Functions[0].Name, decoded.Interface.Functions[0].Name)
}
