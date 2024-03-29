package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	api "github.com/adamdyszy/sportsnews/api/v1"
	"github.com/adamdyszy/sportsnews/internal/poller"
	"github.com/adamdyszy/sportsnews/internal/storage/memory"
	"github.com/adamdyszy/sportsnews/internal/storage/mongo"
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/go-logr/zapr"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// handle args
	var customConfigFile string
	flag.StringVar(&customConfigFile, "customConfigFile", "config/custom.yaml", "Custom config file that will override config/default.yaml")
	flag.Parse()

	// Create a new Viper configuration object.
	v := viper.GetViper()
	v.SetConfigFile("config/default.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file, error: %s\n", err)
		os.Exit(1)
	}

	// Read the custom configuration values if the custom config file exists.
	v.SetConfigFile(customConfigFile)
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Custom config file not found: %v\n", customConfigFile)
			os.Exit(2)
		} else {
			fmt.Printf("Error reading custom config file: %s\n", err)
			os.Exit(2)
		}
	}

	// Create logger
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	// Create a new logger instance from the configuration
	z, err := cfg.Build()
	if err != nil {
		fmt.Printf("Failed to build logger, error: %s\n", err)
		os.Exit(3)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("syncing logger failed %v", err)
		}
	}(z)
	logger := zapr.NewLogger(z)

	var s storage.ArticleStorage
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch v.GetString("storageKind") {
	case "mongo":
		s, err = mongo.NewMongoStorage(v.Sub("mongoStorage"), ctx)
		if err != nil {
			logger.Error(err, "Could not initialize storage.")
			os.Exit(4)
		}
	case "memory", "":
		s = memory.NewMemStorage()
	default:
		err = errors.New("unknown database kind")
	}
	if err != nil {
		logger.Error(err, "Could not initialize storage.")
		os.Exit(4)
	}
	defer func() {
		err := s.Disconnect()
		if err != nil {
			logger.Error(err, "Error during disconnect in storage.")
		}
	}()
	err = poller.StartPollerWithConfigFile(ctx, v.Sub("poller"), logger, s)
	if err != nil {
		logger.Error(err, "Could not start poller.")
		os.Exit(5)
	}
	err = api.ListenAndServe(v.Sub("api"), s, logger)
	if err != nil {
		logger.Error(err, "Could not server api.")
		os.Exit(6)
	}
}
