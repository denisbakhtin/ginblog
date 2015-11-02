package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

type Config struct {
	Public        string `json:"public"`
	SessionSecret string `json:"session_secret"`
	SignupEnabled bool   `json:"signup_enabled"` //always set to false in release mode (config.json)
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
}

var config *Config

//LoadConfig unmarshals config for current GIN_MODE
func LoadConfig(data []byte) {
	configs := &Configs{}
	err := json.Unmarshal(data, configs)
	if err != nil {
		panic(err)
	}
	switch gin.Mode() {
	case gin.DebugMode:
		config = &configs.Debug
	case gin.ReleaseMode:
		config = &configs.Release
	case gin.TestMode:
		config = &configs.Test
	default:
		panic(fmt.Sprintf("Unknown gin mode %s", gin.Mode()))
	}
	if !path.IsAbs(config.Public) {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		config.Public = path.Join(workingDir, config.Public)
	}
}

func GetConfig() *Config {
	return config
}

func PublicPath() string {
	return config.Public
}

func UploadsPath() string {
	return path.Join(config.Public, "uploads")
}
