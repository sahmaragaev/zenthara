package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	v1 "zenthara/api/v1"
	"zenthara/internal/config"
	"zenthara/internal/handlers"
	"zenthara/internal/middleware"
	"zenthara/internal/services/gpt"
	"zenthara/internal/services/prompts"

	speakinghandlers "zenthara/internal/handlers/speaking"
	writinghandlers "zenthara/internal/handlers/writing"
	speakingservices "zenthara/internal/services/speaking"
	writingservices "zenthara/internal/services/writing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type App struct {
	cfg    *config.Config
	logger zerolog.Logger
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.Config, logger zerolog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Setup() error {
	gin.SetMode(a.cfg.App.GinMode)
	a.router = gin.New()
	a.router.Use(gin.Recovery())

	authMiddleware := middleware.NewAuthMiddleware(a.cfg.Auth.APIKey, a.logger)

	systemHandler, speakingQuestionHandler, speakingEvaluationHandler, writingQuestionHandler, writingEvaluationHandler, err := a.initHandlers()
	if err != nil {
		return fmt.Errorf("failed to initialize handlers: %w", err)
	}

	routes := v1.NewRoutes(a.router, a.logger, systemHandler, speakingQuestionHandler, writingQuestionHandler, speakingEvaluationHandler, writingEvaluationHandler)
	routes.Setup(authMiddleware)

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler: a.router,
	}

	return nil
}

func (a *App) initHandlers() (*handlers.SystemHandler, *speakinghandlers.QuestionHandler, *speakinghandlers.EvaluationHandler, *writinghandlers.QuestionHandler, *writinghandlers.EvaluationHandler, error) {
	systemHandler := handlers.NewSystemHandler(a.logger)

	gptClient := gpt.NewClient(&a.cfg.AI, a.logger)
	promptStore := prompts.NewPromptStore("config/prompts")
	if err := promptStore.LoadPrompts(); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to load prompts: %w", err)
	}
	speakingPrompts := prompts.NewSpeakingPrompts(promptStore)
	writingPrompts := prompts.NewWritingPrompts(promptStore)

	speakingQuestionService := speakingservices.NewQuestionService(gptClient, speakingPrompts, a.logger)
	speakingEvaluationService := speakingservices.NewEvaluationService(gptClient, speakingPrompts, a.logger)
	speakingQuestionHandler := speakinghandlers.NewQuestionHandler(speakingQuestionService, a.logger)
	speakingEvaluationHandler := speakinghandlers.NewEvaluationHandler(speakingEvaluationService, a.logger)

	writingQuestionService := writingservices.NewQuestionService(gptClient, writingPrompts, a.logger)
	writingQuestionHandler := writinghandlers.NewQuestionHandler(writingQuestionService, a.logger)
	writingEvaluationService := writingservices.NewEvaluationService(gptClient, writingPrompts, a.logger)
	writingEvaluationHandler := writinghandlers.NewEvaluationHandler(writingEvaluationService, a.logger)

	return systemHandler, speakingQuestionHandler, speakingEvaluationHandler, writingQuestionHandler, writingEvaluationHandler, nil
}

func (a *App) Run() error {
	go func() {
		a.logger.Info().Msgf("Starting server on port %d", a.cfg.App.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}
