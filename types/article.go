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
	Status   string                  `json:"status"`
	Data     *Article                `json:"data,omitempty"`
	Message  string                  `json:"message,omitempty"`
	Metadata ArticleDetailedMetadata `json:"metadata"`
}

type ArticleDetailedMetadata struct {
	CreatedAt time.Time `json:"createdAt"`
}

/*
ArticleList is struct that will be returned by sports news api.

It was generated with the help of:

	go install github.com/twpayne/go-jsonstruct/v2/cmd/gojsonstruct@latest
	cat examples/hullcityArticleList.json | gojsonstruct
*/
type ArticleList struct {
	Status   string              `json:"status"`
	Data     []Article           `json:"data,omitempty"`
	Message  string              `json:"message,omitempty"`
	Metadata ArticleListMetadata `json:"metadata"`
}

type ArticleListMetadata struct {
	CreatedAt  time.Time `json:"createdAt"`
	TotalItems int       `json:"totalItems,omitempty"`
	Sort       string    `json:"sort,omitempty"`
}

type ArticleId string

/*
Article is representing a single article served by API

Use SetGeneratedId after creating the object with empty Id to set it based on ArticleKey.
*/
type Article struct {
	ArticleKey  `json:",inline"`
	Content     string    `json:"content,omitempty"`
	GalleryUrls string    `json:"galleryUrls,omitempty"`
	Id          ArticleId `json:"id"` // id should be generated from ArticleKey fields
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

// SetGeneratedId is used for setting the ID in article after creation
func (a *Article) SetGeneratedId() error {
	if a.Id != "" {
		return nil
	}
	hashStruct, err := rxhash.HashStruct(a.ArticleKey)
	if err != nil {
		return err
	}
	a.Id = ArticleId(fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hashStruct[:8],
		hashStruct[8:12],
		hashStruct[12:16],
		hashStruct[16:20],
		hashStruct[20:],
	))
	return nil
}

func (d ArticleDetailed) GetMessage() string {
	return d.Message
}

func (d ArticleList) GetMessage() string {
	return d.Message
}
