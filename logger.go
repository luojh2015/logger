package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// LoggerLevel0Debug 测试信息
	LoggerLevel0Debug = iota
	// LoggerLevel1Warning 警告信息
	LoggerLevel1Warning
	// LoggerLevel2Error 错误信息
	LoggerLevel2Error
	// LoggerLevel3Fatal 严重信息
	LoggerLevel3Fatal
	// LoggerLevel4Trace 打印信息
	LoggerLevel4Trace
	// LoggerLevel5NoFormat 无格式信息
	LoggerLevel5NoColor
	// LoggerLevel6Off 关闭信息
	LoggerLevel6Off
)

// Logger ...
type Logger struct {
	l          *log.Logger
	fileName   string
	fileReg    string
	fileHandle *os.File
	level      int
	prefix     string
	lock       sync.Mutex
}

// NewLogger ...
// eg: ll := NewLogger(0, "", "")
// eg: ll := NewLogger(logger.LoggerLevelWarning, "demo", "./log/%v.log")
func NewLogger(level int, prefix, file string) *Logger {

	if prefix != "" {
		prefix = "[" + prefix + "]"
	}

	var l *log.Logger
	var fileName string
	var logFile *os.File
	var err error
	if file != "" {
		os.MkdirAll(filepath.Dir(file), 0755)
		fileName = fmt.Sprintf(file, time.Now().Format("2006-01-02"))
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		l = log.New(io.MultiWriter(logFile, os.Stdout), prefix, log.Ltime|log.Lshortfile)
	} else {
		l = log.New(os.Stdout, prefix, log.Ltime|log.Lshortfile)
	}

	return &Logger{l: l, fileName: fileName, fileReg: file, fileHandle: logFile, level: level, prefix: prefix}
}

func (o *Logger) nextLogFile() {
	var _f string
	if strings.Contains(o.fileReg, "%v") {
		_f = fmt.Sprintf(o.fileReg, time.Now().Format("2006-01-02"))
	} else {
		_f = o.fileReg
	}
	// 跨日
	if o.fileReg != "" && _f != o.fileName {
		o.lock.Lock()
		defer o.lock.Unlock()

		oldFileName := o.fileName

		o.fileName = _f
		logFile, err := os.OpenFile(_f, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err == nil {
			if o.fileHandle != nil {
				// 先关闭前一个日志
				o.fileHandle.Close()

				// TODO: 压缩前一月日志
				// 压缩前一天日志
				exec.Command("tar", "czf", oldFileName+".tar.gz", oldFileName).Run()
				os.Remove(oldFileName)
			}
			// 赋值新日志
			o.fileHandle = logFile
			o.l.SetOutput(io.MultiWriter(logFile, os.Stdout))
		}
	}
}

// LogCalldepth ...
func (o *Logger) LogCalldepth(calldepth int, level int, msg ...interface{}) {
	if level < o.level {
		return
	}
	o.nextLogFile()
	o.lock.Lock()
	defer o.lock.Unlock()
	switch level {
	case LoggerLevel0Debug:
		o.l.SetPrefix("\033[32m" + o.prefix)
	case LoggerLevel1Warning:
		o.l.SetPrefix("\033[33m" + o.prefix)
	case LoggerLevel2Error:
		o.l.SetPrefix("\033[31m" + o.prefix)
	case LoggerLevel3Fatal:
		o.l.SetPrefix("\033[31;1;5;7m" + o.prefix)
	case LoggerLevel4Trace:
		o.l.SetPrefix("\033[37m" + o.prefix)
	case LoggerLevel5NoColor:
		o.l.SetPrefix(o.prefix)
		o.l.Output(calldepth, fmt.Sprint(msg...))
		return
	default:
		o.l.SetPrefix(o.prefix)
	}

	o.l.Output(calldepth, fmt.Sprint(msg...)+"\033[m")
}

// SetFile ...
func (o *Logger) SetFile(fileReg string) {
	os.MkdirAll(filepath.Dir(fileReg), 0755)
	o.fileReg = fileReg
	o.nextLogFile()
}

// SetFlag ...
func (o *Logger) SetFlag(flag int) {
	o.l.SetFlags(flag)
}

// SetLevel ...
func (o *Logger) SetLevel(level int) {
	o.level = level
}

// SetPrefix ...
func (o *Logger) SetPrefix(prefix string) {
	if prefix != "" {
		o.prefix = "[" + prefix + "] "
	} else {
		o.prefix = prefix
	}
}

// Println ...
func (o *Logger) Println(v ...interface{}) {
	o.LogCalldepth(3, LoggerLevel5NoColor, fmt.Sprintln(v...))
}

// Printf ...
func (o *Logger) Printf(format string, v ...interface{}) {
	o.LogCalldepth(3, LoggerLevel5NoColor, fmt.Sprintf(format, v...))
}

// Log0Debug ...
func (o *Logger) Log0Debug(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel0Debug, fmt.Sprintf(format, v...))
}

// Log1Warn ...
func (o *Logger) Log1Warn(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel1Warning, fmt.Sprintf(format, v...))
}

// Log2Error ...
func (o *Logger) Log2Error(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel2Error, fmt.Sprintf(format, v...))
}

// Log3Fatal ...
func (o *Logger) Log3Fatal(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel3Fatal, fmt.Sprintf(format, v...))
}

// Log4Trace ...
func (o *Logger) Log4Trace(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel4Trace, fmt.Sprintf(format, v...))
}

// Log5NoFormat ...
func (o *Logger) Log5NoColor(format string, v ...interface{}) {
	if !strings.Contains(format, "%v") && len(v) > 0 {
		format += strings.Repeat("%v", len(v))
	}
	o.LogCalldepth(3, LoggerLevel5NoColor, fmt.Sprintf(format, v...))
}
