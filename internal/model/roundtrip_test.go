package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soroworks/sorovault/internal/model"
)

// roundTrip marshals v to JSON and back into a fresh value of the same type.
func roundTrip[T any](t *testing.T, v T) T {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	var out T
	require.NoError(t, json.Unmarshal(data, &out))
	return out
}

// TestEnumJSONRoundTrip checks that an Enum survives a JSON encode → decode
// cycle unchanged.
func TestEnumJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Enum{
		Name: "Status",
		Doc:  "Contract status codes.",
		Cases: []model.EnumCase{
			{Name: "Active", Value: 0},
			{Name: "Paused", Value: 1},
		},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestEnumCaseJSONRoundTrip checks that an EnumCase survives a JSON
// encode → decode cycle unchanged.
func TestEnumCaseJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.EnumCase{Name: "Expired", Doc: "The entry has expired.", Value: 42}
	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestEnvMetaJSONRoundTripFull checks the full EnvMeta struct including a
// non-zero PreRelease field.
func TestEnvMetaJSONRoundTripFull(t *testing.T) {
	t.Parallel()

	original := model.NewEnvMeta(25, 3)
	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestErrorEnumJSONRoundTrip checks that an ErrorEnum survives a JSON
// encode → decode cycle unchanged.
func TestErrorEnumJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.ErrorEnum{
		Name: "Error",
		Cases: []model.EnumCase{
			{Name: "NotFound", Value: 1},
			{Name: "Unauthorized", Value: 2},
		},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestEventJSONRoundTrip checks that an Event survives a JSON
// encode → decode cycle unchanged.
func TestEventJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Event{
		Name:       "Voted",
		Doc:        "Emitted when a vote is cast.",
		Prefix:     []string{"voted"},
		DataFormat: "map",
		Params: []model.EventParam{
			{
				Name:     "voter",
				Type:     model.Scalar(model.KindAddress),
				Location: "topic-list",
			},
			{
				Name:     "proposal_index",
				Type:     model.Scalar(model.KindU32),
				Location: "data",
			},
		},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestEventParamJSONRoundTrip checks that an EventParam survives a JSON
// encode → decode cycle unchanged.
func TestEventParamJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.EventParam{
		Name:     "nullifier",
		Doc:      "Prevents double-voting.",
		Type:     model.Type{Kind: model.KindBytesN, Display: "BytesN<32>", N: 32},
		Location: "topic-list",
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestFunctionJSONRoundTrip checks that a Function survives a JSON
// encode → decode cycle unchanged.
func TestFunctionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Function{
		Name: "transfer",
		Doc:  "Move tokens between accounts.",
		Inputs: []model.Param{
			{Name: "from", Type: model.Scalar(model.KindAddress)},
			{Name: "to", Type: model.Scalar(model.KindAddress)},
			{Name: "amount", Type: model.Scalar(model.KindI128)},
		},
		Outputs: []model.Type{model.Scalar(model.KindVoid)},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestInterfaceJSONRoundTrip checks that an Interface survives a JSON
// encode → decode cycle unchanged.
func TestInterfaceJSONRoundTrip(t *testing.T) {
	t.Parallel()

	envMeta := model.NewEnvMeta(25, 0)
	original := model.Interface{
		ABIVersion: model.ABIVersion,
		EnvMeta:    &envMeta,
		Meta:       map[string]string{"rsver": "1.93.0", "rssdkver": "25.3.0"},
		Functions: []model.Function{
			{Name: "vote", Inputs: []model.Param{{Name: "index", Type: model.Scalar(model.KindU32)}}},
		},
		Types:  model.Types{Structs: []model.Struct{}, Unions: []model.Union{}, Enums: []model.Enum{}, ErrorEnums: []model.ErrorEnum{}},
		Events: []model.Event{},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original.ABIVersion, got.ABIVersion)
	assert.Equal(t, original.Meta, got.Meta)
	require.NotNil(t, got.EnvMeta)
	assert.Equal(t, *original.EnvMeta, *got.EnvMeta)
	assert.Equal(t, original.Functions, got.Functions)
}

// TestParamJSONRoundTrip checks that a Param survives a JSON
// encode → decode cycle unchanged.
func TestParamJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Param{
		Name: "nullifier",
		Doc:  "Unique per voter.",
		Type: model.Type{Kind: model.KindBytesN, Display: "BytesN<32>", N: 32},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestStructJSONRoundTrip checks that a Struct survives a JSON
// encode → decode cycle unchanged.
func TestStructJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Struct{
		Name:    "Proposal",
		Doc:     "A single voting proposal.",
		IsTuple: false,
		Fields: []model.StructField{
			{Name: "description", Type: model.Scalar(model.KindString)},
			{Name: "vote_count", Type: model.Scalar(model.KindU64)},
		},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestStructFieldJSONRoundTrip checks that a StructField survives a JSON
// encode → decode cycle unchanged.
func TestStructFieldJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.StructField{
		Name: "vote_count",
		Doc:  "Number of votes cast for this proposal.",
		Type: model.Scalar(model.KindU64),
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestTypeJSONRoundTrip checks that a Type survives a JSON
// encode → decode cycle unchanged, covering several variants.
func TestTypeJSONRoundTrip(t *testing.T) {
	t.Parallel()

	scalar := model.Scalar(model.KindAddress)
	inner := model.Scalar(model.KindU32)
	tests := []struct {
		name string
		in   model.Type
	}{
		{"scalar", scalar},
		{"bytesN", model.Type{Kind: model.KindBytesN, Display: "BytesN<32>", N: 32}},
		{"udt", model.Type{Kind: model.KindUDT, Display: "Proposal", Name: "Proposal"}},
		{"vec", model.Type{Kind: model.KindVec, Display: "Vec<u32>", Element: &inner}},
		{"option", model.Type{Kind: model.KindOption, Display: "Option<Address>", Inner: &scalar}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := roundTrip(t, tt.in)
			assert.Equal(t, tt.in, got)
		})
	}
}

// TestTypesJSONRoundTrip checks that a Types container survives a JSON
// encode → decode cycle unchanged.
func TestTypesJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Types{
		Structs:    []model.Struct{{Name: "Proposal", Fields: []model.StructField{{Name: "id", Type: model.Scalar(model.KindU32)}}}},
		Unions:     []model.Union{{Name: "DataKey", Cases: []model.UnionCase{{Name: "Proposal"}}}},
		Enums:      []model.Enum{{Name: "Status", Cases: []model.EnumCase{{Name: "Active", Value: 0}}}},
		ErrorEnums: []model.ErrorEnum{{Name: "Error", Cases: []model.EnumCase{{Name: "NotFound", Value: 1}}}},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestUnionJSONRoundTrip checks that a Union survives a JSON
// encode → decode cycle unchanged.
func TestUnionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.Union{
		Name: "DataKey",
		Doc:  "Keys used in contract storage.",
		Cases: []model.UnionCase{
			{Name: "Params"},
			{Name: "Proposal", Values: []model.Type{model.Scalar(model.KindU32)}},
		},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}

// TestUnionCaseJSONRoundTrip checks that a UnionCase survives a JSON
// encode → decode cycle unchanged.
func TestUnionCaseJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := model.UnionCase{
		Name:   "Proposal",
		Doc:    "A proposal entry.",
		Values: []model.Type{model.Scalar(model.KindU32)},
	}

	got := roundTrip(t, original)
	assert.Equal(t, original, got)
}
