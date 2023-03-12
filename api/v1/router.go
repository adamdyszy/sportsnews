package v1

import (
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

func ListenAndServe(v *viper.Viper, s storage.ArticleStorage) error {
	address := v.GetString("address")

	r := mux.NewRouter()
	r.HandleFunc("/articles/{id}", GetArticleByIdHandler(s)).Methods("GET")
	r.HandleFunc("/articles", GetAllArticlesHandler(s)).Methods("GET")
	return http.ListenAndServe(address, r)
}
