package data

type WithPaging struct {
	List   interface{} `json:"list"`
	Paging Paging      `json:"paging"`
}

type Paging struct {
	Amount int64 `json:"amount"`
	Page   int64 `json:"page"`
	Limit  int64 `json:"limit"`
}

type BasePage struct {
	Page  int64 `json:"page,omitempty"`
	Limit int64 `json:"limit,omitempty"`
}

type TwMetricInfo struct {
	StatisticAt int64           `json:"statistic_at"`
	Items       []*TwMetricItem `json:"items"`
}

type TwMetricItem struct {
	UserId      string `json:"user_id"`
	FollowCount int    `json:"follow_count"`
}

type TwUserRateLimitReq struct {
	UserId string `json:"user_id" binding:"min=1"`
}

type TwUploadMediaStatusReq struct {
	UserId   string   `json:"user_id" binding:"min=1"`
	MediaIds []string `json:"media_ids" binding:"required,dive,numeric"`
}

type CreateTweetReq struct {
	UserId string             `json:"user_id" binding:"min=1"`
	Tweets []*CreateTweetItem `json:"tweets" binding:"required,dive,required"`
}

type CreateTweetItem struct {
	SortId   string   `json:"sort_id" binding:"min=1"`
	Text     string   `json:"text" binding:"required_without=MediaIds"`
	MediaIds []string `json:"media_ids" binding:"required_without=Text,dive,numeric"`
}

type TwUploadMediaReq struct {
	UserId      string                 `json:"user_id" binding:"min=1"`
	UploadFiles []*TwUploadFileReqItem `json:"upload_files" binding:"required,dive,required"`
}

type TwUploadFileReqItem struct {
	Id       string `json:"id" binding:"min=1"`
	MediaUrl string `json:"media_url" binding:"required,url"`
}

type TwUploadMediaResp struct {
	MediaItems []*TwMediaRespItem `json:"media_items"`
}

type TwMediaRespItem struct {
	Id      string `json:"id"`
	MediaId string `json:"media_id"`
	Status  string `json:"status"`
	ErrMsg  string `json:"err_msg"`
}

type TwAddTweetScheduleReq struct {
	UserId          string                       `json:"user_id" binding:"min=1"`
	TwScheduleLibId *int64                       `json:"tw_schedule_lib_id,omitempty" binding:"omitempty,min=1"`
	LoopUnit        int                          `json:"loop_unit" binding:"min=1,max=4"`
	LoopCount       int                          `json:"loop_count" binding:"min=0,max=10"`
	WeekDay         int                          `json:"week_day" binding:"min=1,max=7"`
	Hour            int                          `json:"hour" binding:"min=0,max=23"`
	Minute          int                          `json:"minute" binding:"min=0,max=59"`
	ThreadList      []*TwAddTweetScheduleReqItem `json:"thread_list" binding:"required_without=TwScheduleLibId,dive"`
}

type TwAddTweetScheduleReqItem struct {
	SortId    string   `json:"sort_id" binding:"min=1"`
	Text      string   `json:"text" binding:"required_without=MediaUrls"`
	MediaUrls []string `json:"media_urls" binding:"required_without=Text,dive,url"`
}

type TwUpdateTweetScheduleReq struct {
	ScheduleId int64 `json:"schedule_id" binding:"min=1"`
	LoopUnit   *int  `json:"loop_unit,omitempty" binding:"omitempty,min=1,max=4"`
	LoopCount  *int  `json:"loop_count,omitempty" binding:"omitempty,min=0,max=10"`
	WeekDay    *int  `json:"week_day,omitempty" binding:"omitempty,min=1,max=7"`
	Hour       *int  `json:"hour,omitempty" binding:"omitempty,min=0,max=23"`
	Minute     *int  `json:"minute,omitempty" binding:"omitempty,min=0,max=59"`
}

type TwDelTweetScheduleReq struct {
	ScheduleId int64 `json:"schedule_id" binding:"min=1"`
}

type TwGetTweetScheduleListReq struct {
	UserId string `json:"user_id" binding:"min=1"`
	*BasePage
}

type TwGetTweetScheduleListItemResp struct {
	Id          int64                         `json:"id"`
	UserId      string                        `json:"user_id"`
	Type        int                           `json:"type"`
	TotalCount  int                           `json:"total_count"`
	RemainCount int                           `json:"remain_count"`
	LoopUnit    int                           `json:"loop_unit,omitempty"`
	WeekDay     int                           `json:"week_day,omitempty"`
	Hour        int                           `json:"hour,omitempty"`
	Minute      int                           `json:"minute,omitempty"`
	NextRunAt   int64                         `json:"next_run_at"`
	CreatedAt   int64                         `json:"created_at"`
	ThreadList  []*TwAddTweetScheduleRespItem `json:"thread_list,omitempty"`
}

type TwAddTweetScheduleRespItem struct {
	SortId    string   `json:"sort_id"`
	Text      string   `json:"text,omitempty" `
	MediaUrls []string `json:"media_urls,omitempty"`
}

type GetThreadWithTweetUrlByApiReq struct {
	TweetUrl string `json:"tweet_url" binding:"required,url"`
}

type GetAuthUserListReq struct {
	Name    string `json:"name,omitempty"`
	Account string `json:"account,omitempty"`
}

type GetAuthUserListResp struct {
	UserId          string `json:"user_id"`
	Name            string `json:"name"`
	Account         string `json:"account"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type UploadFiles2S3RespItem struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Status string `json:"status"`
	ErrMsg string `json:"err_msg"`
}

type TwAddTweetScheduleWithTimeReq struct {
	UserId          string                       `json:"user_id" binding:"min=1"`
	TwScheduleLibId *int64                       `json:"tw_schedule_lib_id,omitempty" binding:"omitempty,min=1"`
	LoopCount       *int                         `json:"loop_count" binding:"omitempty,min=1"`
	ScheduleTime    int64                        `json:"schedule_time" binding:"min=1"`
	SourceType      *int                         `json:"source_type,omitempty"`
	ThreadList      []*TwAddTweetScheduleReqItem `json:"thread_list" binding:"required_without=TwScheduleLibId,dive"`
}

type TwUpdateTweetScheduleWithTimeReq struct {
	ScheduleId   int64  `json:"schedule_id" binding:"min=1"`
	ScheduleTime *int64 `json:"schedule_time,omitempty" binding:"omitempty,min=1"`
}
