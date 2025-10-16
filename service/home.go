package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/model"
	"twitter_user_news/utils"
)

// HomeTwitterService 封装Twitter相关服务
type HomeTwitterService struct {
	seenTweets sync.Map // 使用sync.Map替代全局变量，支持并发安全
	client     *http.Client
}

// NewHomeTwitterService 创建新的Twitter服务实例
func NewHomeTwitterService() *HomeTwitterService {
	return &HomeTwitterService{
		client: &http.Client{
			Timeout: 30 * time.Second, // 设置超时时间
		},
	}
}

func (s *HomeTwitterService) Search() {
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

func (s *HomeTwitterService) generateUrl(seenIds []string) string {
	// 用结构体定义搜索条件
	variablesStruct := model.HomeVariables{
		Count:                  10,
		IncludePromotedContent: false,
		LatestControlAvailable: true,
		RequestContext:         "launch",
	}
	featuresStruct := getDefaultFeatures()

	// 序列化成 JSON
	variablesJSON, _ := json.Marshal(variablesStruct)
	featuresJSON, _ := json.Marshal(featuresStruct)

	params := url.Values{}
	params.Set("variables", string(variablesJSON))
	params.Set("features", string(featuresJSON))
	params.Set("queryId", "4GypC700dutgf7ekTk2cyA")

	return common.HomeLatestTimeLine + "?" + params.Encode()
}

func (s *HomeTwitterService) fetchPosts(authToken, ct0 string) error {
	reqURL := s.generateUrl(nil)

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
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Printf("响应解析失败: %v", err)
		return err
	}

	// 判断是否有 errors 字段
	//if errs, ok := raw["errors"]; ok {
	//	log.Printf("接口返回错误: %v", errs)
	//	return fmt.Errorf("接口返回错误: %v", errs)
	//}

	var result model.HomeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("响应解析失败%v,响应:[%s],状态:[%d]\n", err, string(body), resp.StatusCode)
		return fmt.Errorf("authToken:[%s],err:[%s]", authToken, err.Error())
	}

	return s.processTimeline(result)
}

func (s *HomeTwitterService) processTimeline(result model.HomeResponse) error {
	for _, instruction := range result.Data.Home.HomeTimeLineUrt.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				// 输出评论
				for _, item := range entry.Content.Items {
					s.processTweetOrComment(item.Item.ItemContent.TweetResults.Result)
				}
				// 输出推文
				s.processTweetOrComment(entry.Content.ItemContent.TweetResults.Result)
			}
		}
	}
	return nil
}

func (s *HomeTwitterService) processTweetOrComment(tweet model.Tweet) {
	// 检查是否已处理过
	if _, exists := s.seenTweets.Load(tweet.RestId); exists {
		return
	}
	// 检查是否超时
	if tweet.Legacy.CreatedAt == "" {
		return
	}
	t, err := utils.UtcToShanghai(tweet.Legacy.CreatedAt)
	if err != nil {
		return
	}
	publishTime := utils.GetTimeStamp(t)
	fetchTime := utils.GetTimeStamp(time.Now())
	if fetchTime-publishTime > 60 {
		return
	}
	fmt.Println(tweet)
	// 提取推文内容和媒体
	content, mediaMap := extractTweetContent(tweet)
	if content == "" {
		return
	}
	// 去除转义符
	content = strings.ReplaceAll(content, "\n", "")
	// 翻译为中文
	contentZh, err := utils.Translate(content)
	if err != nil {
		logs.Error("翻译失败", zap.Any("err", err))
		return
	}
	// 构建日志对象
	posts := model.LogPosts{
		UserName:    tweet.Core.UserResults.Result.Core.ScreenName,
		UserId:      tweet.Core.UserResults.Result.RestId,
		RestId:      tweet.RestId,
		ContentEn:   content,
		ContentZh:   contentZh,
		PublishTime: publishTime,
		FetchTime:   fetchTime,
		Media:       mediaMap,
	}
	logs.InfoPosts(posts)
	s.seenTweets.Store(tweet.RestId, struct{}{})
}
