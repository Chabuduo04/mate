package services

import "go.uber.org/zap"

type Services struct {
	RoleService  *RoleService
	SessionStore SessionStore
	LLM          LLMService
	ASR          ASRService
	TTS          TTSService
	Storage		 StorageService
	Logger       *zap.Logger

	// config hints (optional)
	LLMModel    string
	ASREndpoint string
	TTSEndpoint string
}
