package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

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
			commands.Handlers[fmt.Sprintf("%s-%s", i.ApplicationCommandData().Name, i.ApplicationCommandData().Options[0].Name)].Fn(s, i)
		// Handle button clicks.
		case discordgo.InteractionMessageComponent:
			commands.Responders[i.Interaction.Data.(discordgo.MessageComponentInteractionData).CustomID].Fn(s, i)
		// Handle modal submissions.
		case discordgo.InteractionModalSubmit:
			commands.Responders[i.Interaction.Data.(discordgo.ModalSubmitInteractionData).CustomID].Fn(s, i)
		}
	})

	// IDs of the commands we've registered with discord for tear down on exit.
	createdCommandIDs := make([]string, len(commands.Commands))

	// Register the commands with discord.
	for _, cmd := range commands.Commands {
		res, err := s.ApplicationCommandCreate(AppID, GuildID, &cmd)
		if err != nil {
			multilog.Error("main", "register discord command", map[string]interface{}{
				"error": err,
			})
		}
		createdCommandIDs = append(createdCommandIDs, res.ID)
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
	for _, id := range createdCommandIDs {
		err := s.ApplicationCommandDelete(AppID, GuildID, id)
		if err != nil {
			multilog.Error("main", "delete discord command", map[string]interface{}{
				"error": err,
			})
		}
	}

}
