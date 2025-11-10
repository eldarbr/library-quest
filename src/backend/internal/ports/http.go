package ports

import (
	"errors"
	"log/slog"

	"github.com/eldarbr/library-quest/backend/generated/models"
	"github.com/eldarbr/library-quest/backend/generated/restapi/operations/version1"
	"github.com/eldarbr/library-quest/backend/internal/app"
	"github.com/eldarbr/library-quest/backend/internal/app/command"
	"github.com/eldarbr/library-quest/backend/internal/domain/quest"
	"github.com/eldarbr/library-quest/backend/internal/domain/word"
	"github.com/go-openapi/runtime/middleware"
)

// HTTPHandler holds a reference to the application logic.
type HTTPHandler struct {
	app app.Application
}

// NewHTTPHandler creates a new HTTP handler with the given application instance.
func NewHTTPHandler(application app.Application) *HTTPHandler {
	return &HTTPHandler{app: application}
}

// GetQuest handles the GET /api/v1/quest endpoint.
func (h *HTTPHandler) GetQuest(params version1.GetAPIV1QuestParams) middleware.Responder {
	slog.Debug("GetQuest received", slog.Int64("teamID", params.TeamID))

	ques, err := h.app.Queries.GetQuest(params.HTTPRequest.Context(), int(params.TeamID))
	if err != nil {
		slog.Error("getting quest for team", slog.Int64("teamID", params.TeamID), slog.Any("err", err))

		return version1.NewGetAPIV1QuestInternalServerError().WithPayload(
			&models.InternalServerError{Error: "internal server error"})
	}

	return version1.NewGetAPIV1QuestOK().WithPayload(questToQuestResponse(ques))
}

func (h *HTTPHandler) GetRandomQuest(params version1.GetAPIV1QuestRandomParams) middleware.Responder {
	slog.Debug("GetRandomQuest received")

	ques, err := h.app.Queries.GetRandomQuest(params.HTTPRequest.Context())
	if err != nil {
		slog.Error("getting random quest", slog.Any("err", err))

		return version1.NewGetAPIV1QuestRandomInternalServerError().WithPayload(
			&models.InternalServerError{Error: "internal server error"})
	}

	return version1.NewGetAPIV1QuestRandomOK().WithPayload(questToQuestResponse(ques))
}

func (h *HTTPHandler) ValidateAnswerRandom(params version1.PostAPIV1ValidateRandomParams) middleware.Responder {
	slog.Debug(
		"ValidateAnswerRandom received", slog.Int64("questID", *params.Body.QuestID))

	answer := command.Answer{
		QuestID: int(*params.Body.QuestID),
		Words:   params.Body.Answer,
	}

	keyword, err := h.app.Commands.ValidateAnswerRandom(params.HTTPRequest.Context(), answer)

	var wrongAnswerErr quest.WrongAnswerError
	if errors.As(err, &wrongAnswerErr) {
		return version1.NewPostAPIV1ValidateRandomBadRequest().WithPayload(&models.ValidationError{
			Mistakes: wrongAnswerErr.Mistakes,
		})
	}

	if err != nil {
		slog.Error("validating answer", slog.Any("err", err))

		return version1.NewPostAPIV1ValidateRandomInternalServerError().WithPayload(
			&models.InternalServerError{Error: "internal server error"})
	}

	return version1.NewPostAPIV1ValidateRandomOK().WithPayload(&models.CorrectAnswerResponse{Keyword: keyword})
}

// ValidateAnswer handles the POST /api/v1/validate endpoint.
func (h *HTTPHandler) ValidateAnswer(params version1.PostAPIV1ValidateParams) middleware.Responder {
	slog.Debug(
		"ValidateAnswer received", slog.Int64("teamID", *params.Body.TeamID),
		slog.Int64("questID", *params.Body.AnswerValidation.QuestID))

	answer := command.Answer{
		QuestID: int(*params.Body.AnswerValidation.QuestID),
		Words:   params.Body.AnswerValidation.Answer,
	}

	keyword, err := h.app.Commands.ValidateAnswer(params.HTTPRequest.Context(), int(*params.Body.TeamID), answer)

	// Handle specific domain errors and map them to appropriate HTTP responses.
	if errors.Is(err, command.ErrQuestIDMismatch) {
		return version1.NewPostAPIV1ValidateUnauthorized().WithPayload(&models.AuthorizationError{
			Error: err.Error(),
		})
	}

	var wrongAnswerErr quest.WrongAnswerError
	if errors.As(err, &wrongAnswerErr) {
		return version1.NewPostAPIV1ValidateBadRequest().WithPayload(&models.ValidationError{
			Mistakes: wrongAnswerErr.Mistakes,
		})
	}

	if err != nil {
		slog.Error("validating answer", slog.Any("err", err))

		return version1.NewPostAPIV1ValidateInternalServerError().WithPayload(
			&models.InternalServerError{Error: "internal server error"})
	}

	return version1.NewPostAPIV1ValidateOK().WithPayload(&models.CorrectAnswerResponse{Keyword: keyword})
}

// questToQuestResponse converts a domain Quest object to a generated QuestResponse model.
func questToQuestResponse(ques quest.Quest) *models.QuestResponse {
	respWords := make([]*models.QuestWord, len(ques.GetWords()))
	for i, w := range ques.GetWords() {
		respWords[i] = wordToQuestWord(w)
	}

	return &models.QuestResponse{
		QuestID: int64(ques.GetID()),
		Words:   respWords,
	}
}

// wordToQuestWord converts a domain Word object to a generated QuestWord model.
func wordToQuestWord(wor word.Word) *models.QuestWord {
	return &models.QuestWord{
		QuestWordIdx: wor.QuestWordIDx,
		Position: &models.WordPosition{
			Book: wor.Position.Book,
			Page: wor.Position.Page,
			Line: wor.Position.Line,
			Word: wor.Position.Word,
		},
	}
}
