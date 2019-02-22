package terminal

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

// ErrTermNotSupported is an error returned when the current terminal
// is not supported.
var ErrTermNotSupported = errors.New("current terminal is not supported")

// Color constant type
type Color int

const (
	// BlackColor black termimal color
	BlackColor Color = iota
	// RedColor red terminal color
	RedColor
	// GreenColor green terminal color
	GreenColor
	// YellowColor yellow terminal color
	YellowColor
	// BlueColor blue terminal color
	BlueColor
	// MagentaColor magenta terminal color
	MagentaColor
	// CyanColor cyan terminal color
	CyanColor
	// WhiteColor white terminal color
	WhiteColor
)

// Terminfo represents all the escape sequence supported by the application
type Terminfo struct {
	FmtResetAttributes string
	FmtClearLine       string
	FmtCursorUp1       string
	FmtSetForeground   string
}

// NewTerminfo creates a new Terminfo based on the current TERM environment variable.
// It returns an error if the terminal is not supported
func NewTerminfo() (*Terminfo, error) {
	return NewTerminfoWithName(os.Getenv("TERM"))
}

// NewTerminfoWithName creates a new Terminfo based on a name
// It returns an error if the terminal is not supported
func NewTerminfoWithName(name string) (*Terminfo, error) {
	terminfo := fetchTerminfo(name)
	if terminfo == nil {
		return nil, ErrTermNotSupported
	}

	return terminfo, nil
}

// ClearLine clears the current line
func (t Terminfo) ClearLine() string {
	return t.FmtClearLine
}

// CursorUp1 moves the cursor up one line
func (t Terminfo) CursorUp1() string {
	return t.FmtCursorUp1
}

// ResetAttributes reset all formatting to use default terminal configuration
func (t Terminfo) ResetAttributes() string {
	return t.FmtResetAttributes
}

// ForegroundColor changes foreground color
func (t Terminfo) ForegroundColor(color Color) string {
	return fmt.Sprintf(t.FmtSetForeground, color)
}

// Columns returns the number of columns of the terminal. It returns infinity
// if it can't find the information
func (t Terminfo) Columns() int {
	rawColumns := os.Getenv("COLUMNS")

	columns, err := strconv.Atoi(rawColumns)
	if err != nil {
		return int(math.Inf(1))
	}

	return columns
}

func fetchTerminfo(name string) *Terminfo {
	return map[string]*Terminfo{
		"xterm-256color": &Terminfo{
			FmtClearLine:       "\x1b[K",
			FmtCursorUp1:       "\x1b[A",
			FmtSetForeground:   "\x1b[3%dm",
			FmtResetAttributes: "\x1b[0m",
		},
	}[name]
}
