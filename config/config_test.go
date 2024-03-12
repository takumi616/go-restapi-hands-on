package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	//Set port number to environment variable PORT
	t.Setenv("PORT", fmt.Sprint(wantPort))

	//Get config
	got, err := New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	//Compare got to want
	if got.Port != wantPort {
		t.Errorf("want %d, but %d", wantPort, got.Port)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
}
