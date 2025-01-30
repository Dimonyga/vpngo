package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Main     MainConfig `yaml:"main"`
	Peers    []Peer     `yaml:"peers"`
	Me       *Peer
	Neighbor *Peer
}

type MainConfig struct {
	Timeout   int    `yaml:"timeout"`
	Keepalive int    `yaml:"keepalive"`
	Secret    string `yaml:"secret"`
}

type Peer struct {
	Name     string `yaml:"name"`
	EAddress string `yaml:"eaddress"`
	IAddress string `yaml:"iaddress"`
	IMask    string `yaml:"imask"`
	Nat      bool   `yaml:"nat"`
	Port     int    `yaml:"port"`
}

// getConfig загружает конфигурацию из файла и возвращает её
func getConfig(name string) Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	for i, peer := range config.Peers {
		if peer.Name == name {
			config.Me = &config.Peers[i]
		} else {
			config.Neighbor = &config.Peers[i]
		}
	}
	if config.Neighbor.Nat {
		config.Neighbor.Port = 0
	}

	return config
}

func (c *Config) SetNeighborPort(port int) {
	c.Neighbor.Port = port
}

func (c *Config) SetNeighborEaddress(eaddress string) {
	c.Neighbor.EAddress = eaddress
}
