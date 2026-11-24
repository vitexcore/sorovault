package mcp_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soroworks/sorovault/internal/mcp"
)

// TestParamsJSONRoundTrip checks that a Params value survives a JSON
// encode → decode cycle unchanged.
func TestParamsJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := mcp.Params{
		Method:  mcp.MethodSearch,
		ID:      42,
		JSONRPC: "2.0",
		Data:    json.RawMessage(`{"query":"vote","limit":10}`),
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded mcp.Params
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.Method, decoded.Method)
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.JSONRPC, decoded.JSONRPC)
	assert.JSONEq(t, string(original.Data), string(decoded.Data))
}

// TestResultJSONRoundTrip checks that a Result value survives a JSON
// encode → decode cycle unchanged.
func TestResultJSONRoundTrip(t *testing.T) {
	t.Parallel()

	payload := json.RawMessage(`{"status":"ok","contracts":[]}`)
	original := mcp.Result{
		JSONRPC: "2.0",
		ID:      7,
		Result:  payload,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded mcp.Result
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.JSONRPC, decoded.JSONRPC)
	assert.Equal(t, original.ID, decoded.ID)
	assert.JSONEq(t, string(original.Result), string(decoded.Result))
	assert.Nil(t, decoded.Error, "no error field present")
}

// TestResultWithErrorJSONRoundTrip checks that a Result carrying an error
// survives a JSON encode → decode cycle unchanged.
func TestResultWithErrorJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := mcp.Result{
		JSONRPC: "2.0",
		ID:      3,
		Error:   &mcp.MCPError{Code: -32601, Message: "method not found"},
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded mcp.Result
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.NotNil(t, decoded.Error)
	assert.Equal(t, original.Error.Code, decoded.Error.Code)
	assert.Equal(t, original.Error.Message, decoded.Error.Message)
}
