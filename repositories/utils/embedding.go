package utils

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/pablosanchi/datastore/core/ports/secondary"
    "os"
)

var (
    ENDPOINT string = "http://" + os.Getenv("OLLAMA_HOST") + ":" + os.Getenv("OLLAMA_PORT") + os.Getenv("OLLAMA_ENDPOINT") //  /api/embeddings"
    MODEL    string = "mistral"
)

type Encoder struct {}

func NewEncoder() secondary.TextEncoder {
	return &Encoder{}
}

type EmbeddingResponse struct {
    Embedding []float32 `json:"embedding"`
}

func (e *Encoder) Encode(text string) ([]float32, error) {
    postBody, _ := json.Marshal(map[string]string{
        "model":  MODEL,
        "prompt": text,
    })

    requestBody := bytes.NewBuffer(postBody)
    resp, err := http.Post(ENDPOINT, "application/json", requestBody)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Read error: %v", err)
        return nil, err
    }

    var embeddingResponse EmbeddingResponse
    err = json.Unmarshal(body, &embeddingResponse)
    if err != nil {
        return nil, err
    }

    return embeddingResponse.Embedding, nil
}
