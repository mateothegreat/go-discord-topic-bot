package commands

import "github.com/bwmarrin/discordgo"

var Commands map[string]discordgo.ApplicationCommand = make(map[string]discordgo.ApplicationCommand)
var Handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

func Add(c discordgo.ApplicationCommand) {
	Commands[c.Name] = c
}

func AddHandler(commandName string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	Handlers[commandName] = handler
}
