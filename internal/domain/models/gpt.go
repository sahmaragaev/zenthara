package models

type GPTRole string

const (
	RoleSystem    GPTRole = "system"
	RoleUser      GPTRole = "user"
	RoleAssistant GPTRole = "assistant"
)

type GPTMessage struct {
	Role    GPTRole `json:"role"`
	Content string  `json:"content"`
}

type GPTRequest struct {
	Model       string       `json:"model"`
	Messages    []GPTMessage `json:"messages"`
	Temperature float64      `json:"temperature",omitempty`
	MaxTokens   int          `json:"max_tokens,omitempty"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}
