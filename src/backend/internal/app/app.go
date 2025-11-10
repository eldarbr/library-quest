package app

import (
	"github.com/eldarbr/library-quest/backend/internal/app/command"
	"github.com/eldarbr/library-quest/backend/internal/app/query"
	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
)

type Application struct {
	Queries  query.Query
	Commands command.Commands
}

func NewApplication(validationKeyworder command.KeywordProvider, questRepo quest.Repository,
	shuffler quest.QuestShuffler, randomer quest.RandomQuestShuffler,
) Application {
	return Application{
		Queries:  query.NewQuery(questRepo, shuffler, randomer),
		Commands: command.NewCommands(validationKeyworder, questRepo, shuffler),
	}
}
