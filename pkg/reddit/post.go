package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"silverbroccoli/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	limit    int    = 10
	category string = "top" // top, new, rising, hot
	endpoint string = "https://www.reddit.com/r/%s/%s.json?limit=%d"
)

var (
	httpClient http.Client = http.Client{Timeout: 3 * time.Second}
	userAgent  string      = "golang http client"
)

// Listing struct
type Listing struct {
	Data struct {
		Children []struct {
			Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// Post struct
type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Subreddit string `json:"subreddit"`
	URL       string `json:"url"`
	IsVideo   bool   `json:"is_video"`
	Media     struct {
		Video `json:"reddit_video,omitempty"`
	} `json:"media,omitempty"`
	Created   float32 `json:"created"`
	Published bool
}

// Video struct
type Video struct {
	Height    int `json:"height"`
	VideoURL  string
	AudioURL  string
	VideoPath string
	AudioPath string
	FilePath  string
}

// Fetch performs a requests to the endpoint and expects a JSON which will be converted into slice of Post structs
func Fetch(subreddit string) []Post {
	url := fmt.Sprintf(endpoint, subreddit, category, limit)
	log.Println(fmt.Sprintf("[+] Making request to %s ...", url))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", userAgent)

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode != 200 {
		log.Fatal("Error", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var results Listing
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var posts []Post
	for _, p := range results.Data.Children {
		posts = append(posts, p.Post)
	}

	return posts
}

// GetPosts returns a slice of Post struct with results of stored items in MongoDB collection
func GetPosts(
	ctx context.Context,
	client *mongo.Client,
	cfg config.Config,
	subreddit string,
	limit int64,
) []Post {
	var posts []Post

	col := client.Database(cfg.MongoData.Database).Collection(config.PostsCollection)

	conditions := bsonx.MDoc{
		"subreddit": bsonx.String(subreddit),
		"published": bsonx.Boolean(false),
	}
	options := options.Find().SetLimit(limit)

	cursor, _ := col.Find(ctx, conditions, options)

	for cursor.Next(ctx) {
		var p Post
		cursor.Decode(&p)
		posts = append(posts, p)
	}

	return posts
}

// StorePosts inserts a new document into MongoDB collection with data of Post structs
func StorePosts(ctx context.Context, client *mongo.Client, cfg config.Config, posts []Post) {
	col := client.Database(cfg.MongoData.Database).Collection(config.PostsCollection)
	added := 0

	for _, item := range posts {
		item.Published = false
		_, err := col.InsertOne(ctx, item)
		if err != nil {
			if cfg.Env == "dev" {
				log.Println("[!]", err)
			}
		} else {
			added++
		}
	}

	log.Printf("[*] Added %v to collection", added)
}

// SetPostAsPublished looks for stored record and sets 'published' parameter as true
func SetPostAsPublished(ctx context.Context, client *mongo.Client, cfg config.Config, post Post) bool {
	col := client.Database(cfg.MongoData.Database).Collection(config.PostsCollection)

	filter := bson.M{"id": bsonx.String(post.ID)}
	update := bson.M{
		"$set": bson.M{"published": bsonx.Boolean(true)},
	}

	result := col.FindOneAndUpdate(ctx, filter, update)

	if result.Err() != nil {
		log.Println("[!] Error on update sent post:", result.Err())
		return false
	}

	log.Println("[*] Updated Published flag on post", post.ID)
	return true
}
