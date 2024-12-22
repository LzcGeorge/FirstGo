package store

import "errors"

var (
	ErrExist    = errors.New("book already exists")
	ErrNotFound = errors.New("book not found")
)

type Book struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Authors []string `json:"authors"`
	Press   string   `json:"press"`
}

type Store interface {
	Create(b *Book) error
	Update(b *Book) error
	Get(id string) (Book, error)
	Delete(id string) error
	List() ([]Book, error)
}
