package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Chabuduo04/mate/back_end/models"
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc *services.Services) {
	api := r.Group("/api")
	{
		api.GET("/roles", func(c *gin.Context) {
			roles := svc.RoleService.ListRoles()
			c.JSON(http.StatusOK, roles)
		})
		api.POST("/chat", makeChatHandler(svc))
		api.POST("/asr", makeASRHandler(svc))
		api.POST("/tts", makeTTSHandler(svc))
	}
}

func makeChatHandler(svc *services.Services) gin.HandlerFunc {
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
		sessKey := req.UserID
		if sessKey == "" {
			sessKey = "anon:" + req.RoleID
		}
		hist, _ := svc.SessionStore.Get(ctx, sessKey)

		// call LLM
		reply, err := svc.LLM.Chat(ctx, role.Prompt, hist, req.Message)
		if err != nil {
			svc.Logger.Sugar().Errorf("llm error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "llm error"})
			return
		}

		// append to session history (naive)
		newHist := hist + "\nUser: " + req.Message + "\nRole: " + reply
		svc.SessionStore.Set(ctx, sessKey, newHist, 30*time.Minute)

		c.JSON(http.StatusOK, models.ChatResponse{
			ReplyText: reply,
		})
	}
}

func makeASRHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "audio part required"})
			return
		}
		defer file.Close()

		ctx := context.Background()
		// call ASR
		text, err := svc.ASR.Transcribe(ctx, file)
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
		audioBase64, err := svc.TTS.Synthesize(context.Background(), body.Text, body.Voice)
		if err != nil {
			svc.Logger.Sugar().Errorf("tts error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tts error"})
			return
		}

		// return base64 payload for frontend to decode/play
		c.JSON(http.StatusOK, gin.H{"audio_base64": audioBase64})
	}
}
