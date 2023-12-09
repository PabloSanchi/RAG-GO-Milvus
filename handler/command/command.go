package command

import (
	"github.com/pablosanchi/datastore/core/domain"
)

type CollectionCommand struct {
	CollectionName string `json:"collection_name"`
}

type UpsertDocumentsCommand struct {
	CollectionName string `json:"collection_name"`
	Documents []domain.Document `json:"documents"`
}

type SearchCommand struct {
	CollectionName string `json:"collection_name"`
	Query string `json:"query"`
}