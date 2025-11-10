package query

import (
	"context"
	"fmt"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/myerrros"
)

type Query struct {
	questRepo quest.Repository
	shuffler  quest.QuestShuffler
	randomer  quest.RandomQuestShuffler
}

func NewQuery(questRepo quest.Repository, shuffler quest.QuestShuffler, randomer quest.RandomQuestShuffler) Query {
	return Query{questRepo: questRepo, shuffler: shuffler, randomer: randomer}
}

func (q Query) GetQuest(ctx context.Context, teamID int) (quest.Quest, error) {
	if teamID < 0 {
		return quest.Quest{}, myerrros.ErrNotFound
	}

	questID := q.shuffler.DetermineQuestForTeam(teamID)

	questObj, err := q.questRepo.GetQuest(ctx, questID)
	if err != nil {
		return quest.Quest{}, fmt.Errorf("questRepo.GetQuest: %w", err)
	}

	return questObj, nil
}

func (q Query) GetRandomQuest(ctx context.Context) (quest.Quest, error) {
	questID := q.randomer.GetRandomQuestID()

	questObj, err := q.questRepo.GetQuest(ctx, questID)
	if err != nil {
		return quest.Quest{}, fmt.Errorf("questRepo.GetQuest: %w", err)
	}

	return questObj, nil
}
