package model

type HomeResponse struct {
	Data struct {
		Home struct {
			HomeTimeLineUrt struct {
				Instructions []struct {
					Type    string  `json:"type"`
					Entries []Entry `json:"entries"`
				} `json:"instructions"`
			} `json:"home_timeline_urt"`
		} `json:"home"`
	} `json:"data"`
}

type HomeVariables struct {
	Count                  int      `json:"count"`
	IncludePromotedContent bool     `json:"includePromotedContent"`
	LatestControlAvailable bool     `json:"latestControlAvailable"`
	RequestContext         string   `json:"requestContext"`
	SeenTweetIds           []string `json:"seenTweetIds"`
}
