package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/mitchellh/go-homedir"

	"github.com/kdisneur/gitup/pkg/config"
	"github.com/kdisneur/gitup/pkg/terminal"
	"github.com/kdisneur/gitup/pkg/vcs"
)

func main() {
	rawConfigParh := "~/.gituprc"
	versionControl := vcs.NewGit()
	ctx := context.Background()

	configPath, err := homedir.Expand(rawConfigParh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't expand path %s: %s", rawConfigParh, err.Error())
		os.Exit(1)
	}

	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read %s: %s", configPath, err.Error())
		os.Exit(1)
	}

	configuration, err := config.Parse(configContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't evaluate configuration file %s: %s", configPath, err.Error())
		os.Exit(1)
	}

	screen, err := terminal.NewScreen(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	var group sync.WaitGroup
	repositories := make([]terminal.Repository, len(configuration.Repositories))

	for i, repository := range configuration.Repositories {
		repositories[i] = terminal.NewRepository(repository.RemoteURL, repository.Path)
	}
	screen.PrintStatus(repositories)

	for i := range configuration.Repositories {
		group.Add(1)
		go func(repository *terminal.Repository, configuredRepository *config.Repository, group *sync.WaitGroup) {
			defer group.Done()
			repository.StartDownload()

			screen.PrintStatus(repositories)

			err = versionControl.Clone(ctx, configuredRepository)

			if err != nil && err != vcs.ErrAlreadyCloned {
				repository.Fail(err.Error())
				screen.PrintStatus(repositories)
			} else {
				repository.FinishDownload()
				screen.PrintStatus(repositories)
			}
		}(&repositories[i], &configuration.Repositories[i], &group)
	}
	group.Wait()
}
