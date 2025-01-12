package twitterapi

import (
	"github.com/g8rswimmer/go-twitter/v2"
)

type userCallbackFunc func(pageToken string) (*twitter.UserRaw, *TWResponseMeta, error)

type TWResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TWResponseMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}
