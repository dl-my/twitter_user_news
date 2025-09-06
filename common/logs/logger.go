package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/config"
)

var Logger *zap.Logger

func InitLogger(config config.LogConfig) {
	infoCore := getEncoderCore(fmt.Sprintf("%s/%s-info.log", config.LogDir, today()), config, zapcore.InfoLevel)
	warnCore := getEncoderCore(fmt.Sprintf("%s/%s-warn.log", config.LogDir, today()), config, zapcore.WarnLevel)
	errorCore := getEncoderCore(fmt.Sprintf("%s/%s-error.log", config.LogDir, today()), config, zapcore.ErrorLevel)

	cores := []zapcore.Core{infoCore, warnCore, errorCore}

	Logger = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zap.ReplaceGlobals(Logger)
}

func getEncoderCore(filename string, cfg config.LogConfig, level zapcore.Level) zapcore.Core {
	writer := getLogWriter(filename, cfg)
	loc, _ := time.LoadLocation(common.TimeLocation)
	shanghaiTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format(common.TimeFormat))
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = shanghaiTimeEncoder
	encoderConfig.StacktraceKey = zapcore.OmitKey
	encoderConfig.TimeKey = zapcore.OmitKey
	encoderConfig.LevelKey = zapcore.OmitKey    // 不输出日志级别
	encoderConfig.MessageKey = zapcore.OmitKey  // 不输出消息字段
	encoderConfig.CallerKey = zapcore.OmitKey   // 不输出调用者信息
	encoderConfig.NameKey = zapcore.OmitKey     // 不输出logger名称
	encoderConfig.FunctionKey = zapcore.OmitKey // 不输出函数名

	// 根据日志等级调整 Caller / Stacktrace
	switch level {
	case zapcore.WarnLevel:
		encoderConfig.CallerKey = "caller" // 显示 Caller
	case zapcore.ErrorLevel:
		encoderConfig.CallerKey = "caller"         // 显示 Caller
		encoderConfig.StacktraceKey = "stacktrace" // 显示堆栈
	default:
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	}))
}

func getLogWriter(filename string, cfg config.LogConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func today() string {
	loc, _ := time.LoadLocation(common.TimeLocation)
	return time.Now().In(loc).Format(common.DateFormat)
}

//func StartLoggerRotation() {
//	// 初始化一次
//	initLogger()
//
//	c := cron.New(cron.WithSeconds()) // 支持秒级 cron 表达式
//	// 每天 0 点执行一次
//	_, err := c.AddFunc("0 0 0 * * *", func() {
//		initLogger()
//	})
//	if err != nil {
//		log.Printf("%s日志初始化失败,err:%v\n", today(), err)
//		return
//	}
//	c.Start()
//}
