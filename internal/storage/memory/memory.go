package memory

import (
	"fmt"
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/adamdyszy/sportsnews/types"
	"sync"
)

type innerStorage struct {
	articles          map[types.ArticleId]types.Article
	newsIdsForDetails map[string]struct{}
	mx                *sync.RWMutex
}

func (i innerStorage) Delete(id types.ArticleId) error {
	i.mx.Lock()
	defer i.mx.Unlock()
	delete(i.newsIdsForDetails, string(id))
	return nil
}

func (i innerStorage) Disconnect() error {
	return nil
}

func (i innerStorage) GetNewsWithoutDetailsIDs() ([]string, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()
	v := make([]string, 0, len(i.newsIdsForDetails))
	for newsId := range i.newsIdsForDetails {
		v = append(v, newsId)
	}
	return v, nil
}

func NewMemStorage() storage.ArticleStorage {
	s := innerStorage{articles: make(map[types.ArticleId]types.Article), mx: &sync.RWMutex{}, newsIdsForDetails: make(map[string]struct{})}
	return s
}

func (i innerStorage) Get(id types.ArticleId) (types.Article, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()
	val, found := i.articles[id]
	if found {
		return val, nil
	}
	return types.Article{}, fmt.Errorf("article with id: %v doesn't exist in memory storage", id)
}

func (i innerStorage) List() ([]types.Article, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()
	v := make([]types.Article, 0, len(i.articles))
	for _, value := range i.articles {
		v = append(v, value)
	}
	return v, nil
}

func (i innerStorage) Write(article types.Article) error {
	i.mx.Lock()
	defer i.mx.Unlock()
	id := article.Id
	if !article.HasDetails {
		_, found := i.articles[id]
		if found {
			return fmt.Errorf("%w with id: %v", storage.ArticleAlreadyExists, id)
		}
		i.newsIdsForDetails[article.NewsId] = struct{}{}
	} else {
		delete(i.newsIdsForDetails, article.NewsId)
	}
	i.articles[id] = article
	return nil
}
