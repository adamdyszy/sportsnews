package poller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPublishedDate(t *testing.T) {
	n := NewsElement{
		PublishDate: "2023-02-17 14:20:33",
	}
	a, err := GetArticleFromNewsElement(n, "t94", false)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	createdDate := a.Published
	assert.Equal(t, 2023, createdDate.Year())
	assert.Equal(t, time.Month(2), createdDate.Month())
	assert.Equal(t, 17, createdDate.Day())
	assert.Equal(t, 14, createdDate.Hour())
	assert.Equal(t, 20, createdDate.Minute())
	assert.Equal(t, 33, createdDate.Second())
	n = NewsElement{
		PublishDate: "2024-11-22 19:44:51",
	}
	a, err = GetArticleFromNewsElement(n, "t94", false)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	createdDate = a.Published
	assert.Equal(t, 2024, createdDate.Year())
	assert.Equal(t, time.Month(11), createdDate.Month())
	assert.Equal(t, 22, createdDate.Day())
	assert.Equal(t, 19, createdDate.Hour())
	assert.Equal(t, 44, createdDate.Minute())
	assert.Equal(t, 51, createdDate.Second())
}

func TestTaxonomies(t *testing.T) {
	n := NewsElement{
		PublishDate: "2023-02-17 14:20:33",
		Taxonomies:  "Club News",
	}
	a, err := GetArticleFromNewsElement(n, "t94", false)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, []string{"Club News"}, a.Type)

	n = NewsElement{
		PublishDate: "2023-02-17 14:20:33",
		Taxonomies:  "Club News,Something",
	}
	a, err = GetArticleFromNewsElement(n, "t94", false)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, []string{"Club News", "Something"}, a.Type)
}
