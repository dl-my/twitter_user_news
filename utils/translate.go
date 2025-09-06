package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/config"
)

type Constant struct {
	TagText string `json:"tgtText"`
}

// 生成权限字符串
func generateAuthStr(params map[string]string) string {
	// 按 key 排序
	keys := make([]string, 0, len(params)+1)
	for k := range params {
		keys = append(keys, k)
	}
	keys = append(keys, "apikey") // 加入 apikey
	sort.Strings(keys)

	// 拼接参数字符串
	var paramPairs []string
	for _, k := range keys {
		v := ""
		if k == "apikey" {
			v = config.GlobalConfig.Document.AppKey
		} else {
			v = params[k]
		}
		paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", k, v))
	}
	paramStr := strings.Join(paramPairs, "&")

	// MD5
	hash := md5.Sum([]byte(paramStr))
	return hex.EncodeToString(hash[:])
}

func Translate(text string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	data := map[string]string{
		"from":      "en",
		"to":        "zh",
		"appId":     config.GlobalConfig.Document.AppId,
		"timestamp": timestamp,
		"srcText":   text,
	}

	authStr := generateAuthStr(data)
	data["authStr"] = authStr

	formData := url.Values{}
	for k, v := range data {
		formData.Set(k, v)
	}

	resp, err := http.PostForm(common.TranslateUrl, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := Constant{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return "", err
	}
	return content.TagText, nil
}
