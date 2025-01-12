package conf

import "fmt"

// error codes

const (

	// API error codes
	ApiCodeSuccess          int32 = 200
	ApiCodeBadRequest       int32 = 400
	ApiCodeNotFound         int32 = 404
	ApiCodeParamErr         int32 = 401
	ApiCodeNoAuth           int32 = 403
	ApiCodeMethodNotAllowed int32 = 405 // Method Not Allowed
	ApiCodeErrMsg           int32 = 500

	// Ext error codes
	ApiCodeAuthWrongUid int32 = 40301 // The uid in the token does not match to the uid that user commited.

	// frontend
	ApiCodeTokenExpired            int32 = 10001
	ApiCodeTokenNotVerified        int32 = 10002
	ApiCode2FAVerifyMaxCountExceed int32 = 10003
	ApiCode2FATempSecretExpired    int32 = 10004
	ApiCodeAddressNotEqual         int32 = 10007

	LoginTokenExpired int32 = 604800 // second
)

var (
	ErrLoginTokenExpired = fmt.Errorf("login token expired")
	ErrIndexOutOfBound   = fmt.Errorf("index out of bound")

	ErrRecordNotFound      = fmt.Errorf("record not found")
	ErrRecordAlreadyExists = fmt.Errorf("record already exists")

	ErrNoAuthForTwitterGetMediaStatus = fmt.Errorf("unable to retrieve uploaded media status from Twitter API due to missing authentication")
)

type ErrSystem struct {
	Msg string
}

func (e *ErrSystem) Error() string {
	return fmt.Sprintf("msg:%s", e.Msg)
}

type ErrLogic struct {
	Code int32
	Msg  string
}

func (e *ErrLogic) Error() string {
	return fmt.Sprintf("code:%d, msg:%s", e.Code, e.Msg)
}
