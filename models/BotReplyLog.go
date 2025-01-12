package models

import (
	"time"

	"github.com/project-miko/miko/tools"
)

const (
	BotReplyStatusSuccess = 1
	BotReplyStatusFail    = 2
)

type BotReplyLog struct {
	Id             int64  `json:"id"`
	BotId          string `json:"bot_id"`
	MediaUrl       string `json:"media_url"`
	AuthorId       string `json:"author_id"`
	RepliedTweetId string `json:"replied_tweet_id"`
	Status         int    `json:"status"`
	ErrorMsg       string `json:"error_msg"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

func (*BotReplyLog) TableName() string {
	return "bot_reply_log"
}

func (m *BotReplyLog) Save() error {
	return GetDbInst().Save(m).Error
}

func (m *BotReplyLog) Update() error {
	m.UpdatedAt = tools.GetMillisecond(time.Now())
	return m.Save()
}
