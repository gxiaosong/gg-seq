package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Database []struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
		URL    string `json:"url"`
	} `json:"database"`
}

func NewConfigFromFile(path string) *Config {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	c := Config{}
	_ = json.Unmarshal(b, &c)
	return &c
}
