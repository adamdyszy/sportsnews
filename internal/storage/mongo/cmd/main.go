package main

import (
	"context"
	"fmt"
	"github.com/adamdyszy/sportsnews/internal/poller"
	storage "github.com/adamdyszy/sportsnews/internal/storage/mongo"
	"github.com/adamdyszy/sportsnews/types"
	"github.com/spf13/viper"
	"reflect"
)

func main() {
	v := viper.New()
	v.SetConfigFile("config/debug.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	s, err := storage.NewMongoStorage(v.Sub("db"), context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		err := s.Disconnect()
		if err != nil {
			fmt.Printf("Error during disconnect in storage: %s\n", err)
		}
	}()
	a, _ := poller.GetArticleFromNewsElement(poller.NewsElement{
		Text:          "test",
		NewsArticleID: "1",
		PublishDate:   "2023-02-17 14:20:33",
		Taxonomies:    "T1,T2",
		Title:         "Title",
	}, "t94", false)

	// delete if was already present
	err = s.Delete(a.Id)
	if err != nil {
		panic(err)
	}

	// write without details
	err = s.Write(a)
	if err != nil {
		panic(err)
	}

	// check if not detailed is available in list
	withoutDetailsIDs, err := s.GetNewsWithoutDetailsIDs()
	if err != nil {
		panic(err)
	}

	found := false
	for _, id := range withoutDetailsIDs {
		if id == "1" {
			found = true
		}
	}
	if !found {
		panic(fmt.Sprintf("not found 1 in withoutDetailsIDs: %v", withoutDetailsIDs))
	}

	// do override with details
	a.Content = "we now have details!"
	a.HasDetails = true
	err = s.Write(a)
	if err != nil {
		panic(err)
	}

	// check override worked
	aFromDB, err := s.Get(a.Id)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(a, aFromDB) {
		panic(fmt.Sprintf("a: %v , fromDB: %v", a, aFromDB))
	}

	// check list
	list, err := s.List()
	if err != nil {
		panic(err)
	}

	found = false
	var lastFound types.Article
	for _, aIter := range list {
		if aIter.NewsId == a.NewsId || aIter.Id == a.Id {
			if found {
				panic(fmt.Sprintf("should only find article once: %v , lastFound: %v, iteratedNow: %v", a, lastFound, aIter))
			}
			found = true
			if !reflect.DeepEqual(a, aIter) {
				panic(fmt.Sprintf("a: %v , iteratedNow: %v", a, aFromDB))
			}
		}
		lastFound = aIter
	}

	if !found {
		panic(fmt.Sprintf("not found 1 in list: %v", list))
	}

	// check if after override we no longer have it without ids
	withoutDetailsIDs, err = s.GetNewsWithoutDetailsIDs()
	if err != nil {
		panic(err)
	}

	for _, id := range withoutDetailsIDs {
		if id == "1" {
			panic(fmt.Sprintf("news Id 1 should be no longer in withoutDetailsIDs: %v", list))
		}
	}

	fmt.Println("Success!")
}
