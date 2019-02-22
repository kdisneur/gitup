package config_test

import (
	"os"
	"testing"

	"github.com/kdisneur/gitup/pkg/config"
)

func overrideEnv(name string, value string) func() {
	oldValue := os.Getenv(name)
	os.Setenv(name, value)

	return func() { os.Setenv(name, oldValue) }
}

func TestValidConfiguration(t *testing.T) {
	resetEnv := overrideEnv("HOME", "/home/testuser")
	defer resetEnv()

	expectedRepositories := []struct {
		URL  string
		Path string
	}{
		{URL: "git@github.com:kdisneur/gitup", Path: "/home/testuser/Workspace/kdisneur/gitup"},
		{URL: "git@github.com:kdisneur/dotfiles", Path: "/home/testuser/Workspace/kdisneur/somewhereelse"},
	}

	configFile := []byte(`
repositories:
  - url: git@github.com:kdisneur/gitup
    path: ~/Workspace/kdisneur/gitup
  - url: git@github.com:kdisneur/dotfiles
    path: ~/Workspace/kdisneur/somewhereelse
`)

	config, err := config.Parse(configFile)
	if err != nil {
		t.Fatalf("Wrong configuration: %s", err.Error())
	}

	if len(config.Repositories) != len(expectedRepositories) {
		t.Fatalf("Wrong number of repositories. Want: %d, Got: %d", len(expectedRepositories), len(config.Repositories))
	}

	for i := range expectedRepositories {
		if config.Repositories[i].RemoteURL != expectedRepositories[i].URL {
			t.Fatalf("Wrong repository. Want: %+v, Got: %+v", config.Repositories[i], expectedRepositories[i])
		}

		if config.Repositories[i].Path != expectedRepositories[i].Path {
			t.Fatalf("Wrong repository. Want: %+v, Got: %+v", config.Repositories[i], expectedRepositories[i])
		}
	}
}
