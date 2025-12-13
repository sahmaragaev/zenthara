package writing

import (
	dto "zenthara/internal/domain/dto/writing"
	"zenthara/internal/services/writing"
	"zenthara/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type EvaluationHandler struct {
	evaluationService writing.EvaluationService
	logger            zerolog.Logger
	response          response.Response
}

func NewEvaluationHandler(es writing.EvaluationService, logger zerolog.Logger) *EvaluationHandler {
	return &EvaluationHandler{
		evaluationService: es,
		logger:            logger.With().Str("handler", "evaluation").Logger(),
		response:          response.Response{},
	}
}

func (h *EvaluationHandler) EvaluateAnswer(c *gin.Context) {
	var req dto.EvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to bind request")
		response.BadRequest(c, err, nil)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request")
		response.BadRequest(c, err, nil)
		return
	}

	h.logger.Debug().
		Str("testType", string(req.TestType)).
		RawJSON("data", req.Data).
		Msg("Processing evaluation request")

	result, err := h.evaluationService.EvaluateAnswer(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to evaluate answer")
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, result, "Answer evaluated successfully")
}
