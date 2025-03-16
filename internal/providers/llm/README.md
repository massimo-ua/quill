# LLM Providers for Quill

This package implements Large Language Model (LLM) providers that can be used to analyze messages, generate documentation, categorize content, and detect references in text.

## Supported Providers

### OpenAI

The OpenAI provider uses the OpenAI API to access models like GPT-4, GPT-3.5-turbo, etc.

Features:
- Full implementation of the AiAgentProvider interface
- Support for message analysis with confidence scoring
- Documentation generation with metadata integration
- Content categorization
- Reference detection in text
- Error handling and fallbacks for unstructured responses

### Ollama

The Ollama provider integrates with local Ollama instances to use open-source models like Llama2, Mistral, Gemma, etc.

Features:
- Same interface as the OpenAI provider
- Local deployment support
- Configurable model parameters
- Optimized for lower latency
- No API key requirements

## Usage

### Creating a Provider

```go
import (
    "github.com/massimo-ua/quill/internal/providers/llm"
    "github.com/massimo-ua/quill/internal/providers/llm/openai"
    "github.com/massimo-ua/quill/internal/providers/llm/ollama"
)

// For OpenAI
openAIConfig := &openai.Config{
    APIKey:      "your-openai-api-key",
    Model:       "gpt-4",
    Temperature: 0.7,
    MaxTokens:   1024,
}

config := &llm.Config{
    Type:   llm.ProviderTypeOpenAI,
    OpenAI: openAIConfig,
}

// OR for Ollama
ollamaConfig := &ollama.Config{
    ServerURL:   "http://localhost:11434",
    Model:       "llama2",
    Temperature: 0.7,
    MaxTokens:   1024,
}

config := &llm.Config{
    Type:   llm.ProviderTypeOllama,
    Ollama: ollamaConfig,
}

// Create provider
provider, err := llm.NewLLMProvider(config)
if err != nil {
    // Handle error
}

// Inject into services
botService := services.NewBotService(chatProvider, docStore, provider, projectService, docService)
```

### Analyzing Messages

```go
result, err := provider.AnalyzeMessage(ctx, "Let's update the deployment strategy to use blue-green deployments.")
if err != nil {
    // Handle error
}

fmt.Printf("Type: %s\n", result.MessageType())
fmt.Printf("Category: %s\n", result.Category())
fmt.Printf("Confidence: %.2f\n", result.ConfidenceScore())
fmt.Printf("Tags: %v\n", result.SuggestedTags())
```

### Generating Documentation

```go
metadata := map[string]interface{}{
    "type":     "decision",
    "category": "operations",
}

doc, err := provider.GenerateDocumentation(ctx, "We decided to use Kubernetes for our deployment infrastructure.", metadata)
if err != nil {
    // Handle error
}

fmt.Println(doc)
```

### Detecting References

```go
references, err := provider.DetectReferences(ctx, "As discussed in message thread-123, we need to review doc://architecture/overview.md.")
if err != nil {
    // Handle error
}

for _, ref := range references {
    fmt.Printf("Type: %s, Value: %s\n", ref.Type(), ref.Value())
}
```

## Configuration

### OpenAI Configuration

| Parameter    | Description                                       | Default   |
|--------------|---------------------------------------------------|-----------|
| APIKey       | OpenAI API key                                    | Required  |
| Model        | Model name (e.g., "gpt-4", "gpt-3.5-turbo")      | Required  |
| Temperature  | Controls randomness (0-2)                         | 0.7       |
| MaxTokens    | Maximum tokens to generate                        | 1024      |
| BaseURL      | Custom API endpoint                               | OpenAI API|
| Organization | OpenAI organization ID                            | None      |

### Ollama Configuration

| Parameter    | Description                                       | Default   |
|--------------|---------------------------------------------------|-----------|
| ServerURL    | Ollama server URL (e.g., "http://localhost:11434")| Required  |
| Model        | Model name (e.g., "llama2", "mistral")           | Required  |
| Temperature  | Controls randomness (0-2)                         | 0.7       |
| MaxTokens    | Maximum tokens to generate                        | 1024      |
| SystemPrompt | Default system prompt                             | None      |