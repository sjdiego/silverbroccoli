package reddit

import (
	"context"
	"fmt"
	"log"
	"silverbroccoli/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Subreddit data
type Subreddit struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

// GetSubreddits returns a slice of Subreddit structs stored in MongoDB
func GetSubreddits(ctx context.Context, client *mongo.Client, cfg config.Config) []Subreddit {
	var subreddits []Subreddit

	col := client.Database(cfg.MongoData.Database).Collection(config.SubredditsCollection)

	cursor, _ := col.Find(ctx, bson.D{})

	for cursor.Next(ctx) {
		var record Subreddit
		err := cursor.Decode(&record)

		if err != nil {
			log.Println(err)
		} else {
			subreddits = append(subreddits, record)
		}
	}

	return subreddits
}

// AddSubreddit inserts a new document with data of Subreddit struct
func AddSubreddit(ctx context.Context, client *mongo.Client, cfg config.Config) bool {
	var name string

	fmt.Print("- Write the name of the subreddit: ")
	fmt.Scanln(&name)

	col := client.Database(cfg.MongoData.Database).Collection(config.SubredditsCollection)

	subreddit := new(Subreddit)
	subreddit.Name = name

	_, err := col.InsertOne(ctx, subreddit)

	if err != nil {
		log.Println("[!] Error adding subreddit:", err)
		return false
	}

	log.Println("[+] Subreddit added successfully")
	return true
}
