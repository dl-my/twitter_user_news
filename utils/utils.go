package utils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/config"
)

type Token struct {
	Auth string `json:"Auth"`
	Ct0  string `json:"Ct0"`
}

var (
	tokens []Token
	index  int
	mu     sync.RWMutex
)

var userMap = make(map[string]string)

// LoadTokens 从 JSON 文件加载 tokens
func LoadTokens() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(config.GlobalConfig.CookieFile)
	if err != nil {
		return fmt.Errorf("读取cookie.json文件失败: %w", err)
	}

	var tmp []Token
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}
	if len(tmp) == 0 {
		return fmt.Errorf("cookie 文件为空")
	}

	tokens = tmp
	index = 0
	return nil
}

// GetAuthAndCt0 返回下一个 token
func GetAuthAndCt0() (string, string) {
	mu.Lock()
	defer mu.Unlock()

	token := tokens[index]
	index = (index + 1) % len(tokens) // 轮转
	return token.Auth, token.Ct0
}

func UtcToShanghai(XTime string) (time.Time, error) {
	// 解析 UTC 时间
	t, err := time.Parse(common.XTimeParse, XTime)
	if err != nil {
		logs.Error("时间转换失败", zap.Any("err", err))
		return time.Time{}, err
	}
	return t, nil
}

func GetTimeStamp(t time.Time) int64 {
	loc, err := time.LoadLocation(common.TimeLocation)
	if err != nil {
		logs.Error("上号时区加载失败", zap.Any("err", err))
		return 0
	}
	// 转换到上海时区
	shanghaiTime := t.In(loc)

	// 返回时间戳（秒）
	return shanghaiTime.Unix()
}

func AddUser(username, userId string) {
	mu.Lock()
	defer mu.Unlock()
	userMap[username] = userId
}

func DelUser(username string) {
	mu.Lock()
	defer mu.Unlock()
	delete(userMap, username)
}

func GetUser(username string) string {
	mu.RLock()
	defer mu.RUnlock()
	return userMap[username]
}

// LoadUsers 从 JSON 文件加载
func LoadUsers() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(config.GlobalConfig.UserFile)
	if err != nil {
		return fmt.Errorf("读取user.json文件失败: %w", err)
	}
	err = json.Unmarshal(data, &userMap)
	if err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return nil
}

// SaveUsers 保存到 JSON 文件
func SaveUsers() error {
	mu.RLock()
	defer mu.RUnlock()

	f, err := os.Create(config.GlobalConfig.UserFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(userMap)
}
