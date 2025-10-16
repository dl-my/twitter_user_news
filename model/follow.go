package model

type Create struct {
	IncludeProfileInterstitialType int `form:"include_profile_interstitial_type"`
	IncludeBlocking                int `form:"include_blocking"`
	IncludeBlockedBy               int `form:"include_blocked_by"`
	IncludeFollowedBy              int `form:"include_followed_by"`
	IncludeWantRetweets            int `form:"include_want_retweets"`
	IncludeMuteEdge                int `form:"include_mute_edge"`
	IncludeCanDM                   int `form:"include_can_dm"`
	IncludeCanMediaTag             int `form:"include_can_media_tag"`
	IncludeExtIsBlueVerified       int `form:"include_ext_is_blue_verified"`
	IncludeExtVerifiedType         int `form:"include_ext_verified_type"`
	IncludeExtProfileImageShape    int `form:"include_ext_profile_image_shape"`
	SkipStatus                     int `form:"skip_status"`
	UserID                         int `form:"user_id"`
}

type Update struct {
	IncludeProfileInterstitialType int  `form:"include_profile_interstitial_type"`
	IncludeBlocking                int  `form:"include_blocking"`
	IncludeBlockedBy               int  `form:"include_blocked_by"`
	IncludeFollowedBy              int  `form:"include_followed_by"`
	IncludeWantRetweets            int  `form:"include_want_retweets"`
	IncludeMuteEdge                int  `form:"include_mute_edge"`
	IncludeCanDM                   int  `form:"include_can_dm"`
	IncludeCanMediaTag             int  `form:"include_can_media_tag"`
	IncludeExtIsBlueVerified       int  `form:"include_ext_is_blue_verified"`
	IncludeExtVerifiedType         int  `form:"include_ext_verified_type"`
	IncludeExtProfileImageShape    int  `form:"include_ext_profile_image_shape"`
	SkipStatus                     int  `form:"skip_status"`
	Cursor                         int  `form:"cursor"`
	UserID                         int  `form:"id"`
	Device                         bool `form:"device"`
}
