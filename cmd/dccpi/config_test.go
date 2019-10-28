package main

import "testing"

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("config.json")
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

	_, err = LoadConfig("bad-config.json")
	if err == nil {
		t.Error("should have returned an error parsing")
	}
}

func TestSave(t *testing.T) {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal("cannot load config")
	}
	err = SaveConfig(cfg, "config.json")
	if err != nil {
		t.Error("error saving config")
	}
}
