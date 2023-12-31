package helpers

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"context"
	"log"
	"os"
)

var (
    MilvusAddress = "http://" + os.Getenv("MILVUS_HOST") + ":" + os.Getenv("MILVUS_PORT")
)

func NewMilvusClient() client.Client {
    milvusClient, err := client.NewClient(context.Background(), client.Config{Address: MilvusAddress})
	if err != nil {
		log.Fatal("Failed to connect to Milvus:", err)
	}

	return milvusClient
}