package v1

import (
	"encoding/json"
	"errors"
	"github.com/adamdyszy/sportsnews/types"
	"net/http"

	"github.com/adamdyszy/sportsnews/storage"
	"github.com/gorilla/mux"
)

func GetArticleByIdHandler(s storage.ArticleStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleId := vars["id"]
		article, err := s.Get(types.ArticleId(articleId))
		if err != nil {
			if errors.Is(err, storage.ArticleNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//a := types.ArticleDetailed{} // todo
		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func GetAllArticlesHandler(s storage.ArticleStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := s.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(articles)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
