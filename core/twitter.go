package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/models"
	"github.com/project-miko/miko/models/data"
	"github.com/project-miko/miko/sdk/twitterapi"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/log"
	"github.com/shopspring/decimal"
)

const (
	TwitterBaseUrl = "https://twitter.com"

	RouteDeFlashGetThreadByApi = "/security/twitter/getthreadwithtweeturlbyapi"

	RouteDeFlashGetThread  = "/twitter/getthreadwithtweeturl"
	RouteRSSHubTweetDetail = "/twitter/tweet"
)

// getUserMap user_id -> UserObj
func getUserMap(userIds []string, token string, totalReqCount *int) (map[string]*twitter.UserObj, error) {
	maxResults := 100
	twAPI, err := twitterapi.NewTwitterAPI(token, maxResults)
	if err != nil {
		return nil, fmt.Errorf("twitterapi.NewTwitterAPI(%s, %d) error %s", token, maxResults, err.Error())
	}

	userRaw, err := twAPI.GetFollowerCount(userIds)
	if err != nil {
		return nil, fmt.Errorf("twAPI.GetFollowerCount() error %s", err.Error())
	}

	userMap := make(map[string]*twitter.UserObj, 0)
	for _, v := range userRaw.Users {
		userMap[v.ID] = v
	}

	*totalReqCount++
	return userMap, nil
}

func splitTask(amount, limit int64) int64 {
	var result int64 = 0
	if amount%limit == 0 {
		result = amount / limit
	} else {
		result = amount/limit + 1
	}
	return result
}

// concatUserId eg: (from:1234 OR from:4567) -is:quote -is:reply -is:retweet
func concatUserId(ids []string) string {
	handledIds := make([]string, 0)
	prefix := "from:"
	for _, v := range ids {
		v = prefix + v
		handledIds = append(handledIds, v)
	}
	sep := " OR "
	exclude := "-is:quote -is:reply -is:retweet"
	return fmt.Sprintf("(%s) %s", strings.Join(handledIds, sep), exclude)
}

func RefreshTwMetricInfo() error {
	log.Info("", "RefreshTwMetricInfo() start")
	userInfoList, err := models.GetAllTwUserInfoList(-1)
	if err != nil {
		return err
	}

	userIds := make([]string, 0)
	for _, v := range userInfoList {
		userIds = append(userIds, v.UserId)
	}

	temp := 0
	// use user_id to get user's follower count
	userMap, err := getUserMap(userIds, conf.TwitterAPIToken, &temp)
	if err != nil {
		return err
	}

	items := make([]*data.TwMetricItem, 0)
	for _, v := range userInfoList {
		userObj, ok := userMap[v.UserId]
		if !ok {
			log.Info("tw-monitor-in-time", "the query data not contain the user's userId. user_id :%s", v.UserId)
			continue
		}

		item := new(data.TwMetricItem)
		item.UserId = userObj.ID
		item.FollowCount = userObj.PublicMetrics.Followers
		items = append(items, item)
	}

	twMetricInfo := new(data.TwMetricInfo)
	twMetricInfo.StatisticAt = tools.GetMillisecond(time.Now().UTC())
	twMetricInfo.Items = items

	err = saveTwMetricInfoToRedis(twMetricInfo)
	if err != nil {
		return err
	}

	log.Info("", "RefreshTwMetricInfo() done")
	return nil
}

func saveTwMetricInfoToRedis(value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	client := models.GetRdbInst()
	err = client.SetString(conf.AISERTwMetricInfoString, string(b), conf.TwMetricRefreshInterval)
	if err != nil {
		return err
	}

	return nil
}

func GetTwitterStatisticalChartInfoAll(userType, start, end int64) (map[string]interface{}, error) {
	userInfoList, err := models.GetAllTwUserInfoList(userType)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	for _, v := range userInfoList {
		userIds = append(userIds, v.UserId)
	}

	dataList, err := models.GetListByUserIds(userIds, start, end)
	if err != nil {
		return nil, err
	}

	return doGetTwitterStatisticalChartInfo(dataList, start, end)
}

func GetTwitterStatisticalChartInfoByUserId(userId string, start, end int64) (map[string]interface{}, error) {
	dataList, err := models.GetListByUserIds([]string{userId}, start, end)
	if err != nil {
		return nil, err
	}

	return doGetTwitterStatisticalChartInfo(dataList, start, end)
}

func doGetTwitterStatisticalChartInfo(dataList []*models.TwDailyData, start, end int64) (map[string]interface{}, error) {
	followCountMap := make(map[int64]int64)
	for _, v := range dataList { // count follower count by time point
		followCountMap[v.StatisticAt] += int64(v.FollowerCount)
	}

	type Chart struct {
		X []string `json:"x"`
		Y []string `json:"y"`
	}

	chart := new(Chart)
	var step int64 = conf.CrawlInterval * 60 * 1000 // step, 1 hour, unit ms
	for x := start; x <= end; x += step {
		y := ""
		followCount, ok := followCountMap[x] // x axis time point follower count
		if ok {                              // if has follower count, then use it
			y = strconv.FormatInt(followCount, 10)
		}

		chart.X = append(chart.X, strconv.FormatInt(x, 10))
		chart.Y = append(chart.Y, y)
	}

	m := map[string]interface{}{
		"data": chart,
	}

	return m, nil
}

func GetTwitterStatisticInfo(page, limit, userType, start, end int64) (map[string]interface{}, error) {
	amount, userInfoList, err := models.GetTwUserInfoListByType(page, limit, userType, start)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	userInfoMap := make(map[string]*models.TwUserInfo, 0)
	for _, v := range userInfoList {
		userIds = append(userIds, v.UserId)
		userInfoMap[v.UserId] = v
	}

	allList, err := models.GetListByUserIds(userIds, start, end)
	if err != nil {
		return nil, err
	}

	userDataMap := make(map[string][]*models.TwDailyData, 0)
	for _, v := range allList {
		ll, ok := userDataMap[v.UserId]
		if !ok {
			ll = make([]*models.TwDailyData, 0)
		}
		ll = append(ll, v)
		userDataMap[v.UserId] = ll
	}

	// sort by statistic_at asc
	for uId, dataList := range userDataMap {
		sort.Slice(dataList, func(i, j int) bool {
			return dataList[i].StatisticAt < dataList[j].StatisticAt
		})
		userDataMap[uId] = dataList
	}

	monitoringDTO := new(data.TwMonitoringDTO)
	userDataList := make([]*data.TwUserDataDTO, 0)
	for _, id := range userIds {
		userInfo := userInfoMap[id]
		userData := userDataMap[id]
		dto := &data.TwUserDataDTO{}

		dto.UserId = userInfo.UserId
		dto.UserName = userInfo.Name
		dto.UserAccount = userInfo.Account
		dto.ProfileImageUrl = userInfo.ProfileImageUrl

		// if not found user data, return user basic info, other statistic data is 0, then return
		if userData == nil {
			userDataList = append(userDataList, dto)
			continue
		}

		// if only one record, it means new added user, only count one day data. then use 0 follower count and 0 increase follower rate
		followChange := 0
		if len(userData) > 1 {
			followChange = userData[len(userData)-1].FollowerCount - userData[0].FollowerCount
		}
		// use latest follower count
		dto.FollowCount = userData[len(userData)-1].FollowerCount
		dto.FollowChange = followChange

		// when query, to count follower count change, the start time boundary is also queried. when count, exclude the start boundary
		startStatisticAt := userData[0].StatisticAt
		for _, v := range userData {
			if v.StatisticAt == startStatisticAt {
				continue
			}
			dto.LikeCount += v.LikeCount
			dto.ReplyCount += v.ReplyCount
			dto.RetweetCount += v.RetweetCount
			dto.TweetCount += v.TweetCount
		}

		dto.IncreaseFollowRate = Div(followChange, dto.FollowCount)
		dto.LikeFollowRate = Div(dto.LikeCount, dto.FollowCount)

		monitoringDTO.LikeCount += dto.LikeCount
		monitoringDTO.ReplyCount += dto.ReplyCount
		monitoringDTO.ReTweetCount += dto.RetweetCount
		monitoringDTO.FollowCount += dto.FollowCount
		monitoringDTO.FollowChange += dto.FollowChange

		userDataList = append(userDataList, dto)
	}
	monitoringDTO.IncreaseFollowRate = Div(monitoringDTO.FollowChange, monitoringDTO.FollowCount)
	monitoringDTO.UserDataList = userDataList

	m := map[string]interface{}{
		"list": monitoringDTO,
		"paging": data.Paging{
			Amount: amount,
			Page:   page,
			Limit:  limit,
		},
	}

	return m, nil
}

func Div(dividend, divisor int) decimal.Decimal {
	if divisor == 0 || dividend == 0 {
		return decimal.NewFromInt(0)
	}
	return decimal.NewFromInt(int64(dividend) * 100).Div(decimal.NewFromInt(int64(divisor))).Truncate(4)
}

type GetTweetResp struct {
	Account   string   `json:"account"`
	Avatar    string   `json:"avatar"`
	Tweets    []string `json:"tweets"`
	PublishAt int64    `json:"publish_at"`
}

func mapToStruct(m map[string]interface{}, out interface{}) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, out)
}

type ErrGetMediaUploadStatusInProgress struct {
	PInfo *twitterapi.ProcessingInfo
}

func (e *ErrGetMediaUploadStatusInProgress) Error() string {
	return fmt.Sprintf("progress percent:%d, check after seconds:%d.", e.PInfo.ProgressPercent, e.PInfo.CheckAfterSecs)
}

type ErrGetMediaUploadStatusFailed struct {
	Err error
}

func (e *ErrGetMediaUploadStatusFailed) Error() string {
	return e.Err.Error()
}

type ErrCreateTweet struct {
	Err error
}

func (e *ErrCreateTweet) Error() string {
	return e.Err.Error()
}

func RefreshAccessToken(account *models.TwAccount) error {
	now := time.Now()
	if account.ExpiredAt > tools.GetMillisecond(now.Add(10*time.Minute)) { // refresh 10 minutes before expired
		return nil
	}

	log.Info("", "refresh token start, userId:%s", account.UserId)
	resp, err := twitterapi.RefreshToken(account.RefreshToken)
	if err != nil {
		return err
	}
	account.AccessToken = resp.AccessToken
	account.RefreshToken = resp.RefreshToken
	account.ExpiredAt = tools.GetMillisecond(now.Add(time.Duration(resp.ExpiresIn) * time.Second))
	account.UpdatedAt = tools.GetMillisecond(now)

	err = account.SaveOrUpdateWithLog()
	if err != nil {
		return err
	}

	log.Info("", "refresh token success, userId:%s", account.UserId)
	return nil
}

var JobHandleFunc = func(userId string, twScheduleLibId int64) {
	log.Info("", "job callback function execution start")

	now := tools.GetMillisecond(time.Now())
	err := doJobHandle(userId, twScheduleLibId)

	jobId := fmt.Sprintf("%s-%d", userId, twScheduleLibId)
	if e := models.SaveScheduleLog(jobId, now, err); e != nil {
		errMsg := e.Error()
		log.Error("", "doJobHandle() error %s", errMsg)
		return
	}

	log.Info("", "job callback function execution success")
}

func doJobHandle(userId string, twScheduleLibId int64) error {
	// todo add lock
	twSchedule, err := models.GetTwScheduleByUserIdAndTwLibId(userId, twScheduleLibId)
	if err != nil {
		return err
	}
	if twSchedule == nil {
		return conf.ErrRecordNotFound
	}

	nextRunAt := new(int64)
	err = doUploadTwMediaAndCreateTweet(userId, twScheduleLibId, nextRunAt)
	if err != nil {
		err = fmt.Errorf("doUploadTwMediaAndCreateTweet() error, %w", err)
		twSchedule.Status = models.TwScheduleStatusError
	} else {
		twSchedule.RemainCount--
		if twSchedule.RemainCount == 0 {
			twSchedule.Status = models.TwScheduleStatusFinished
		}
		twSchedule.NextRunAt = *nextRunAt
	}

	if e := twSchedule.Update(); e != nil {
		return e
	}

	return err
}

func doUploadTwMediaAndCreateTweet(userId string, twScheduleLibId int64, nextRunAt *int64) error {
	twScheduleLib, err := models.GetTwScheduleLibById(twScheduleLibId)
	if err != nil {
		return err
	}

	items := make([]*data.TwAddTweetScheduleReqItem, 0)
	err = json.Unmarshal([]byte(twScheduleLib.Content), &items)
	if err != nil {
		return err
	}

	uploadFiles := make([]*data.TwUploadFileReqItem, 0)
	for _, item := range items {
		for _, v2 := range item.MediaUrls {
			reqItem := new(data.TwUploadFileReqItem)
			reqItem.Id = item.SortId
			reqItem.MediaUrl = v2
			uploadFiles = append(uploadFiles, reqItem)
		}
	}

	uReq := new(data.TwUploadMediaReq)
	uReq.UserId = userId
	uReq.UploadFiles = uploadFiles

	mediaResp, err := UploadTwMedia(uReq)
	if err != nil {
		return err
	}

	m := make(map[string][]string)
	for _, item := range mediaResp.MediaItems {
		if len(item.ErrMsg) != 0 {
			return fmt.Errorf("UploadTwMedia() error %s", item.ErrMsg)
		}
		m[item.Id] = append(m[item.Id], item.MediaId)
	}

	tweets := make([]*data.CreateTweetItem, 0)
	for _, v := range items {
		mediaIds := m[v.SortId]
		ci := new(data.CreateTweetItem)
		ci.MediaIds = mediaIds
		ci.SortId = v.SortId
		ci.Text = v.Text
		tweets = append(tweets, ci)
	}

	req := new(data.CreateTweetReq)
	req.UserId = userId
	req.Tweets = tweets

	time.Sleep(1 * time.Second)
	const maxRetryCount = 6
	const waitSec = 10                   // use fixed time to wait
	for i := 0; i < maxRetryCount; i++ { // max wait 1 minute
		err = CreateTweet(req)
		if err == nil { // send success, exit loop
			break
		}

		if _, ok := err.(*ErrGetMediaUploadStatusInProgress); ok { // if send failed, and return Twitter is uploading status, then wait and retry
			if i == maxRetryCount-1 {
				return fmt.Errorf("create tweet error after max retry")
			}
			log.Error("", "call core.CreateTweet() retry")
			//waitSec := e.PInfo.CheckAfterSecs
			time.Sleep(time.Duration(waitSec) * time.Second)
		} else { // if not uploading status, return directly
			return err
		}
	}

	tag := GetTag(userId, twScheduleLibId)
	jobs, err := scheduler.Scheduler.FindJobsByTag(tag)
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return fmt.Errorf("can not found job by tag. %s", tag)
	}

	*nextRunAt = tools.GetMillisecond(jobs[0].NextRun())

	return nil
}

func GetTag(userId string, twScheduleLibId int64) string {
	return fmt.Sprintf("%s-%d", userId, twScheduleLibId)
}

func ParseTagString(tag string) (string, int64) {
	s := strings.Split(tag, "-")
	i, _ := strconv.ParseInt(s[1], 10, 64)
	return s[0], i
}

func InitTwCreateTweetJobs() error {
	log.Info("", "init twitter create tweet cron jobs start")
	return ReloadTwCreateTweetJobsFromDB()
}

func ReloadTwCreateTweetJobsFromDB() error {
	log.Info("", "reload twitter create tweet cron jobs start")

	list, err := models.GetAllTwScheduleList(models.TwScheduleStatusUnFinished)
	if err != nil {
		return err
	}

	s := GetScheduler()
	for _, v := range list {
		scheduler.SetJobFuncAndParams(JobHandleFunc, v.UserId, v.TwScheduleLibId)
		tag := GetTag(v.UserId, v.TwScheduleLibId)
		j, err := s.Add(v.CronExpression, tag, v.RemainCount, v.NextRunAt)
		if err != nil {
			return fmt.Errorf("scheduler.Add() error %s", err.Error())
		}

		if !j.NextRun().IsZero() {
			v.NextRunAt = tools.GetMillisecond(j.NextRun())
		}
		if e := v.Update(); e != nil {
			return fmt.Errorf("twSchedule.Update() error %s", e.Error())
		}
	}

	log.Info("", "reload twitter create tweet cron jobs success. length of jobs :%d", len(list))

	return nil
}

func UploadTwMedia(req *data.TwUploadMediaReq) (*data.TwUploadMediaResp, error) {
	log.Info("", "upload media start. user_id: %s, the length of media: %d", req.UserId, len(req.UploadFiles))
	uploadFiles := req.UploadFiles
	userId := req.UserId

	exists, err := models.GetTwOAuth1ByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("models.GetTwOAuth1ByUserId() error %s", err.Error())
	}
	if exists == nil {
		return nil, conf.ErrRecordNotFound
	}

	twApi := twitterapi.NewTwitterAPIV1(exists.AccessToken, exists.AccessSecret)

	resp := new(data.TwUploadMediaResp)
	items := make([]*data.TwMediaRespItem, 0)
	for _, v := range uploadFiles {
		item := new(data.TwMediaRespItem)
		item.Id = v.Id
		item.Status = twitterapi.UploadMediaSucceeded

		mediaId, err := twApi.UploadMediaFromUrl(v.MediaUrl)
		if err != nil {
			log.Error("", "twitterapi.UploadMediaFromUrl() error %s", err.Error())
			item.ErrMsg = err.Error()
			item.Status = twitterapi.UploadMediaFailed
		}
		item.MediaId = mediaId
		items = append(items, item)
	}

	resp.MediaItems = items

	return resp, nil
}

func CreateTweet(req *data.CreateTweetReq) error {
	log.Info("", "create tweet start, userId:%s. the length of tweets: %d", req.UserId, len(req.Tweets))
	userId := req.UserId
	tweetItems := req.Tweets

	//const maxTweetCharCount = 10000
	//const showCharCount = 20
	//errMsg := ""
	//for _, v := range tweetItems {
	//	l := utf8.RuneCountInString(v.Text)
	//	if l > maxTweetCharCount {
	//		chars := []rune(v.Text)[:showCharCount]
	//		errMsg += fmt.Sprintf("tweet text: %s, length: %d. ", string(chars)+"...", l)
	//	}
	//}
	//if len(errMsg) != 0 {
	//	ctrl.JsonError(c, conf.ApiCodeParamErr, errMsg)
	//	return
	//}

	account, err := models.GetTwAccountByUserId(userId)
	if err != nil {
		return fmt.Errorf("models.GetTwAccountByUserId() error %s", err.Error())
	}
	if account == nil {
		return conf.ErrRecordNotFound
	}

	twOAuth1, err := models.GetTwOAuth1ByUserId(userId)
	if err != nil {
		return fmt.Errorf("models.GetTwOAuth1ByUserId() error %s", err.Error())
	}
	if twOAuth1 == nil {
		return conf.ErrRecordNotFound
	}

	twApiV1 := twitterapi.NewTwitterAPIV1(twOAuth1.AccessToken, twOAuth1.AccessSecret)
	for _, v := range tweetItems {
		for _, mediaId := range v.MediaIds {
			mediaResp, err := twApiV1.GetMediaUploadStatus(mediaId)
			if err != nil {
				if e, ok := err.(*anaconda.ApiError); ok {
					if e.StatusCode == http.StatusNotFound {
						continue
					}
				}
				return fmt.Errorf("twitter.GetMediaUploadStatus() error %w", err)
			}

			pInfo := mediaResp.ProcessingInfo
			switch pInfo.State {
			case twitterapi.UploadMediaInProgress:
				if pInfo.Error != nil {
					return &ErrGetMediaUploadStatusFailed{Err: pInfo.Error}
				}

				return &ErrGetMediaUploadStatusInProgress{PInfo: pInfo}
			case twitterapi.UploadMediaFailed:
				return &ErrGetMediaUploadStatusFailed{Err: pInfo.Error}
			}
		}
	}

	if err = RefreshAccessToken(account); err != nil {
		return fmt.Errorf("RefreshAccessToken() error %s", err.Error())
	}

	twApi, err := twitterapi.NewTwitterAPI(account.AccessToken, -1)
	if err != nil {
		return err
	}

	sort.SliceStable(tweetItems, func(i, j int) bool {
		ai, _ := strconv.Atoi(tweetItems[i].SortId)
		aj, _ := strconv.Atoi(tweetItems[j].SortId)
		return ai < aj
	})

	successTweetIds := make([]string, 0)
	tempInReplyToTweetID := ""
	for _, v := range tweetItems {
		createTweetReq := &twitter.CreateTweetRequest{
			Text: v.Text,
		}
		if len(v.MediaIds) != 0 {
			createTweetReq.Media = &twitter.CreateTweetMedia{
				IDs: v.MediaIds,
			}
		}
		if len(tempInReplyToTweetID) != 0 {
			createTweetReq.Reply = &twitter.CreateTweetReply{
				InReplyToTweetID: tempInReplyToTweetID,
			}
		}

		log.Info("", "create tweet start, userId:%s", userId)
		resp, err := twApi.CreateTweet(createTweetReq)
		SaveUserRateLimitCreateTweet(userId, resp, err)
		if err != nil {
			return &ErrCreateTweet{Err: err}
		}

		respTweetId := resp.Tweet.ID
		tempInReplyToTweetID = respTweetId

		log.Info("", "create tweet success, userId:%s, tweetId:%s, inReplyToTweetID:%s", userId, respTweetId, tempInReplyToTweetID)
		successTweetIds = append(successTweetIds, respTweetId)
	}

	return nil
}

type AddTweetScheduleParams struct {
	UserId     string
	SourceType int
	LoopCount  int
	CronExp    string
	ThreadList []*data.TwAddTweetScheduleReqItem
}
