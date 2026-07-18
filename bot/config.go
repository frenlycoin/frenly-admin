package bot

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Config stores configuration vars
type Config struct {
	Dev         bool   `yaml:"dev"`
	TelegramKey string `yaml:"telegram_key"`
	DbURI       string `yaml:"db_uri"`
	GroqAPIKey  string `yaml:"groq_api_key"`
}

// Load method loads configuration file to Config struct
func (c *Config) load(configFile string) {
	file, err := os.Open(configFile)

	if err != nil {
		log.Println(err.Error())

		configFile = os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "/persistent/frenly-admin.config.yaml"
		}

		file, err = os.Open(configFile)
		if err != nil {
			log.Println(err.Error())
		}
	}

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&c)

	if err != nil {
		log.Println(err.Error())
	}
}

// Initializes configuration
func initConfig() *Config {
	c := &Config{}
	c.load("config.yaml")
	return c
}
