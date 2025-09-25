package config

import (
	"os"
	"strconv"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	RedisAddr   string
	ApiKey		string
	ApiUrl		string
	LLMModel    string
	ASREndpoint string
	TTSEndpoint string
	KodoHost	string
	AccessKey 	string
	SecretKey 	string
	Bucket	  	string
}

var AppConfig *Config

func InitConfig() {
	_ = godotenv.Load()
	AppConfig = &Config{
		Port:        getEnvAsInt("PORT", 8080),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		ApiKey:	 os.Getenv("API_KEY"),
		ApiUrl:	 os.Getenv("API_URL"),
		LLMModel:    os.Getenv("LLM_MODEL"),
		ASREndpoint: os.Getenv("ASR_ENDPOINT"),
		TTSEndpoint: os.Getenv("TTS_ENDPOINT"),
		KodoHost:	os.Getenv("KODO_HOST"),
		AccessKey:	os.Getenv("ACCESS_KEY"),
		SecretKey:	os.Getenv("SECRET_KEY"),
		Bucket:		os.Getenv("BUCKET"),
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
