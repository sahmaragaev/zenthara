package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"zenthara/internal/config"
	"zenthara/internal/domain/models"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type Client interface {
	Complete(ctx context.Context, prompt string) (json.RawMessage, error)
}

type client struct {
	config *config.AIConfig
	http   *resty.Client
	logger zerolog.Logger
}

func NewClient(config *config.AIConfig, logger zerolog.Logger) Client {
	httpClient := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.APIKey)).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			},
		)

	httpClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		logger.Debug().
			Str("url", r.URL).
			Str("method", r.Method).
			Msg("Starting GPT request")
		return nil
	})

	httpClient.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		logger.Debug().
			Str("url", r.Request.URL).
			Str("method", r.Request.Method).
			Int("status", r.StatusCode()).
			Dur("duration", r.Time()).
			Msg("Completed GPT request")
		return nil
	})

	return &client{
		config: config,
		http:   httpClient,
		logger: logger.With().Str("component", "gpt_client").Logger(),
	}
}

func (c *client) Complete(ctx context.Context, prompt string) (json.RawMessage, error) {
	logger := c.logger.With().
		Str("method", "Complete").
		Int("promptLength", len(prompt)).
		Logger()

	req := models.GPTRequest{
		Model: c.config.Model,
		Messages: []models.GPTMessage{
			{
				Role:    models.RoleSystem,
				Content: "You are an IELTS examiner. Provide responses in JSON format only.",
			},
			{
				Role:    models.RoleUser,
				Content: prompt,
			},
		},
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}

	var resp models.GPTResponse
	httpResp, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&resp).
		Post("/chat/completions")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to make GPT request")
		return nil, fmt.Errorf("GPT request failed: %w", err)
	}

	if httpResp.IsError() {
		logger.Error().
			Int("statusCode", httpResp.StatusCode()).
			RawJSON("response", httpResp.Body()).
			Msg("AI API returned error")
		return nil, fmt.Errorf("API error: %s", httpResp.String())
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := resp.Choices[0].Message.Content
	var jsonContent json.RawMessage
	if err := json.Unmarshal([]byte(content), &jsonContent); err != nil {
		logger.Error().
			Err(err).
			Str("content", content).
			Msg("Failed to parse GPT response as JSON")
		return nil, fmt.Errorf("GPT response is not valid JSON: %w", err)
	}

	return jsonContent, nil
}
