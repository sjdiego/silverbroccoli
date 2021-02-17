package telegram

import (
	"log"
	"silverbroccoli/config"
	"silverbroccoli/pkg/reddit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendMessage func
func SendMessage(bot tgbotapi.BotAPI, cfg config.Config, channel Channel, post reddit.Post) bool {
	log.Printf("[*] Sending message %s to channel %s (%v)", post.Title, channel.Name, channel.ChannelID)

	message := tgbotapi.NewMessage(channel.ChannelID, post.URL)

	_, err := bot.Send(message)

	if err != nil {
		log.Println("[!] Error sending message:", err)
		return false
	}

	log.Println("[+] Message sent successfully")
	return true
}
