package logs

import (
	"go.uber.org/zap"
	"twitter_user_news/model"
)

// Info 封装 info 日志
func Info(fields ...zap.Field) {
	Logger.Info("", fields...)
}

func InfoPosts(posts model.LogPosts) {
	Logger.Info("",
		zap.String("username", posts.UserName),
		zap.String("user_id", posts.UserId),
		zap.String("rest_id", posts.RestId),
		zap.String("content_en", posts.ContentEn),
		zap.String("content_zh", posts.ContentZh),
		zap.Int64("publish_time", posts.PublishTime),
		zap.Int64("fetch_time", posts.FetchTime),
		zap.Any("media", posts.Media),
	)
}

// Warn 封装 warn 日志
func Warn(msg any, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Any("message", msg)}, fields...)
	Logger.Warn("", allFields...)
}

// Error 封装 error 日志
func Error(msg any, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Any("message", msg)}, fields...)
	Logger.Error("", allFields...)
}
