package poller

import (
	"encoding/xml"
	"github.com/adamdyszy/sportsnews/types"
	"strings"
	"time"
)

/*
NewsList is struct for getting news from feed.

It was generated with the help of:

	go install github.com/miku/zek/cmd/zek@latest
	cat examples/hullcityList.xml | zek -e

The inner news element was refactored into separate struct NewsElement.
*/
type NewsList struct {
	XMLName             xml.Name `xml:"NewListInformation"`
	Text                string   `xml:",chardata"`
	ClubName            string   `xml:"ClubName"`       // Hull City
	ClubWebsiteURL      string   `xml:"ClubWebsiteURL"` // https://www.wearehullcity...
	NewsletterNewsItems struct {
		Text               string        `xml:",chardata"`
		NewsletterNewsItem []NewsElement `xml:"NewsletterNewsItem"`
	} `xml:"NewsletterNewsItems"`
}

/*
NewsElement is inner element merged
from NewsList inner element
and NewsDetailed inner element.
*/
type NewsElement struct {
	Text              string `xml:",chardata"`
	ArticleURL        string `xml:"ArticleURL"`        // https://www.wearehullcity...
	NewsArticleID     string `xml:"NewsArticleID"`     // 444541, 444515, 444506, 4...
	PublishDate       string `xml:"PublishDate"`       // 2023-03-04 18:58:00, 2023...
	Taxonomies        string `xml:"Taxonomies"`        // Academy, Academy, Intervi...
	TeaserText        string `xml:"TeaserText"`        // Midfielder Sincere Hall w...
	ThumbnailImageURL string `xml:"ThumbnailImageURL"` // https://www.wearehullcity...
	Title             string `xml:"Title"`             // Hall: ‘Really happy wit...
	OptaMatchId       string `xml:"OptaMatchId"`       // g2322054, g2322054, g2300...
	LastUpdateDate    string `xml:"LastUpdateDate"`    // 2023-03-05 02:00:11, 2023...
	IsPublished       string `xml:"IsPublished"`       // True, True, True, True, T...
	Subtitle          string `xml:"Subtitle"`
	BodyText          string `xml:"BodyText"` // <p class="x_MsoNormal">Be...
	GalleryImageURLs  string `xml:"GalleryImageURLs"`
	VideoURL          string `xml:"VideoURL"`
}

/*
NewsDetailed is struct for getting detailed single news from feed.

It was generated with the help of:

	go install github.com/miku/zek/cmd/zek@latest
	cat examples/hullcityDetailed.xml | zek -e

The inner news element was refactored into separate struct NewsElement.
*/
type NewsDetailed struct {
	XMLName        xml.Name    `xml:"NewsArticleInformation"`
	Text           string      `xml:",chardata"`
	ClubName       string      `xml:"ClubName"`       // Hull City
	ClubWebsiteURL string      `xml:"ClubWebsiteURL"` // https://www.wearehullcity...
	NewsArticle    NewsElement `xml:"NewsArticle"`
}

// NewsPublishedDateLayout is the layout for date time in gathered news
// remember that layout needs to be pointing to Jan 2, 2006 at 3:04pm (MST) in expected format
const NewsPublishedDateLayout = "2006-01-02 15:04:05"

/*
GetArticleFromNewsElement creates Article
from NewsElement taken as a value so that it can create pointers to its fields
*/
func GetArticleFromNewsElement(n NewsElement, teamId string, hasDetails bool) (types.Article, error) {
	publishedDate, err := time.Parse(NewsPublishedDateLayout, n.PublishDate)
	if err != nil {
		return types.Article{}, err
	}
	article := types.Article{
		ArticleKey: types.ArticleKey{
			Published: publishedDate,
			TeamId:    teamId,
			NewsId:    n.NewsArticleID,
		},
		Content:     n.BodyText,
		GalleryUrls: n.GalleryImageURLs,
		ImageURL:    n.ThumbnailImageURL,
		OptaMatchId: n.OptaMatchId,
		Teaser:      n.TeaserText,
		Title:       n.Title,
		Type:        strings.FieldsFunc(n.Taxonomies, splitTaxonomies),
		URL:         n.ArticleURL,
		VideoURL:    n.VideoURL,
		HasDetails:  hasDetails,
	}
	err = article.SetGeneratedId()
	if err != nil {
		return types.Article{}, err
	}
	return article, nil
}

func splitTaxonomies(r rune) bool {
	return r == ',' || r == ':'
}
