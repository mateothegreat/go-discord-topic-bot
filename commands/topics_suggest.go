package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mateothegreat/go-multilog/multilog"
)

func init() {
	AddHandler("topics-suggest", TopicsSuggest)
}

func TopicsSuggest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "topics_suggest_" + i.Interaction.Member.User.ID,
			Title:    "Suggest a new topic!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "title",
							Label:       "Title",
							Style:       discordgo.TextInputShort,
							Placeholder: "My first topic",
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "description",
							Label:     "Describe your topic",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 4000,
						},
					},
				},
			},
		},
	})
	if err != nil {
		multilog.Error("topics-suggest", "respond to interaction", map[string]interface{}{
			"error": err,
		})
	}
}
