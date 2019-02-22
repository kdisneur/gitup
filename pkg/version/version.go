package version

import (
	"fmt"
	"runtime"
)

var (
	buildDate = "1970-01-01T00:00:00Z"
	gitBranch = "unknown"
	gitCommit = "unknown"
	gitState  = "unknown"
)

// Info represents all the information about the command line, from the commit
// to the destination platform
type Info struct {
	BuildDate string
	Compiler  string
	GitBranch string
	GitCommit string
	GitState  string
	GoVersion string
	Platform  string
}

// GetInfo returns the current versoin of the command line
// The informations come from ldflags. You can look at the Makefile
// for more information
func GetInfo() Info {
	return Info{
		BuildDate: buildDate,
		Compiler:  runtime.Compiler,
		GitBranch: gitBranch,
		GitCommit: gitCommit,
		GitState:  gitState,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
