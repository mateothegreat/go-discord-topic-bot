package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mateothegreat/go-discord-topic-bot/commands"
	"github.com/mateothegreat/go-multilog/multilog"
)

var s *discordgo.Session

func init() {
	flag.Parse()
}

func init() {
	godotenv.Load()
	var err error
	s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
	}))

	ResultsChannel := os.Getenv("DISCORD_CHANNEL_ID")
	AppID := os.Getenv("DISCORD_APP_ID")
	GuildID := os.Getenv("DISCORD_GUILD_ID")

	// Handle the bot being ready.
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		multilog.Info("main", "bot is up", nil)
	})

	// Handle interactions.
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Handle slash commands.
		case discordgo.InteractionApplicationCommand:
			multilog.Debug("main", "interaction", map[string]interface{}{
				"type":  i.Type,
				"name:": fmt.Sprintf("%s-%s", i.ApplicationCommandData().Name, i.ApplicationCommandData().Options[0].Name),
			})
			commands.Handlers[fmt.Sprintf("%s-%s", i.ApplicationCommandData().Name, i.ApplicationCommandData().Options[0].Name)](s, i)
		// Handle modal submissions.
		case discordgo.InteractionModalSubmit:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Thank you for taking your time to suggest a topic, we'll check it out in a bit!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				multilog.Error("main", "respond to interaction", map[string]interface{}{
					"error": err,
				})
			}

			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "topics_suggest") {
				return
			}

			userid := strings.Split(data.CustomID, "_")[2]
			_, err = s.ChannelMessageSend(ResultsChannel, fmt.Sprintf(
				"New topic suggestion received. From <@%s>\n\n**Title**:\n%s\n\n**Description**:\n%s",
				userid,
				data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			))
			if err != nil {
				panic(err)
			}
		}
	})

	// IDs of the commands we've registered with discord for tear down on exit.
	discordCommandIDs := make(map[string]string, len(commands.Commands))

	// Register the commands with discord.
	for _, cmd := range commands.Commands {
		rcmd, err := s.ApplicationCommandCreate(AppID, GuildID, &cmd)
		if err != nil {
			multilog.Error("main", "register discord command", map[string]interface{}{
				"error": err,
			})
		}

		discordCommandIDs[rcmd.ID] = rcmd.Name
	}

	// Open the discord session.
	err := s.Open()
	if err != nil {
		multilog.Fatal("main", "open discord session", map[string]interface{}{
			"error": err,
		})
	}
	defer s.Close()

	// Wait for the user to close the program.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	multilog.Info("main", "gracefully shutting down", map[string]interface{}{})

	// Delete the commands we registered with discord above.
	for id, _ := range discordCommandIDs {
		err := s.ApplicationCommandDelete(AppID, GuildID, id)
		if err != nil {
			multilog.Error("main", "delete discord command", map[string]interface{}{
				"error": err,
			})
		}
	}

}
