package v1

import (
	"encoding/json"
	"errors"
	"github.com/adamdyszy/sportsnews/types"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/adamdyszy/sportsnews/storage"
)

const internalServerErrorMsg = "Internal server error"
const articleIdNotFoundMsg = "ArticleId not found"
const failFromStorageMsg = "Failure while getting articles from storage"
const failJsonEncodeMsg = "Failure during json encoding"

type WithMessage interface {
	GetMessage() string
}

func MakeSuccessArticleList(data []types.Article) types.ArticleList {
	return types.ArticleList{
		Data: data,
		Metadata: types.ArticleListMetadata{
			CreatedAt:  time.Now(),
			Sort:       "random",
			TotalItems: len(data),
		},
		Status: "success",
	}
}

func MakeErrorArticleList(msg string) types.ArticleList {
	return types.ArticleList{
		Data: nil,
		Metadata: types.ArticleListMetadata{
			CreatedAt: time.Now(),
		},
		Status:  "error",
		Message: msg,
	}
}

func MakeSuccessArticleDetailed(data types.Article) types.ArticleDetailed {
	return types.ArticleDetailed{
		Data: &data,
		Metadata: types.ArticleDetailedMetadata{
			CreatedAt: time.Now(),
		},
		Status: "success",
	}
}

func MakeErrorArticleDetailed(msg string) types.ArticleDetailed {
	return types.ArticleDetailed{
		Data: nil,
		Metadata: types.ArticleDetailedMetadata{
			CreatedAt: time.Now(),
		},
		Status:  "error",
		Message: msg,
	}
}

func jsonEncodeSuccessResponse(w http.ResponseWriter, response WithMessage, logger logr.Logger) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error(err, failJsonEncodeMsg, "response", response)
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
	}
}

func jsonEncodeErrorResponse(w http.ResponseWriter, response WithMessage, status int, logger logr.Logger) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error(err, failJsonEncodeMsg, "response", response)
		http.Error(w, response.GetMessage(), status)
	}
}

func GetArticleByIdHandler(s storage.ArticleStorage, logger logr.Logger) http.HandlerFunc {
	logger.WithValues("handler", "GetArticleByIdHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleId := vars["id"]
		logger.WithValues("articleId", articleId)
		article, err := s.Get(types.ArticleId(articleId))
		if err != nil {
			if errors.Is(err, storage.ArticleNotFound) {
				// not logging here since it might happen often, we want to log important errors
				response := MakeErrorArticleDetailed(articleIdNotFoundMsg)
				jsonEncodeErrorResponse(w, response, http.StatusNotFound, logger)
				return
			}
			logger.Error(err, failFromStorageMsg)
			response := MakeErrorArticleDetailed(internalServerErrorMsg)
			jsonEncodeErrorResponse(w, response, http.StatusInternalServerError, logger)
			return
		}
		response := MakeSuccessArticleDetailed(article)
		jsonEncodeSuccessResponse(w, response, logger)
	}
}

func GetAllArticlesHandler(s storage.ArticleStorage, logger logr.Logger) http.HandlerFunc {
	logger.WithValues("handler", "GetAllArticlesHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := s.List()
		if err != nil {
			logger.Error(err, failFromStorageMsg)
			response := MakeErrorArticleList(internalServerErrorMsg)
			jsonEncodeErrorResponse(w, response, http.StatusInternalServerError, logger)
			return
		}
		response := MakeSuccessArticleList(articles)
		jsonEncodeSuccessResponse(w, response, logger)
	}
}
