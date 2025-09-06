package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/model"
	"twitter_user_news/utils"
)

var seenTweets = make(map[string]struct{})

func Search(userId string) {
	const maxRetries = 5
	retryCount := 0

	authToken, ct0 := utils.GetAuthAndCt0()

	for {
		err := Posts(userId, authToken, ct0)
		if err == nil {
			// 成功就直接退出
			return
		}

		retryCount++
		logs.Error(fmt.Sprintf("搜索失败 [userId: %s]，错误: %v，第 %d 次重试...", userId, err, retryCount))

		if retryCount >= maxRetries {
			// 超过 5 次，重新获取 AuthAndCt0
			logs.Warn("连续失败 5 次，切换 AuthAndCt0 ...")
			authToken, ct0 = utils.GetAuthAndCt0()
			retryCount = 0
		}

		time.Sleep(2 * time.Second) // 等待一会再试
	}
}

func generateUrl(userId string) string {
	// 用结构体定义搜索条件
	variablesStruct := model.PostsVariables{
		UserID:                                 userId,
		Count:                                  10,
		IncludePromotedContent:                 true,
		WithQuickPromoteEligibilityTweetFields: true,
		WithVoice:                              true,
	}
	featuresStruct := model.PostsFeatures{
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
	fieldTogglesStruct := model.PostsFieldToggles{
		WithArticlePlainText: false,
	}

	// 序列化成 JSON
	variablesJSON, _ := json.Marshal(variablesStruct)
	featuresJSON, _ := json.Marshal(featuresStruct)
	fieldTogglesJSON, _ := json.Marshal(fieldTogglesStruct)

	params := url.Values{}
	params.Set("variables", string(variablesJSON))
	params.Set("features", string(featuresJSON))
	params.Set("fieldToggles", string(fieldTogglesJSON))

	reqURL := common.UserTweetsAndRepliesUrl + "?" + params.Encode()
	return reqURL
}

func Posts(userId, authToken, ct0 string) error {
	reqURL := generateUrl(userId)

	// 创建 GET 请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("创建请求失败%v\n", err)
		return err
	}

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
	req.AddCookie(&http.Cookie{Name: "ct0", Value: ct0})

	u, _ := url.Parse(common.UserTweetsAndRepliesUrl)

	// 必要的请求头
	req.Header.Set("Authorization", common.Authorization)
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Transaction-Id", utils.XClientTransactionID("GET", u.Path))
	req.Header.Set("X-Csrf-Token", ct0)

	req.Header.Set("referer", "https://x.com")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-auth-type", "OAuth2Session")
	req.Header.Set("x-twitter-client-language", "en")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败%v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据失败%v\n", err)
		return err
	}

	fmt.Println(resp.StatusCode, string(body))

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Printf("响应解析失败: %v", err)
		return err
	}

	// 判断是否有 errors 字段
	if errs, ok := raw["errors"]; ok {
		log.Printf("接口返回错误: %v", errs)
		return fmt.Errorf("接口返回错误: %v", errs)
	}

	var result model.PostsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("响应解析失败%v,响应:[%s],状态:[%d]\n", err, string(body), resp.StatusCode)
		return fmt.Errorf("authToken:[%s],err:[%s]", authToken, err.Error())
	}

	for _, instruction := range result.Data.User.Result.Timeline.Timeline.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				// 输出评论
				for _, item := range entry.Content.Items {
					logComment(userId, item.Item.ItemContent.TweetResults.Result)
				}
				logComment(userId, entry.Content.ItemContent.TweetResults.Result)
			}
		}
	}

	return nil
}

func logComment(userId string, item model.Tweet) {
	var content string
	var mediaMap = make(map[string]string)
	if item.Legacy.UserIdStr != userId {
		return
	}
	if _, ok := seenTweets[item.RestId]; ok {
		return
	}
	t, err := utils.UtcToShanghai(item.Legacy.CreatedAt)
	if err != nil {
		return
	}
	publishTime := utils.GetTimeStamp(t)
	fetchTime := utils.GetTimeStamp(time.Now())
	if fetchTime-publishTime > 10 {
		return
	}
	if item.Legacy.RetweetedStatusResult != nil {
		content, mediaMap = getRetweet(item.Legacy.RetweetedStatusResult)
	} else {
		content, mediaMap = getTweet(item)
	}
	content = strings.ReplaceAll(content, "\n", "")
	contentZh, err := utils.Translate(content)
	if err != nil {
		logs.Error("翻译失败", zap.Any("err", err))
		return
	}
	posts := model.LogPosts{
		UserName:    "",
		UserId:      userId,
		RestId:      item.RestId,
		ContentEn:   content,
		ContentZh:   contentZh,
		PublishTime: publishTime,
		FetchTime:   fetchTime,
		Media:       mediaMap,
	}
	logs.InfoPosts(posts)
	seenTweets[item.RestId] = struct{}{}
}

func getTweet(tweet model.Tweet) (string, map[string]string) {
	mediaMap := getMedia(tweet.Legacy.Entities.Medias)
	if tweet.NoteTweet != nil {
		return tweet.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return tweet.Legacy.FullText, mediaMap
}

func getRetweet(retweet *model.RetweetedStatusResult) (string, map[string]string) {
	mediaMap := getMedia(retweet.Result.Legacy.Entities.Medias)
	if retweet.Result.NoteTweet != nil {
		return retweet.Result.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return retweet.Result.Legacy.FullText, mediaMap

}

func getMedia(medias []model.Media) map[string]string {
	mediaMap := make(map[string]string)
	for _, media := range medias {
		switch media.Type {
		case common.Photo:
			mediaMap[common.Photo] = media.MediaUrl
		case common.AnimatedGif, common.Video:
			for _, variant := range media.VideoInfo.Variants {
				if strings.Contains(variant.ContentType, common.Video) {
					mediaMap[media.Type] = variant.URL
					break
				}
			}
		}
	}
	return mediaMap
}
