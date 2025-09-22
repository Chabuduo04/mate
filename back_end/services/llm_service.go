package services

import (
	"context"
	"fmt"
)

// LLMService interface (so we can swap real impl)
type LLMService interface {
	Chat(ctx context.Context, rolePrompt string, history string, userMessage string) (string, error)
}

// MockLLMService - simple mock: echo + role prompt + canned behavior.
// Replace this implementation with real API call (OpenAI/Anthropic/Vertex).
type MockLLMService struct {
}

func NewMockLLMService() LLMService {
	return &MockLLMService{}
}

func (m *MockLLMService) Chat(ctx context.Context, rolePrompt string, history string, userMessage string) (string, error) {
	// Very simple mock: combine prompt, history, and user message and return a role-style reply.
	reply := fmt.Sprintf("[角色回答 — 基于 prompt: %s]\n我听到你说：%s\n（这是一个 mock 回复，可替换为真实 LLM 输出）", rolePrompt, userMessage)
	return reply, nil
}
