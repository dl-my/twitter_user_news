package model

type PostsVariables struct {
	UserID                                 string `json:"userId"`
	Count                                  int    `json:"count"`
	IncludePromotedContent                 bool   `json:"includePromotedContent"`
	WithQuickPromoteEligibilityTweetFields bool   `json:"withQuickPromoteEligibilityTweetFields"`
	WithVoice                              bool   `json:"withVoice"`
}

type PostsFeatures struct {
	RwebVideoScreenEnabled                                         bool `json:"rweb_video_screen_enabled"`
	PaymentsEnabled                                                bool `json:"payments_enabled"`
	RwebXchatEnabled                                               bool `json:"rweb_xchat_enabled"`
	ProfileLabelImprovementsPcfLabelInPostEnabled                  bool `json:"profile_label_improvements_pcf_label_in_post_enabled"`
	RwebTipjarConsumptionEnabled                                   bool `json:"rweb_tipjar_consumption_enabled"`
	VerifiedPhoneLabelEnabled                                      bool `json:"verified_phone_label_enabled"`
	CreatorSubscriptionsTweetPreviewAPIEnabled                     bool `json:"creator_subscriptions_tweet_preview_api_enabled"`
	ResponsiveWebGraphqlTimelineNavigationEnabled                  bool `json:"responsive_web_graphql_timeline_navigation_enabled"`
	ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled      bool `json:"responsive_web_graphql_skip_user_profile_image_extensions_enabled"`
	PremiumContentAPIReadEnabled                                   bool `json:"premium_content_api_read_enabled"`
	CommunitiesWebEnableTweetCommunityResultsFetch                 bool `json:"communities_web_enable_tweet_community_results_fetch"`
	C9STweetAnatomyModeratorBadgeEnabled                           bool `json:"c9s_tweet_anatomy_moderator_badge_enabled"`
	ResponsiveWebGrokAnalyzeButtonFetchTrendsEnabled               bool `json:"responsive_web_grok_analyze_button_fetch_trends_enabled"`
	ResponsiveWebGrokAnalyzePostFollowupsEnabled                   bool `json:"responsive_web_grok_analyze_post_followups_enabled"`
	ResponsiveWebJetfuelFrame                                      bool `json:"responsive_web_jetfuel_frame"`
	ResponsiveWebGrokShareAttachmentEnabled                        bool `json:"responsive_web_grok_share_attachment_enabled"`
	ArticlesPreviewEnabled                                         bool `json:"articles_preview_enabled"`
	ResponsiveWebEditTweetAPIEnabled                               bool `json:"responsive_web_edit_tweet_api_enabled"`
	GraphqlIsTranslatableRwebTweetIsTranslatableEnabled            bool `json:"graphql_is_translatable_rweb_tweet_is_translatable_enabled"`
	ViewCountsEverywhereAPIEnabled                                 bool `json:"view_counts_everywhere_api_enabled"`
	LongformNotetweetsConsumptionEnabled                           bool `json:"longform_notetweets_consumption_enabled"`
	ResponsiveWebTwitterArticleTweetConsumptionEnabled             bool `json:"responsive_web_twitter_article_tweet_consumption_enabled"`
	TweetAwardsWebTippingEnabled                                   bool `json:"tweet_awards_web_tipping_enabled"`
	ResponsiveWebGrokShowGrokTranslatedPost                        bool `json:"responsive_web_grok_show_grok_translated_post"`
	ResponsiveWebGrokAnalysisButtonFromBackend                     bool `json:"responsive_web_grok_analysis_button_from_backend"`
	CreatorSubscriptionsQuoteTweetPreviewEnabled                   bool `json:"creator_subscriptions_quote_tweet_preview_enabled"`
	FreedomOfSpeechNotReachFetchEnabled                            bool `json:"freedom_of_speech_not_reach_fetch_enabled"`
	StandardizedNudgesMisinfo                                      bool `json:"standardized_nudges_misinfo"`
	TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled bool `json:"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled"`
	LongformNotetweetsRichTextReadEnabled                          bool `json:"longform_notetweets_rich_text_read_enabled"`
	LongformNotetweetsInlineMediaEnabled                           bool `json:"longform_notetweets_inline_media_enabled"`
	ResponsiveWebGrokImageAnnotationEnabled                        bool `json:"responsive_web_grok_image_annotation_enabled"`
	ResponsiveWebGrokImagineAnnotationEnabled                      bool `json:"responsive_web_grok_imagine_annotation_enabled"`
	ResponsiveWebGrokCommunityNoteAutoTranslationIsEnabled         bool `json:"responsive_web_grok_community_note_auto_translation_is_enabled"`
	ResponsiveWebEnhanceCardsEnabled                               bool `json:"responsive_web_enhance_cards_enabled"`
}

type PostsFieldToggles struct {
	WithArticlePlainText bool `json:"withArticlePlainText"`
}
