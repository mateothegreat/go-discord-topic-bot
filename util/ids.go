package util

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CreateCustomInteractionID(prefix string, interaction *discordgo.Interaction, suffixes ...string) string {
	customID := fmt.Sprintf("%s_%s", prefix, interaction.Member.User.ID)
	for _, suffix := range suffixes {
		customID = fmt.Sprintf("%s_%s", customID, suffix)
	}
	return customID
}

func ParseCustomInteractionID(prefix string, customID string) ([]string, error) {
	trimmedID := strings.TrimPrefix(customID, prefix+"_")
	if trimmedID == customID {
		return nil, fmt.Errorf("invalid prefix")
	}
	return strings.Split(trimmedID, "_"), nil
}
