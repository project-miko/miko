package twitterapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/log"
	"github.com/project-miko/miko/tools/netutils"
)

const (
	twitterDomain = "https://api.twitter.com"

	authorizeUri = "https://twitter.com/i/oauth2/authorize"
	authTokenUri = "https://api.twitter.com/2/oauth2/token"
)

var (
	redirectUri = ""

	codeChallengeMethod = ""

	clientId     = ""
	clientSecret = ""

	challengeKey = ""
)

type Authorize struct {
	Token string
}

func (a Authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func GetGoTwitterClient(token string) *twitter.Client {
	client := &twitter.Client{
		Authorizer: Authorize{
			Token: token,
		},
		Client: &http.Client{
			Timeout: defaultTimeOut * time.Second,
		},
		Host: twitterDomain,
	}
	return client
}

func InitConfig() {
	redirectUri = conf.GetConfigString("twitter", "redirect_uri")
	challengeKey = conf.GetConfigString("twitter", "challenge_key")
	codeChallengeMethod = conf.GetConfigString("twitter", "code_challenge_method")
	clientId = conf.GetConfigString("twitter", "client_id")
	clientSecret = conf.GetConfigString("twitter", "client_secret")
}

func GetAuthCodeUrl(scope string) string {
	responseType := "code"
	state := tools.GetRandStr(16)
	// todo sha256
	codeChallenge := challengeKey

	u := url.Values{}
	u.Add("scope", scope)
	u.Add("state", state)
	u.Add("client_id", clientId)
	u.Add("redirect_uri", redirectUri)
	u.Add("code_challenge", codeChallenge)
	u.Add("response_type", responseType)
	u.Add("code_challenge_method", codeChallengeMethod)
	authCodeReqUrl := fmt.Sprintf("%s?%s", authorizeUri, u.Encode())
	return authCodeReqUrl
}

func GetAuthToken(code, codeVerifier string) (*TWResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("the parameter code can not be null")
	}
	grantType := "authorization_code"

	u := url.Values{}
	u.Add("code", code)
	u.Add("client_id", clientId)
	u.Add("grant_type", grantType)
	u.Add("redirect_uri", redirectUri)
	u.Add("code_verifier", codeVerifier)
	reqUrl := fmt.Sprintf("%s?%s", authTokenUri, u.Encode())

	req := netutils.NewHttpRequest(reqUrl)

	basicHeader := fmt.Sprintf("%s:%s", clientId, clientSecret)
	basicHeader = base64.RawStdEncoding.WithPadding('=').EncodeToString([]byte(basicHeader))

	_ = req.SetMethod("POST")
	req.SetHeader("Accept", "application/json")
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	req.SetHeader("Authorization", fmt.Sprintf("Basic %s", basicHeader))

	resp, _, err := req.Exec(time.Second * 10)
	if err != nil {
		return nil, fmt.Errorf("send get token request error %s", err.Error())
	}

	twResponse := new(TWResponse)
	err = json.Unmarshal(resp, twResponse)

	if err != nil {
		log.Error("", "json.Unmarshal: the byte transfer to struct TWResponse error %s", err.Error())

		result := make(map[string]interface{}, 0)
		err = json.Unmarshal(resp, &result)
		if err != nil {
			log.Error("", "json.Unmarshal: the byte transfer to map error %s", err.Error())
		} else {
			log.Error("", "json.Unmarshal: the raw response map is %v", result)
		}

		return nil, fmt.Errorf("json.Unmarshal: the byte transfer to struct TWResponse error %s", err.Error())
	}

	if len(twResponse.Error) > 0 {
		return nil, fmt.Errorf("refresh token error %s", twResponse.ErrorDescription)
	}

	return twResponse, nil
}

// cursor pagination

func GetAuthUser(token string) (*twitter.UserRaw, error) {
	opts := twitter.UserLookupOpts{UserFields: []twitter.UserField{twitter.UserFieldProfileImageURL}}

	client := GetGoTwitterClient(token)
	resp, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return resp.Raw, err
}

func GetAuthUserIdByToken(token string) (string, error) {
	authUser, err := GetAuthUser(token)
	if err != nil {
		return "", fmt.Errorf("get auth user by access token failed %s", err.Error())
	}
	if len(authUser.Users) == 0 {
		return "", fmt.Errorf("the auth user not found")
	}
	return authUser.Users[0].ID, nil
}

func RefreshToken(refreshToken string) (*TWResponse, error) {
	u := url.Values{}
	u.Add("grant_type", "refresh_token")
	u.Add("refresh_token", refreshToken)
	reqUrl := fmt.Sprintf("%s?%s", authTokenUri, u.Encode())
	req := netutils.NewHttpRequest(reqUrl)

	basicHeader := fmt.Sprintf("%s:%s", clientId, clientSecret)
	basicHeader = base64.RawStdEncoding.WithPadding('=').EncodeToString([]byte(basicHeader))

	_ = req.SetMethod("POST")
	req.SetHeader("Accept", "application/json")
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	req.SetHeader("Authorization", fmt.Sprintf("Basic %s", basicHeader))
	resp, _, err := req.Exec(10 * time.Second)
	if err != nil {
		return nil, err
	}

	twResponse := new(TWResponse)
	err = json.Unmarshal(resp, twResponse)

	if err != nil {
		log.Error("", "json.Unmarshal: the byte transfer to struct TWResponse error %s", err.Error())

		result := make(map[string]interface{}, 0)
		err = json.Unmarshal(resp, &result)
		if err != nil {
			log.Error("", "json.Unmarshal: the byte transfer to map error %s", err.Error())
		} else {
			log.Error("", "json.Unmarshal: the raw response map is %v", result)
		}

		return nil, fmt.Errorf("json.Unmarshal: the byte transfer to struct TWResponse error %s", err.Error())
	}

	if len(twResponse.Error) > 0 {
		return nil, fmt.Errorf("refresh token error %s", twResponse.ErrorDescription)
	}

	return twResponse, nil
}
