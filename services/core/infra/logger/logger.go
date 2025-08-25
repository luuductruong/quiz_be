package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/quiz_be/services/core/context"
)

const (
	// OutputStdout check config(logger.output).
	OutputStdout = "stdout"
	// OutputFile check config(logger.output).
	OutputFile = "file"
	// FormatText check config(logger.format).
	FormatText = "text"
	// FormatJSON check config(logger.format).
	FormatJson = "json"
)

// Config of logger
type Config struct {
	Output        string `config:"output"`
	File          string `config:"file"`
	Level         string `config:"level"`
	Format        string `config:"format"`
	ExpectingLine bool   `config:"expecting_line"`
}

var Default Logger

func SetDefault(logger Logger) {
	if logger == nil {
		return
	}
	Default = logger
}

type Logger interface {
	Fatal(...interface{})
	Panic(...interface{})
	Error(...interface{})
	Debug(...interface{})

	// with appcontext
	ErrorCtx(context.Context, error, ...interface{})
	DebugCtx(context.Context, ...interface{})
}

type logger struct {
	log           *logrus.Logger
	logfile       *os.File
	expectingLine bool
}

func (l *logger) Fatal(i ...interface{}) {
	l.withDefaultFields().Fatal(i...)
}

func (l *logger) Panic(i ...interface{}) {
	l.withDefaultFields().Panic(i...)
}

func (l *logger) Error(i ...interface{}) {
	l.withDefaultFields().Error(i...)
}

func (l *logger) Debug(i ...interface{}) {
	l.withDefaultFields().Debug(i...)
}

func (l *logger) DebugCtx(ctx context.Context, i ...interface{}) {
	l.withDefaultFields(ctx).Debug(i...)
}

func (l *logger) ErrorCtx(ctx context.Context, err error, i ...interface{}) {
	if err != nil {
		l.withDefaultFields(ctx).Error(append(i, " : ", err)...)
	} else {
		l.withDefaultFields(ctx).Error(i...)
	}
}

func (l *logger) withDefaultFields(ctx ...context.Context) *logrus.Entry {
	locName, locPath := l.getLocationField()
	var fields logrus.Fields
	if len(ctx) > 0 {
		fields = logrus.Fields{"trace": ctx[0].GetTracerId(), locName: locPath}
	} else {
		fields = logrus.Fields{locName: locPath}
	}
	return l.log.WithFields(fields)
}

// getLocationField returns a key-value pair ("location", value)
// where the value indicates the caller's function name and optionally the line number.
func (l *logger) getLocationField() (string, string) {
	pc, _, line, _ := runtime.Caller(3)
	frs := runtime.CallersFrames([]uintptr{pc})
	fr, _ := frs.Next()
	var loc string
	if l.expectingLine {
		loc = fmt.Sprintf("%s:%d", fr.Function, line)
	} else {
		loc = fr.Function
	}
	return "location", loc
}

// NewLogger returns new Logger.
// repository: https://github.com/sirupsen/logrus
func NewLogger(c *Config) Logger {
	if c == nil {
		return nil
	}
	var err error
	var file *os.File

	// new logrus.
	log := logrus.New()

	// set output.
	switch c.Output {
	case OutputStdout: // output: stdout
		log.Out = os.Stdout
		logrus.SetOutput(os.Stdout)
	case OutputFile: // output: file
		file, err = os.OpenFile(c.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			_, fileName := path.Split(c.File)
			file, err = os.OpenFile(fmt.Sprintf("%s/%s", os.Getenv("API_DIR"), fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil
			}
		}
		log.Out = file
		logrus.SetOutput(file)
	default:
		log.Out = os.Stdout
		logrus.SetOutput(os.Stdout)
	}

	// set formatter.
	switch c.Format {
	case FormatText:
		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		formatter.TimestampFormat = time.RFC3339Nano
		// Set specific colors for prefix and timestamp
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle:    "blue+b",
			TimestampStyle: "grey+h",
		})
		log.SetFormatter(formatter)
		logrus.SetFormatter(formatter)
	case FormatJson:
		formatter := &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "func",
			},
		}
		logrus.SetFormatter(formatter)
		log.Formatter = formatter
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
		log.Formatter = &logrus.TextFormatter{}
	}

	level := c.Level
	// set level.
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return nil
	}
	if level == "" {
		level = "info"
	}

	log.Info("log level: ", lv)
	log.SetLevel(lv)
	logrus.SetLevel(lv)

	return &logger{
		log:           log,
		logfile:       file,
		expectingLine: c.ExpectingLine,
	}
}
