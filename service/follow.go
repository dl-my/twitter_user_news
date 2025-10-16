package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/model"
	"twitter_user_news/utils"
)

// FollowService 封装Twitter相关服务
type FollowService struct {
	client *http.Client
}

// NewFollowService 创建新的Twitter服务实例
func NewFollowService() *FollowService {
	return &FollowService{
		client: &http.Client{
			Timeout: 30 * time.Second, // 设置超时时间
		},
	}
}

func (s *FollowService) Create(userID int) {
	retryCount := 0
	authToken, ct0 := utils.GetAuthAndCt0()
	for {
		// 1️⃣ 先关注用户
		err := s.doFollow(authToken, ct0, userID)
		if err == nil {
			fmt.Println("关注成功，准备开启通知...")
			// 2️⃣ 关注成功后再开启通知
			err = s.enableNotification(authToken, ct0, userID)
			if err == nil {
				fmt.Println("通知开启成功 ✅")
				return
			}
			logs.Warn("通知开启失败", zap.Error(err))
			return
		}

		retryCount++
		logs.Error("关注失败", zap.Error(err), zap.Int("retryCount", retryCount))

		if retryCount >= common.MaxRetries {
			logs.Error("连续失败达到最大重试次数，切换认证信息",
				zap.String("authToken", authToken))
			authToken, ct0 = utils.GetAuthAndCt0()
			retryCount = 0
		}
		time.Sleep(1 * time.Second)
	}
}

func structToValues(s interface{}) url.Values {
	values := url.Values{}
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("form")
		if tag == "" {
			continue
		}
		value := v.Field(i).Interface()
		values.Set(tag, fmt.Sprintf("%v", value))
	}
	return values
}

func buildTwitterRequest(authToken, ct0, reqURL string, data interface{}) (*http.Request, error) {
	form := structToValues(data)
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
	req.AddCookie(&http.Cookie{Name: "ct0", Value: ct0})

	req.Header.Set("Authorization", common.Authorization)
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Csrf-Token", ct0)
	req.Header.Set("Referer", "https://x.com")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Client-Language", "en")

	return req, nil
}

func followRequest(authToken, ct0 string, userID int) (*http.Request, error) {
	reqData := model.Create{
		IncludeProfileInterstitialType: 1,
		IncludeBlocking:                1,
		IncludeBlockedBy:               1,
		IncludeFollowedBy:              1,
		IncludeWantRetweets:            1,
		IncludeMuteEdge:                1,
		IncludeCanDM:                   1,
		IncludeCanMediaTag:             1,
		IncludeExtIsBlueVerified:       1,
		IncludeExtVerifiedType:         1,
		IncludeExtProfileImageShape:    1,
		SkipStatus:                     1,
		UserID:                         userID,
	}

	return buildTwitterRequest(authToken, ct0, common.FriendshipCreate, reqData)
}

func notificationRequest(authToken, ct0 string, userID int) (*http.Request, error) {
	reqData := model.Update{
		IncludeProfileInterstitialType: 1,
		IncludeBlocking:                1,
		IncludeBlockedBy:               1,
		IncludeFollowedBy:              1,
		IncludeWantRetweets:            1,
		IncludeMuteEdge:                1,
		IncludeCanDM:                   1,
		IncludeCanMediaTag:             1,
		IncludeExtIsBlueVerified:       1,
		IncludeExtVerifiedType:         1,
		IncludeExtProfileImageShape:    1,
		SkipStatus:                     1,
		Cursor:                         -1,
		UserID:                         userID,
		Device:                         true,
	}

	return buildTwitterRequest(authToken, ct0, common.FriendshipUpdate, reqData)
}

func (s *FollowService) doFollow(authToken, ct0 string, userID int) error {
	req, err := followRequest(authToken, ct0, userID)
	if err != nil {
		return fmt.Errorf("创建关注请求失败: %w", err)
	}

	return s.doRequest(req)
}

func (s *FollowService) enableNotification(authToken, ct0 string, userID int) error {
	req, err := notificationRequest(authToken, ct0, userID)
	if err != nil {
		return fmt.Errorf("创建通知请求失败: %w", err)
	}

	return s.doRequest(req)
}

func (s *FollowService) doRequest(req *http.Request) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var raw map[string]interface{}
	if err = json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("响应解析失败: %w", err)
	}
	if errs, ok := raw["errors"]; ok {
		return fmt.Errorf("接口返回错误: %v", errs)
	}
	return nil
}
