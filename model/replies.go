package model

type PostsResponse struct {
	Data struct {
		User struct {
			Result struct {
				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Type    string  `json:"type"`
							Entries []Entry `json:"entries"`
						} `json:"Instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type Entry struct {
	SortIndex string `json:"sortIndex"`
	Content   struct {
		Items []struct {
			Dispensable bool `json:"dispensable"`
			Item        struct {
				ItemContent struct {
					TweetResults struct {
						Result Tweet `json:"result"`
					} `json:"tweet_results"`
				} `json:"itemContent"`
			} `json:"item"`
		} `json:"items"`
		ItemContent struct {
			TweetResults struct {
				Result Tweet `json:"result"`
			} `json:"tweet_results"`
		} `json:"itemContent"`
	} `json:"content"`
}

type Tweet struct {
	RestId    string      `json:"rest_id"`
	NoteTweet *NoteTweet  `json:"note_tweet"`
	Legacy    TweetLegacy `json:"legacy"`
}

type ReTweet struct {
	NoteTweet *NoteTweet `json:"note_tweet"`
	Legacy    struct {
		FullText string `json:"full_text"`
		Entities struct {
			Medias []Media `json:"media"`
		} `json:"entities"`
	} `json:"legacy"`
}
type RetweetedStatusResult struct {
	Result ReTweet `json:"result"`
}

type TweetLegacy struct {
	CreatedAt     string `json:"created_at"`
	FullText      string `json:"full_text"`
	UserIdStr     string `json:"user_id_str"`
	IsQuoteStatus bool   `json:"is_quote_status"`
	Entities      struct {
		Medias []Media `json:"media"`
	} `json:"entities"`
	RetweetedStatusResult *RetweetedStatusResult `json:"retweeted_status_result"`
}

type Media struct {
	Type      string `json:"type"`
	VideoInfo struct {
		Variants []struct {
			ContentType string `json:"content_type"`
			URL         string `json:"url"`
		} `json:"variants"`
	} `json:"video_info"`
	MediaUrl string `json:"media_url_https"`
}

type NoteTweet struct {
	NoteTweetResults struct {
		Result struct {
			Text string `json:"text"`
		} `json:"result"`
	} `json:"note_tweet_results"`
}

type LogPosts struct {
	UserName    string            `json:"username"`
	UserId      string            `json:"user_id"`
	RestId      string            `json:"rest_id"`
	ContentEn   string            `json:"content_en"`
	ContentZh   string            `json:"content_zh"`
	PublishTime int64             `json:"publish_time"`
	FetchTime   int64             `json:"fetch_time"`
	Media       map[string]string `json:"media"`
}
