package suggestions

import (
	"context"

	"github.com/mateothegreat/go-discord-topic-bot/database"
	"github.com/mateothegreat/go-discord-topic-bot/prisma/db"
)

type CreateArgs struct {
	UserID      string
	UserName    string
	Title       string
	Description string
}

func Create(args CreateArgs) (*db.SuggestionModel, error) {
	res, err := database.DatabaseClient.Suggestion.CreateOne(
		db.Suggestion.UserID.Set(args.UserID),
		db.Suggestion.UserName.Set(args.UserName),
		db.Suggestion.Title.Set(args.Title),
		db.Suggestion.Description.Set(args.Description),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
