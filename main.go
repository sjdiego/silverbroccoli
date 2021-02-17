package main

import (
	"flag"
	"silverbroccoli/config"
	"silverbroccoli/pkg/reddit"
	"silverbroccoli/pkg/telegram"
)

func main() {
	getPosts := flag.Bool("fetch", false, "Retrieve and save new posts from Reddit")
	publishPosts := flag.Bool("publish", false, "Send unpublished posts to Telegram")
	publishQty := flag.Int64("num", 3, "Number of posts to publish")
	addSubreddits := flag.Bool("addsub", false, "Add new subreddit to fetch posts")
	addChannel := flag.Bool("addchan", false, "Add new channel to send posts")

	flag.Parse()

	cfg := config.ReadConfig()
	client, ctx := config.InitMongo(cfg)

	if *addSubreddits {
		reddit.AddSubreddit(ctx, client, cfg)
	}

	if *addChannel {
		telegram.AddChannel(ctx, client, cfg)
	}

	if *getPosts {
		for _, subreddit := range reddit.GetSubreddits(ctx, client, cfg) {
			posts := reddit.Fetch(subreddit.Name)
			reddit.StorePosts(ctx, client, cfg, posts)
		}
	}

	if *publishPosts {
		for _, channel := range telegram.GetChannels(ctx, client, cfg) {
			for _, subreddit := range channel.Subreddits {
				posts := reddit.GetPosts(ctx, client, cfg, subreddit.Name, *publishQty)

				if len(posts) > 0 {
					telegram.PublishPosts(ctx, client, cfg, channel, posts)
				}
			}
		}
	}

	client.Disconnect(ctx)
}
