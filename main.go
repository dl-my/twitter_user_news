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
	//if err := utils.LoadUsers(); err != nil {
	//	log.Println("文件加载失败", err)
	//}
}

func main() {
	twitterService := service.NewListTwitterService()
	twitterService.Search("")

	//ticker := time.NewTicker(10 * time.Second)
	//defer ticker.Stop()
	//log.Println("监听任务已启动")
	//for {
	//	select {
	//	case <-ticker.C:
	//		twitterService.Search()
	//	}
	//}

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
