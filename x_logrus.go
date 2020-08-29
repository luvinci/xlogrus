package xlogrus

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

func Init(opts ...Option) {
	opt := newOptions(opts...)
	initLogFormat(opt.EnableLineNum, opt.LogLevel)
	initLogSettings(opt.Address, opt.Name, opt.MaxAge, opt.RotationTime, opt.EnableSaveToFile)
}

// 定义日志格式样式
func initLogFormat(enableLineNum bool, logLevel logrus.Level) {
	formatter := &prefixed.TextFormatter{}
	// 格式化日志样式
	formatter.ForceFormatting = true
	formatter.ForceColors = true
	formatter.DisableColors = false
	// 开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	// 定义时间格式
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	// ---- 设置logrus的formatter ----
	logrus.SetFormatter(formatter)
	logrus.SetOutput(colorable.NewColorableStdout())
	// 日志级别（默认InfoLevel），应该设置在配置文件中！！！
	logrus.SetLevel(logLevel)
	// 记录到文件中的日志 是否包含代码行的信息
	logrus.SetReportCaller(true)
	// 设置文件/函数中代码行信息的输出的hook
	lnHook := NewLineNumLogrusHook(enableLineNum)
	logrus.AddHook(lnHook)
}

// 初始化log配置，配置logrus日志文件滚动生成
func initLogSettings(logPath, logName string, maxAge, rotationTime time.Duration, saveToFile bool) {
	if !saveToFile {
		return
	}
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Println("create log dir failed")
		return
	}
	baseLogPath := path.Join(logPath, logName)
	// 设置滚动日志输出writer
	writer, err := rotatelogs.New(
		// 日志文件名
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H%M",
		// 为最新的日志建立软链接，以方便找到当前日志文件
		rotatelogs.WithLinkName(baseLogPath),
		// 日志文件清理前的最长保存时间
		rotatelogs.WithMaxAge(maxAge),
		// 日志分割时间间隔
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		fmt.Println("failed to create rotatelogs: ", err)
		return
	}

	// 设置日志文件输出的日志格式
	formatter := &logrus.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatter)
	logrus.AddHook(lfHook)
}
