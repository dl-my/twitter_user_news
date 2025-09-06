package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twitter_user_news/common/logs"
	"twitter_user_news/config"
	"twitter_user_news/service"
	"twitter_user_news/utils"
)

func init() {
	// 初始化配置
	config.Init()
	// 初始化日志
	logs.InitLogger(config.GlobalConfig.Log)
	// transaction
	utils.Init()
	// 初始化cookie
	if err := utils.LoadTokens(); err != nil {
		log.Println(err)
	}
	// 初始化user
	if err := utils.LoadUsers(); err != nil {
		log.Println("文件加载失败", err)
	}
}

func main2() {
	userId := "902926941413453824"
	token := "741b41d581fbc6732c1665dde9c8f96ce3fe3617"
	ct0 := "ab6842bbf24d08acc142bd93a3ac5d92e36d2f05062e9ddf9bf3869bcfd6b12486a1d941b6e76a68c5e646b24f6d0f77ca2471634f77d449315f830b34437122336e2c3e429d969aaa7982b0eb62bb33"
	service.Posts(userId, token, ct0)
}

func main() {
	twitterService := service.NewTwitterService()
	twitterService.Search("cz_binance")
}

// 监听信号并优雅关闭
func setupGracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("捕获退出信号，正在保存用户")
		if err := utils.SaveUsers(); err != nil {
			log.Println("用户保存失败")
		} else {
			log.Println("用户保存成功")
		}
		time.Sleep(500 * time.Millisecond) // 确保日志输出完毕
		os.Exit(0)
	}()
}
