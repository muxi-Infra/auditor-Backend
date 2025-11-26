package logger

import (
	"go.uber.org/zap"
	defLog "log"
	"os"
)

const defaultLogPrefix = "[muxi-auditor]"

type ZapLogger struct {
	l *zap.Logger
}

func NewLogger(l *zap.Logger) Logger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toArgs(args)...)
}

// 记录异步,缓存等需要警告的事情
func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toArgs(args)...)
}

// 记录发生的错误,数据库级别
func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toArgs(args)...)
}

func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Val))
	}
	return res
}

// DefaultLogger 用于test和部分插件的初始化临时替代
type DefaultLogger struct {
	log *defLog.Logger
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		log: defLog.New(os.Stderr, defaultLogPrefix, defLog.LstdFlags),
	}
}

func (d *DefaultLogger) WithField(key string, value interface{}) interface{} {
	return d
}

func (d *DefaultLogger) Info(msg string, args ...Field) {
	d.log.Print(msg)
}

func (d *DefaultLogger) Debug(msg string, args ...Field) {
	d.log.Printf(msg)
}

func (d *DefaultLogger) Warn(msg string, args ...Field) {
	d.log.Print(msg)
}

func (d *DefaultLogger) Error(msg string, args ...Field) {
	d.log.Print(msg)
}
