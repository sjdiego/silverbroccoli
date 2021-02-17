package telegram

import (
	"log"
	"os"
	"silverbroccoli/config"
	"silverbroccoli/pkg/files"
	"silverbroccoli/pkg/reddit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendVideomessage function
func SendVideomessage(
	bot tgbotapi.BotAPI,
	cfg config.Config,
	channel Channel,
	post reddit.Post,
) bool {
	video := files.CreateVideo(post)

	log.Printf("[*] Sending video %s to channel %s (%v)", post.Title, channel.Name, channel.ChannelID)
	_, err := bot.Send(createVideoMessage(cfg, channel, post, video))

	if err != nil {
		log.Println("[!] Error sending video:", err)
		return false
	}

	log.Println("[+] Video uploaded successfully")
	return true
}

func createVideoMessage(
	cfg config.Config,
	channel Channel,
	post reddit.Post,
	video reddit.Video,
) tgbotapi.VideoConfig {
	videoConfig := tgbotapi.NewVideoUpload(channel.ChannelID, video.FilePath)
	videoConfig.DisableNotification = true

	if len(channel.Caption) > 0 {
		videoConfig.Caption = post.Title + "\n\n" + channel.Caption
	} else {
		videoConfig.Caption = post.Title
	}

	return videoConfig
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
