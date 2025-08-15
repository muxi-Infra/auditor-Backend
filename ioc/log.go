package ioc

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

func InitLogger(logConfig *config.LogConfig) logger.Logger {
	// 直接使用 zap 本身的配置结构体来处理
	// 配置Lumberjack以支持日志文件的滚动
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logConfig.Path, // 指定日志文件路径
		MaxSize:    50,             // 每个日志文件的最大大小，单位：MB
		MaxBackups: 3,              // 保留旧日志文件的最大个数
		MaxAge:     28,             // 保留旧日志文件的最大天数
		Compress:   true,           // 是否压缩旧的日志文件
	}
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel := getZapLevel(logLevelStr)
	// 创建zap日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		logLevel, // 设置日志级别
	)

	l := zap.New(core, zap.AddCaller())
	res := logger.NewLogger(l)

	return res
}
func getZapLevel(levelStr string) zapcore.Level {
	levelStr = strings.ToLower(levelStr)

	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		// 默认使用Debug级别
		return zapcore.DebugLevel
	}
}
