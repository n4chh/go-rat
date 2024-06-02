package log

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

const DebugLevel = log.DebugLevel
const InfoLevel = log.InfoLevel
const SuccessLevel = 2

const PwnedLevel = 3

const WarnLevel = log.WarnLevel
const ErrorLevel = log.ErrorLevel
const FatalLevel = log.FatalLevel

type logger struct {
	*log.Logger
}

func (l *logger) Pwned(msg string, args ...any) {
	l.Log(PwnedLevel, msg, args...)
}

func (l *logger) Success(msg string, args ...any) {
	l.Log(SuccessLevel, msg, args...)
}

var Log *logger

func InitLogger() *logger {
	styles := log.DefaultStyles()
	styles.Levels[DebugLevel] = lipgloss.NewStyle().
		SetString("[#]").
		Bold(true).
		Foreground(lipgloss.Color("12"))
	styles.Levels[InfoLevel] = lipgloss.NewStyle().
		SetString("[*]").
		Bold(true).
		Foreground(lipgloss.Color("14"))
	styles.Levels[SuccessLevel] = lipgloss.NewStyle().
		SetString("[+]").
		Bold(true).
		Foreground(lipgloss.Color("10"))
	styles.Levels[PwnedLevel] = lipgloss.NewStyle().
		SetString("[$]").
		Bold(true).
		Foreground(lipgloss.Color("#9fef00"))
	styles.Levels[WarnLevel] = lipgloss.NewStyle().
		SetString("[-]").
		Bold(true).
		Foreground(lipgloss.Color("11"))
	styles.Levels[ErrorLevel] = lipgloss.NewStyle().
		SetString("[!]").
		Bold(true).
		Foreground(lipgloss.Color("9"))
	styles.Levels[FatalLevel] = lipgloss.NewStyle().
		SetString("[X]").
		Bold(true).
		Foreground(lipgloss.Color("13"))
	L := log.NewWithOptions(os.Stderr, log.Options{
		Prefix:          "R.A.T. ",
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           InfoLevel,
	})
	Log = new(logger)
	Log.Logger = L
	Log.SetStyles(styles)
	return Log
}
