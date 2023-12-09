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

type DatastoreMilvus struct {
	ID          string
	Title       string
	Content		string
	Category	string
	Embedding	[]float32
}

type DatastoreListMilvus []DatastoreMilvus


func (dsm *DatastoreMilvus) ToDomain() *domain.Document {
	return &domain.Document{
		ID: dsm.ID,
		Title: dsm.Title,
		Content: dsm.Content,
		Category: dsm.Category,
	}
}

func (dsm *DatastoreMilvus) FromDomain(document *domain.Document) {
	if dsm == nil {
		dsm = &DatastoreMilvus{}
	}

	dsm.ID = document.ID
	dsm.Title = document.Title
	dsm.Content = document.Content
	dsm.Category = document.Category
	dsm.Embedding = []float32{} // here goes the embedding
}

var milvusClient, _ = client.NewClient(context.Background(), client.Config{
	Address: "localhost:19530",
})

type DatastoreMilvusRepository struct {
	
}

func NewDatastoreMilvusRepository() ports.DatastoreRepository {
	return &DatastoreMilvusRepository{
		
	}
}

func (m *DatastoreMilvusRepository) CreateCollection(collectionName string) error {

	schema := &entity.Schema{
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

	err := milvusClient.CreateCollection(
		context.Background(),
		schema,
		entity.DefaultShardNumber,
	)

	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
		return err
	}

	err = buildIndex(collectionName)

	if err != nil {
		log.Fatal("failed to build index:", err.Error())
		return err
	}

	log.Println("collection created successfully")
	
	return nil
}

func (m *DatastoreMilvusRepository) DeleteCollection(collectionName string) error {
	
	err := milvusClient.DropCollection(
		context.Background(),
		collectionName,
	)
	
	if err != nil {
		log.Fatal("fail to drop collection:", err.Error())
	}

	log.Println("collection dropped successfully")

	return nil
}

func (m *DatastoreMilvusRepository) List() ([]string, error) {

	listColl, err := milvusClient.ListCollections(context.Background(),)
	if err != nil {
		return nil, err
	}

	var collections []string
	for _, collection := range listColl {
		collections = append(collections, collection.Name)
	}

	log.Println("collections listed successfully")

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
func buildIndex(collectionName string) error {
	idx, err := entity.NewIndexIvfFlat(
		entity.COSINE,
		1024,
	)

	if err != nil {
	  log.Fatal("fail to create ivf flat index parameter:", err.Error())
	  return err
	}

	err = milvusClient.CreateIndex(
		context.Background(),
		collectionName,
		"embedding",
		idx,
		false,
	)
	
	if err != nil {
		log.Fatal("fail to create index:", err.Error())
		return err
	}

	return nil
}