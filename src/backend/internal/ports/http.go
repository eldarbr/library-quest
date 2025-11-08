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

		return version1.NewPostAPIV1ValidateInternalServerError().WithPayload(
			&models.InternalServerError{Error: "internal server error"})
	}

	return version1.NewGetAPIV1QuestOK().WithPayload(questToQuestResponse(ques))
}

// ValidateAnswer handles the POST /api/v1/validate endpoint.
func (h *HTTPHandler) ValidateAnswer(params version1.PostAPIV1ValidateParams) middleware.Responder {
	slog.Debug(
		"ValidateAnswer received", slog.Int64("teamID", *params.Body.TeamID), slog.Int64("questID", *params.Body.QuestID))

	answer := command.Answer{
		QuestID: int(*params.Body.QuestID),
		TeamID:  int(*params.Body.TeamID),
		Words:   params.Body.Answer,
	}

	err := h.app.Commands.ValidateAnswer(params.HTTPRequest.Context(), answer)

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

	return version1.NewPostAPIV1ValidateOK()
}

// questToQuestResponse converts a domain Quest object to a generated QuestResponse model.
func questToQuestResponse(ques quest.Quest) *models.QuestResponse {
	respWords := make([]*models.QuestWord, len(ques.GetWords()))
	for i, w := range ques.GetWords() {
		respWords[i] = wordToQuestWord(w)
	}

	questID := int64(ques.GetID())

	return &models.QuestResponse{
		QuestID: &questID,
		Words:   respWords,
	}
}

// wordToQuestWord converts a domain Word object to a generated QuestWord model.
func wordToQuestWord(wor word.Word) *models.QuestWord {
	return &models.QuestWord{
		QuestWordIdx: int64(wor.QuestWordIDx),
		Position: &models.WordPosition{
			Book: wor.Position.Book,
			Page: wor.Position.Page,
			Line: wor.Position.Line,
			Word: wor.Position.Word,
		},
	}
}
