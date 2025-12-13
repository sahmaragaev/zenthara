package writing

import (
	"fmt"
	models "zenthara/internal/domain/dto/writing"
	services "zenthara/internal/services/writing"
	"zenthara/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type QuestionHandler struct {
	questionService services.QuestionService
	logger          zerolog.Logger
}

func NewQuestionHandler(qs services.QuestionService, logger zerolog.Logger) *QuestionHandler {
	return &QuestionHandler{
		questionService: qs,
		logger:          logger.With().Str("handler", "question").Logger(),
	}
}

func (h *QuestionHandler) GenerateQuestions(c *gin.Context) {
	var req models.QuestionRequest
	fmt.Println("req", req)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "Invalid request payload")
		return
	}

	questions, err := h.questionService.GenerateQuestions(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate question")
		response.InternalServerError(c, err, "Failed to generate question")
		return
	}

	response.Success(c, questions, "Writing question generated successfully")
}
