package service

import (
	"context"
	"log"
	"sync"
	"time"
	"twitter_user_news/utils"
)

// 管理所有用户名任务
var (
	taskMap = make(map[string]context.CancelFunc)
	mu      sync.Mutex
)

// StartSearchTask 启动指定用户名的定时搜索任务
func StartSearchTask(userName, userId string, interval time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	// 如果任务已经存在，先停止旧任务
	if cancel, exists := taskMap[userName]; exists {
		cancel()
		delete(taskMap, userName)
		utils.DelUser(userName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	taskMap[userName] = cancel
	utils.AddUser(userName, userId)

	twitterService := NewTwitterService()
	twitterService.Search(userName)
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		log.Printf("监听任务已启动, 用户名: %s, 间隔: %v", userName, interval)
		for {
			select {
			case <-ticker.C:
				twitterService.Search(userName)
			case <-ctx.Done():
				log.Printf("监听任务已停止, 用户名: %s", userName)
				return
			}
		}
	}()
}

// StopSearchTask 停止指定用户名的定时搜索任务
func StopSearchTask(userName string) {
	mu.Lock()
	defer mu.Unlock()

	if cancel, exists := taskMap[userName]; exists {
		cancel()
		delete(taskMap, userName)
		utils.DelUser(userName)
	} else {
		log.Printf("没有找到用户名 [%s] 对应的任务", userName)
	}
}

// StopAllTasks 停止所有定时搜索任务
func StopAllTasks() {
	mu.Lock()
	defer mu.Unlock()

	for key, cancel := range taskMap {
		cancel()
		log.Printf("停止任务: %s", key)
	}
	taskMap = make(map[string]context.CancelFunc)
}

func GetAllTasks() (result []string) {
	mu.Lock()
	defer mu.Unlock()

	for userName, _ := range taskMap {
		result = append(result, userName)
	}
	return
}
