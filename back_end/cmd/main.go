package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/Chabuduo04/mate/back_end/handlers"
	"github.com/Chabuduo04/mate/back_end/services"
)

func main() {
	// logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.LoadConfigFromEnv()

	// init session store (Redis or in-memory)
	var sessionStore services.SessionStore
	if cfg.RedisAddr != "" {
		redisStore, err := services.NewRedisSessionStore(cfg.RedisAddr)
		if err != nil {
			logger.Sugar().Warnf("connect redis fail: %v, fallback to in-memory", err)
			sessionStore = services.NewMemorySessionStore()
		} else {
			sessionStore = redisStore
		}
	} else {
		sessionStore = services.NewMemorySessionStore()
	}

	// init role service (load roles.json)
	roleSvc, err := services.NewRoleService("roles.json")
	if err != nil {
		logger.Sugar().Fatalf("load roles error: %v", err)
	}

	// create API clients (currently mock implementations)
	llm := services.NewMockLLMService()
	asr := services.NewMockASRService()
	tts := services.NewMockTTSService()

	// aggregate services
	svc := &services.Services{
		RoleService:  roleSvc,
		SessionStore: sessionStore,
		LLM:          llm,
		ASR:          asr,
		TTS:          tts,
		Logger:       logger,
		LLMModel:     cfg.LLMModel,
		ASREndpoint:  cfg.ASREndpoint,
		TTSEndpoint:  cfg.TTSEndpoint,
	}

	router := gin.Default()
	handlers.RegisterRoutes(router, svc)

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Sugar().Infof("starting server at %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Sugar().Fatalf("server failed: %v", err)
	}
}
