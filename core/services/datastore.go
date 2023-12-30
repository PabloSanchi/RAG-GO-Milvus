package services

import (
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/pablosanchi/datastore/core/domain"
)

type DatastoreService struct {
	datastoreRepository ports.DatastoreRepository
}

func NewDatastoreService(datastoreRepository ports.DatastoreRepository) ports.DatastoreService {
	return &DatastoreService{
		datastoreRepository: datastoreRepository,
	}
}

func (dsSrv *DatastoreService) CreateCollection(collectionName string) error {
	return dsSrv.datastoreRepository.CreateCollection(collectionName)
}

func (dsSrv *DatastoreService) DeleteCollection(collectionName string) error {
	return dsSrv.datastoreRepository.DeleteCollection(collectionName)
}

func (dsSrv *DatastoreService) List() ([]string, error) {
	collectionList, err := dsSrv.datastoreRepository.List()
	
	if err != nil {
		return nil, err
	}

	return collectionList, nil
}

func (dsSrv *DatastoreService) UpsertDocuments(collectionName string, documents []domain.Document) error {
	return dsSrv.datastoreRepository.UpsertDocuments(collectionName, documents)
}

func (dsSrv *DatastoreService) Search(collectionName string, query string) ([]domain.Document, error) {
	documents, err := dsSrv.datastoreRepository.Search(collectionName, query)
	
	if err != nil {
		return nil, err
	}

	return documents, nil
}
