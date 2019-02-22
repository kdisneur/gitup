package config

// Config represents the top level of the YAML configuration
// file
type Config struct {
	Repositories []Repository
}

// Repository represents a repository to synchronize
type Repository struct {
	RemoteURL string
	Path      string
}
