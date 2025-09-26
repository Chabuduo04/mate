package services

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "encoding/json"
    "strings"

    "github.com/Chabuduo04/mate/back_end/config"
)

// TTSRequest 定义请求结构体
type TTSRequest struct {
	Audio   Audio   `json:"audio"`
	Request Request `json:"request"`
}

type Audio struct {
	VoiceType string  `json:"voice_type"`
	Encoding  string  `json:"encoding"`
	SpeedRatio float64 `json:"speed_ratio"`
}

type Request struct {
	Text string `json:"text"`
}

// TTSResponse 定义响应结构体
type TTSResponse struct {
	ReqID     string         `json:"reqid"`
	Operation string         `json:"operation"`
	Sequence  int            `json:"sequence"`
	Data      string         `json:"data"`
	Addition  ResponseAddition `json:"addition"`
	AudioData []byte         `json:"-"` // 音频二进制数据
}

type ResponseAddition struct {
	Duration string `json:"duration"`
}

// TTSService interface: text -> audio (return base64 for simplicity)
type TTSService interface {
	Synthesize(text string, voice string) (string, error) // returns base64 audio
    ListVoicesRaw() ([]byte, error)
}

type QiniuTTSService struct {
	APIKey string
	URL    string
}

func NewQiniuTTSService() *QiniuTTSService {
	return &QiniuTTSService{
		APIKey: config.AppConfig.ApiKey,
		URL:	config.AppConfig.TTSEndpoint,
	}
}

// ListVoicesRaw 调用提供方的 /voice/list 接口，返回原始响应体
func (s *QiniuTTSService) ListVoicesRaw() ([]byte, error) {
    base := s.URL
    endpoint := strings.TrimRight(base, "/") + "/list"
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, fmt.Errorf("创建请求失败: %v", err)
    }
    if s.APIKey != "" {
        req.Header.Set("Authorization", "Bearer "+s.APIKey)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("发送请求失败: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API返回错误: %s, 响应体: %s", resp.Status, string(body))
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("读取响应体失败: %v", err)
    }
    return body, nil
}

func (s *QiniuTTSService) Synthesize(text string, voice string) (string, error) {
	// 构建请求体
	// voice 为空时使用默认值
	defaultVoice := "qiniu_zh_female_tmjxxy"
	if voice == "" {
		voice = defaultVoice
	}
	requestBody := TTSRequest{
		Audio: Audio{
			VoiceType: voice,
			Encoding:  "mp3",
			SpeedRatio: 1.0,
		},
		Request: Request{
			Text: text,
		},
	}

	// 序列化请求体
	base := s.URL
    endpoint := strings.TrimRight(base, "/") + "/tts"
	jsonData, err := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
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

	// 读取完整的响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 解析JSON响应
	var ttsResponse TTSResponse
	err = json.Unmarshal(responseBody, &ttsResponse)
	if err != nil {
		return "", fmt.Errorf("解析JSON响应失败: %v", err)
	}

	// 解码base64音频数据
	// if ttsResponse.Data != "" {
	// 	audioData, err := base64.StdEncoding.DecodeString(ttsResponse.Data)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("base64解码失败: %v", err)
	// 	}
	// 	ttsResponse.AudioData = audioData
	// } else {
	// 	return nil, fmt.Errorf("响应中未包含音频数据")
	// }

	return ttsResponse.Data, nil
}

// MockTTSService 用于本地演示/测试
type MockTTSService struct{}

func NewMockTTSService() TTSService {
    return &MockTTSService{}
}

func (m *MockTTSService) Synthesize(text string, voice string) (string, error) {
    if text == "" {
        return "", fmt.Errorf("empty text")
    }
    return "mock_audio_base64_data_for_" + text, nil
}

func (m *MockTTSService) ListVoicesRaw() ([]byte, error) {
    return []byte(`["female_zh","male_zh","neutral_en"]`), nil
}
