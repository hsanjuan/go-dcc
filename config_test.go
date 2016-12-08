package dcc

import "testing"

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("./test/config.json")
	if err != nil {
		t.Error("error loading valid config")
	}

	if len(cfg.Locomotives) != 3 {
		t.Error("config not parsed correctly")
	}

	_, err = LoadConfig("non-existent.json")
	if err == nil {
		t.Error("should have returned an error")
	}

	_, err = LoadConfig("./test/bad-config.json")
	if err == nil {
		t.Error("should have returned an error parsing")
	}
}
