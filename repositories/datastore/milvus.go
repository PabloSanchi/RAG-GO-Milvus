package datastore

import (
	"github.com/pablosanchi/datastore/core/domain"
	"github.com/pablosanchi/datastore/core/ports"
	"fmt"
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



type DatastoreMilvusRepository struct {
	
}

func NewDatastoreMilvusRepository() ports.DatastoreRepository {
	return &DatastoreMilvusRepository{
		
	}
}

func (m *DatastoreMilvusRepository) CreateCollection(collectionName string) error {
	fmt.Printf("Creating collection | %s", collectionName)
	return nil
}

func (m *DatastoreMilvusRepository) DeleteCollection(collectionName string) error {
	fmt.Println("Deleting collection", collectionName)
	return nil
}

func (m *DatastoreMilvusRepository) List() ([]string, error) {
	fmt.Println("Listing collections")
	return []string{}, nil
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