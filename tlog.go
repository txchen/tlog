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

// These are the log levels defined in tlog.
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

// tlog exposes *log.Logger instances which can be used to do actual logging.
var (
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger

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

// arrangeLoggers will use the current settings to create log.logger instances.
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

// LogfileLogLevel returns current logfile threashold
func LogfileLogLevel() LogLevel {
	return logfileThreshold
}

// ConsoleLogLevel returns current console threashold.
func ConsoleLogLevel() LogLevel {
	return consoleThreshold
}

// SetLogfileLogLevel sets the logfile loglevel, matching or above will be logged.
// set to LevelOff to turn off logfile logging.
func SetLogfileLogLevel(level LogLevel) {
	logfileThreshold = levelCheck(level)
	arrangeLoggers()
}

// SetConsoleLogLevel sets the console loglevel, matching or above will be logged.
// set to LevelOff to turn off console logging.
func SetConsoleLogLevel(level LogLevel) {
	consoleThreshold = levelCheck(level)
	arrangeLoggers()
}

// SetLogFlags takes "log" package's log flags, sets output flags for the logger.
func SetLogFlags(flag int) {
	logFlag = flag
	arrangeLoggers()
}

// SetLogFile takes filename as input and enables logfile logging, log will be appended to the file.
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

// UseTempLogFile creates a temporary file and set it as the logfile name.
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
