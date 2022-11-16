package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

// почему хуки? потому что у них более высокая призводительность
// Для того чтобы мы во множество райторов писали любое количество уровней логирования
type writerHook struct {
	Writer   []io.Writer
	LogLevel []logrus.Level
}

// Fire вызывается когда мы пишем что то кудато
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevel
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

// почему так? потому что таким образом мы можем завести второй логер с новым полем например GetLoggerWithField
func GetLogger() *Logger {
	return &Logger{e}
}

//func (l *Logger) GetLoggerWithField(key string, v interface{}) *Logger {
//	return &Logger{l.WithField(key, v)}
//}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileName := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", fileName, f.Line)
		},
		FullTimestamp: true,
		//DisableColors: true,

	}
	err := os.MkdirAll("logs", 0744)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		panic(err)
	}
	l.SetOutput(io.Discard) //отключаем запись логруса
	l.AddHook(&writerHook{
		Writer:   []io.Writer{allFile, os.Stdout},
		LogLevel: logrus.AllLevels,
	})
	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
