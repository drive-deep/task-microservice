package config

import (
    "log"
    "os"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
    Port int `yaml:"port"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type RedisConfig struct {
    Addr     string `yaml:"addr"`
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
}

type KafkaConfig struct {
    Broker string `yaml:"broker"`
}

func LoadConfig() (*Config, error) {
    file, err := os.Open("config/config.yaml")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func GetConfig() *Config {
    config, err := LoadConfig()
    if err != nil {
        log.Fatalf("could not load config: %v", err)
    }
    return config
}