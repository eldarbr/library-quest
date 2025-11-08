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

func NewApplication(questRepo quest.Repository) Application {
	return Application{
		Queries:  query.NewQuery(questRepo),
		Commands: command.NewCommands(questRepo),
	}
}
