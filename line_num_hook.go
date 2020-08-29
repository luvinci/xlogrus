package xlogrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"strings"
)

/*
	logrus 配置日志输出 文件/函数 中代码行位置输出的Hook
	使用方式：
	var enableLineNum = true // enableLineLog应该设置在配置文件中！！！
	lnHook := NewLineNumLogrusHook(enableLineNum)
	logrus.AddHook(lnHook)
*/
type lineNumLogrusHook struct {
	EnableFileNameLog bool
}

func NewLineNumLogrusHook(enable bool) *lineNumLogrusHook {
	return &lineNumLogrusHook{
		EnableFileNameLog: enable,
	}
}

func (hook *lineNumLogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *lineNumLogrusHook) Fire(entry *logrus.Entry) error {
	if entry.HasCaller() {
		var (
			function string
			line     int
		)
		frame := entry.Caller
		line = frame.Line
		function = frame.Function
		arr := strings.Split(function, ".")
		dir := arr[0]
		_, filename := path.Split(frame.File)
		if hook.EnableFileNameLog {
			entry.Message = fmt.Sprintf("%s/%s:%d %s", dir, filename, line, entry.Message)
		}
	}
	return nil
}
