package twitterapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
	"strconv"
)

const (
	defaultTimeOut = 30 // s

	tweetCreateEndpoint endpoint = "2/tweets"

	rateLimit     = "x-user-limit-24hour-limit"
	rateRemaining = "x-user-limit-24hour-remaining"
	rateReset     = "x-user-limit-24hour-reset"
)

type endpoint string

type TwitterAPI struct {
	Token        string
	MaxResults   int
	UserId       string
	Client       *twitter.Client
	ErrorHandler ErrorHandler
}

func NewTwitterAPI(token string, maxResults int) (*TwitterAPI, error) {
	obj := new(TwitterAPI)
	client := GetGoTwitterClient(token)
	obj.MaxResults = maxResults
	obj.Token = token
	obj.Client = client
	obj.ErrorHandler = defaultErrorHandler

	return obj, nil
}

type ErrorHandler func(err error) error

func WithErrorHandler(handler ErrorHandler) func(*TwitterAPI) {
	return func(ta *TwitterAPI) {
		ta.ErrorHandler = handler
	}
}

func defaultErrorHandler(err error) error {
	var rtnErr error
	switch err.(type) {
	case *twitter.ErrorResponse:
		e := err.(*twitter.ErrorResponse)
		errMsg := ""
		for _, e2 := range e.Errors {
			errMsg += fmt.Sprintf(" detail message: %s", e2.Message)
		}
		errMsg = fmt.Sprintf("twitter callout status %d %s:%s %s", e.StatusCode, e.Title, e.Detail, errMsg)
		rtnErr = fmt.Errorf("%s", errMsg)
	default:
		rtnErr = fmt.Errorf("%w", err)
	}

	return rtnErr
}

// SetUserId the TwitterAPI struct field userId default is the auth user's id . By invoke this method change the userId
// when call the GetFollowersByUserId、GetFollowingByUserId、GetAllFollowerList、GetAllFollowingList require call this method
func (ta *TwitterAPI) SetUserId(userId string) {
	ta.UserId = userId
}

func (ta *TwitterAPI) GetFollowersByUserId(pageToken string) (*twitter.UserRaw, *TWResponseMeta, error) {
	opts := twitter.UserFollowersLookupOpts{
		MaxResults:      ta.MaxResults,
		PaginationToken: pageToken,
	}

	resp, err := ta.Client.UserFollowersLookup(context.Background(), ta.UserId, opts)
	if err != nil {
		return nil, nil, err
	}

	obj := &TWResponseMeta{
		ResultCount:   resp.Meta.ResultCount,
		NextToken:     resp.Meta.NextToken,
		PreviousToken: resp.Meta.PreviousToken,
	}

	return resp.Raw, obj, err
}

func (ta *TwitterAPI) GetFollowingByUserId(pageToken string) (*twitter.UserRaw, *TWResponseMeta, error) {
	opts := twitter.UserFollowingLookupOpts{
		MaxResults:      ta.MaxResults,
		PaginationToken: pageToken,
	}

	resp, err := ta.Client.UserFollowingLookup(context.Background(), ta.UserId, opts)
	if err != nil {
		return nil, nil, err
	}

	obj := &TWResponseMeta{
		ResultCount:   resp.Meta.ResultCount,
		NextToken:     resp.Meta.NextToken,
		PreviousToken: resp.Meta.PreviousToken,
	}
	return resp.Raw, obj, err
}

func (ta *TwitterAPI) GetAuthUser() (*twitter.UserRaw, error) {
	opts := twitter.UserLookupOpts{}

	resp, err := ta.Client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return resp.Raw, err
}

// getFollowInfo Twitter API Docs https://developer.twitter.com/en/docs/twitter-api/users/follows
func (ta *TwitterAPI) getAllFollowInfo(cb userCallbackFunc) ([]*twitter.UserObj, error) {
	results := make([]*twitter.UserObj, 0)
	pageToken := ""
	for {
		userRaw, meta, err := cb(pageToken)
		if err != nil {
			return nil, err
		}

		results = append(results, userRaw.Users...)
		if meta.NextToken == "" {
			break
		}
		pageToken = meta.NextToken
	}
	return results, nil
}

func (ta *TwitterAPI) SearchTweets(queryString string, opts twitter.TweetRecentSearchOpts) (*twitter.TweetRaw, *TWResponseMeta, error) {
	resp, err := ta.Client.TweetRecentSearch(context.Background(), queryString, opts)
	if err != nil {
		return nil, nil, err
	}

	obj := &TWResponseMeta{
		ResultCount: resp.Meta.ResultCount,
		NextToken:   resp.Meta.NextToken,
	}

	return resp.Raw, obj, err
}

func (ta *TwitterAPI) GetFollowerCount(userIds []string) (*twitter.UserRaw, error) {
	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldPublicMetrics,
			twitter.UserFieldProfileImageURL,
		},
	}
	resp, err := ta.Client.UserLookup(context.Background(), userIds, opts)
	if err != nil {
		return nil, err
	}
	return resp.Raw, err
}

func (ta *TwitterAPI) GetUserByAccount(account string) (*twitter.UserRaw, error) {
	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldProfileImageURL,
		},
	}
	resp, err := ta.Client.UserNameLookup(context.Background(), []string{account}, opts)
	return resp.Raw, err
}

func (ta *TwitterAPI) Like(tweetId string) (bool, error) {
	resp, err := ta.Client.UserLikes(context.Background(), ta.UserId, tweetId)
	if err != nil {
		return false, err
	}
	return resp.Data.Liked, nil
}

func (ta *TwitterAPI) Retweet(tweetId string) (bool, error) {
	resp, err := ta.Client.UserRetweet(context.Background(), ta.UserId, tweetId)
	if err != nil {
		return false, err
	}
	return resp.Data.Retweeted, nil
}

func (ta *TwitterAPI) CreateTweet(req *twitter.CreateTweetRequest) (*twitter.CreateTweetResponse, error) {
	//resp, err := ta.Client.CreateTweet(context.Background(), *req)
	resp, err := ta.MyCreateTweet(context.Background(), *req)
	if err != nil {
		return nil, ta.ErrorHandler(err)
	}

	return resp, nil
}

func (ta *TwitterAPI) UserMentionTimeline(userId string, opts *twitter.UserMentionTimelineOpts) (*twitter.UserMentionTimelineResponse, error) {
	resp, err := ta.Client.UserMentionTimeline(context.Background(), userId, *opts)
	if err != nil {
		return nil, ta.ErrorHandler(err)
	}

	return resp, nil
}

func (ta *TwitterAPI) TweetLookup(tweetIds []string) (*twitter.TweetLookupResponse, error) {
	opts := twitter.TweetLookupOpts{
		Expansions: []twitter.Expansion{twitter.ExpansionAttachmentsMediaKeys},
		MediaFields: []twitter.MediaField{
			twitter.MediaFieldMediaKey,
			twitter.MediaFieldURL,
			twitter.MediaFieldType,
			twitter.MediaFieldVariants,
		},
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldInReplyToUserID,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldAttachments,
			twitter.TweetFieldEntities,
			twitter.TweetFieldCreatedAt,
		},
	}

	resp, err := ta.Client.TweetLookup(context.Background(), tweetIds, opts)
	if err != nil {
		return nil, ta.ErrorHandler(err)
	}

	return resp, nil
}

func (ta *TwitterAPI) UserLookup(ids []string) (*twitter.UserLookupResponse, error) {
	opt := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldID,
			twitter.UserFieldName,
			twitter.UserFieldUserName,
			twitter.UserFieldProfileImageURL,
		},
	}

	resp, err := ta.Client.UserLookup(context.Background(), ids, opt)
	if err != nil {
		return nil, ta.ErrorHandler(err)
	}

	return resp, nil
}

func (ta *TwitterAPI) MyCreateTweet(ctx context.Context, tweet twitter.CreateTweetRequest) (*twitter.CreateTweetResponse, error) {
	if tweet.Media != nil {
		if len(tweet.Media.IDs) == 0 {
			return nil, fmt.Errorf("twitter input parameter error. media ids is required")
		}

		for _, id := range tweet.Media.IDs {
			if len(id) == 0 {
				return nil, fmt.Errorf("twitter input parameter error. media id is required")
			}
		}
	}
	if (tweet.Media == nil || len(tweet.Media.IDs) == 0) && len(tweet.Text) == 0 {
		return nil, fmt.Errorf("twitter input parameter error. create tweet text is required if no media ids")
	}

	body, err := json.Marshal(tweet)
	if err != nil {
		return nil, fmt.Errorf("create tweet marshal error %w", err)
	}

	ep := tweetCreateEndpoint.url(ta.Client.Host)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create tweet request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	ta.Client.Authorizer.Add(req)

	resp, err := ta.Client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create tweet response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusCreated {
		e := &twitter.ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &twitter.HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.Detail += fmt.Sprintf("\n Content: %s", tweet.Text)
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &twitter.CreateTweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &twitter.ResponseDecodeError{
			Name:      "create tweet",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

func (e endpoint) url(host string) string {
	return fmt.Sprintf("%s/%s", host, string(e))
}

func rateFromHeader(header http.Header) *twitter.RateLimit {
	limit, err := strconv.Atoi(header.Get(rateLimit))
	if err != nil {
		return nil
	}
	remaining, err := strconv.Atoi(header.Get(rateRemaining))
	if err != nil {
		return nil
	}
	reset, err := strconv.Atoi(header.Get(rateReset))
	if err != nil {
		return nil
	}
	return &twitter.RateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     twitter.Epoch(reset),
	}
}
