package model

type LogPosts struct {
	UserName    string              `json:"username"`
	UserId      string              `json:"user_id"`
	RestId      string              `json:"rest_id"`
	ContentEn   string              `json:"content_en"`
	ContentZh   string              `json:"content_zh"`
	PublishTime int64               `json:"publish_time"`
	FetchTime   int64               `json:"fetch_time"`
	Media       map[string][]string `json:"media"`
}
