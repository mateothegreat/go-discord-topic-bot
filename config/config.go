package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Config Conf

type Conf struct {
	DiscordToken               string
	TopicsSuggestionsChannelID string
}

func init() {
	godotenv.Load()
	Config = Conf{
		DiscordToken:               os.Getenv("DISCORD_TOKEN"),
		TopicsSuggestionsChannelID: os.Getenv("TOPICS_SUGGESTIONS_CHANNEL_ID"),
	}
}
