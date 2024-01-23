package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	want := "127.0.0.1:8080"
	os.Setenv("RENUMBER_ADDRESS", want)

	config := LoadConfig()

	if want != config.Address {
		t.Fatalf(`config.Address = %v, want %v`, config.Address, want)
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	want := "0.0.0.0:8000"
	os.Unsetenv("RENUMBER_ADDRESS")

	config := LoadConfig()

	if want != config.Address {
		t.Fatalf(`config.Address = %v, want %v`, config.Address, want)
	}
}
