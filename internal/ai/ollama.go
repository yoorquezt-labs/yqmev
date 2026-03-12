package ai

// NewOllama creates an Ollama provider using the OpenAI-compatible API.
// Default model: llama3.1, default base URL: http://localhost:11434/v1.
func NewOllama(model, baseURL string) *OpenAI {
	if model == "" {
		model = "llama3.1"
	}
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}
	return &OpenAI{
		apiKey:  "",
		model:   model,
		baseURL: baseURL,
	}
}
