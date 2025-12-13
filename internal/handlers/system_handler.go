package handlers

import (
	"zenthara/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type SystemHandler struct {
	logger zerolog.Logger
}

func NewSystemHandler(logger zerolog.Logger) *SystemHandler {
	return &SystemHandler{
		logger: logger.With().Str("handler", "system").Logger(),
	}
}

func (h *SystemHandler) HealthCheck(c *gin.Context) {
	systemInfo := map[string]string{
		"service": "zenthara",
		"status":  "healthy",
	}
	response.Success(c, systemInfo, "System is up and running")
}
