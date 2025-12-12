package prompts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"zenthara/internal/domain/enums/prompttype"

	"gopkg.in/yaml.v3"
)

type PromptStore struct {
	mu       sync.RWMutex
	prompts  map[prompttype.PromptType]map[string]*template.Template
	basePath string
}

func NewPromptStore(promptsPath string) *PromptStore {
	return &PromptStore{
		prompts:  make(map[prompttype.PromptType]map[string]*template.Template),
		basePath: promptsPath,
	}
}

func (s *PromptStore) LoadPrompts() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.loadPromptFile("speaking_questions.yaml", prompttype.SpeakingQuestions); err != nil {
		return fmt.Errorf("failed to load speaking questions prompts: %w", err)
	}

	if err := s.loadPromptFile("speaking_evaluations.yaml", prompttype.SpeakingEvaluation); err != nil {
		return fmt.Errorf("failed to load speaking evaluation prompts: %w", err)
	}

	if err := s.loadPromptFile("writing_questions.yaml", prompttype.WritingQuestions); err != nil {
		return fmt.Errorf("failed to load writing questions prompts: %w", err)
	}

	if err := s.loadPromptFile("writing_evaluations.yaml", prompttype.WritingEvaluation); err != nil {
		return fmt.Errorf("failed to load writing evaluation prompts: %w", err)
	}

	return nil
}

func (s *PromptStore) loadPromptFile(filename string, promptType prompttype.PromptType) error {
	path := filepath.Join(s.basePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var promptsMap map[string]string
	if err := yaml.Unmarshal(data, &promptsMap); err != nil {
		return err
	}

	templates := make(map[string]*template.Template)
	for name, promptText := range promptsMap {
		tmpl, err := template.New(name).Parse(promptText)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		templates[name] = tmpl
	}

	s.prompts[promptType] = templates
	return nil
}

func (s *PromptStore) GetPrompt(promptType prompttype.PromptType, name string, data any) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	templates, ok := s.prompts[promptType]
	if !ok {
		return "", fmt.Errorf("prompt type %s not found", promptType)
	}

	tmpl, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("prompt %s not found for type %s", name, promptType)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
