package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/adamdyszy/sportsnews/storage"
	"github.com/adamdyszy/sportsnews/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongoStorage struct {
	client       *mongo.Client
	database     string
	articlesColl *mongo.Collection
	timeout      time.Duration
}

func NewMongoStorage(v *viper.Viper, ctx context.Context) (storage.ArticleStorage, error) {
	// Load the configuration values into variables.
	dbURI := v.GetString("uri")
	dbName := v.GetString("name")
	user := v.GetString("user")
	password := v.GetString("password")
	articlesCollName := v.GetString("articlesColl")
	timeoutSeconds := v.GetInt("timeoutSeconds")
	timeout := time.Duration(timeoutSeconds) * time.Second

	// Create a new MongoDB client.
	clientOptions := options.Client().ApplyURI(dbURI)
	if user != "" {
		clientOptions.SetAuth(options.Credential{
			Username: user,
			Password: password,
		})
	}
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Connect to the MongoDB server.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB server: %w", err)
	}

	// Ping the MongoDB server to check the connection.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB client: %w", err)
	}

	// Check if collection exists, if not create it
	collection := client.Database(dbName).Collection(articlesCollName)
	if collection == nil {
		return nil, errors.New("error creating mongo collection")
	}

	return &mongoStorage{
		client:       client,
		database:     dbName,
		articlesColl: collection,
		timeout:      timeout,
	}, nil
}

func (m mongoStorage) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (m mongoStorage) Delete(id types.ArticleId) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	_, err := m.articlesColl.DeleteMany(ctx, bson.M{"id": id})
	return err
}

func (m mongoStorage) GetNewsWithoutDetailsIDs() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	filter := bson.M{"hasDetails": false}
	cur, err := m.articlesColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting news without details: %w", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			fmt.Printf("error closing cursor %s", err)
		}
	}(cur, ctx)

	var newsIds []string
	for cur.Next(ctx) {
		var article articleBson
		if err := cur.Decode(&article); err != nil {
			return nil, fmt.Errorf("error decoding news without details: %w", err)
		}
		newsIds = append(newsIds, article.NewsId)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error iterating news without details: %w", err)
	}
	return newsIds, nil
}

func (m mongoStorage) List() ([]types.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var articles []types.Article
	cur, err := m.articlesColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting articles: %w", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			fmt.Printf("error closing cursor %s", err)
		}
	}(cur, ctx)

	for cur.Next(ctx) {
		var article articleBson
		if err := cur.Decode(&article); err != nil {
			return nil, fmt.Errorf("error decoding article: %w", err)
		}
		articles = append(articles, article.ToArticle())
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %w", err)
	}
	return articles, nil
}

func (m mongoStorage) Get(id types.ArticleId) (types.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	filter := bson.M{"id": id}
	var article articleBson
	err := m.articlesColl.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Article{}, fmt.Errorf("%w with id: %v", storage.ArticleNotFound, id)
		}
		return types.Article{}, err
	}

	return article.ToArticle(), nil
}

func (m mongoStorage) Write(article types.Article) error {
	foundArticle, err := m.Get(article.Id)
	override := false
	if err != nil {
		if !errors.Is(err, storage.ArticleNotFound) {
			return err
		}
		// if not found we are ok with creation
	} else {
		// if found
		override = true
		if !article.HasDetails {
			// if our new article doesn't have details
			return fmt.Errorf("%w with articleID %v with NewsId %v", storage.ArticleAlreadyExists, foundArticle.Id, foundArticle.NewsId)
		}
		if foundArticle.HasDetails {
			// if our old article already had details
			return fmt.Errorf("%w with articleID %v with NewsId %v, but wanted to override with NewsId %v", storage.ArticleAlreadyExists, foundArticle.Id, foundArticle.NewsId, article.NewsId)
		}
	}

	// Delete previous articles with that Id
	if override {
		err = m.Delete(article.Id)
		if err != nil {
			return err
		}
	}

	// Insert the article into the collection
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	_, err = m.articlesColl.InsertOne(ctx, fromArticle(article))
	if err != nil {
		return fmt.Errorf("%w with id: %v", storage.ArticleWriteFailed, article.Id)
	}
	return nil
}
