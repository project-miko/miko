package data

import "github.com/shopspring/decimal"

type TwMonitoringDTO struct {
	FollowCount        int              `json:"follow_count"`
	FollowChange       int              `json:"follow_change"`
	LikeCount          int              `json:"like_count"`
	ReplyCount         int              `json:"reply_count"`
	ReTweetCount       int              `json:"re_tweet_count"`
	IncreaseFollowRate decimal.Decimal  `json:"increase_follow_rate"`
	StatisticAt        int64            `json:"statistic_at"`
	UserDataList       []*TwUserDataDTO `json:"user_data_list"`
}

type TwUserDataDTO struct {
	UserId             string          `json:"user_id"`
	UserName           string          `json:"user_name"`
	UserAccount        string          `json:"user_account"`
	ProfileImageUrl    string          `json:"profile_image_url"`
	FollowCount        int             `json:"follow_count"`
	FollowChange       int             `json:"follow_change"`
	LikeCount          int             `json:"like_count"`
	ReplyCount         int             `json:"reply_count"`
	RetweetCount       int             `json:"retweet_count"`
	TweetCount         int             `json:"tweet_count"`
	LikeFollowRate     decimal.Decimal `json:"like_follow_rate"`
	IncreaseFollowRate decimal.Decimal `json:"increase_follow_rate"`
}
