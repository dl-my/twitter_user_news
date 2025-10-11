package main

import (
	"log"
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
	// 初始化cookie
	if err := utils.LoadTokens(); err != nil {
		log.Println(err)
	}
	// transaction-id
	//utils.Init()
}

func main() {
	twitterService := service.NewNotificationService()
	twitterService.Search()
	ticker := time.NewTicker(6000 * time.Millisecond)
	defer ticker.Stop()
	log.Println("监听任务已启动")
	for {
		select {
		case <-ticker.C:
			//time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			twitterService.Search()
		}
	}
}
