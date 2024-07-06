package topics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mateothegreat/go-discord-topic-bot/commands"
)

func init() {
	commands.Add(discordgo.ApplicationCommand{
		Name:        "topics",
		Description: "Topics",
	})
}
