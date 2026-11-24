package refresher_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soroworks/sorovault/internal/refresher"
)

// TestChangeJSONRoundTrip checks that a Change value survives a JSON
// encode → decode cycle unchanged.
func TestChangeJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 11, 24, 10, 0, 0, 0, time.UTC)
	original := refresher.Change{
		ContractID:   "CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
		Network:      "testnet",
		OldWasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		NewWasmHash:  "0f5bd7b2e93efd1c07d2f5b1f0f2f2a1c0d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7",
		InterfaceURL: "/api/contracts/CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN?network=testnet",
		ChangedAt:    ts,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded refresher.Change
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.ContractID, decoded.ContractID)
	assert.Equal(t, original.Network, decoded.Network)
	assert.Equal(t, original.OldWasmHash, decoded.OldWasmHash)
	assert.Equal(t, original.NewWasmHash, decoded.NewWasmHash)
	assert.Equal(t, original.InterfaceURL, decoded.InterfaceURL)
	assert.Equal(t, original.ChangedAt.UTC(), decoded.ChangedAt.UTC())
}

// TestSummaryJSONRoundTrip checks that a Summary value survives a JSON
// encode → decode cycle unchanged.
func TestSummaryJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := refresher.Summary{
		Checked: 148,
		Changed: 2,
		Failed:  1,
		Busy:    0,
		Missing: 3,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded refresher.Summary
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original, decoded)
}

// TestSweepRecordJSONRoundTrip checks that a SweepRecord value survives a
// JSON encode → decode cycle unchanged.
func TestSweepRecordJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 11, 24, 10, 5, 0, 0, time.UTC)
	original := refresher.SweepRecord{
		Summary: refresher.Summary{
			Checked: 50,
			Changed: 3,
			Failed:  0,
			Busy:    1,
			Missing: 0,
		},
		CompletedAt: ts,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded refresher.SweepRecord
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.Checked, decoded.Checked)
	assert.Equal(t, original.Changed, decoded.Changed)
	assert.Equal(t, original.Failed, decoded.Failed)
	assert.Equal(t, original.Busy, decoded.Busy)
	assert.Equal(t, original.Missing, decoded.Missing)
	assert.Equal(t, original.CompletedAt.UTC(), decoded.CompletedAt.UTC())
}
