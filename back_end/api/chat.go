package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
)

type ChatApi struct {
	Service services.Services
}

func NewChatApi(service services.Services) ChatApi {
	return ChatApi{Service: service}
}

func (c *ChatApi) ChatWithText(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	role, ok := c.Service.RoleService.GetRole(req.RoleID)
	if !ok {
		ctx.JSON(404, gin.H{"error": "role not found"})
		return
	}

	//todo: get session context using req.UserID if provided

	var messages []services.ChatMessage
	messages = append(messages, services.ChatMessage{Role: "system", Content: role.Prompt})

	// call LLM
	reply, err := c.Service.LLM.Chat(messages)
	if err != nil {
		c.Service.Logger.Sugar().Errorf("llm error: %v", err)
		ctx.JSON(500, gin.H{"error": "llm error"})
		return
	}

	voiceType := req.Voice
	if voiceType == "" {
		voiceType = role.VoiceType
	}
	audioBase64, err := c.Service.TTS.Synthesize(reply, voiceType)
	if err != nil {
		c.Service.Logger.Sugar().Errorf("tts error: %v", err)
		// TTS失败不影响文字回复，继续返回文字
		audioBase64 = ""
	} else if audioBase64 != "" {
		audioBase64 = "data:audio/mp3;base64," + audioBase64
		c.Service.Logger.Sugar().Infof("TTS success, audio length: %d", len(audioBase64))
	}

	ctx.JSON(200, gin.H{
		"reply_text":   reply,
		"audio_base64": audioBase64,
	})
}

func (c *ChatApi) ChatWithAudio(ctx *gin.Context) {
	// 获取上传的音频文件
	file, header, err := ctx.Request.FormFile("audio")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
		return
	}
	defer file.Close()

	// 上传音频文件
	audioURL, err := c.Service.Storage.UploadAudio(file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload audio"})
		return
	}
	c.Service.Logger.Sugar().Infof("audio uploaded to: %s", audioURL)
	// 将音频文件转换为文本
	audioText, err := c.Service.ASR.TranscribeFromURL(audioURL)
	if err != nil {
		c.Service.Logger.Sugar().Errorf("ASR error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "audio to text error"})
		return
	}

	// 使用转换后的文本进行聊天
	ctx.Request.Body = io.NopCloser(strings.NewReader(audioText))
	c.ChatWithText(ctx)
	ctx.JSON(200, gin.H{"transcribed_text": audioText})
}
