// Package types has structs for the sports news api and for its poller of news
package types

import (
	"fmt"
	"github.com/rxwycdh/rxhash"
	"time"
)

/*
ArticleDetailed is struct that will be returned by sports news api.

It was generated with the help of:

	go install github.com/twpayne/go-jsonstruct/v2/cmd/gojsonstruct@latest
	cat examples/hullcityArticleDetailed.json | gojsonstruct
*/
type ArticleDetailed struct {
	Data     Article `json:"data"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"` // 2023-03-06T10:26:56.560Z
	} `json:"metadata"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ArticleList struct {
	Data     []Article `json:"data"`
	Metadata struct {
		CreatedAt  time.Time `json:"createdAt"`
		Sort       string    `json:"sort"`
		TotalItems int       `json:"totalItems"` // 2023-03-06T10:26:56.560Z
	} `json:"metadata"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ArticleId string

/*
Article is defined as a contract for the API and hashing function meaning it should not be changed.
*/
type Article struct {
	ArticleKey  `json:",inline"`
	Content     string    `json:"content,omitempty"`
	GalleryUrls string    `json:"galleryUrls,omitempty"`
	Id          ArticleId `json:"id"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	OptaMatchId string    `json:"optaMatchId,omitempty"`
	Teaser      string    `json:"teaser,omitempty"`
	Title       string    `json:"title,omitempty"`
	Type        []string  `json:"type,omitempty"`
	URL         string    `json:"url,omitempty"`
	VideoURL    string    `json:"videoUrl,omitempty"`
	HasDetails  bool      `json:"hasDetails"`
}

// ArticleKey represents fields that are used to generate hash from article.
type ArticleKey struct {
	TeamId    string    `json:"teamId"`
	NewsId    string    `json:"-"`
	Published time.Time `json:"published"`
}

func (a Article) GetOrGenerateID() (ArticleId, error) {
	if a.Id != "" {
		return a.Id, nil
	}
	hashStruct, err := rxhash.HashStruct(a.ArticleKey)
	if err != nil {
		return "", err
	}
	return ArticleId(fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hashStruct[:8],
		hashStruct[8:12],
		hashStruct[12:16],
		hashStruct[16:20],
		hashStruct[20:],
	)), nil
}
