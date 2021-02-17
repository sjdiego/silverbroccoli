package config

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v2"
)

const (
	// PostsCollection name collection for posts documents
	PostsCollection string = "posts"
	// SubredditsCollection name collection for subreddits documents
	SubredditsCollection string = "subreddits"
	// ChannelsCollection name collection for channels documents
	ChannelsCollection string = "channels"
)

// Config data
type Config struct {
	Env        string `yaml:"env"`
	Subreddits string `yaml:"subreddits"`
	MongoData  struct {
		ConnectionString string `yaml:"connectionString"`
		Database         string `yaml:"database"`
	} `yaml:"mongoData"`
	Telegram struct {
		BotKey string `yaml:"botKey"`
	} `yaml:"telegram"`
}

// ReadConfig data
func ReadConfig() Config {
	filename, _ := filepath.Abs("config/config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

// InitMongo data
func InitMongo(cfg Config) (client *mongo.Client, ctx context.Context) {
	log.Println("[*] Connecting to MongoDB...")

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoData.ConnectionString))
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("[!] Could not ping MongoDB service:", err)
		os.Exit(1)
	}
	log.Println("[+] Connected to MongoDB!")

	createIndexes(ctx, client, cfg)

	return client, ctx
}

func createIndexes(ctx context.Context, client *mongo.Client, cfg Config) {

	idxPosts := []mongo.IndexModel{
		{
			Keys:    bson.M{"id": bson.TypeInt32},
			Options: options.Index().SetBackground(true).SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.M{"isvideo": bson.TypeBoolean, "published": bson.TypeBoolean},
			Options: options.Index().SetBackground(true),
		},
	}

	idxSubreddits := []mongo.IndexModel{
		{
			Keys:    bson.M{"id": bson.TypeInt64},
			Options: options.Index().SetBackground(true).SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.M{"subreddits": bson.TypeObjectID},
			Options: options.Index().SetBackground(true),
		},
	}

	idxChannels := mongo.IndexModel{
		Keys:    bson.M{"name": bson.TypeString},
		Options: options.Index().SetBackground(true).SetUnique(true).SetSparse(true),
	}

	client.Database(cfg.MongoData.Database).
		Collection(PostsCollection).
		Indexes().CreateMany(ctx, idxPosts)

	client.Database(cfg.MongoData.Database).
		Collection(SubredditsCollection).
		Indexes().CreateMany(ctx, idxSubreddits)

	client.Database(cfg.MongoData.Database).
		Collection(ChannelsCollection).
		Indexes().CreateOne(ctx, idxChannels)
}
