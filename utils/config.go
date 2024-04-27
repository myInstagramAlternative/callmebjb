package utils

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Configs struct {
	Modem struct {
		Port     string `yaml:"port"`
		BaudRate int    `yaml:"baudrate"`
	} `yaml:"modem"`
	Server struct {
		Listen string `yaml:"listen"`
		Port   string `yaml:"port"`
	} `yaml:"server"`
}

func ReadConfig(cfg *Configs, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

var Config Configs

func InitConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	ReadConfig(&Config, configPath)
}
