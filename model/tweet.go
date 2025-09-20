package model

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
	Core      Core        `json:"core"`
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

type Core struct {
	UserResults struct {
		Result struct {
			RestId string `json:"rest_id"`
			Core   struct {
				ScreenName string `json:"screen_name"`
			} `json:"core"`
		} `json:"result"`
	} `json:"user_results"`
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
