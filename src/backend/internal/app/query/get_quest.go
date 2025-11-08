package query

import (
	"context"
	"fmt"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/myerrros"
)

type Query struct {
	questRepo quest.Repository
}

func NewQuery(questRepo quest.Repository) Query {
	return Query{questRepo: questRepo}
}

func (q Query) GetQuest(ctx context.Context, teamID int) (quest.Quest, error) {
	if teamID < 0 {
		return quest.Quest{}, myerrros.ErrNotFound
	}

	questID := quest.DetermineQuestForTeam(teamID)

	questObj, err := q.questRepo.GetQuest(ctx, questID)
	if err != nil {
		return quest.Quest{}, fmt.Errorf("questRepo.GetQuest: %w", err)
	}

	return questObj, nil
}
