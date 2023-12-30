package domain

import "fmt"

type Document struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Category string `json:"category"`
}

func NewDocument(id string, title string, content string, category string) *Document {
	return &Document{
		ID: id,
		Title: title,
		Content: content,
		Category: category,
	}
}

func (doc *Document) String() string {
	return fmt.Sprintf("Document(\nID: %s,\nTitle: %s,\nContent: %s,\nCategory: %s)", doc.ID, doc.Title, doc.Content, doc.Category)
}