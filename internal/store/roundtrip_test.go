package store_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soroworks/sorovault/internal/model"
	"github.com/soroworks/sorovault/internal/store"
)

// storeRoundTrip marshals v to JSON and back into a fresh value of the same
// type. It is shared by all round-trip tests in this package.
func storeRoundTrip[T any](t *testing.T, v T) T {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	var out T
	require.NoError(t, json.Unmarshal(data, &out))
	return out
}

// newContract returns a fully-populated store.Contract for round-trip tests.
func newContract() store.Contract {
	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	verifiedAt := ts.Add(time.Hour)
	missingSince := ts.Add(2 * time.Hour)
	deprecatedSince := ts.Add(3 * time.Hour)
	ledger := uint32(1478234)
	return store.Contract{
		ContractID:       "CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
		Network:          "testnet",
		CurrentWasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		Name:             "zkVote",
		SourceURL:        "https://github.com/example/zkvote",
		SourceCommit:     "abc123",
		FirstSeen:        ts,
		LastRefreshed:    ts,
		ObservedLedger:   nil,
		Creator:          "GDUZ7Q4K...creator",
		CreatedLedger:    ledger,
		VerifiedBy:       "GDUZ7Q4K...publisher",
		VerifiedAt:       &verifiedAt,
		VerifiedWasmHash: "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		MissingSince:     &missingSince,
		Implements:       []string{"sep41"},
		DeprecatedSince:  &deprecatedSince,
		Successor:        "CDZZZPYQALFRBQAZD5KCW2SZK4ZF455PJVD7H53AGHXTWAJO4BYVPP2V",
		DeprecationNote:  "Move to v2.",
	}
}

// TestContractJSONRoundTrip checks that a Contract value survives a JSON
// encode → decode cycle unchanged.
func TestContractJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := newContract()
	got := storeRoundTrip(t, original)

	assert.Equal(t, original.ContractID, got.ContractID)
	assert.Equal(t, original.Network, got.Network)
	assert.Equal(t, original.CurrentWasmHash, got.CurrentWasmHash)
	assert.Equal(t, original.Name, got.Name)
	assert.Equal(t, original.SourceURL, got.SourceURL)
	assert.Equal(t, original.Creator, got.Creator)
	assert.Equal(t, original.CreatedLedger, got.CreatedLedger)
	assert.Equal(t, original.VerifiedBy, got.VerifiedBy)
	require.NotNil(t, got.VerifiedAt)
	assert.Equal(t, original.VerifiedAt.UTC(), got.VerifiedAt.UTC())
	require.NotNil(t, got.MissingSince)
	assert.Equal(t, original.MissingSince.UTC(), got.MissingSince.UTC())
	assert.Equal(t, original.Implements, got.Implements)
	require.NotNil(t, got.DeprecatedSince)
	assert.Equal(t, original.DeprecatedSince.UTC(), got.DeprecatedSince.UTC())
	assert.Equal(t, original.Successor, got.Successor)
	assert.Equal(t, original.DeprecationNote, got.DeprecationNote)
}

// TestContractPageJSONRoundTrip checks that a ContractPage survives a JSON
// encode → decode cycle unchanged.
func TestContractPageJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.ContractPage{
		Contracts: []store.Contract{newContract()},
		Total:     1,
		Next:      "cursor-abc",
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Total, got.Total)
	assert.Equal(t, original.Next, got.Next)
	require.Len(t, got.Contracts, 1)
	assert.Equal(t, original.Contracts[0].ContractID, got.Contracts[0].ContractID)
}

// TestContractSnapshotJSONRoundTrip checks that a ContractSnapshot survives a
// JSON encode → decode cycle unchanged.
func TestContractSnapshotJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{ABIVersion: model.ABIVersion, Functions: []model.Function{{Name: "vote"}}, Types: model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}}, Events: []model.Event{}}
	original := store.ContractSnapshot{
		Contract: newContract(),
		Versions: []store.SpecSnapshot{
			{
				WasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
				Interface: iface,
				DecodedAt: ts,
			},
		},
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Contract.ContractID, got.Contract.ContractID)
	require.Len(t, got.Versions, 1)
	assert.Equal(t, original.Versions[0].WasmHash, got.Versions[0].WasmHash)
}

// TestNetworkStatsJSONRoundTrip checks that a NetworkStats value survives a
// JSON encode → decode cycle unchanged.
func TestNetworkStatsJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.NetworkStats{Network: "testnet", Contracts: 42, Versions: 60}
	got := storeRoundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestSnapshotJSONRoundTrip checks that a Snapshot survives a JSON
// encode → decode cycle unchanged.
func TestSnapshotJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{ABIVersion: model.ABIVersion, Functions: []model.Function{{Name: "vote"}}, Types: model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}}, Events: []model.Event{}}
	original := store.Snapshot{
		Version: store.SnapshotVersion,
		Contracts: []store.ContractSnapshot{
			{
				Contract: newContract(),
				Versions: []store.SpecSnapshot{{WasmHash: "abc", Interface: iface, DecodedAt: ts}},
			},
		},
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Version, got.Version)
	require.Len(t, got.Contracts, 1)
	assert.Equal(t, original.Contracts[0].Contract.ContractID, got.Contracts[0].Contract.ContractID)
}

// TestSpecJSONRoundTrip checks that a Spec survives a JSON
// encode → decode cycle unchanged.
func TestSpecJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{ABIVersion: model.ABIVersion, Functions: []model.Function{}, Types: model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}}, Events: []model.Event{}}
	ledger := uint32(1478234)
	original := store.Spec{
		ID:             42,
		ContractID:     "CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
		Network:        "testnet",
		WasmHash:       "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		Interface:      iface,
		DecodedAt:      ts,
		ObservedLedger: &ledger,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.ID, got.ID)
	assert.Equal(t, original.ContractID, got.ContractID)
	assert.Equal(t, original.Network, got.Network)
	assert.Equal(t, original.WasmHash, got.WasmHash)
	assert.Equal(t, original.DecodedAt.UTC(), got.DecodedAt.UTC())
	require.NotNil(t, got.ObservedLedger)
	assert.Equal(t, *original.ObservedLedger, *got.ObservedLedger)
}

// TestSpecSnapshotJSONRoundTrip checks that a SpecSnapshot survives a JSON
// encode → decode cycle unchanged.
func TestSpecSnapshotJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{ABIVersion: model.ABIVersion, Functions: []model.Function{}, Types: model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}}, Events: []model.Event{}}
	original := store.SpecSnapshot{
		WasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		Interface: iface,
		DecodedAt: ts,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.WasmHash, got.WasmHash)
	assert.Equal(t, original.DecodedAt.UTC(), got.DecodedAt.UTC())
	require.NotNil(t, got.Interface)
}

// TestSpecVersionJSONRoundTrip checks that a SpecVersion survives a JSON
// encode → decode cycle unchanged.
func TestSpecVersionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	ledger := uint32(1478234)
	original := store.SpecVersion{
		WasmHash:       "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		DecodedAt:      ts,
		Current:        true,
		ObservedLedger: &ledger,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.WasmHash, got.WasmHash)
	assert.Equal(t, original.Current, got.Current)
	assert.Equal(t, original.DecodedAt.UTC(), got.DecodedAt.UTC())
	require.NotNil(t, got.ObservedLedger)
	assert.Equal(t, *original.ObservedLedger, *got.ObservedLedger)
}

// TestStatsJSONRoundTrip checks that a Stats value survives a JSON
// encode → decode cycle unchanged.
func TestStatsJSONRoundTrip(t *testing.T) {
	t.Parallel()

	c := newContract()
	original := store.Stats{
		Networks:  []store.NetworkStats{{Network: "testnet", Contracts: 10, Versions: 15}},
		Contracts: 10,
		Versions:  15,
		Latest:    &c,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Contracts, got.Contracts)
	assert.Equal(t, original.Versions, got.Versions)
	require.Len(t, got.Networks, 1)
	assert.Equal(t, original.Networks[0], got.Networks[0])
	require.NotNil(t, got.Latest)
	assert.Equal(t, c.ContractID, got.Latest.ContractID)
}

// TestSubscribeRequestJSONRoundTrip checks that a SubscribeRequest survives
// a JSON encode → decode cycle unchanged.
func TestSubscribeRequestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.SubscribeRequest{
		ContractIDs: []string{
			"CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
		},
		URL: "https://hooks.example.com/sorovault",
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.URL, got.URL)
	assert.Equal(t, original.ContractIDs, got.ContractIDs)
}

// TestSubscriptionJSONRoundTrip checks that a Subscription survives a JSON
// encode → decode cycle unchanged.
func TestSubscriptionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 11, 24, 10, 0, 0, 0, time.UTC)
	original := store.Subscription{
		ID:          "deadbeefcafe0001",
		Network:     "testnet",
		ContractIDs: []string{"CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN"},
		URL:         "https://hooks.example.com/sorovault",
		Active:      true,
		CreatedAt:   ts,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.ID, got.ID)
	assert.Equal(t, original.Network, got.Network)
	assert.Equal(t, original.ContractIDs, got.ContractIDs)
	assert.Equal(t, original.URL, got.URL)
	assert.Equal(t, original.Active, got.Active)
	assert.Equal(t, original.CreatedAt.UTC(), got.CreatedAt.UTC())
	assert.Empty(t, got.Secret, "Secret has json:\"-\" and must not round-trip")
}

// TestSubscriptionListItemJSONRoundTrip checks that a SubscriptionListItem
// survives a JSON encode → decode cycle unchanged.
func TestSubscriptionListItemJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.SubscriptionListItem{
		ID:        "deadbeefcafe0001",
		Network:   "testnet",
		Contracts: 2,
		Active:    true,
		CreatedAt: "2026-11-24T10:00:00Z",
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestUpgradeJSONRoundTrip checks that an Upgrade survives a JSON
// encode → decode cycle unchanged.
func TestUpgradeJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 9, 1, 10, 2, 0, 0, time.UTC)
	original := store.Upgrade{
		ContractID:  "CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
		Network:     "testnet",
		OldWasmHash: "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		NewWasmHash: "0f5bd7b2e93efd1c07d2f5b1f0f2f2a1c0d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7",
		ObservedAt:  ts,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.ContractID, got.ContractID)
	assert.Equal(t, original.Network, got.Network)
	assert.Equal(t, original.OldWasmHash, got.OldWasmHash)
	assert.Equal(t, original.NewWasmHash, got.NewWasmHash)
	assert.Equal(t, original.ObservedAt.UTC(), got.ObservedAt.UTC())
}

// TestUpgradePageJSONRoundTrip checks that an UpgradePage survives a JSON
// encode → decode cycle unchanged.
func TestUpgradePageJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 9, 1, 10, 2, 0, 0, time.UTC)
	original := store.UpgradePage{
		Upgrades: []store.Upgrade{{
			ContractID:  "CDZZXVGUCLBUTE3WDD6TF2NXPTBPOYTVPPBBIZMSPVIQLUG226XJ4PAN",
			Network:     "testnet",
			OldWasmHash: "c618...",
			NewWasmHash: "0f5b...",
			ObservedAt:  ts,
		}},
		Total: 1,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Total, got.Total)
	require.Len(t, got.Upgrades, 1)
	assert.Equal(t, original.Upgrades[0].ContractID, got.Upgrades[0].ContractID)
}

// TestWasmContractJSONRoundTrip checks that a WasmContract survives a JSON
// encode → decode cycle unchanged.
func TestWasmContractJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.WasmContract{
		Contract: newContract(),
		Current:  true,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.ContractID, got.ContractID)
	assert.Equal(t, original.Current, got.Current)
}

// TestWasmContractPageJSONRoundTrip checks that a WasmContractPage survives
// a JSON encode → decode cycle unchanged.
func TestWasmContractPageJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := store.WasmContractPage{
		Contracts: []store.WasmContract{{Contract: newContract(), Current: true}},
		Total:     1,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Total, got.Total)
	require.Len(t, got.Contracts, 1)
	assert.Equal(t, original.Contracts[0].ContractID, got.Contracts[0].ContractID)
	assert.True(t, got.Contracts[0].Current)
}

// TestWasmSpecJSONRoundTrip checks that a WasmSpec survives a JSON
// encode → decode cycle unchanged.
func TestWasmSpecJSONRoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 8, 7, 9, 14, 11, 0, time.UTC)
	iface := &model.Interface{ABIVersion: model.ABIVersion, Functions: []model.Function{}, Types: model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}}, Events: []model.Event{}}
	original := store.WasmSpec{
		Network:   "testnet",
		WasmHash:  "c618dae264864ccf446a3c7db27da80c7c83e840131242cc2fe9cd32a2a20781",
		Interface: iface,
		DecodedAt: ts,
	}

	got := storeRoundTrip(t, original)
	assert.Equal(t, original.Network, got.Network)
	assert.Equal(t, original.WasmHash, got.WasmHash)
	assert.Equal(t, original.DecodedAt.UTC(), got.DecodedAt.UTC())
	require.NotNil(t, got.Interface)
}
