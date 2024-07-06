package topics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mateothegreat/go-discord-topic-bot/commands"
)

func init() {
	commands.AddCommand(discordgo.ApplicationCommand{
		Name:        "topics",
		Description: "Topics",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "suggest",
				Description: "Suggest a topic",
			},
		},
	})
}
