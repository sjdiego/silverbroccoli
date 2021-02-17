package telegram

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"silverbroccoli/config"
	"silverbroccoli/pkg/reddit"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Channel struct
type Channel struct {
	ID         primitive.ObjectID `bson:"_id"`
	ChannelID  int64              `bson:"id"`
	Name       string             `bson:"name"`
	Caption    string             `bson:"caption"`
	IsActive   bool               `bson:"is_active"`
	Subreddits []reddit.Subreddit `bson:"subreddits"`
}

// GetChannels function
func GetChannels(ctx context.Context, client *mongo.Client, cfg config.Config) []Channel {
	var channels []Channel
	col := client.Database(cfg.MongoData.Database).Collection(config.ChannelsCollection)

	cursor, _ := col.Find(ctx, bson.M{"is_active": bsonx.Boolean(true)})

	for cursor.Next(ctx) {
		var record Channel
		err := cursor.Decode(&record)

		if err != nil {
			log.Println(err)
		} else {
			channels = append(channels, record)
		}
	}

	return channels
}

// AddChannel function
func AddChannel(ctx context.Context, client *mongo.Client, cfg config.Config) {
	var name string
	fmt.Print("- Write the name of the channel: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		name = scanner.Text()
	}

	var caption string
	fmt.Print("- Write the caption of the channel: ")
	if scanner.Scan() {
		caption = scanner.Text()
	}

	var id int64
	fmt.Print("- Write the ID of the channel: ")
	fmt.Scanln(&id)

	selectedSubreddits := appendSubreddits(reddit.GetSubreddits(ctx, client, cfg))

	col := client.Database(cfg.MongoData.Database).Collection(config.ChannelsCollection)

	_, err := col.InsertOne(ctx, &Channel{
		ID:         primitive.NewObjectID(),
		ChannelID:  id,
		Name:       name,
		Caption:    caption,
		IsActive:   true,
		Subreddits: selectedSubreddits,
	})

	if err != nil {
		log.Println("[!] Error adding channel:", err)
	} else {
		log.Println("[+] Channel added successfully")
	}
}

func appendSubreddits(storedSubreddits []reddit.Subreddit) []reddit.Subreddit {
	var subreddits []reddit.Subreddit
	var create bool = false

	for create == false {
		fmt.Println("- Add subreddit subscription: ")
		for key, sub := range storedSubreddits {
			fmt.Printf("\t%d. %s\n", key+1, sub.Name)
		}

		var response string
		fmt.Scanln(&response)
		idx, _ := strconv.ParseInt(response, 10, 0)

		fmt.Println("Selected subreddit is:", storedSubreddits[idx-1].Name)
		subreddits = append(subreddits, storedSubreddits[idx-1])

		fmt.Println("- Add another? (y/n)")
		if !askForConfirmation() {
			create = true
		}
	}

	return subreddits
}

func askForConfirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Please write (y)es or (n)o")
		return askForConfirmation()
	}
}
