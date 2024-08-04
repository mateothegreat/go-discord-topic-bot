package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mateothegreat/go-discord-topic-bot/config"
	"github.com/mateothegreat/go-discord-topic-bot/suggestions"
	"github.com/mateothegreat/go-discord-topic-bot/util"
	"github.com/mateothegreat/go-multilog/multilog"
)

func init() {
	AddHandler("topics-suggest", TopicsSuggestCreator)
}

func TopicsSuggestCreator(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: util.CreateCustomInteractionID("topics_suggest", i.Interaction),
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
							MinLength:   1,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "description",
							Label:     "Describe your topic",
							Style:     discordgo.TextInputParagraph,
							Required:  false,
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

	AddResponder(util.CreateCustomInteractionID("topics_suggest", i.Interaction), Responder{
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour * 24),
		Fn:      TopicsSuggestResponder,
	})
}

func TopicsSuggestResponder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you for taking your time to suggest a topic, we'll check it out in a bit!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		multilog.Error("TopicsSuggestResponder", "respond to interaction", map[string]interface{}{
			"error": err,
		})
		return
	}

	data := i.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, "topics_suggest") {
		return
	}

	suggestion, err := suggestions.Create(suggestions.CreateArgs{
		UserID:      i.Member.User.ID,
		UserName:    i.Member.User.Username,
		Title:       data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		Description: data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	})
	if err != nil {
		multilog.Error("TopicsSuggestResponder", "create suggestion", map[string]interface{}{
			"error": err,
		})
		return

	}
	_, err = s.ChannelMessageSendComplex(config.Config.TopicsSuggestionsChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("New topic suggestion received from <@%s>", i.Member.User.ID),
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: "New Topic Suggestion",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Title",
						Value:  suggestion.Title,
						Inline: false,
					},
					{
						Name:   "Description",
						Value:  suggestion.Description,
						Inline: false,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Suggested at %s", time.Now().Format(time.RFC1123)),
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Approve",
						Style:    discordgo.PrimaryButton,
						CustomID: `{"foo": "bar"}`,
					},
					discordgo.Button{
						Label:    "Reject",
						Style:    discordgo.DangerButton,
						CustomID: "reject_button",
					},
				},
			},
		},
	})
	if err != nil {
		multilog.Error("TopicsSuggestResponder", "send message", map[string]interface{}{
			"error": err,
		})
	}
}
