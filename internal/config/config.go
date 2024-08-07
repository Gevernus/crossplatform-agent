package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIURL        string `yaml:"api_url"`
	PollInterval  int    `yaml:"poll_interval"`
	LogLevel      string `yaml:"log_level"`
	LogPath       string `yaml:"log_path"`
	AgentID       string `yaml:"agent_id"`
	AgentPassword string `yaml:"agent_password"`
	ConfigPath    string `yaml:"-"`
}

func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.ConfigPath = path
	return &cfg, nil
}
