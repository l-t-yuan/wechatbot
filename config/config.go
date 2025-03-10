package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type FeishuConfig struct {
	AppId  string `json:"appId"`
	Token  string `json:"token"`
	Secret string `json:"secret"`
	Encrpy string `json:"encrpy"`
}

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`

	JwtToken string `json:"jwt_token"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`

	TeleToken string         `json:"tele_token"`
	FeiShu    []FeishuConfig `json:"feishu"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		AutoPass := os.Getenv("AutoPass")
		JwtToken := os.Getenv("JwtToken")
		if ApiKey != "" {
			config.JwtToken = JwtToken
		}
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
	})
	return config
}
