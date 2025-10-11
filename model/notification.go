package model

import (
	"net/url"
	"reflect"
	"strconv"
)

type NotificationResponse struct {
	GlobalObjects struct {
		Tweets map[string]NotificationTweet `json:"tweets"`
	} `json:"globalObjects"`
}

type NotificationTweet struct {
	CreatedAt        string `json:"created_at"`
	ID               string `json:"id_str"`
	FullText         string `json:"full_text"`
	UserID           string `json:"user_id_str"`
	ExtendedEntities struct {
		Medias []Media `json:"media"`
	} `json:"extended_entities"`
}

type NotificationParams struct {
	IncludeProfileInterstitialType   int    `url:"include_profile_interstitial_type"`
	IncludeBlocking                  int    `url:"include_blocking"`
	IncludeBlockedBy                 int    `url:"include_blocked_by"`
	IncludeFollowedBy                int    `url:"include_followed_by"`
	IncludeWantRetweets              int    `url:"include_want_retweets"`
	IncludeMuteEdge                  int    `url:"include_mute_edge"`
	IncludeCanDM                     int    `url:"include_can_dm"`
	IncludeCanMediaTag               int    `url:"include_can_media_tag"`
	IncludeExtIsBlueVerified         int    `url:"include_ext_is_blue_verified"`
	IncludeExtVerifiedType           int    `url:"include_ext_verified_type"`
	IncludeExtProfileImageShape      int    `url:"include_ext_profile_image_shape"`
	SkipStatus                       int    `url:"skip_status"`
	CardsPlatform                    string `url:"cards_platform"`
	IncludeCards                     int    `url:"include_cards"`
	IncludeExtAltText                bool   `url:"include_ext_alt_text"`
	IncludeExtLimitedActionResults   bool   `url:"include_ext_limited_action_results"`
	IncludeQuoteCount                bool   `url:"include_quote_count"`
	IncludeReplyCount                int    `url:"include_reply_count"`
	TweetMode                        string `url:"tweet_mode"`
	IncludeExtViews                  bool   `url:"include_ext_views"`
	IncludeEntities                  bool   `url:"include_entities"`
	IncludeUserEntities              bool   `url:"include_user_entities"`
	IncludeExtMediaColor             bool   `url:"include_ext_media_color"`
	IncludeExtMediaAvailability      bool   `url:"include_ext_media_availability"`
	IncludeExtSensitiveMediaWarning  bool   `url:"include_ext_sensitive_media_warning"`
	IncludeExtTrustedFriendsMetadata bool   `url:"include_ext_trusted_friends_metadata"`
	SendErrorCodes                   bool   `url:"send_error_codes"`
	SimpleQuotedTweet                bool   `url:"simple_quoted_tweet"`
	Count                            int    `url:"count"`
	Ext                              string `url:"ext"`
}

func (p NotificationParams) Encode() string {
	values := url.Values{}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("url")
		if tag == "" {
			continue
		}

		fieldValue := v.Field(i)
		var str string

		switch fieldValue.Kind() {
		case reflect.Int:
			str = strconv.Itoa(int(fieldValue.Int()))
		case reflect.String:
			str = fieldValue.String()
		case reflect.Bool:
			str = strconv.FormatBool(fieldValue.Bool())
		default:
			continue
		}

		values.Set(tag, str)
	}
	return values.Encode()
}
