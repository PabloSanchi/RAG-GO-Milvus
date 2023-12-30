package utils_test

import (
    "encoding/json"
    "net/http"
    "testing"
    "github.com/jarcoal/httpmock"
    "github.com/stretchr/testify/assert"
	"github.com/pablosanchi/datastore/repositories/utils"
)


func TestEncodeSuccess(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    embedding := []float32{0.1, 0.2, 0.3}
    mockResponse := utils.EmbeddingResponse{Embedding: embedding}
    responseBody, _ := json.Marshal(mockResponse)
    httpmock.RegisterResponder("POST", utils.ENDPOINT,
        httpmock.NewBytesResponder(http.StatusOK, responseBody))

    encoder := utils.NewEncoder()

    result, err := encoder.Encode("test text")

    assert.Nil(t, err)
    assert.Equal(t, embedding, result)

    info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST "+ utils.ENDPOINT])
}
