package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/Chabuduo04/mate/back_end/models"
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(r *gin.Engine, svc *services.Services) {
	api := r.Group("/api")
	{
		api.GET("/roles", func(c *gin.Context) {
			roles := svc.RoleService.ListRoles()
			c.JSON(http.StatusOK, roles)
		})
		api.POST("/llm", makeLLMHandler(svc))
		api.POST("/asr", makeASRHandler(svc))
		api.POST("/tts", makeTTSHandler(svc))
		api.POST("/upload", makeKodoHandler(svc))
		api.POST("/voice-chat", makeVoiceChatHandler(svc))
		api.GET("/voice/list", makeVoiceListHandler(svc))
	}
}

func makeLLMHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		role, ok := svc.RoleService.GetRole(req.RoleID)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role not found"})
			return
		}

		// get session context
		ctx := context.Background()
		// sessKey应包含userID和roleID，确保同一用户与不同角色的会话分开
		var sessKey string
		if req.UserID != "" {
			sessKey = req.UserID + ":" + req.RoleID
		} else {
			sessKey = "anon:" + req.RoleID
		}
		hist, _ := svc.SessionStore.Get(ctx, sessKey)

		// 解析历史为ChatMessage数组
		var messages []services.ChatMessage
		messages = append(messages, services.ChatMessage{Role: "system", Content: role.Prompt})
		// 解析历史
		if hist != "" {
			lines := splitHistory(hist)
			messages = append(messages, lines...)
		}
		// 当前用户消息
		messages = append(messages, services.ChatMessage{Role: "user", Content: req.Message})

		// call LLM
		reply, err := svc.LLM.Chat(messages)
		if err != nil {
			svc.Logger.Sugar().Errorf("llm error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "llm error"})
			return
		}

		// call TTS to generate audio for the reply (respect selected voice if provided)
		svc.Logger.Sugar().Infof("calling TTS with text: %s", reply)
		// 优先用请求 voice，否则用角色默认 voice_type
		voiceType := req.Voice
		if voiceType == "" {
			voiceType = role.VoiceType
		}
		fmt.Println("voice_type:", voiceType)
		audioBase64, err := svc.TTS.Synthesize(reply, voiceType)
		if err != nil {
			svc.Logger.Sugar().Errorf("tts error: %v", err)
			// TTS失败不影响文字回复，继续返回文字
			audioBase64 = ""
		} else if audioBase64 != "" {
			audioBase64 = "data:audio/mp3;base64," + audioBase64
			svc.Logger.Sugar().Infof("TTS success, audio length: %d", len(audioBase64))
		}

		// append to session history (naive)
		newHist := hist + "\nUser: " + req.Message + "\nRole: " + reply
		svc.SessionStore.Set(ctx, sessKey, newHist, 30*time.Minute)

		// 返回文字和语音
		c.JSON(http.StatusOK, gin.H{
			"reply_text":   reply,
			"audio_base64": audioBase64,
		})
	}
}

func makeASRHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
			return
		}
		defer file.Close()

		// 生成唯一key
		fileExt := filepath.Ext(header.Filename)
		if fileExt == "" {
			fileExt = ".mp3" // 默认扩展名
		}
		uniqueKey := fmt.Sprintf("asr/%s%s", uuid.New().String(), fileExt)

		// 上传到云存储
		err = svc.Storage.Upload(file, uniqueKey, header.Filename)
		if err != nil {
			svc.Logger.Sugar().Errorf("upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload error"})
			return
		}

		// 构建音频URL
		audioURL := fmt.Sprintf("%s/%s", config.AppConfig.KodoHost, uniqueKey)

		// 调用ASR
		text, err := svc.ASR.TranscribeFromURL(audioURL)
		if err != nil {
			svc.Logger.Sugar().Errorf("asr error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "asr error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"text": text})
	}
}

func makeTTSHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Text  string `json:"text" binding:"required"`
			Voice string `json:"voice"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call TTS
		audioBase64, err := svc.TTS.Synthesize(body.Text, body.Voice)
		if err != nil {
			svc.Logger.Sugar().Errorf("tts error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tts error"})
			return
		}
		if audioBase64 != "" {
			audioBase64 = "data:audio/mp3;base64," + audioBase64
		}
		// return base64 payload for frontend to decode/play
		c.JSON(http.StatusOK, gin.H{"audio_base64": audioBase64})
	}
}

func makeKodoHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
			return
		}
		defer file.Close()
		// call Upload
		err = svc.Storage.Upload(file, "user123", time.Now().Format("202506010000"))
		if err != nil {
			svc.Logger.Sugar().Errorf("upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"upload": true})
	}
}

func makeVoiceChatHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取上传的音频文件
		file, header, err := c.Request.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
			return
		}
		defer file.Close()

		// 1.1 读取文件内容到内存，避免多次读取导致无声
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "read audio error"})
			return
		}

		// 2. 获取角色ID和用户ID
		roleID := c.PostForm("role_id")
		userID := c.PostForm("user_id")
		if roleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role_id is required"})
			return
		}

		// 3. 验证角色是否存在
		role, ok := svc.RoleService.GetRole(roleID)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role not found"})
			return
		}

		// 4. 生成唯一的文件key
		fileExt := filepath.Ext(header.Filename)
		if fileExt == "" {
			fileExt = ".wav" // 默认扩展名
		}
		uniqueKey := fmt.Sprintf("voice-chat/%s%s", uuid.New().String(), fileExt)

		// 5. 上传音频文件到云存储（用新 reader）
		uploadReader := bytes.NewReader(fileBytes)
		err = svc.Storage.Upload(uploadReader, uniqueKey, header.Filename)
		if err != nil {
			svc.Logger.Sugar().Errorf("upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload error"})
			return
		}

		// 6. 构建音频URL
		audioURL := fmt.Sprintf("%s%s/%s", "http://", config.AppConfig.KodoHost, uniqueKey)
		fmt.Println("audioURL:", audioURL)
		// 7. 调用ASR将音频转换为文字
		transcribedText, err := svc.ASR.TranscribeFromURL(audioURL)
		if err != nil {
			svc.Logger.Sugar().Errorf("asr error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "asr error"})
			return
		}

		// 8. 获取会话上下文
		ctx := context.Background()
		sessKey := userID
		if sessKey == "" {
			sessKey = "anon:" + roleID
		}
		hist, _ := svc.SessionStore.Get(ctx, sessKey)

		// 9. 组装历史消息
		var messages []services.ChatMessage
		messages = append(messages, services.ChatMessage{Role: "system", Content: role.Prompt})
		if hist != "" {
			lines := splitHistory(hist)
			messages = append(messages, lines...)
		}
		messages = append(messages, services.ChatMessage{Role: "user", Content: transcribedText})

		// 9. 调用LLM处理文字
		reply, err := svc.LLM.Chat(messages)

		if err != nil {
			svc.Logger.Sugar().Errorf("llm error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "llm error"})
			return
		}

		// 10. 调用TTS将回复转换为语音（允许通过query/form传入voice选择）
		selectedVoice := c.Query("voice")
		if selectedVoice == "" {
			selectedVoice = c.PostForm("voice")
		}
		audioBase64, err := svc.TTS.Synthesize(reply, selectedVoice)
		if err != nil {
			svc.Logger.Sugar().Errorf("tts error: %v", err)
			// TTS失败不影响文字回复，继续返回文字
			audioBase64 = ""
		} else if audioBase64 != "" {
			audioBase64 = "data:audio/mp3;base64," + audioBase64
		}

		// 11. 更新会话历史
		newHist := hist + "\nUser: " + transcribedText + "\nRole: " + reply
		svc.SessionStore.Set(ctx, sessKey, newHist, 30*time.Minute)

		// 12. 返回转录文字、AI回复文字和语音
		c.JSON(http.StatusOK, gin.H{
			"transcribed_text": transcribedText,
			"reply_text":       reply,
			"audio_base64":     audioBase64,
		})
	}
}

func makeVoiceListHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := svc.TTS.ListVoicesRaw()
		if err != nil {
			svc.Logger.Sugar().Errorf("list voices error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list voices error"})
			return
		}
		c.Data(http.StatusOK, "application/json", body)
	}
}

// splitHistory 将历史字符串解析为ChatMessage数组
func splitHistory(hist string) []services.ChatMessage {
	var msgs []services.ChatMessage
	lines := make([]string, 0)
	curr := ""
	for _, r := range hist {
		if r == '\n' {
			lines = append(lines, curr)
			curr = ""
		} else {
			curr += string(r)
		}
	}
	if curr != "" {
		lines = append(lines, curr)
	}
	for _, line := range lines {
		if len(line) > 6 && line[:6] == "User: " {
			msgs = append(msgs, services.ChatMessage{Role: "user", Content: line[6:]})
		} else if len(line) > 6 && line[:6] == "Role: " {
			msgs = append(msgs, services.ChatMessage{Role: "assistant", Content: line[6:]})
		}
	}
	return msgs
}
