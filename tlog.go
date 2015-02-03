package tlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// LogLevel 1-5, from debug to critical, 6 is off
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelOff
	DefaultLogfileThreshold = LevelInfo
	DefaultConsoleThreshold = LevelWarn
)

type tlogger struct {
	Level  LogLevel
	Prefix string
	Logger **log.Logger
}

var (
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger

	// by default, not logging to log file
	LogFileWriter = ioutil.Discard
	StdoutWriter  = os.Stdout
	StderrWriter  = os.Stderr

	debug    = &tlogger{Level: LevelDebug, Logger: &DEBUG, Prefix: "DEBUG: "}
	info     = &tlogger{Level: LevelInfo, Logger: &INFO, Prefix: "INFO: "}
	warn     = &tlogger{Level: LevelWarn, Logger: &WARN, Prefix: "WARN: "}
	err      = &tlogger{Level: LevelError, Logger: &ERROR, Prefix: "ERROR: "}
	critical = &tlogger{Level: LevelCritical, Logger: &CRITICAL, Prefix: "CRITICAL: "}

	tloggers         = []*tlogger{debug, info, warn, err, critical}
	logfileThreshold = DefaultLogfileThreshold
	consoleThreshold = DefaultConsoleThreshold

	logFlag = log.Ldate | log.Ltime
)

func init() {
	arrangeLoggers()
}

// initialize will setup the jWalterWeatherman standard approach of providing the user
// some feedback and logging a potentially different amount based on independent log and output thresholds.
// By default the output has a lower threshold than logged
// Don't use if you have manually set the Handles of the different levels as it will overwrite them.
func arrangeLoggers() {
	for _, tl := range tloggers {
		var wrtier io.Writer
		if tl.Level < consoleThreshold && tl.Level < logfileThreshold {
			wrtier = ioutil.Discard
		} else if tl.Level >= consoleThreshold && tl.Level >= logfileThreshold {
			wrtier = io.MultiWriter(getConsoleWriter(tl.Level), LogFileWriter)
		} else if tl.Level >= consoleThreshold && tl.Level < logfileThreshold {
			wrtier = getConsoleWriter(tl.Level)
		} else {
			wrtier = LogFileWriter
		}
		*tl.Logger = log.New(wrtier, tl.Prefix, logFlag)
	}
}

func getConsoleWriter(level LogLevel) io.Writer {
	if level < LevelError {
		return StdoutWriter
	}
	return StderrWriter
}

// Ensures that the level provided is within the bounds of available levels
func levelCheck(level LogLevel) LogLevel {
	switch {
	case level <= LevelDebug:
		return LevelDebug
	case level >= LevelOff:
		return LevelOff
	default:
		return level
	}
}

// Get the current logfile threashold
func LogfileLogLevel() LogLevel {
	return logfileThreshold
}

// Get the current console threashold
func ConsoleLogLevel() LogLevel {
	return consoleThreshold
}

// Establishes a threshold where anything matching or above will be logged
func SetLogfileLogLevel(level LogLevel) {
	logfileThreshold = levelCheck(level)
	arrangeLoggers()
}

// Establishes a threshold where anything matching or above will be output
func SetConsoleLogLevel(level LogLevel) {
	consoleThreshold = levelCheck(level)
	arrangeLoggers()
}

func SetLogFlags(flag int) {
	logFlag = flag
	arrangeLoggers()
}

// Conveniently Sets the Log Handle to a io.writer created for the file behind the given filepath
// Will only append to this file
func SetLogFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		CRITICAL.Println("Failed to open log file:", path, err)
		return err
	}
	fmt.Println("Logging to file", file.Name())

	LogFileWriter = file
	arrangeLoggers()
	return nil
}

// Conveniently Creates a temporary file and sets the Log Handle to a io.writer created for it
func UseTempLogFile(prefix string) error {
	file, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		CRITICAL.Println(err)
		return err
	}

	fmt.Println("Logging to file:", file.Name())

	LogFileWriter = file
	arrangeLoggers()
	return nil
}
