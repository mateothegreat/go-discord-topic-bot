package suggestions

import (
	"context"

	"github.com/mateothegreat/go-discord-topic-bot/database"
	"github.com/mateothegreat/go-discord-topic-bot/prisma/db"
)

type ListArgs struct {
	TopicID string
	Page    int
}

func List(args ListArgs) ([]db.TopicModel, error) {
	res, err := database.DatabaseClient.Topic.FindMany(
		db.Topic.ID.Equals(args.TopicID),
	).Take(10).Skip(args.Page * 10).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
