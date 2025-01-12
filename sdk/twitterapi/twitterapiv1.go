package twitterapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/tools/log"
	"github.com/project-miko/miko/tools/mediautils"
)

const (
	maxChunkSize = 5 * 1024 * 1024 // 5MB

	MimePrefixImage = "image/"
	MimePrefixVideo = "video/"

	UploadMediaInProgress = "in_progress"
	UploadMediaSucceeded  = "succeeded"
	UploadMediaFailed     = "failed"

	uploadEndpoint = "https://upload.twitter.com/1.1/media/upload.json"

	maxRetryCount = 3
)

var (
	consumerKey    = ""
	consumerSecret = ""

	callbackUrl = ""
	tempAuthMap = sync.Map{}

	oauthCredentials oauth.Credentials
)

type V1 struct {
	Client      *anaconda.TwitterApi
	OauthClient *oauth.Client
}

type ProcessingInfo struct {
	State           string `json:"state"`
	CheckAfterSecs  int    `json:"check_after_secs"`
	ProgressPercent int    `json:"progress_percent,omitempty"`
	Error           *Error `json:"error,omitempty"`
}

type ErrUnsupportedMimeType struct {
	MimeType string
}

func (e *ErrUnsupportedMimeType) Error() string {
	return fmt.Sprintf("unsupported mime type: %s", e.MimeType)
}

type Error struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d, name:%s, message:%s.", e.Code, e.Name, e.Message)
}

type VideoInfo struct {
	VideoType string `json:"video_type"`
}

type MediaData struct {
	MediaID          int64           `json:"media_id"`
	MediaIDString    string          `json:"media_id_string"`
	ExpiresAfterSecs int64           `json:"expires_after_secs,omitempty"`
	Video            *VideoInfo      `json:"video,omitempty"`
	ProcessingInfo   *ProcessingInfo `json:"processing_info"`
}

func InitTwitterAPIV1() {
	consumerKey = conf.GetConfigString("twitter_v1", "consumer_Key")
	if len(consumerKey) == 0 {
		panic("twitter consumer_Key is not config")
	}
	consumerSecret = conf.GetConfigString("twitter_v1", "consumer_secret")
	if len(consumerSecret) == 0 {
		panic("twitter consumer_secret is not config")
	}
	callbackUrl = conf.GetConfigString("twitter_v1", "callback_url")
	if len(callbackUrl) == 0 {
		panic("twitter callback_url is not config")
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	oauthCredentials = oauth.Credentials{
		Token:  consumerKey,
		Secret: consumerSecret,
	}
}

func NewTwitterAPIV1(accessToken string, accessSecret string) *V1 {
	client := anaconda.NewTwitterApi(accessToken, accessSecret)

	client.HttpClient = &http.Client{
		Timeout: defaultTimeOut * time.Second,
	}
	oauthClient := &oauth.Client{
		Credentials: oauthCredentials,
	}

	return &V1{
		Client:      client,
		OauthClient: oauthClient,
	}
}

func (ta *V1) SetCredentials(accessToken string, accessSecret string) {
	ta.Client.Credentials.Token = accessToken
	ta.Client.Credentials.Secret = accessSecret
}

func (ta *V1) GetAuthorizationUrl() (string, error) {
	authUrl, tempCred, err := ta.Client.AuthorizationURL(callbackUrl)
	if err != nil {
		return "", err
	}

	tempAuthMap.Store(tempCred.Token, tempCred.Secret)

	return authUrl, nil
}

func (ta *V1) ExchangeAccessToken(requestToken, verifier string) (string, string, error) {
	val, ok := tempAuthMap.Load(requestToken)
	if !ok {
		return "", "", fmt.Errorf("get requestSecret by requestToken from temp auth map error")
	}

	requestSecret := val.(string)

	cred := &oauth.Credentials{
		Token:  requestToken,
		Secret: requestSecret,
	}

	cred, _, err := ta.Client.GetCredentials(cred, verifier)
	if err != nil {
		return "", "", err
	}

	return cred.Token, cred.Secret, err
}

func (ta *V1) UploadMediaFromUrl(mediaUrl string) (string, error) {
	log.Info("", "UploadMediaFromUrl start, media url: %s", mediaUrl)
	contentType, d, err := mediautils.DownloadFileFromURL(mediaUrl)
	if err != nil {
		return "", fmt.Errorf("mediautils.DownloadFileFromURL() error %s", err.Error())
	}
	log.Info("", "download file success, length of bytes: %d", len(d))

	mediaId, err := ta.UploadMediaBinary(contentType, d)
	if err != nil {
		return "", fmt.Errorf("twitterapi.UploadMediaBinary() error %w. mediaUrl: %s", err, mediaUrl)
	}
	log.Info("", "upload media success, mediaId:%s", mediaId)

	return mediaId, nil
}

func (ta *V1) UploadMediaBase64(mediaData []byte) (string, error) {
	client := ta.Client
	mimeType := http.DetectContentType(mediaData)

	if strings.HasPrefix(mimeType, MimePrefixImage) {
		media, err := client.UploadMedia(base64.StdEncoding.EncodeToString(mediaData))
		if err != nil {
			return "", err
		}

		log.Info("", "upload image success. mediaId:%s", media.MediaIDString)
		return media.MediaIDString, nil
	} else if strings.HasPrefix(mimeType, MimePrefixVideo) {
		// initialize video upload
		videoInit, err := ta.UploadVideoInit(len(mediaData), mimeType)
		if err != nil {
			return "", err
		}

		mediaId := videoInit.MediaIDString
		log.Info("", "upload video success. command:init, mediaId:%s, totalBytes:%d", mediaId, len(mediaData))

		segmentIdx := 0
		for i := 0; i < len(mediaData); i += maxChunkSize {
			end := i + maxChunkSize
			if end > len(mediaData) {
				end = len(mediaData)
			}

			chunk := mediaData[i:end]
			base64Chunk := base64.RawURLEncoding.EncodeToString(chunk)

			if err = client.UploadVideoAppend(mediaId, segmentIdx, base64Chunk); err != nil {
				return "", err
			}

			log.Info("", "upload video success. command:append, mediaId:%s, segmentIndex:%d", mediaId, segmentIdx)
			segmentIdx++
		}

		videoFinalize, err := client.UploadVideoFinalize(mediaId)
		if err != nil {
			return "", err
		}

		log.Info("", "upload video success. command:finalize, mediaId:%s", mediaId)

		return videoFinalize.MediaIDString, nil
	} else {
		return "", fmt.Errorf("unsupported mime type: %s", mimeType)
	}
}

func (ta *V1) UploadMediaBinary(contentType string, mediaData []byte) (string, error) {
	// DetectContentType always returns a valid MIME type: if it cannot determine a more specific one,
	// it returns "application/octet-stream"
	mimeType := http.DetectContentType(mediaData)
	//const trimIdx = 4
	//// if the mime type detected by the standard library does not match the mime type in the http response header,
	//// then use the mime type in the http response header as the mime type
	//if len(contentType) > trimIdx && len(mimeType) > trimIdx && contentType[:trimIdx] != mimeType[:trimIdx] {
	//	mimeType = contentType
	//}

	if strings.HasPrefix(mimeType, MimePrefixImage) { // image
		media, err := ta.UploadImage(mediaData)
		if err != nil {
			return "", err
		}

		log.Info("", "upload image success. mediaId:%s", media.MediaIDString)
		return media.MediaIDString, nil
	}

	const MimeTypeOctetStream = "application/octet-stream" // use binary stream type
	// initialize video upload
	videoInit, err := ta.UploadVideoInit(len(mediaData), MimeTypeOctetStream)
	if err != nil {
		return "", err
	}

	mediaId := videoInit.MediaIDString
	log.Info("", "upload video success. command:init, mediaId:%s, totalBytes:%d", mediaId, len(mediaData))

	segmentIdx := 0
	for i := 0; i < len(mediaData); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(mediaData) {
			end = len(mediaData)
		}

		chunk := mediaData[i:end]
		if err = ta.UploadVideoAppend(mediaId, segmentIdx, chunk); err != nil {
			return "", err
		}

		log.Info("", "upload video success. command:append, mediaId:%s, segmentIndex:%d", mediaId, segmentIdx)
		segmentIdx++
	}

	videoFinalize, err := ta.UploadVideoFinalize(mediaId)
	if err != nil {
		return "", err
	}

	log.Info("", "upload video success. command:finalize, mediaId:%s", mediaId)

	return videoFinalize.MediaIDString, nil
}

func (ta *V1) UploadImage(mediaData []byte) (*anaconda.Media, error) {
	v := url.Values{}

	media := new(anaconda.Media)
	if err := ta.doFormRequest(v, mediaData, media); err != nil {
		return nil, err
	}

	return media, nil
}

func (ta *V1) UploadVideoInit(totalBytes int, mimeType string) (*anaconda.ChunkedMedia, error) {
	// initialize video upload
	chunkMedia := new(anaconda.ChunkedMedia)
	v := url.Values{}
	v.Set("total_bytes", strconv.Itoa(totalBytes))
	v.Set("command", "INIT")
	v.Set("media_type", mimeType)
	v.Set("media_category", "tweet_video")

	if err := ta.doFormRequest(v, nil, chunkMedia); err != nil {
		return nil, err
	}

	return chunkMedia, nil
}

func (ta *V1) UploadVideoAppend(mediaId string, segmentIndex int, chunkBytes []byte) error {
	v := url.Values{}
	v.Set("command", "APPEND")
	v.Set("media_id", mediaId)
	v.Set("segment_index", strconv.Itoa(segmentIndex))

	var emptyResponse interface{}

	if err := ta.doFormRequest(v, chunkBytes, &emptyResponse); err != nil {
		return err
	}

	return nil
}

func (ta *V1) UploadVideoFinalize(mediaId string) (*MediaData, error) {
	v := url.Values{}
	v.Set("command", "FINALIZE")
	v.Set("media_id", mediaId)

	mediaResponse := new(MediaData)

	if err := ta.doFormRequest(v, nil, mediaResponse); err != nil {
		return nil, err
	}

	return mediaResponse, nil
}

func (ta *V1) GetMediaUploadStatus(mediaId string) (*MediaData, error) {
	v := url.Values{}
	v.Set("command", "STATUS")
	v.Set("media_id", mediaId)

	mediaResponse := new(MediaData)

	if err := ta.doRequest(http.MethodGet, "", &v, nil, mediaResponse); err != nil {
		return nil, err
	}

	return mediaResponse, nil
}

func (ta *V1) doFormRequest(fields url.Values, b []byte, data interface{}) error {
	// create a multipart form with fields and files
	body, contentType, err := createMultipartForm(fields, b)
	if err != nil {
		return err
	}

	return ta.doRequest(http.MethodPost, contentType, nil, body, data)
}

func (ta *V1) doRequest(method string, contentType string, params *url.Values, body io.Reader, data interface{}) error {
	reqUrl := uploadEndpoint
	if params != nil {
		reqUrl = reqUrl + "?" + params.Encode()
	}

	req, err := http.NewRequest(method, reqUrl, body)
	if err != nil {
		return err
	}

	// set request header
	req.Header.Set("Content-Type", contentType)
	// Sign the request.
	if err := ta.OauthClient.SetAuthorizationHeader(req.Header, ta.Client.Credentials, method, req.URL, url.Values{}); err != nil {
		return err
	}

	// send request
	resp, err := ta.Client.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	if strings.HasSuffix(resp.Request.URL.String(), "upload.json") {
		if resp.StatusCode == 204 {
			// empty response, don't decode
			return nil
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return anaconda.NewApiError(resp)
		}
	} else if resp.StatusCode != 200 {
		return anaconda.NewApiError(resp)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request failed: %s", err.Error())
	}

	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		return err
	}

	return nil
}

func createMultipartForm(fields url.Values, b []byte) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add fields
	for key, value := range fields {
		_ = writer.WriteField(key, value[0])
	}

	if len(b) != 0 {
		part, err := writer.CreateFormFile("media", "filename.txt")
		if err != nil {
			return nil, "", err
		}
		_, err = io.Copy(part, bytes.NewReader(b))
		if err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	contentType := writer.FormDataContentType()
	return body, contentType, nil
}
