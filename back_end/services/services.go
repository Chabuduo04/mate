package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Services struct {
	RoleService  *RoleService
	SessionStore SessionStore
	UserService  *UserService
	LLM          LLMService
	ASR          ASRService
	TTS          TTSService
	Storage      StorageService
	Logger       *zap.Logger

	// config hints (optional)
	LLMModel    string
	ASREndpoint string
	TTSEndpoint string
}

func NewServices(
	roleService *RoleService,
	sessionStore SessionStore,
	userService *UserService,
	llm LLMService,
	asr ASRService,
	tts TTSService,
	storage StorageService,
	logger *zap.Logger,
	llmModel string,
	asrEndpoint string,
	ttsEndpoint string,
) Services {
	return Services{
		RoleService:  roleService,
		SessionStore: sessionStore,
		UserService:  userService,
		LLM:          llm,
		ASR:          asr,
		TTS:          tts,
		Storage:      storage,
		Logger:       logger,
		LLMModel:     llmModel,
		ASREndpoint:  asrEndpoint,
		TTSEndpoint:  ttsEndpoint,
	}
}

func (s *Services) AudioToText(file io.Reader, header *multipart.FileHeader) (string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read audio file: %w", err)
	}

	fileExt := filepath.Ext(header.Filename)
	if fileExt == "" {
		fileExt = ".wav" // 默认扩展名
	}
	uniqueKey := fmt.Sprintf("voice-chat/%s%s", uuid.New().String(), fileExt)

	uploadReader := bytes.NewReader(fileBytes)
	err = s.Storage.Upload(uploadReader, uniqueKey, header.Filename)
	if err != nil {
		s.Logger.Sugar().Errorf("upload error: %v", err)
		return "", fmt.Errorf("upload error: %w", err)
	}

	audioURL := fmt.Sprintf("%s%s/%s", "http://", config.GetConfig().Kodo.Host, uniqueKey)

	return s.ASR.TranscribeFromURL(audioURL)
}
