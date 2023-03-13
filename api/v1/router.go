package v1

import (
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

func ListenAndServe(v *viper.Viper, s storage.ArticleStorage, logger logr.Logger) error {
	address := v.GetString("address")
	r := mux.NewRouter()
	// Serve api handlers
	r.HandleFunc("/articles/{id}", GetArticleByIdHandler(s, logger)).Methods("GET")
	r.HandleFunc("/articles", GetAllArticlesHandler(s, logger)).Methods("GET")

	return http.ListenAndServe(address, r)
}
