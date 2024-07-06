package commands

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Commands is a map of command names to commands.
var Commands map[string]discordgo.ApplicationCommand = make(map[string]discordgo.ApplicationCommand)

// Handlers is a map of command names to handlers.
var Handlers map[string]Handler = make(map[string]Handler)

// Responders is a map of custom IDs to responders.
var Responders map[string]Responder = make(map[string]Responder)

// Handler maps a command to a handler.
type Handler struct {
	// Fn is the method to create the command.
	Fn func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Responder is a handler that responds to an interaction.
type Responder struct {
	// Created is the time the responder was created.
	Created time.Time
	// Expires is the time the responder expires.
	Expires time.Time
	// Fn is the method to handle the interaction.
	Fn func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// CustomID is a custom ID for mapping interactions to a specific handler.
type CustomID struct {
	// Action is the action to map to.
	Action string `json:"a,omitempty"`
	// Context is typically a user id.
	Context string `json:"c,omitempty"`
	// Meta is any additional data.
	Meta interface{} `json:"m,omitempty"`
}

// JSON returns the custom ID as a JSON string.
func (c *CustomID) JSON() string {
	json, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(json)
}

// AddCommand adds a command to the commands map.
//
// Arguments:
//   - c: The command to add.
func AddCommand(c discordgo.ApplicationCommand) {
	Commands[c.Name] = c
}

// AddHandler adds a handler to the handlers map.
//
// Arguments:
//   - commandName: The name of the command to add.
//   - creator: The method to create the command.
func AddHandler(commandName string, Fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	Handlers[commandName] = Handler{
		Fn: Fn,
	}
}

// AddResponder adds a responder to the responders map.
//
// Arguments:
//   - customID: The custom ID to add.
//   - responder: The responder to add.
func AddResponder(customID string, responder Responder) {
	Responders[customID] = responder
}
