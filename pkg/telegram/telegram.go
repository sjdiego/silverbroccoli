package telegram

import (
	"context"
	"log"
	"silverbroccoli/config"
	"silverbroccoli/pkg/files"
	"silverbroccoli/pkg/reddit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

// PublishPosts function
func PublishPosts(
	ctx context.Context,
	client *mongo.Client,
	cfg config.Config,
	channel Channel,
	posts []reddit.Post,
) {
	bot, err := tgbotapi.NewBotAPI(config.ReadConfig().Telegram.BotKey)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("[+] Telegram API authorized on account '%s'", bot.Self.UserName)

	for _, post := range posts {
		if post.IsVideo {
			if SendVideomessage(*bot, cfg, channel, post) {
				files.Delete(post)
				reddit.SetPostAsPublished(ctx, client, cfg, post)
			}
		} else {
			if SendMessage(*bot, cfg, channel, post) {
				reddit.SetPostAsPublished(ctx, client, cfg, post)
			}
		}
	}
}
