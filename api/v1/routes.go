package v1

import (
	"zenthara/internal/handlers"
	speakinghandlers "zenthara/internal/handlers/speaking"
	writinghandlers "zenthara/internal/handlers/writing"
	"zenthara/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Routes struct {
	router                    *gin.Engine
	logger                    zerolog.Logger
	speakingQuestionHandler   *speakinghandlers.QuestionHandler
	writingQuestionHandler    *writinghandlers.QuestionHandler
	speakingEvaluationHandler *speakinghandlers.EvaluationHandler
	writingEvaluationHandler  *writinghandlers.EvaluationHandler
	systemHandler             *handlers.SystemHandler
}

func NewRoutes(
	router *gin.Engine,
	logger zerolog.Logger,
	systemHandler *handlers.SystemHandler,
	speakingQuestionHandler *speakinghandlers.QuestionHandler,
	writingQuestionHandler *writinghandlers.QuestionHandler,
	speakingEvaluationHandler *speakinghandlers.EvaluationHandler,
	writingEvaluationHandler *writinghandlers.EvaluationHandler,
) *Routes {
	return &Routes{
		router:                    router,
		logger:                    logger,
		systemHandler:             systemHandler,
		speakingQuestionHandler:   speakingQuestionHandler,
		writingQuestionHandler:    writingQuestionHandler,
		speakingEvaluationHandler: speakingEvaluationHandler,
		writingEvaluationHandler:  writingEvaluationHandler,
	}
}

func (r *Routes) Setup(authMiddleware *middleware.AuthMiddleware) {
	r.router.GET("/health", r.systemHandler.HealthCheck)

	v1 := r.router.Group("/v1")
	v1.Use(authMiddleware.Authenticate())
	{
		speaking := v1.Group("/speaking")
		{
			speaking.POST("/generate-questions", r.speakingQuestionHandler.GenerateQuestions)
			speaking.POST("/evaluate", r.speakingEvaluationHandler.EvaluateAnswer)
		}

		writing := v1.Group("/writing")
		{
			writing.POST("/generate-question", r.writingQuestionHandler.GenerateQuestions)
			writing.POST("/evaluate", r.writingEvaluationHandler.EvaluateAnswer)
		}
	}
}
