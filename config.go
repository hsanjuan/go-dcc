package dcc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config allows to store configuration settings to initialize go-dcc.
type Config struct {
	Locomotives []*Locomotive `json:"locomotives"`
}

// LoadConfig parses a configuration file and returns a Config object.
func LoadConfig(path string) (*Config, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(conf, &cfg)
	if err != nil {
		return nil, err
	}
	log.Printf("XLoaded configuration for %d locomotive(s)", len(cfg.Locomotives))
	return &cfg, nil
}
