package filedb

import (
	"context"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/myerrros"
)

var _ quest.Repository = FileToMemoDB{}

func (db FileToMemoDB) GetQuest(_ context.Context, id int) (quest.Quest, error) {
	q, ok := db.quests[id]
	if !ok {
		return quest.Quest{}, myerrros.ErrNotFound
	}

	return q, nil
}
