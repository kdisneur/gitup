package vcs

import (
	"context"
	"fmt"
	"os"
	"path"

	"gopkg.in/src-d/go-git.v4"

	"github.com/kdisneur/gitup/pkg/config"
)

// Git represents a Git Version Control System
type Git struct{}

// NewGit creates a new Git Version Control System
func NewGit() Behavior {
	return &Git{}
}

// Clone copies a remote project and its history to a local path
func (g Git) Clone(ctx context.Context, repository *config.Repository) error {
	gitPath := path.Join(repository.Path, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		return ErrAlreadyCloned
	}

	_, err := git.PlainCloneContext(ctx, repository.Path, false, &git.CloneOptions{
		URL:               repository.RemoteURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return fmt.Errorf("can't clone %s to %s: %s", repository.RemoteURL, repository.Path, err.Error())
	}

	return nil
}
