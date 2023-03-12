package storage

import (
	"errors"
	"github.com/adamdyszy/sportsnews/types"
)

type ArticleStorage interface {
	ArticleReader
	ArticleWriter
	Disconnect() error
}

type ArticleReader interface {
	Get(types.ArticleId) (types.Article, error)
	GetNewsWithoutDetailsIDs() ([]string, error)
	List() ([]types.Article, error)
}

var ArticleNotFound = errors.New("article not found")

type ArticleWriter interface {
	// Write takes article to save
	Write(types.Article) error
	// Delete takes articleID and tries to delete it from the storage
	Delete(id types.ArticleId) error
}

var ArticleAlreadyExists = errors.New("tried to write to already existing article id")
var ArticleWriteFailed = errors.New("could not write article")
