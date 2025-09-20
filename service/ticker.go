package service

import (
	"context"
	"log"
	"sync"
	"time"
)

// 管理所有用户名任务
var (
	taskMap = make(map[string]context.CancelFunc)
	mu      sync.Mutex
)

// StartSearchTask 启动指定用户名的定时搜索任务
func StartSearchTask(listId string, interval time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	// 如果任务已经存在，先停止旧任务
	if cancel, exists := taskMap[listId]; exists {
		cancel()
		delete(taskMap, listId)
	}

	ctx, cancel := context.WithCancel(context.Background())
	taskMap[listId] = cancel

	twitterService := NewListTwitterService()
	twitterService.Search(listId)
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		log.Printf("监听任务已启动, listId: %s, 间隔: %v", listId, interval)
		for {
			select {
			case <-ticker.C:
				twitterService.Search(listId)
			case <-ctx.Done():
				log.Printf("监听任务已停止, listId: %s", listId)
				return
			}
		}
	}()
}

// StopSearchTask 停止指定用户名的定时搜索任务
func StopSearchTask(listId string) {
	mu.Lock()
	defer mu.Unlock()

	if cancel, exists := taskMap[listId]; exists {
		cancel()
		delete(taskMap, listId)
	} else {
		log.Printf("没有找到listId [%s] 对应的任务", listId)
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
