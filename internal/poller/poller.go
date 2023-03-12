package poller

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/go-logr/logr"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"
)

type Config struct {
	TeamId        string `mapstructure:"teamId"`
	RunOnceAtBoot bool   `mapstructure:"runOnceAtBoot"`
	List          struct {
		URL      string `mapstructure:"url"`
		Count    int    `mapstructure:"count"`
		Schedule string `mapstructure:"schedule"`
	} `mapstructure:"list"`
	Details struct {
		URL      string `mapstructure:"url"`
		Schedule string `mapstructure:"schedule"`
	} `mapstructure:"details"`
}

/*
StartPollerWithConfigFile starts cron jobs based on config file.
The jobs then queries news provider for list of new news and their details
*/
func StartPollerWithConfigFile(
	ctx context.Context,
	v *viper.Viper,
	logger logr.Logger,
	s storage.ArticleStorage,
) error {
	// Unmarshal the poller config
	var pollerConfig Config
	err := v.Unmarshal(&pollerConfig)
	if err != nil {
		return fmt.Errorf("error unmarshaling poller config: %w", err)
	}
	logger = logger.WithValues("workerKind", "NewsPoller")
	logger.Info("Starting poller with this config.", "config", pollerConfig)
	c := cron.New()
	err = c.AddFunc(pollerConfig.List.Schedule, func() {
		PollNewsListIntoStorage(ctx, pollerConfig.List.URL, pollerConfig.List.Count, logger, s)
	})
	if err != nil {
		return fmt.Errorf("error adding PollNewsListIntoStorage to cron: %w", err)
	}
	err = c.AddFunc(pollerConfig.Details.Schedule, func() {
		PollNewsDetailsIntoStorage(ctx, pollerConfig.Details.URL, logger, s)
	})
	if err != nil {
		return fmt.Errorf("error adding PollNewsDetailsIntoStorage to cron: %w", err)
	}
	if pollerConfig.RunOnceAtBoot {
		logger.Info("Running jobs for the first time.")
		for _, e := range c.Entries() {
			e.Job.Run()
		}
	}
	logger.Info("Starting the scheduler.")
	c.Start()
	return nil
}

func PollNewsDetailsIntoStorage(ctx context.Context, url string, logger logr.Logger, s storage.ArticleStorage) {
	logger = logger.WithValues("workerJob", "DetailsPolling")
	logger = logger.WithValues("url", url)
	logger.Info("Getting news IDs that don't have details filled in.")
	ids, err := s.GetNewsWithoutDetailsIDs()
	if err != nil {
		logger.Error(err, "Could not get IDs of news that needs to get details from storage.")
		return
	}
	if ids == nil {
		logger.Info("There are no news IDs to get details of.")
		return
	}
	for _, id := range ids {
		err := PollNewsDetailsIntoStorageOfGivenID(ctx, url, logger, s, id)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Fail when polling details of newsId %v", id))
			return
		}
	}
	logger.Info("Finished polling details of all newses.")
}

func PollNewsDetailsIntoStorageOfGivenID(ctx context.Context, url string, logger logr.Logger, s storage.ArticleStorage, newsId string) error {
	logger = logger.WithValues("newsId", newsId)
	logger.Info("Starting to poll detailed news.")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error(err, "Failed to create new GET request.")
		return nil
	}
	q := req.URL.Query()
	q.Add("id", newsId)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err, "Failed to do http request.")
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err, "Error during close of response.")
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Errorf("%v", resp.StatusCode), "Status error")
		return nil
	}
	dec := xml.NewDecoder(resp.Body)
	dec.Strict = false
	var news NewsDetailed
	err = dec.Decode(&news)
	if err != nil {
		logger.Error(err, "Could not decode response.")
		return nil
	}
	article, err := GetArticleFromNewsElement(news.NewsArticle, "t94", true)
	if err != nil {
		logger.Error(err, fmt.Sprintf("Could not parse article from news %v", news.NewsArticle))
		return nil
	}
	err = s.Write(article)
	if err != nil {
		if errors.Is(err, storage.ArticleAlreadyExists) || errors.Is(err, storage.ArticleWriteFailed) {
			logger.Error(err, fmt.Sprintf("Could not write article with newsId %v", article.Id))
			return nil
		}
		return err
	}
	logger.Info("Saved article from detailed news.", "articleID", article.Id, "newsId", newsId)
	return nil
}

func PollNewsListIntoStorage(ctx context.Context, url string, count int, logger logr.Logger, s storage.ArticleStorage) {
	logger = logger.WithValues("workerJob", "ListPolling")
	logger = logger.WithValues("pollerNewsListCount", count)
	logger = logger.WithValues("url", url)
	logger.Info("Starting to poll news.")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error(err, "Failed to create new GET request.")
		return
	}
	q := req.URL.Query()
	q.Add("Count", strconv.Itoa(count))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err, "Failed to do http request.")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err, "Error during close of response.")
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Errorf("%v", resp.StatusCode), "Status error")
		return
	}
	dec := xml.NewDecoder(resp.Body)
	dec.Strict = false
	var news NewsList
	err = dec.Decode(&news)
	if err != nil {
		logger.Error(err, "Could not decode response.")
		return
	}
	logger.Info("Polled news.", "newsAmount", len(news.NewsletterNewsItems.NewsletterNewsItem))
	for _, v := range news.NewsletterNewsItems.NewsletterNewsItem {
		article, err := GetArticleFromNewsElement(v, "t94", false)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Could not parse article from news %v", v))
			continue
		}
		err = s.Write(article)
		if err != nil {
			if errors.Is(err, storage.ArticleAlreadyExists) {
				continue
			}
			if errors.Is(err, storage.ArticleWriteFailed) {
				logger.Error(err, fmt.Sprintf("Could not write article with id %v", article.Id))
				continue
			}
			logger.Error(err, fmt.Sprintf("Fail when processing article %v", v))
			return
		} else {
			logger.Info("Saved article from listed news.", "articleID", article.Id, "newsId", v.NewsArticleID)
		}
	}
	logger.Info("Finished polling and saving news.")
}
