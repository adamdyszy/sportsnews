package main

import (
	"context"
	"fmt"
	api "github.com/adamdyszy/sportsnews/api/v1"
	"github.com/adamdyszy/sportsnews/internal/poller"
	storageImpl "github.com/adamdyszy/sportsnews/internal/storage/mongo"
	"github.com/go-logr/zapr"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// Create a new Viper configuration object.
	v := viper.GetViper()
	v.SetConfigFile("config/default.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file, error: %s\n", err)
		os.Exit(1)
	}

	// Read the custom configuration values if the custom config file exists.
	v.SetConfigFile("config/custom.yaml")
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Custom config file not found")
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := storageImpl.NewMongoStorage(v.Sub("db"), ctx)
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
	err = api.ListenAndServe(v.Sub("api"), s)
	if err != nil {
		logger.Error(err, "Could not server api.")
		os.Exit(6)
	}
}
