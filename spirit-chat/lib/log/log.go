package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"time"
)

type fileHook struct {
	levels        []logrus.Level
	formatter     logrus.Formatter
	defaultWriter io.Writer
}

func newFileHook(writer io.Writer) *fileHook {
	hook := &fileHook{
		formatter:     &logrus.TextFormatter{DisableColors: true},
		defaultWriter: writer,
	}
	return hook
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	file, line := getCallerIgnoringLogMulti(1)
	entry.Data["line"] = fmt.Sprintf("%s:%d", file, line)

	return hook.ioWrite(entry)
}

// Write a log line to an io.Writer.
func (hook *fileHook) ioWrite(entry *logrus.Entry) error {
	var (
		writer io.Writer
		msg    []byte
		err    error
	)

	writer = hook.defaultWriter

	msg, err = hook.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	_, err = writer.Write(msg)
	return err
}

func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLog(log_path string, log_level string) {
	log_file, err := rotatelogs.New(log_path,
		rotatelogs.WithRotationTime(24 * time.Hour),
		rotatelogs.WithRotationSize(100 * 1024 ))
	if err != nil {
		panic(err)
	}

	lv, _ := logrus.ParseLevel(log_level)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true}) //防止linux输出没有时间之类的东西
	logrus.SetLevel(lv)
	logrus.AddHook(newFileHook(log_file))
}
