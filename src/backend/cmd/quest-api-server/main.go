package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"

	"github.com/eldarbr/library-quest/backend/generated/restapi"
	"github.com/eldarbr/library-quest/backend/generated/restapi/operations"
	"github.com/eldarbr/library-quest/backend/generated/restapi/operations/version1"
	"github.com/eldarbr/library-quest/backend/internal/app"
	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/ports"
	"github.com/eldarbr/library-quest/backend/internal/repository/filedb"
)

type options struct {
	DBFilePath         string   `long:"db-file-path" short:"f" description:"Path to the tab-separated database file" required:"true"`
	CORSAllowedOrigins []string `long:"cors-allowed-origins" description:"A list of allowed origins for CORS" required:"false"`
}

const (
	envKeyCORS         = `CORS_ALLOWED_ORIGINS`
	envKeyQuestKeyword = `QUEST_KEYWORD`
	envKeyQuestsSeed   = `QUEST_SHUFFLE_SEED`
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	var opts options

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		slog.Error("loads embedded swagger", slog.Any("err", err))

		return
	}

	api := operations.NewQuestAPIAPI(swaggerSpec)
	api.Logger = func(format string, args ...any) {
		logger.Info(fmt.Sprintf(format, args...))
	}

	server := restapi.NewServer(api)
	defer server.Shutdown()

	parseFlags(server, api, &opts)

	questRepo, err := filedb.NewFileToMemoDB(opts.DBFilePath)
	if err != nil {
		slog.Error("failed to initialize quest repository", slog.Any("err", err))

		return
	}

	questKeyword, questKeywordPresent := os.LookupEnv(envKeyQuestKeyword)
	if !questKeywordPresent {
		slog.Error("no quest keyword provided")

		return
	}

	envSeed, envSeedPresent := os.LookupEnv(envKeyQuestsSeed)
	if !envSeedPresent {
		slog.Error("no seed to shuffle quests")
		return
	}
	seed, err := strconv.Atoi(envSeed)
	if err != nil {
		slog.Error("parse shuffle seed", slog.Any("err", err))
		return
	}

	randshuffler := quest.NewDeterminedTeamToQuestShuffler(questRepo.GetTotalQuestsCnt(), int64(seed))
	application := app.NewApplication(quest.ConstantKeyworder{Keyword: questKeyword},
		questRepo, randshuffler, randshuffler)
	httpHandler := ports.NewHTTPHandler(application)

	api.Version1GetAPIV1QuestHandler = version1.GetAPIV1QuestHandlerFunc(httpHandler.GetQuest)
	api.Version1GetAPIV1QuestRandomHandler = version1.GetAPIV1QuestRandomHandlerFunc(httpHandler.GetRandomQuest)
	api.Version1PostAPIV1ValidateHandler = version1.PostAPIV1ValidateHandlerFunc(httpHandler.ValidateAnswer)
	api.Version1PostAPIV1ValidateRandomHandler = version1.PostAPIV1ValidateRandomHandlerFunc(httpHandler.ValidateAnswerRandom)

	server.ConfigureAPI()

	allowedOrigins := opts.CORSAllowedOrigins
	if envCors, present := os.LookupEnv(envKeyCORS); len(opts.CORSAllowedOrigins) <= 0 && present {
		allowedOrigins = strings.Split(envCors, ",")
	}

	if len(allowedOrigins) == 0 {
		slog.Error("CORS allowed origins must be provided via --cors-allowed-origins flag or CORS_ALLOWED_ORIGINS env var")

		return
	}

	server.SetHandler(applyCors(allowedOrigins, api.Serve(nil)))

	if err := server.Serve(); err != nil {
		slog.Error("server serve", slog.Any("err", err))

		return
	}
}

func applyCors(allowedOrigins []string, handler http.Handler) http.Handler {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler = corsOptions.Handler(handler)

	return handler
}

func parseFlags(server *restapi.Server, api *operations.QuestAPIAPI, opts *options) {
	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Quest API"
	parser.LongDescription = "API for retrieving and validating team quests."

	_, err := parser.AddGroup("Application Options", "Options for the application", opts)
	if err != nil {
		slog.Error("parser add group", slog.Any("err", err))
		os.Exit(1)
	}

	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err = parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			slog.Error("parser add group", slog.Any("err", err))
			os.Exit(1)
		}
	}

	if _, err = parser.Parse(); err != nil {
		flagErr := &flags.Error{}
		if errors.As(err, &flagErr) && flagErr.Type != flags.ErrHelp {
			slog.Error("parse flags", slog.Any("err", flagErr))
			os.Exit(1)
		}
	}
}
