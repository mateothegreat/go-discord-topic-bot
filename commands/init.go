package commands

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	AddCommand(discordgo.ApplicationCommand{
		Name:        "topics",
		Description: "Topics",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "topics-suggest",
				Description: "Suggest a topic",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "topics-list",
				Description: "List topics",
			},
		},
	})
}
