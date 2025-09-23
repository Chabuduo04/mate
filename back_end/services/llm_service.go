package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 请求体
type ChatRequest struct {
	Stream   bool          `json:"stream"`
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse 响应体
type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

type ChatChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
}

// LLMService 接口
type LLMService interface {
	Chat(messages []ChatMessage) (string, error)
}

// QiniuLLMService 实现
type QiniuLLMService struct {
	APIKey string
	URL    string
	Model  string
}

func NewQiniuLLMService() *QiniuLLMService {
	return &QiniuLLMService{
		APIKey: os.Getenv("QINIU_API_KEY"),
		URL:    os.Getenv("QINIU_API_URL"),
		Model:  os.Getenv("QINIU_API_MODEL"),
	}
}

func (s *QiniuLLMService) Chat(messages []ChatMessage) (string, error) {
	reqBody := ChatRequest{
		Stream:   false,
		Model:    s.Model,
		Messages: messages,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v, raw: %s", err, string(respBody))
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from LLM: %s", string(respBody))
}
