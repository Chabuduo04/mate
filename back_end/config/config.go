package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        int
	RedisAddr   string
	LLMModel    string
	ASREndpoint string
	TTSEndpoint string
}

func LoadConfigFromEnv() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}
	return &Config{
		Port:        port,
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		LLMModel:    os.Getenv("LLM_MODEL"),
		ASREndpoint: os.Getenv("ASR_ENDPOINT"),
		TTSEndpoint: os.Getenv("TTS_ENDPOINT"),
	}
}
