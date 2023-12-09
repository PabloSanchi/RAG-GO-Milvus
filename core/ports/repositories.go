package ports

import (
	"github.com/pablosanchi/datastore/core/domain"
)

type DatastoreRepository interface {
	CreateCollection(collectionName string) error
	DeleteCollection(collectionName string) error
	List() ([]string, error)
	UpsertDocuments(collectionName string, documents []domain.Document) error
	Search(collectionName string, query string) ([]domain.Document, error)
}