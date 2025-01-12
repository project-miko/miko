package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-miko/miko/tools"
)

const (
	TweetLibQueryModeSorted = 1 // sorted
	TweetLibQueryModeRandom = 2 // random
)

var (
	AllowTweetLibQueryMode = map[int64]struct{}{
		TweetLibQueryModeSorted: {},
		TweetLibQueryModeRandom: {},
	}
)

type TweetLib struct {
	Id           uint64 `json:"id"`
	Account      string `json:"account"`
	Avatar       string `json:"avatar"`
	Category     string `json:"category"`
	SourceUrl    string `json:"source_url"`
	LikeCount    int64  `json:"like_count"`
	RetweetCount int64  `json:"retweet_count"`
	Content      string `json:"content"`
	PublishAt    int64  `json:"publish_at"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

func (m *TweetLib) TableName() string {
	return "tweet_lib"
}

func (m *TweetLib) Save() error {
	return GetDbInst().Save(m).Error
}

func (m *TweetLib) Update() error {
	m.UpdatedAt = tools.GetMillisecond(time.Now())
	return m.Save()
}

func (m *TweetLib) Del() error {
	return GetDbInst().Delete(m).Error
}

func GeTweetById(id int64) (*TweetLib, error) {
	result := new(TweetLib)
	err := GetDbInst().Where("id=?", id).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return result, err
}

func GeTweetBySourceUrl(sourceUrl string) (*TweetLib, error) {
	result := new(TweetLib)
	err := GetDbInst().Where("source_url=?", sourceUrl).Find(result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return result, err
}

func GetTweetCategoryList() ([]string, error) {
	type result struct {
		Category string
	}

	_results := make([]result, 0)
	sql := "select category from tweet_lib group by category"
	err := GetDbInst().Model(TweetLib{}).Raw(sql).Scan(&_results).Error
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, v := range _results {
		results = append(results, v.Category)
	}

	return results, nil
}

func GetTweetListByPage(category, keywords string, likeCount, retweetCount, page, limit int64) (int64, []*TweetLib, error) {
	var amount int64
	results := make([]*TweetLib, 0)
	db := GetDbInst()

	if len(category) > 0 {
		db = db.Where("category = ?", category)
	}
	if len(keywords) > 0 {
		db = db.Where("content like ?", "%"+keywords+"%")
	}
	if likeCount >= 0 {
		db = db.Where("like_count >= ?", likeCount)
	}
	if retweetCount >= 0 {
		db = db.Where("retweet_count >= ?", retweetCount)
	}

	err := db.Model(TweetLib{}).Count(&amount).Error
	if err != nil {
		return 0, nil, err
	}
	if amount == 0 {
		return 0, results, nil
	}

	offset := (page - 1) * limit
	err = db.Offset(offset).Limit(limit).Order("like_count desc, retweet_count desc").Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, results, nil
	}

	return amount, results, err
}
