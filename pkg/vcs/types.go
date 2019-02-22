package vcs

import (
	"context"
	"errors"

	"github.com/kdisneur/gitup/pkg/config"
)

// ErrAlreadyCloned represents an error when the folder is already a Git repository
var ErrAlreadyCloned = errors.New("repository already exists")

// Behavior defines how a Version Control System should behave
type Behavior interface {
	Clone(context.Context, *config.Repository) error
}
