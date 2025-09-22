package services

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// ASRService interface for audio->text
type ASRService interface {
	Transcribe(ctx context.Context, audio io.Reader) (string, error)
}

// MockASRService reads incoming bytes and returns a canned text or file-size.
type MockASRService struct{}

func NewMockASRService() ASRService {
	return &MockASRService{}
}

func (m *MockASRService) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	// a naive mock: count bytes and return a placeholder text
	buf := make([]byte, 1024)
	n, err := audio.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	if n == 0 {
		return "", errors.New("empty audio")
	}
	return fmt.Sprintf("（mock 转写）检测到 %d 字节音频，返回示例文本：你好，我想和角色聊天。", n), nil
}
