# Zenthara

An IELTS test API service for generating questions and evaluating answers for speaking and writing tests using AI.

## Features

- **Speaking Test Support**: Generate questions and evaluate speaking test answers
- **Writing Test Support**: Generate questions and evaluate writing test answers
- **AI-Powered**: Uses GPT models for intelligent question generation and evaluation
- **REST API**: Clean RESTful API with authentication middleware
- **Configurable Prompts**: YAML-based prompt templates for customization

## API Endpoints

### Speaking
- `POST /v1/speaking/generate-questions` - Generate speaking test questions
- `POST /v1/speaking/evaluate` - Evaluate speaking test answers

### Writing
- `POST /v1/writing/generate-question` - Generate writing test questions
- `POST /v1/writing/evaluate` - Evaluate writing test answers

### System
- `GET /health` - Health check endpoint

## Configuration

Copy `example.config.yaml` to `config.yaml` and configure:
- Application settings (port, environment, etc.)
- Authentication API key
- AI service configuration (GPT API key, model, etc.)
- Logger settings

## Running

```bash
go run cmd/zenthara/main.go
```

## Project Structure

- `api/v1/` - API routes and handlers
- `internal/domain/` - Domain models, DTOs, and enums
- `internal/services/` - Business logic services
- `internal/handlers/` - HTTP request handlers
- `config/prompts/` - YAML prompt templates
