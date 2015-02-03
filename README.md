# tlog [![wercker status](https://app.wercker.com/status/5cf56d48878565a6da12cafd85e1ac8c/s "wercker status")](https://app.wercker.com/project/bykey/5cf56d48878565a6da12cafd85e1ac8c)
Package tlog provides simple APIs to write log to console and file simultaneously.

Tlog is a wrapper around golang's standard [log](http://golang.org/pkg/log) package.

Features:
* Predefined log levels, from debug to critical.
* Console and logfile loglevel can be set separately.
* Console log writes to stderr or stdout according to log level.

Log levels:
* DEBUG - 1
* INFO - 2
* WARN - 3
* ERROR - 4
* CRITICAL - 5
* OFF - 6

### Usage

```go
package main

import (
  "log"
  "github.com/txchen/tlog"
)

func main() {
  // By default, only log to console
  tlog.WARN.Println("warn")

  // Default console loglevel is WARN, we can change it
  tlog.SetConsoleLogLevel(LevelInfo)
  tlog.INFO.Println("now I can show")

  // Set logfile to enable file based logging
  tlog.SetLogFile("z.log")
  tlog.CRITICAL.Printf("2 + 2 = %d", 4)

  // By default, logfile loglevel is INFO, let's change it
  tlog.SetLogfileLogLevel(LevelDebug)
  tlog.Debug.Println("you can only see me in file")

  // By default, tlog use Ldate | Ltime flag, you can tune it
  tlog.SetLogFlags(log.Ldate | log.Ltime | log.Lshortfile)
  // Now you can see different logging style
  tlog.ERROR.Println("different style")

  // And you can turn off the flag
  tlog.SetLogFlags(0)
  // Now the output will be "WARN: clean"
  tlog.WARN.Println("clean")

  // To turn off logging, set the level to off
  tlog.SetLogfileLogLevel(LevelOff)
  tlog.SetConsoleLogLevel(LevelOff)
  tlog.CRITICAL.Println("no one can see me")
}
```
