package model

type SearchResponse struct {
	Data struct {
		SearchByRawQuery struct {
			SearchTimeline struct {
				Timeline struct {
					Instructions []struct {
						Type    string  `json:"type"`
						Entries []Entry `json:"entries"`
					} `json:"Instructions"`
				} `json:"timeline"`
			} `json:"search_timeline"`
		} `json:"search_by_raw_query"`
	} `json:"data"`
}

type SearchVariables struct {
	RawQuery              string `json:"rawQuery"`
	Count                 int    `json:"count"`
	QuerySource           string `json:"querySource"`
	Product               string `json:"product"`
	WithGrokTranslatedBio bool   `json:"withGrokTranslatedBio"`
}
