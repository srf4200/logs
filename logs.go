// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/issue9/logs/internal/config"
	"github.com/issue9/logs/writers"
)

// 定义了一些日志的类型
const (
	LevelInfo = iota
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelCritical
	levelSize
)

var levels = map[string]int{
	"info":     LevelInfo,
	"trace":    LevelTrace,
	"debug":    LevelDebug,
	"warn":     LevelWarn,
	"error":    LevelError,
	"critical": LevelCritical,
}

var loggers [levelSize]*logger

type logger struct {
	flush writers.Flusher
	log   *log.Logger // 要确保这些值不能为空，因为要保证对应的 ERROR() 等函数的返回值是始终可用的。
}

func (l *logger) set(w io.Writer, prefix string, flag int) {
	if f, ok := w.(writers.Flusher); ok {
		l.flush = f
	}
	l.log = log.New(w, prefix, flag)
}

func init() {
	for index := range loggers {
		loggers[index] = &logger{}
	}

	setDefaultLog()
}

func setDefaultLog() {
	for _, l := range loggers {
		l.set(ioutil.Discard, "", log.LstdFlags)
	}
}

// InitFromXMLFile 从一个 XML 文件中初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLFile(path string) error {
	cfg, err := config.ParseXMLFile(path)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// InitFromXMLString 从一个 XML 字符串初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLString(xml string) error {
	cfg, err := config.ParseXMLString(xml)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// 从 config.Config 中初始化整个 logs 系统
func initFromConfig(cfg *config.Config) error {
	if loggers[0].log != nil { // 加载新配置文件。先输出旧的内容。
		Flush()
		setDefaultLog()
	}

	for name, c := range cfg.Items {
		flag := 0
		flagStr, found := c.Attrs["flag"]
		if found && (len(flagStr) > 0) {
			var err error
			flag, err = parseFlag(flagStr)
			if err != nil {
				return err
			}
		}

		w, err := toWriter(c)
		if err != nil {
			return err
		}
		index, found := levels[name]
		if !found {
			return fmt.Errorf("未知道的二级元素名称:[%v]", name)
		}

		loggers[index].set(w, c.Attrs["prefix"], flag)
	}

	return nil
}

// Flush 输出所有的缓存内容。
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Flush() 输出可能缓存的日志内容。
func Flush() {
	for _, l := range loggers {
		if l.flush != nil {
			l.flush.Flush()
		}
	}
}

// INFO 获取 INFO 级别的 log.Logger 实例，在未指定 info 级别的日志时，该实例返回一个 nil。
func INFO() *log.Logger {
	return loggers[LevelInfo].log
}

// Info 相当于 INFO().Println(v...) 的简写方式
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func Info(v ...interface{}) {
	loggers[LevelInfo].log.Output(2, fmt.Sprintln(v...))
}

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func Infof(format string, v ...interface{}) {
	loggers[LevelInfo].log.Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例，在未指定 debug 级别的日志时，该实例返回一个 nil。
func DEBUG() *log.Logger {
	return loggers[LevelDebug].log
}

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func Debug(v ...interface{}) {
	loggers[LevelDebug].log.Output(2, fmt.Sprintln(v...))
}

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func Debugf(format string, v ...interface{}) {
	loggers[LevelDebug].log.Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例，在未指定 trace 级别的日志时，该实例返回一个 nil。
func TRACE() *log.Logger {
	return loggers[LevelTrace].log
}

// Trace 相当于 TRACE().Println(v...) 的简写方式
func Trace(v ...interface{}) {
	loggers[LevelTrace].log.Output(2, fmt.Sprintln(v...))
}

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func Tracef(format string, v ...interface{}) {
	loggers[LevelTrace].log.Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例，在未指定 warn 级别的日志时，该实例返回一个 nil。
func WARN() *log.Logger {
	return loggers[LevelWarn].log
}

// Warn 相当于 WARN().Println(v...) 的简写方式
func Warn(v ...interface{}) {
	loggers[LevelWarn].log.Output(2, fmt.Sprintln(v...))
}

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func Warnf(format string, v ...interface{}) {
	loggers[LevelWarn].log.Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例，在未指定 error 级别的日志时，该实例返回一个 nil。
func ERROR() *log.Logger {
	return loggers[LevelError].log
}

// Error 相当于 ERROR().Println(v...) 的简写方式
func Error(v ...interface{}) {
	loggers[LevelError].log.Output(2, fmt.Sprintln(v...))
}

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func Errorf(format string, v ...interface{}) {
	loggers[LevelError].log.Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例，在未指定 critical 级别的日志时，该实例返回一个 nil。
func CRITICAL() *log.Logger {
	return loggers[LevelCritical].log
}

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) {
	loggers[LevelCritical].log.Output(2, fmt.Sprintln(v...))
}

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func Criticalf(format string, v ...interface{}) {
	loggers[LevelCritical].log.Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容。
func All(v ...interface{}) {
	all(v...)
}

// Allf 向所有的日志输出内容。
func Allf(format string, v ...interface{}) {
	allf(format, v...)
}

// Fatal 输出错误信息，然后退出程序。
func Fatal(v ...interface{}) {
	all(v...)
	Flush()
	os.Exit(2)
}

// Fatalf 输出错误信息，然后退出程序。
func Fatalf(format string, v ...interface{}) {
	allf(format, v...)
	Flush()
	os.Exit(2)
}

// Panic 输出错误信息，然后触发 panic。
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	all(s)
	Flush()
	panic(s)
}

// Panicf 输出错误信息，然后触发 panic。
func Panicf(format string, v ...interface{}) {
	allf(format, v...)
	Flush()
	panic(fmt.Sprintf(format, v...))
}

func all(v ...interface{}) {
	for _, l := range loggers {
		l.log.Output(3, fmt.Sprintln(v...))
	}
}

func allf(format string, v ...interface{}) {
	for _, l := range loggers {
		l.log.Output(3, fmt.Sprintf(format, v...))
	}
}
