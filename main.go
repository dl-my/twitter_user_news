package main

import (
	"log"
	"math/rand"
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
	//twitterSearchService := service.NewSearchTwitterService()
	//twitterSearchService.Search()
	//twitterListService := service.NewListTwitterService()
	//twitterListService.Search("1967411334661988555")

	//ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()
	//log.Println("监听任务已启动")
	//for {
	//	select {
	//	case <-ticker.C:
	//		//twitterSearchService.Search()
	//		twitterListService.Search("1967411334661988555")
	//	}
	//}

	//for {
	//	// 执行任务
	//	twitterListService.Search("1967411334661988555")
	//
	//	// 随机延迟 3-5 秒
	//	delay := time.Duration(rand.Intn(3)+3) * time.Second
	//	log.Printf("下一次执行将在 %v 后", delay)
	//	time.Sleep(delay)
	//}

	//twitterHomeService := service.NewHomeTwitterService()
	//twitterHomeService.Search()
	//ticker := time.NewTicker(3 * time.Second)
	//defer ticker.Stop()
	//log.Println("监听任务已启动")
	//for {
	//	select {
	//	case <-ticker.C:
	//		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	//		twitterHomeService.Search()
	//	}
	//}
	twitterHomeService := service.NewHomeTwitterService()
	baseDelay := 2500 * time.Millisecond

	for {
		// 先发请求
		twitterHomeService.Search()

		// 请求间隔 = 基础周期 ± 随机扰动(0~1s)
		jitter := time.Duration(rand.Intn(1500)) * time.Millisecond
		sleep := baseDelay + jitter

		log.Printf("下一次请求将在 %v 后执行", sleep)
		time.Sleep(sleep)
	}
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
