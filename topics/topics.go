package topcis

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

func Create(args CreateArgs) (*db.TopicModel, error) {
	res, err := database.DatabaseClient.Topic.CreateOne(
		db.Topic.
			db.Topic.Title.Set(args.Title),
		db.Topic.Description.Set(args.Description),
		db.Topic.Title.Set(args.Title),
		db.Topic.Description.Set(args.Description),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}
