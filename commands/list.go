package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mateothegreat/go-multilog/multilog"
)

type ActionMessage struct {
	Created   time.Time
	Expires   time.Time
	AuthorID  string
	MessageID string
	Page      int
}

var actionMessages = make(map[string]ActionMessage)

func init() {
	AddHandler("topics-list", TopicsListCreator)
}

func TopicsListCreator(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("New topic suggestion received from <@%s>", i.Member.User.ID),
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: "List of topics",

				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Description",
						Value:  "Descriptionasdfasdfasdf",
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
						Label:    "Page 1",
						Style:    discordgo.PrimaryButton,
						CustomID: `topics-list`,
					},
					discordgo.Button{
						Label:    "Page 2",
						Style:    discordgo.SecondaryButton,
						CustomID: `topics-listage`,
					},
				},
			},
		},
	})
	if err != nil {
		multilog.Error("TopicsListCreator", "send message", map[string]interface{}{
			"error": err,
		})
	}

	actionMessages[i.Interaction.ID] = ActionMessage{
		Created:   time.Now(),
		Expires:   time.Now().Add(time.Hour * 24),
		AuthorID:  i.Member.User.ID,
		MessageID: i.Interaction.ID,
		Page:      1,
	}

	AddResponder("topics-list", Responder{
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour * 24),
		Fn:      TopicsListResponder,
	})
}

func TopicsListResponder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	actionMessage, ok := actionMessages[i.Interaction.ID]
	if !ok {
		return
	}

	content := fmt.Sprintf("asdf topic suggestion received from <@%s>", actionMessage.AuthorID)

	_, err := s.ChannelMessageEditComplex(discordgo.NewMessageEdit(i.ChannelID, actionMessage.MessageID).SetContent(content))
	if err != nil {
		multilog.Error("TopicsListResponder", "edit message", map[string]interface{}{
			"error": err,
		})
	}
}
