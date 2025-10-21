package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Chabuduo04/mate/back_end/config"
)

// ASRRequest 定义ASR请求结构体
type ASRRequest struct {
	Model string   `json:"model"`
	Audio AudioAsr `json:"audio"`
}

type AudioAsr struct {
	Format string `json:"format"`
	URL    string `json:"url"`
}

// ASRResponse 定义ASR响应结构体
type ASRResponse struct {
	ReqID     string  `json:"reqid"`
	Operation string  `json:"operation"`
	Data      ASRData `json:"data"`
}

type ASRData struct {
	AudioInfo AudioInfo `json:"audio_info"`
	Result    Result    `json:"result"`
}

type AudioInfo struct {
	Duration int `json:"duration"`
}

type Result struct {
	Additions Additions `json:"additions"`
	Text      string    `json:"text"`
}

type Additions struct {
	Duration string `json:"duration"`
}

// ASRService interface for audio->text
type ASRService interface {
	Transcribe(ctx context.Context, audio io.Reader) (string, error)
	TranscribeFromURL(audioURL string) (string, error)
	// ListVoicesRaw is not part of ASR, but user said same base URL; ASR unchanged
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

func (m *MockASRService) TranscribeFromURL(audioURL string) (string, error) {
	// Mock implementation for URL-based transcription
	return fmt.Sprintf("（mock 转写）从URL %s 转写音频，返回示例文本：你好，我想和角色聊天。", audioURL), nil
}

type QiniuASRService struct {
	APIKey string
	URL    string
}

func NewQiniuASRService() *QiniuASRService {
	cfg := config.GetConfig()
	return &QiniuASRService{
		APIKey: cfg.ASR.ApiKey,
		URL:    cfg.ASR.Endpoint,
	}
}
func (s *QiniuASRService) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	// For file-based transcription, we need to upload first
	// This is a simplified implementation - in practice you might want to handle this differently
	return "（文件转写）需要先上传文件到云存储", nil
}

func (s *QiniuASRService) TranscribeFromURL(audioURL string) (string, error) {
	// 构建请求体
	requestBody := ASRRequest{
		Model: "asr",
		Audio: AudioAsr{
			Format: "mp3",
			URL:    audioURL,
		},
	}

	// 序列化请求体
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API返回错误: %s, 响应体: %s", resp.Status, string(body))
	}

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 解析JSON响应
	var asrResponse ASRResponse
	err = json.Unmarshal(responseBody, &asrResponse)
	if err != nil {
		return "", fmt.Errorf("解析JSON响应失败: %v", err)
	}

	return asrResponse.Data.Result.Text, nil
}
