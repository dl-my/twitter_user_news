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
	"sync"
	"time"
	"twitter_user_news/common"
	"twitter_user_news/common/logs"
	"twitter_user_news/model"
	"twitter_user_news/utils"
)

// SearchTwitterService 封装Twitter相关服务
type SearchTwitterService struct {
	seenTweets sync.Map // 使用sync.Map替代全局变量，支持并发安全
	client     *http.Client
}

// NewSearchTwitterService 创建新的Twitter服务实例
func NewSearchTwitterService() *SearchTwitterService {
	return &SearchTwitterService{
		client: &http.Client{
			Timeout: 30 * time.Second, // 设置超时时间
		},
	}
}

func (s *SearchTwitterService) Search() {
	retryCount := 0
	authToken, ct0 := utils.GetAuthAndCt0()
	for {
		err := s.fetchPosts(authToken, ct0)
		if err == nil {
			fmt.Println("获取成功")
			return
		}
		retryCount++
		logs.Error("搜索失败",
			zap.Error(err),
			zap.Int("retryCount", retryCount))

		if retryCount >= common.MaxRetries {
			// 超过 5 次，重新获取 AuthAndCt0
			logs.Warn("连续失败达到最大重试次数，切换认证信息",
				zap.String("authToken", authToken))
			authToken, ct0 = utils.GetAuthAndCt0()
			retryCount = 0
		}

		time.Sleep(2 * time.Second) // 等待一会再试
	}
}

func (s *SearchTwitterService) generateUrl() string {
	// 用结构体定义搜索条件
	variablesStruct := model.SearchVariables{
		RawQuery:              "filter:replies filter:follows",
		Count:                 20,
		QuerySource:           "typed_query",
		Product:               "Latest",
		WithGrokTranslatedBio: false,
	}
	featuresStruct := s.getDefaultFeatures()

	// 序列化成 JSON
	variablesJSON, _ := json.Marshal(variablesStruct)
	featuresJSON, _ := json.Marshal(featuresStruct)

	params := url.Values{}
	params.Set("variables", string(variablesJSON))
	params.Set("features", string(featuresJSON))

	return common.SearchTimeline + "?" + params.Encode()
}

func (s *SearchTwitterService) getDefaultFeatures() model.PostsFeatures {
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

func (s *SearchTwitterService) createRequest(reqURL, authToken, ct0 string) (*http.Request, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置cookies
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
	req.AddCookie(&http.Cookie{Name: "ct0", Value: ct0})

	// 设置请求头
	u, _ := url.Parse(common.SearchTimeline)
	req.Header.Set("Authorization", common.Authorization)
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Transaction-Id", utils.XClientTransactionID("GET", u.Path))
	req.Header.Set("X-Csrf-Token", ct0)
	req.Header.Set("Referer", "https://x.com")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Client-Language", "en")

	return req, nil
}

func (s *SearchTwitterService) fetchPosts(authToken, ct0 string) error {
	reqURL := s.generateUrl()

	req, err := s.createRequest(reqURL, authToken, ct0)
	if err != nil {
		log.Printf("创建请求失败%v\n", err)
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("发送请求失败%v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
		return fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据失败%v\n", err)
		return err
	}

	//fmt.Println(resp.StatusCode, string(body))

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Printf("响应解析失败: %v", err)
		return err
	}

	// 判断是否有 errors 字段
	//if errs, ok := raw["errors"]; ok {
	//	log.Printf("接口返回错误: %v", errs)
	//	return fmt.Errorf("接口返回错误: %v", errs)
	//}

	var result model.SearchResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("响应解析失败%v,响应:[%s],状态:[%d]\n", err, string(body), resp.StatusCode)
		return fmt.Errorf("authToken:[%s],err:[%s]", authToken, err.Error())
	}
	//fmt.Println(result)

	return s.processTimeline(result)
}

func (s *SearchTwitterService) processTimeline(result model.SearchResponse) error {
	for _, instruction := range result.Data.SearchByRawQuery.SearchTimeline.Timeline.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				// 输出评论
				for _, item := range entry.Content.Items {
					s.processTweetOrComment(item.Item.ItemContent.TweetResults.Result)
				}
				// 输出推文
				s.processTweetOrComment(entry.Content.ItemContent.TweetResults.Result)
			}
		}
	}
	return nil
}

func (s *SearchTwitterService) processTweetOrComment(tweet model.Tweet) {
	// 检查是否已处理过
	if _, exists := s.seenTweets.Load(tweet.RestId); exists {
		return
	}
	// 检查是否超时
	t, err := utils.UtcToShanghai(tweet.Legacy.CreatedAt)
	if err != nil {
		return
	}
	publishTime := utils.GetTimeStamp(t)
	fetchTime := utils.GetTimeStamp(time.Now())
	if fetchTime-publishTime > 60*60*1 {
		return
	}
	// 提取推文内容和媒体
	content, mediaMap := s.extractTweetContent(tweet)
	if content == "" {
		return
	}
	// 去除转义符
	content = strings.ReplaceAll(content, "\n", "")
	// 翻译为中文
	contentZh, err := utils.Translate(content)
	if err != nil {
		logs.Error("翻译失败", zap.Any("err", err))
		return
	}
	// 构建日志对象
	posts := model.LogPosts{
		UserName:    tweet.Core.UserResults.Result.Core.ScreenName,
		UserId:      tweet.Core.UserResults.Result.RestId,
		RestId:      tweet.RestId,
		ContentEn:   content,
		ContentZh:   contentZh,
		PublishTime: publishTime,
		FetchTime:   fetchTime,
		Media:       mediaMap,
	}
	logs.InfoPosts(posts)
	s.seenTweets.Store(tweet.RestId, struct{}{})
}

func (s *SearchTwitterService) extractTweetContent(tweet model.Tweet) (string, map[string]string) {
	if tweet.Legacy.RetweetedStatusResult != nil {
		return s.extractRetweetContent(tweet.Legacy.RetweetedStatusResult)
	}
	return s.extractOriginalContent(tweet)
}

func (s *SearchTwitterService) extractOriginalContent(tweet model.Tweet) (string, map[string]string) {
	mediaMap := s.getMedia(tweet.Legacy.Entities.Medias)
	if tweet.NoteTweet != nil {
		return tweet.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return tweet.Legacy.FullText, mediaMap
}

func (s *SearchTwitterService) extractRetweetContent(retweet *model.RetweetedStatusResult) (string, map[string]string) {
	mediaMap := s.getMedia(retweet.Result.Legacy.Entities.Medias)
	if retweet.Result.NoteTweet != nil {
		return retweet.Result.NoteTweet.NoteTweetResults.Result.Text, mediaMap
	}
	return retweet.Result.Legacy.FullText, mediaMap

}

func (s *SearchTwitterService) getMedia(medias []model.Media) map[string]string {
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
