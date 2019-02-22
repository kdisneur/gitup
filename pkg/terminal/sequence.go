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
	BlackColor Color = iota
	RedColor
	GreenColor
	YellowColor
	BlueColor
	MagentaColor
	CyanColor
	WhiteColor
)

type Terminfo struct {
	FmtResetAttributes string
	FmtClearLine       string
	FmtCursorUp1       string
	FmtSetForeground   string
}

func NewTerminfo() (*Terminfo, error) {
	return NewTerminfoWithName(os.Getenv("TERM"))
}

func NewTerminfoWithName(name string) (*Terminfo, error) {
	terminfo := fetchTerminfo(name)
	if terminfo == nil {
		return nil, ErrTermNotSupported
	}

	return terminfo, nil
}

func (t Terminfo) ClearLine() string {
	return t.FmtClearLine
}

func (t Terminfo) CursorUp1() string {
	return t.FmtCursorUp1
}

func (t Terminfo) ResetAttributes() string {
	return t.FmtResetAttributes
}

func (t Terminfo) ForegroundColor(color Color) string {
	return fmt.Sprintf(t.FmtSetForeground, color)
}

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
