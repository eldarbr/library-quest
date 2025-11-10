package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
)

type Commands struct {
	questRepo     quest.Repository
	keywordSource KeywordProvider
	shuffler      quest.QuestShuffler
}

type Answer struct {
	QuestID int
	Words   []string
}

func NewCommands(keywordSource KeywordProvider, questRepo quest.Repository, shuffler quest.QuestShuffler) Commands {
	return Commands{
		questRepo:     questRepo,
		keywordSource: keywordSource,
		shuffler:      shuffler,
	}
}

func (c Commands) ValidateAnswer(ctx context.Context, teamID int, answer Answer) (string, error) {
	if c.shuffler.DetermineQuestForTeam(teamID) != answer.QuestID {
		return "", ErrQuestIDMismatch
	}

	return c.ValidateAnswerRandom(ctx, answer)
}

func (c Commands) ValidateAnswerRandom(ctx context.Context, answer Answer) (string, error) {
	questObj, err := c.questRepo.GetQuest(ctx, answer.QuestID)
	if err != nil {
		return "", fmt.Errorf("questRepo.GetQuest: %w", err)
	}

	err = questObj.ValidateAnswer(answer.Words)
	if err != nil {
		return "", fmt.Errorf("answer validation: %w", err)
	}

	keyword := c.keywordSource.GetKeyword()

	return keyword, nil
}

type KeywordProvider interface {
	GetKeyword() string
}

var (
	ErrQuestIDMismatch = errors.New("this team is not allowed to do the quest")
)
