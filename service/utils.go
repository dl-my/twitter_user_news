package service

import (
	"fmt"
	"net/http"
	"strings"
	"twitter_user_news/common"
	"twitter_user_news/model"
)

func getDefaultFeatures() model.PostsFeatures {
	return model.PostsFeatures{
		RwebVideoScreenEnabled: false,
		PaymentsEnabled:        false,
		RwebXchatEnabled:       false,
		ProfileLabelImprovementsPcfLabelInPostEnabled:                  true,
		RwebTipjarConsumptionEnabled:                                   true,
		VerifiedPhoneLabelEnabled:                                      false,
		CreatorSubscriptionsTweetPreviewAPIEnabled:                     true,
		ResponsiveWebGraphqlTimelineNavigationEnabled:                  true,
		ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled:      false,
		PremiumContentAPIReadEnabled:                                   false,
		CommunitiesWebEnableTweetCommunityResultsFetch:                 true,
		C9STweetAnatomyModeratorBadgeEnabled:                           true,
		ResponsiveWebGrokAnalyzeButtonFetchTrendsEnabled:               false,
		ResponsiveWebGrokAnalyzePostFollowupsEnabled:                   true,
		ResponsiveWebJetfuelFrame:                                      true,
		ResponsiveWebGrokShareAttachmentEnabled:                        true,
		ArticlesPreviewEnabled:                                         true,
		ResponsiveWebEditTweetAPIEnabled:                               true,
		GraphqlIsTranslatableRwebTweetIsTranslatableEnabled:            true,
		ViewCountsEverywhereAPIEnabled:                                 true,
		LongformNotetweetsConsumptionEnabled:                           true,
		ResponsiveWebTwitterArticleTweetConsumptionEnabled:             true,
		TweetAwardsWebTippingEnabled:                                   false,
		ResponsiveWebGrokShowGrokTranslatedPost:                        false,
		ResponsiveWebGrokAnalysisButtonFromBackend:                     true,
		CreatorSubscriptionsQuoteTweetPreviewEnabled:                   false,
		FreedomOfSpeechNotReachFetchEnabled:                            true,
		StandardizedNudgesMisinfo:                                      true,
		TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled: true,
		LongformNotetweetsRichTextReadEnabled:                          true,
		LongformNotetweetsInlineMediaEnabled:                           true,
		ResponsiveWebGrokImageAnnotationEnabled:                        true,
		ResponsiveWebGrokImagineAnnotationEnabled:                      true,
		ResponsiveWebGrokCommunityNoteAutoTranslationIsEnabled:         false,
		ResponsiveWebEnhanceCardsEnabled:                               false,
	}
}

func createRequest(reqURL, authToken, ct0 string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置cookies
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
	req.AddCookie(&http.Cookie{Name: "ct0", Value: ct0})

	// 设置请求头
	//u, _ := url.Parse(reqURL)
	//req.Header.Set("X-Client-Transaction-Id", utils.XClientTransactionID(http.MethodGet, u.Path))
	req.Header.Set("Authorization", common.Authorization)
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Csrf-Token", ct0)
	req.Header.Set("Referer", "https://x.com")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Client-Language", "en")

	return req, nil
}

func extractTweetContent(tweet model.Tweet) (string, map[string][]string) {
	if tweet.Legacy.RetweetedStatusResult != nil {
		return extractRetweetContent(tweet.Legacy.RetweetedStatusResult)
	}
	return extractOriginalContent(tweet)
}

func extractOriginalContent(tweet model.Tweet) (string, map[string][]string) {
	mediaMap := getMedia(tweet.Legacy.Entities.Medias)
	if tweet.NoteTweet != nil {
		return tweet.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return tweet.Legacy.FullText, mediaMap
}

func extractRetweetContent(retweet *model.RetweetedStatusResult) (string, map[string][]string) {
	mediaMap := getMedia(retweet.Result.Legacy.Entities.Medias)
	if retweet.Result.NoteTweet != nil {
		return retweet.Result.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return retweet.Result.Legacy.FullText, mediaMap

}

func getMedia(medias []model.Media) map[string][]string {
	mediaMap := make(map[string][]string)
	for _, media := range medias {
		switch media.Type {
		case common.Photo:
			mediaMap[common.Photo] = append(mediaMap[common.Photo], media.MediaUrl)
		case common.AnimatedGif, common.Video:
			for _, variant := range media.VideoInfo.Variants {
				if strings.Contains(variant.ContentType, common.Video) {
					mediaMap[media.Type] = append(mediaMap[media.Type], variant.URL)
					break
				}
			}
		}
	}
	return mediaMap
}

func getNotificationParams() model.NotificationParams {
	return model.NotificationParams{
		IncludeProfileInterstitialType:   1,
		IncludeBlocking:                  1,
		IncludeBlockedBy:                 1,
		IncludeFollowedBy:                1,
		IncludeWantRetweets:              1,
		IncludeMuteEdge:                  1,
		IncludeCanDM:                     1,
		IncludeCanMediaTag:               1,
		IncludeExtIsBlueVerified:         1,
		IncludeExtVerifiedType:           1,
		IncludeExtProfileImageShape:      1,
		SkipStatus:                       1,
		CardsPlatform:                    "Web-12",
		IncludeCards:                     1,
		IncludeExtAltText:                true,
		IncludeExtLimitedActionResults:   true,
		IncludeQuoteCount:                true,
		IncludeReplyCount:                1,
		TweetMode:                        "extended",
		IncludeExtViews:                  true,
		IncludeEntities:                  true,
		IncludeUserEntities:              true,
		IncludeExtMediaColor:             true,
		IncludeExtMediaAvailability:      true,
		IncludeExtSensitiveMediaWarning:  true,
		IncludeExtTrustedFriendsMetadata: true,
		SendErrorCodes:                   true,
		SimpleQuotedTweet:                true,
		Count:                            20,
		Ext: "mediaStats,highlightedLabel,parodyCommentaryFanLabel,voiceInfo," +
			"birdwatchPivot,superFollowMetadata,unmentionInfo,editControl,article",
	}
}
