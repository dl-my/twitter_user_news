package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/model"
	"twitter_user_news/utils"
)

// NotificationService 封装Twitter相关服务
type NotificationService struct {
	seenTweets sync.Map
	client     *http.Client
}

// NewNotificationService 创建新的Twitter服务实例
func NewNotificationService() *NotificationService {
	return &NotificationService{
		client: &http.Client{
			Timeout: 30 * time.Second, // 设置超时时间
		},
	}
}

func (s *NotificationService) Search() {
	retryCount := 0
	authToken, ct0 := utils.GetAuthAndCt0()
	for {
		err := s.fetchPosts(authToken, ct0)
		if err == nil {
			fmt.Println("获取成功")
			return
		}
		retryCount++
		logs.Warn("获取失败",
			zap.Error(err),
			zap.Int("retryCount", retryCount))

		if retryCount >= common.MaxRetries {
			// 超过 5 次，重新获取 AuthAndCt0
			logs.Error("连续失败达到最大重试次数，切换认证信息",
				zap.String("authToken", authToken))
			authToken, ct0 = utils.GetAuthAndCt0()
			retryCount = 0
		}

		time.Sleep(1 * time.Second) // 等待一会再试
	}
}

func (s *NotificationService) generateUrl() string {
	params := getNotificationParams()

	return common.Notification + "?" + params.Encode()
}

func (s *NotificationService) fetchPosts(authToken, ct0 string) error {
	reqURL := s.generateUrl()

	req, err := createRequest(reqURL, authToken, ct0)
	if err != nil {
		log.Printf("创建请求失败%v\n", err)
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("发送请求失败%v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
		return fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("x-rate-limit-limit: %s, x-rate-limit-remaining: %s\n", resp.Header.Get("x-rate-limit-limit"), resp.Header.Get("x-rate-limit-remaining"))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据失败%v\n", err)
		return err
	}
	//fmt.Println(resp.StatusCode, string(body))
	fmt.Println(resp.StatusCode, string(body[:600]))

	var raw map[string]interface{}
	if err = json.Unmarshal(body, &raw); err != nil {
		log.Printf("响应解析失败: %v", err)
		return err
	}

	// 判断是否有 errors 字段
	if errs, ok := raw["errors"]; ok {
		log.Printf("接口返回错误: %v", errs)
		return fmt.Errorf("接口返回错误: %v", errs)
	}

	var result model.NotificationResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("响应数据解析失败%v\n", err)
		return err
	}
	for _, tweet := range result.GlobalObjects.Tweets {
		if _, exists := s.seenTweets.Load(tweet.ID); exists {
			continue
		}
		t, err := utils.UtcToShanghai(tweet.CreatedAt)
		if err != nil {
			continue
		}
		publishTime := utils.GetTimeStamp(t)
		fetchTime := utils.GetTimeStamp(time.Now())
		if fetchTime-publishTime > 60 {
			continue
		}
		fmt.Println(time.Now().Unix())
		// 构建日志对象
		posts := model.LogPosts{
			UserId:      tweet.UserID,
			RestId:      tweet.ID,
			ContentEn:   tweet.FullText,
			PublishTime: publishTime,
			FetchTime:   fetchTime,
			Media:       getMedia(tweet.ExtendedEntities.Medias),
		}
		logs.InfoPosts(posts)
		fmt.Println(tweet)
		s.seenTweets.Store(tweet.ID, struct{}{})
	}
	return nil
}
