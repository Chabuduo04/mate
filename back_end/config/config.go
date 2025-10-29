package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type LogConfig struct {
	LogPath string `toml:"logPath"`
}

type JWTConfig struct {
	Secret string `toml:"secret"`
}

type MysqlConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type LLMConfig struct {
	ApiKey string `toml:"apiKey"`
	ApiUrl string `toml:"apiUrl"`
	Model  string `toml:"model"`
}

type ASRConfig struct {
	ApiKey   string `toml:"apiKey"`
	ApiUrl   string `toml:"apiUrl"`
	Endpoint string `toml:"endpoint"`
}

type TTSConfig struct {
	ApiKey   string `toml:"apiKey"`
	ApiUrl   string `toml:"apiUrl"`
	Endpoint string `toml:"endpoint"`
}

type KodoConfig struct {
	Host      string `toml:"host"`
	AccessKey string `toml:"accessKey"`
	SecretKey string `toml:"secretKey"`
	Bucket    string `toml:"bucket"`
}

type Config struct {
	Main  MainConfig  `toml:"main"`
	Log   LogConfig   `toml:"log"`
	JWT   JWTConfig   `toml:"jwt"`
	Mysql MysqlConfig `toml:"mysql"`
	Redis RedisConfig `toml:"redis"`
	LLM   LLMConfig   `toml:"llm"`
	ASR   ASRConfig   `toml:"asr"`
	TTS   TTSConfig   `toml:"tts"`
	Kodo  KodoConfig  `toml:"kodo"`
}

var config *Config

// LoadConfig 从指定的 TOML 文件加载配置，优先使用环境变量 MATE_CONFIG_PATH 指定路径
func LoadConfig() error {
	// 允许通过环境变量覆盖配置路径
	path := os.Getenv("MATE_CONFIG_PATH")
	if path == "" {
		path = "./config/config.toml"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", path)
		return err
	}

	if _, err := toml.DecodeFile(path, config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
		return err
	}

	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = LoadConfig()
	}
	return config
}
