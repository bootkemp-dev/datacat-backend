package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
	} `yaml:"database"`
	Jwt struct {
		JwtKey string `yaml:"jwtKey"`
		Issuer string `yaml:"issuer"`
	} `yaml:"jwt"`
	Accounts struct {
		ResetPasswordTokenLength     int `yaml:"resetPasswordTokenLength"`
		ResetPasswordTokenExpiration int `yaml:"resetPasswordTokenExpiration"`
	}
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func validatePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
