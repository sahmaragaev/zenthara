# Zenthara

An IELTS test API service for generating questions and evaluating answers for speaking and writing tests using AI.

## Features

- **Speaking Test Support**: Generate questions and evaluate speaking test answers
- **Writing Test Support**: Generate questions and evaluate writing test answers
- **AI-Powered**: Uses GPT models for intelligent question generation and evaluation
- **REST API**: Clean RESTful API with authentication middleware
- **Configurable Prompts**: YAML-based prompt templates for customization

## Prerequisites

- Go 1.21.5 or higher
- OpenAI API key (or compatible GPT API)
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd zenthara
```

2. Install dependencies:
```bash
go mod download
```

3. Set up configuration:
```bash
cp example.config.yaml config.yaml
```

4. Edit `config.yaml` with your settings:
   - Set your OpenAI API key
   - Configure authentication API key
   - Adjust port and other settings as needed

## Building

### Development Build
```bash
go build -o bin/zenthara cmd/zenthara/main.go
```

### Production Build
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/zenthara cmd/zenthara/main.go
```

The binary will be created in the `bin/` directory.

## Running

### Development Mode
Run directly with Go:
```bash
go run cmd/zenthara/main.go
```

### Production Mode
Run the compiled binary:
```bash
./bin/zenthara
```

The server will start on the port specified in `config.yaml` (default: 8081).

### Environment Variables
You can also use environment variables in your `config.yaml`:
```yaml
auth:
  api_key: ${API_KEY}

gpt:
  api_key: ${OPENAI_API_KEY}
```

Then set them before running:
```bash
export API_KEY=your-api-key
export OPENAI_API_KEY=your-openai-key
go run cmd/zenthara/main.go
```

## API Endpoints

All endpoints except `/health` require authentication via API key header:
```
X-API-Key: <your-api-key>
```

### Speaking
- `POST /v1/speaking/generate-questions` - Generate speaking test questions
- `POST /v1/speaking/evaluate` - Evaluate speaking test answers

### Writing
- `POST /v1/writing/generate-question` - Generate writing test questions
- `POST /v1/writing/evaluate` - Evaluate writing test answers

### System
- `GET /health` - Health check endpoint (no authentication required)

## Configuration

The `config.yaml` file contains all application settings:

- **app**: Application name, environment, port, and Gin mode
- **auth**: API key for endpoint authentication
- **http**: HTTP client timeout and retry settings
- **gpt**: AI service configuration (API key, model, temperature, etc.)
- **logger**: Logging level and format

See `example.config.yaml` for a complete example.

## Development

### Project Structure

```
zenthara/
├── api/v1/              # API routes and handlers
├── cmd/zenthara/         # Application entry point
├── config/               # Configuration files
│   └── prompts/         # YAML prompt templates
├── internal/
│   ├── domain/          # Domain models, DTOs, and enums
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, etc.)
│   ├── services/       # Business logic services
│   ├── shared/         # Shared utilities (logger, response)
│   └── utils/          # Helper utilities
└── config.yaml         # Application configuration
```

### Testing the API

Test the health endpoint:
```bash
curl http://localhost:8081/health
```

Test with authentication:
```bash
curl -X POST http://localhost:8081/v1/speaking/generate-questions \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"test_type": "part1", "include_topics": [], "exclude_topics": []}'
```

## License

[Add your license here]
