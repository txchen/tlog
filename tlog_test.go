package tlog

import (
	"bytes"
	"log"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogLevels(t *testing.T) {
	SetConsoleLogLevel(LevelOff)
	assert.Equal(t, LevelOff, ConsoleLogLevel())
	SetConsoleLogLevel(LevelDebug)
	assert.Equal(t, LevelDebug, ConsoleLogLevel())

	SetLogfileLogLevel(LevelError)
	assert.Equal(t, LevelError, LogfileLogLevel())
	SetLogfileLogLevel(LevelCritical)
	assert.Equal(t, LevelCritical, LogfileLogLevel())
}

func TestSetLogFlags(t *testing.T) {
	SetConsoleLogLevel(LevelDebug)
	consoleBuf := new(bytes.Buffer)
	StderrWriter = consoleBuf
	SetLogFlags(0)
	ERROR.Println("error logging 1")
	assert.Equal(t, "ERROR: error logging 1\n", consoleBuf.String())

	SetLogFlags(log.Ldate)
	CRITICAL.Println("critical logging 1")
	assert.Regexp(t, regexp.MustCompile("CRITICAL: \\d{4}/\\d{2}/\\d{2} critical logging 1"), consoleBuf.String())
}

func TestLogfileAndConsole(t *testing.T) {
	consoleBuf := new(bytes.Buffer)
	logfileBuf := new(bytes.Buffer)
	SetConsoleLogLevel(LevelInfo)
	SetLogfileLogLevel(LevelWarn)
	StdoutWriter = consoleBuf
	LogFileWriter = logfileBuf
	SetLogFlags(0)
	WARN.Println("log line 1")
	assert.Equal(t, "WARN: log line 1\n", consoleBuf.String())
	assert.Equal(t, "WARN: log line 1\n", logfileBuf.String())

	INFO.Println("log line 2")
	assert.Equal(t, "WARN: log line 1\nINFO: log line 2\n", consoleBuf.String())
	assert.Equal(t, "WARN: log line 1\n", logfileBuf.String())
	assert.Fail(t, "fail")
}
