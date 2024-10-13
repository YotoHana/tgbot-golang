package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

var configPath = "config.yml"

type Cfg struct {
	ApiKey string `yaml:"bot_api"`
	AdminChatID int64 `yaml:"admin_chat_id"`
}

func ReadCfg() (*Cfg, error){
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)

	var AppConfig Cfg

	err = decoder.Decode(&AppConfig)

	if err != nil {
		return nil, err
	}
	log.Println("Config parsing is success")
	return &AppConfig, nil
}