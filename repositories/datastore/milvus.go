package datastore

import (
	"github.com/pablosanchi/datastore/core/domain"
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"fmt"
	"context"
	"log"
)

const (
    MilvusAddress = "localhost:19530"
)

type DatastoreMilvus struct {
    ID        string
    Title     string
    Content   string
    Category  string
    Embedding []float32
}

func (dsm *DatastoreMilvus) ToDomain() *domain.Document {
    return &domain.Document{
        ID:       dsm.ID,
        Title:    dsm.Title,
        Content:  dsm.Content,
        Category: dsm.Category,
    }
}

func (dsm *DatastoreMilvus) FromDomain(document *domain.Document) {
    dsm.ID = document.ID
    dsm.Title = document.Title
    dsm.Content = document.Content
    dsm.Category = document.Category
}

type DatastoreMilvusRepository struct{
	client client.Client
}

func NewDatastoreMilvusRepository() ports.DatastoreRepository {
    milvusClient, err := client.NewClient(context.Background(), client.Config{Address: MilvusAddress})
	if err != nil {
		log.Fatal("Failed to connect to Milvus:", err)
	}

	return &DatastoreMilvusRepository{client: milvusClient}
}

func (m *DatastoreMilvusRepository) CreateCollection(collectionName string) error {
    schema := defineSchema(collectionName)
    if err := m.client.CreateCollection(context.Background(), schema, entity.DefaultShardNumber); err != nil {
        return fmt.Errorf("failed to create collection: %w", err)
    }

    if err := m.buildIndex(collectionName); err != nil {
        return fmt.Errorf("failed to build index: %w", err)
    }

    log.Println("Collection created successfully")
    return nil
}

func (m *DatastoreMilvusRepository) DeleteCollection(collectionName string) error {
    if err := m.client.DropCollection(context.Background(), collectionName); err != nil {
        return fmt.Errorf("failed to drop collection: %w", err)
    }

    log.Println("Collection dropped successfully")
    return nil
}

func (m *DatastoreMilvusRepository) List() ([]string, error) {
    listColl, err := m.client.ListCollections(context.Background())
    if err != nil {
        return nil, fmt.Errorf("failed to list collections: %w", err)
    }

    var collections []string
    for _, collection := range listColl {
        collections = append(collections, collection.Name)
    }

    log.Println("Collections listed successfully")
    return collections, nil
}

func (m *DatastoreMilvusRepository) UpsertDocuments(collectionName string, documents []domain.Document) error {
	fmt.Println("Upserting documents", collectionName)

	for _, document := range documents {
		fmt.Println(document.String())
	}

	return nil
}

func (m *DatastoreMilvusRepository) Search(collectionName string, query string) ([]domain.Document, error) {
	fmt.Printf("Searching documents in collection %s with query %s", collectionName, query)
	return []domain.Document{}, nil
}

// private
func (m *DatastoreMilvusRepository) buildIndex(collectionName string) error {
    idx, err := entity.NewIndexIvfFlat(entity.COSINE, 1024)
    if err != nil {
        return fmt.Errorf("fail to create IVF flat index parameter: %w", err)
    }

    err = m.client.CreateIndex(context.Background(), collectionName, "embedding", idx, false)
    if err != nil {
        return fmt.Errorf("fail to create index: %w", err)
    }

    return nil
}

func defineSchema(collectionName string) *entity.Schema {
	return &entity.Schema{
		CollectionName: collectionName,
		Description:    "",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				AutoID:     false,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 255),
				},
			},
			{
				Name:       "title",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: false,
				AutoID:     false,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 255),
				},
			},
			{
				Name:       "content",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: false,
				AutoID:     false,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 3000),
				},
			},
			{
				Name:       "category",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: false,
				AutoID:     false,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 100),
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					entity.TypeParamDim: fmt.Sprintf("%d", 768),
				},
			},
		},
	}
}
