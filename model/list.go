package model

type ListResponse struct {
	Data struct {
		List struct {
			TweetsTimeline struct {
				Timeline struct {
					Instructions []struct {
						Type    string  `json:"type"`
						Entries []Entry `json:"entries"`
					} `json:"Instructions"`
				} `json:"timeline"`
			} `json:"tweets_timeline"`
		} `json:"list"`
	} `json:"data"`
}

type ListVariables struct {
	ListID string `json:"listId"`
	Count  int    `json:"count"`
}
