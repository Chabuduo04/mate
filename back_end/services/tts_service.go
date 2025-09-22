package services

import (
	"context"
	"encoding/base64"
)

// TTSService interface: text -> audio (return base64 for simplicity)
type TTSService interface {
	Synthesize(ctx context.Context, text string, voice string) (string, error) // returns base64 audio
}

type MockTTSService struct{}

func NewMockTTSService() TTSService {
	return &MockTTSService{}
}

func (m *MockTTSService) Synthesize(ctx context.Context, text string, voice string) (string, error) {
	// mock: return base64 of text bytes (NOT real audio) so front-end can still receive some payload
	b := []byte("AUDIO-MOCK:" + text)
	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}
