package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	apiKey string
	logger zerolog.Logger
}

func NewAuthMiddleware(apiKey string, logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		apiKey: apiKey,
		logger: logger.With().Str("middleware", "auth").Logger(),
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientAPIKey := c.GetHeader("X-API-Key")

		if clientAPIKey == "" {
			m.logger.Warn().Msg("Missing API key in request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			return
		}

		if subtle.ConstantTimeCompare([]byte(clientAPIKey), []byte(m.apiKey)) != 1 {
			m.logger.Warn().Msg("Invalid API key provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}

		c.Next()
	}
}
