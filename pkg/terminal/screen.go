package terminal

import (
	"fmt"
	"io"
	"sync"
)

type downloadState rune

const (
	// AwaitingDownload represents a repository not yet downloaded
	AwaitingDownload downloadState = '\u25EF'
	// Downloading represents a repository currently in download
	Downloading = '\u262F'
	// Downloaded represents a repository already downloaded
	Downloaded = '\u25C9'
	// DownloadFailed represents a repository that can't be downloaded
	DownloadFailed = '\u26A0'
)

// Repository represents a Repository
type Repository struct {
	URL     string
	Path    string
	Status  downloadState
	Failure string
}

// Screen represents a download screen
type Screen struct {
	Lines    int
	Writer   io.Writer
	mutex    sync.Mutex
	terminfo Terminfo
}

// NewScreen reprensents a new terminal canva where we can write
func NewScreen(w io.Writer) (*Screen, error) {
	terminfo, err := NewTerminfo()
	if err != nil {
		return nil, err
	}

	return &Screen{
		Writer:   w,
		Lines:    0,
		mutex:    sync.Mutex{},
		terminfo: *terminfo,
	}, nil
}

// NewRepository builds a new Repository in AwaitingDownload state
func NewRepository(url string, path string) Repository {
	return Repository{
		URL:    url,
		Path:   path,
		Status: AwaitingDownload,
	}
}

// StartDownload changes the state of the Repository to Downloading
func (r *Repository) StartDownload() {
	r.Status = Downloading
}

// Fail changes the state of the Repository to DownloadFailed
func (r *Repository) Fail(cause string) {
	r.Status = DownloadFailed
	r.Failure = cause
}

// FinishDownload changes the state of the Repository to Downloaded
func (r *Repository) FinishDownload() {
	r.Status = Downloaded
}

func (s *Screen) clearScreen() string {
	lines := ""
	for i := 0; i < s.Lines; i++ {
		lines += s.terminfo.CursorUp1() + "\r" + s.terminfo.ClearLine()
	}

	return lines
}

func (s *Screen) defaultColor() string {
	return s.terminfo.ResetAttributes()
}

func (s *Screen) errorColor() string {
	return s.terminfo.ForegroundColor(RedColor)
}

// PrintStatus display the download status
func (s *Screen) PrintStatus(repositories []Repository) {
	s.mutex.Lock()

	fmt.Fprint(s.Writer, s.clearScreen())
	s.Lines = 0
	for _, repository := range repositories {
		statusline := fmt.Sprintf("%s  %s", string(repository.Status), repository.URL)
		errorLine := ""

		fmt.Fprintln(s.Writer, statusline)
		s.Lines++

		if repository.Status == DownloadFailed {
			errorLine = fmt.Sprintf("   %s%s%s", s.errorColor(), repository.Failure, s.defaultColor())
			fmt.Fprintln(s.Writer, errorLine)
			s.Lines++
		}
	}
	s.mutex.Unlock()
}
