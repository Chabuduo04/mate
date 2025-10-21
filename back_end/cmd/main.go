package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/Chabuduo04/mate/back_end/db"
	"github.com/Chabuduo04/mate/back_end/handlers"
	"github.com/Chabuduo04/mate/back_end/services"
)

func main() {
	// logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// load config
	cfg := config.GetConfig()

	// initialize database (GORM)
	db.InitDB()

	// init session store (Redis or in-memory)
	var sessionStore services.SessionStore
	if cfg.Redis.Host != "" {
		redisStore, err := services.NewRedisSessionStore()
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
	llm := services.NewQiniuLLMService()
	asr := services.NewQiniuASRService()
	tts := services.NewQiniuTTSService()
	kodo := services.NewKodoService()

	// aggregate services
	svc := &services.Services{
		RoleService:  roleSvc,
		SessionStore: sessionStore,
		UserService:  services.NewUserService(),
		LLM:          llm,
		ASR:          asr,
		TTS:          tts,
		Storage:      kodo,
		Logger:       logger,
		LLMModel:     cfg.LLM.Model,
		ASREndpoint:  cfg.ASR.Endpoint,
		TTSEndpoint:  cfg.TTS.Endpoint,
	}

	router := gin.Default()
	handlers.RegisterRoutes(router, svc)

	addr := fmt.Sprintf(":%d", cfg.Main.Port)
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
