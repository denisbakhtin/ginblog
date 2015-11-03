package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

//Configs contains application configurations for all gin modes
type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

//Config contains application configuration for active gin mode
type Config struct {
	Public        string `json:"public"`
	Domain        string `json:"domain"`
	SessionSecret string `json:"session_secret"`
	SignupEnabled bool   `json:"signup_enabled"` //always set to false in release mode (config.json)
	Database      DatabaseConfig
}

//DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
}

//current loaded config
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

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}

//PublicPath returns path to application public folder
func PublicPath() string {
	return config.Public
}

//UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return path.Join(config.Public, "uploads")
}
