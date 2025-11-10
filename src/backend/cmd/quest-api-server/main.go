package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"

	"github.com/eldarbr/library-quest/backend/generated/restapi"
	"github.com/eldarbr/library-quest/backend/generated/restapi/operations"
	"github.com/eldarbr/library-quest/backend/generated/restapi/operations/version1"
	"github.com/eldarbr/library-quest/backend/internal/app"
	"github.com/eldarbr/library-quest/backend/internal/ports"
	"github.com/eldarbr/library-quest/backend/internal/repository/filedb"
)

type options struct {
	DBFilePath         string   `long:"db-file-path" short:"f" description:"Path to the tab-separated database file" required:"true"`
	CORSAllowedOrigins []string `long:"cors-allowed-origins" description:"A list of allowed origins for CORS" required:"false"`
}

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

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Quest API"
	parser.LongDescription = "API for retrieving and validating team quests."

	_, err = parser.AddGroup("Application Options", "Options for the application", &opts)
	if err != nil {
		slog.Error("parser add group", slog.Any("err", err))

		return
	}

	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err = parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			slog.Error("parser add group", slog.Any("err", err))

			return
		}
	}

	if _, err = parser.Parse(); err != nil {
		code := 1

		flagErr := &flags.Error{}
		if errors.As(err, &flagErr) {
			if flagErr.Type == flags.ErrHelp {
				code = 0
			}
		}

		os.Exit(code)
	}

	questRepo, err := filedb.NewFileToMemoDB(opts.DBFilePath)
	if err != nil {
		slog.Error("failed to initialize quest repository", slog.Any("err", err))

		return
	}

	application := app.NewApplication(questRepo)
	httpHandler := ports.NewHTTPHandler(application)

	api.Version1GetAPIV1QuestHandler = version1.GetAPIV1QuestHandlerFunc(httpHandler.GetQuest)
	api.Version1PostAPIV1ValidateHandler = version1.PostAPIV1ValidateHandlerFunc(httpHandler.ValidateAnswer)

	server.ConfigureAPI()

	var allowedOrigins []string
	if len(opts.CORSAllowedOrigins) > 0 {
		allowedOrigins = opts.CORSAllowedOrigins
	} else {
		originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
		if originsStr != "" {
			allowedOrigins = strings.Split(originsStr, ",")
		}
	}

	if len(allowedOrigins) == 0 {
		slog.Error("CORS allowed origins must be provided via --cors-allowed-origins flag or CORS_ALLOWED_ORIGINS env var")

		return
	}

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(api.Serve(nil))
	server.SetHandler(handler)

	if err := server.Serve(); err != nil {
		slog.Error("server serve", slog.Any("err", err))

		return
	}
}
