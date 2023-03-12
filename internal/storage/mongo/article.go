package mongo

import (
	"github.com/adamdyszy/sportsnews/types"
	"time"
)

/*
ArticleBson is a copy of Article but with bson tags
*/
type articleBson struct {
	Published   time.Time `bson:"published"`
	TeamId      string    `bson:"teamId"`
	NewsId      string    `bson:"newsId"`
	Content     string    `bson:"content,omitempty"`
	GalleryUrls string    `bson:"galleryUrls,omitempty"`
	Id          string    `bson:"id"`
	ImageURL    string    `bson:"imageUrl,omitempty"`
	OptaMatchID string    `bson:"optaMatchId,omitempty"`
	Teaser      string    `bson:"teaser,omitempty"`
	Title       string    `bson:"title,omitempty"`
	URL         string    `bson:"url,omitempty"`
	VideoURL    string    `bson:"videoUrl,omitempty"`
	Type        []string  `bson:"type,omitempty"`
	HasDetails  bool      `bson:"hasDetails"`
}

func fromArticle(a types.Article) articleBson {
	return articleBson{
		TeamId:      a.TeamId,
		NewsId:      a.NewsId,
		Published:   a.Published,
		Content:     a.Content,
		GalleryUrls: a.GalleryUrls,
		Id:          string(a.Id),
		ImageURL:    a.ImageURL,
		OptaMatchID: a.OptaMatchId,
		Teaser:      a.Teaser,
		Title:       a.Title,
		Type:        a.Type,
		URL:         a.URL,
		VideoURL:    a.VideoURL,
		HasDetails:  a.HasDetails,
	}
}

func (a *articleBson) ToArticle() types.Article {
	return types.Article{
		ArticleKey: types.ArticleKey{
			TeamId:    a.TeamId,
			NewsId:    a.NewsId,
			Published: a.Published,
		},
		Content:     a.Content,
		GalleryUrls: a.GalleryUrls,
		Id:          types.ArticleId(a.Id),
		ImageURL:    a.ImageURL,
		OptaMatchId: a.OptaMatchID,
		Teaser:      a.Teaser,
		Title:       a.Title,
		Type:        a.Type,
		URL:         a.URL,
		VideoURL:    a.VideoURL,
		HasDetails:  a.HasDetails,
	}
}
