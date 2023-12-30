package datastore

import (
	"github.com/pablosanchi/datastore/core/domain"
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/pablosanchi/datastore/core/ports/secondary"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"fmt"
	"context"
	"log"
)

const (
    MilvusAddress = "localhost:19530"
)

type DatastoreMilvusRepository struct{
	client client.Client
	encoder secondary.TextEncoder
}

func NewDatastoreMilvusRepository(encoder secondary.TextEncoder) ports.DatastoreRepository {
    milvusClient, err := client.NewClient(context.Background(), client.Config{Address: MilvusAddress})
	if err != nil {
		log.Fatal("Failed to connect to Milvus:", err)
	}

	return &DatastoreMilvusRepository{client: milvusClient, encoder: encoder}
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
	nEntities := len(documents)
	idList:= make([]string, 0, nEntities)
	titleList:= make([]string, 0, nEntities)
	contentList:= make([]string, 0, nEntities)
	categoryList:= make([]string, 0, nEntities)
	embeddingList := make([][]float32, 0, nEntities)

	for _, document := range documents {
		encodedContent, err := m.encoder.Encode(document.Content)

		if err != nil {
			log.Fatal("fail to encode content:", err.Error())
		}
		
		idList = append(idList, document.ID)
		titleList = append(titleList, document.Title)
		contentList = append(contentList, document.Content)
		categoryList = append(categoryList, document.Category)
		embeddingList = append(embeddingList, encodedContent)
	}

	idColumn := entity.NewColumnVarChar("id", idList)
	titleColumn := entity.NewColumnVarChar("title", titleList)
	contentColumn := entity.NewColumnVarChar("content", contentList)
	categoryColumn := entity.NewColumnVarChar("category", categoryList)
	embeddingColumn := entity.NewColumnFloatVector("embedding", 4096, embeddingList)

	if _, err := m.client.Upsert(
		context.Background(), 
		collectionName, 
		"",
		idColumn,
		titleColumn,
		contentColumn,
		categoryColumn,
		embeddingColumn,	
	);

	err != nil {
			log.Fatalf("failed to upsert data, err: %v", err)
	}

	return nil
}

func (m *DatastoreMilvusRepository) Search(collectionName string, query string) ([]domain.Document, error) {

	encodedQuery, err := m.encoder.Encode(query)

	if err != nil {
		log.Fatal("fail to encode query:", err.Error())
	}

	err = m.client.LoadCollection(
		context.Background(),
		collectionName,
		false,
	)

	if err != nil {
		log.Fatal("failed to load collection:", err.Error())
	}

	sp, _ := entity.NewIndexIvfFlatSearchParam( // NewIndex*SearchParam func
		10,                                  // searchParam
	)
	
	opt := client.SearchQueryOptionFunc(func(option *client.SearchQueryOption) {
		option.Limit = 3
		option.Offset = 0
		option.ConsistencyLevel = entity.ClStrong
		option.IgnoreGrowing = false
	})

	searchResult, err := m.client.Search(
		context.Background(),
		collectionName,
		[]string{},
		"",
		[]string{"title", "content", "category"},
		[]entity.Vector{entity.FloatVector(encodedQuery)},
		"embedding",
		entity.COSINE,
		10,
		sp,
		opt,
	)

	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}

	fields := searchResult[0].Fields
	titleList := fields.GetColumn("title")
	contentList := fields.GetColumn("content")
	categoryList := fields.GetColumn("category")

	var documents []domain.Document
	for i := 0; i < titleList.Len(); i++ {
		
		title, _ := titleList.GetAsString(i);
		content, _ := contentList.GetAsString(i);
		category, _ := categoryList.GetAsString(i);

		documents = append(documents, domain.Document{
			ID: "",
			Title: title,
			Content: content,
			Category: category,
		})
	}
	
	err = m.client.ReleaseCollection(
		context.Background(),                            // ctx
		collectionName,                                          // CollectionName
	)

	if err != nil {
		log.Fatal("failed to release collection:", err.Error())
	}

	return documents, nil
}

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
					entity.TypeParamDim: fmt.Sprintf("%d", 4096),
				},
			},
		},
	}
}
