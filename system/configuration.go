package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//Configs contains application configurations for all gin modes
type Configs struct {
	Debug   Config
	Release Config
}

//Config contains application configuration for active gin mode
type Config struct {
	Public        string `json:"public"`
	Domain        string `json:"domain"`
	SessionSecret string `json:"session_secret"`
	SignupEnabled bool   `json:"signup_enabled"` //always set to false in release mode (config.json)
	Database      DatabaseConfig
	Oauth         OauthConfig
}

//DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
}

//OauthConfig contains oauth client ids and secrets
type OauthConfig struct {
	GoogleClientID string `json:"google_client_id"`
	GoogleSecret   string `json:"google_secret"`
}

//current loaded config
var config *Config

//LoadConfig unmarshals config for current GIN_MODE
func LoadConfig() {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	configs := &Configs{}
	err = json.Unmarshal(data, configs)
	if err != nil {
		panic(err)
	}
	switch gin.Mode() {
	case gin.DebugMode:
		config = &configs.Debug
	case gin.ReleaseMode:
		config = &configs.Release
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

//GetConnectionString returns a database connection string
func GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)
}
