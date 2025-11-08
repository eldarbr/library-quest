package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
)

type Commands struct {
	questRepo quest.Repository
}

type Answer struct {
	QuestID int
	TeamID  int
	Words   []string
}

func NewCommands(questRepo quest.Repository) Commands {
	return Commands{questRepo: questRepo}
}

func (c Commands) ValidateAnswer(ctx context.Context, answer Answer) error {
	if quest.DetermineQuestForTeam(answer.TeamID) != answer.QuestID {
		return ErrQuestIDMismatch
	}

	questObj, err := c.questRepo.GetQuest(ctx, answer.QuestID)
	if err != nil {
		return fmt.Errorf("questRepo.GetQuest: %w", err)
	}

	err = questObj.ValidateAnswer(answer.Words)
	if err != nil {
		return fmt.Errorf("answer validation: %w", err)
	}

	return nil
}

var (
	ErrQuestIDMismatch = errors.New("this team is not allowed to do the quest")
)
