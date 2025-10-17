package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	RedisAddr   string
	RedisPass   string
	ApiKey      string
	ApiUrl      string
	LLMModel    string
	ASREndpoint string
	TTSEndpoint string
	KodoHost    string
	AccessKey   string
	SecretKey   string
	Bucket      string
}

var AppConfig *Config

func InitConfig() {
	_ = godotenv.Load()
	AppConfig = &Config{
		Port:        getEnvAsInt("PORT", 8080),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   os.Getenv("REDIS_PASS"),
		ApiKey:      getEnvRequired("API_KEY"),
		ApiUrl:      getEnvRequired("API_URL"),
		LLMModel:    getEnvRequired("LLM_MODEL"),
		ASREndpoint: getEnvRequired("ASR_ENDPOINT"),
		TTSEndpoint: getEnvRequired("TTS_ENDPOINT"),
		KodoHost:    getEnvRequired("KODO_HOST"),
		AccessKey:   getEnvRequired("ACCESS_KEY"),
		SecretKey:   getEnvRequired("SECRET_KEY"),
		Bucket:      getEnvRequired("BUCKET"),
	}
}

// getEnvAsInt 读取整数环境变量
func getEnvAsInt(key string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			log.Fatalf("环境变量 %s 必须是整数，但得到: %s", key, valueStr)
		}
		return value
	}
	return defaultVal
}

// getEnv 读取字符串环境变量，支持默认值
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvRequired 读取必须的环境变量（敏感数据），未设置时退出
func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("必须设置环境变量 %s", key)
	}
	return value
}
