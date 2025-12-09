package ioc

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"

	"github.com/muxi-Infra/auditor-Backend/config"
	"github.com/muxi-Infra/auditor-Backend/pkg/logger"
)

func InitLogger(logConfig *config.LogConfig) logger.Logger {
	// 直接使用 zap 本身的配置结构体来处理
	// 配置Lumberjack以支持日志文件的滚动
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logConfig.Path, // 指定日志文件路径
		MaxSize:    5,              // 每个日志文件的最大大小，单位：MB
		MaxBackups: 3,              // 保留旧日志文件的最大个数
		MaxAge:     28,             // 保留旧日志文件的最大天数
		Compress:   true,           // 是否压缩旧的日志文件
	}
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel := getZapLevel(logLevelStr)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := &PrettyJSONEncoder{Encoder: zapcore.NewJSONEncoder(encoderCfg)}

	// 创建zap日志核心
	core := zapcore.NewCore(
		encoder, //zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), // 开发环境用encoder可以看到格式化日志
		zapcore.AddSync(lumberjackLogger),
		logLevel, // 设置日志级别
	)

	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
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

type PrettyJSONEncoder struct {
	zapcore.Encoder
}

func (p *PrettyJSONEncoder) Clone() zapcore.Encoder {
	return &PrettyJSONEncoder{Encoder: p.Encoder.Clone()}
}

func (p *PrettyJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := p.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return buf, err
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &tmp); err != nil {
		return buf, err
	}

	prettyBytes, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return buf, err
	}

	buf.Reset()
	buf.Write(prettyBytes)
	buf.WriteByte('\n')
	return buf, nil
}
