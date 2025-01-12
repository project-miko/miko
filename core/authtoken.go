package core

import (
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/crypt"
)

const (
	TokenVerified = 1
	TokenUnVerify = 2
)

type LoginTokenJsonReq struct {
	LoginToken string `json:"login_token"`
}

type LoginToken struct {
	Uid       int64  `json:"uid"`
	CreatedAt int64  `json:"created_at"`
	ExpiredAt int64  `json:"expired_at"`
	Nonce     string `json:"nonce"`
}

type AdminToken struct {
	Uid       int64  `json:"uid"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	ExpiredAt int64  `json:"expired_at"`
	Status    int    `json:"status"`
	Nonce     string `json:"nonce"`
}

func CheckPassword(src, dist, salt string) bool {
	pwdCrypted := crypt.Md5(src)
	pwdCrypted = crypt.Md5(fmt.Sprintf("%s.%s", pwdCrypted, salt))

	return pwdCrypted == dist
}

func HashPassword(pwd string) (salt, pwdCrypted string) {
	salt = tools.GetRandStr(6)

	pwdCrypted = crypt.Md5(pwd)
	pwdCrypted = crypt.Md5(fmt.Sprintf("%s.%s", pwdCrypted, salt))

	return salt, pwdCrypted
}

func ParseParamToken(paramToken string) (map[string]string, error) {
	b, e := base64.RawURLEncoding.DecodeString(paramToken)
	if e != nil {
		return nil, e
	}

	deced, e := crypt.DeDes(b, ParamKey)
	if e != nil {
		return nil, e
	}

	deced, e = crypt.PKCS5Unpadding(deced)
	if e != nil {
		return nil, e
	}

	param := make(map[string]string)

	e = json.Unmarshal(deced, &param)
	if e != nil {
		return nil, e
	}

	return param, nil
}

func GenerateParamToken(param map[string]string) (string, error) {

	b, e := json.Marshal(param)
	if e != nil {
		return "", e
	}

	b = crypt.PKCS5Padding(b, des.BlockSize)
	r, e := crypt.EnDes(b, ParamKey)
	if e != nil {
		return "", e
	}

	return base64.RawURLEncoding.EncodeToString(r), nil
}

func ParseAdminToken(loginToken string) (*AdminToken, error) {
	b, e := base64.RawURLEncoding.DecodeString(loginToken)
	if e != nil {
		return nil, e
	}

	deced, e := crypt.DeDes(b, LoginKey)
	if e != nil {
		return nil, e
	}

	deced, e = crypt.PKCS5Unpadding(deced)
	if e != nil {
		return nil, e
	}

	info := new(AdminToken)

	e = json.Unmarshal(deced, info)
	if e != nil {
		return nil, e
	}

	now := time.Now().Unix()
	if now > info.ExpiredAt {
		return nil, conf.ErrLoginTokenExpired
	}

	return info, nil
}

func GenerateAdminToken(uid int64, username string, status int) (string, error) {

	nonce := tools.GetRandStr(6)
	now := time.Now().Unix()
	info := &AdminToken{
		Uid:       uid,
		Username:  username,
		CreatedAt: now,
		ExpiredAt: now + int64(conf.LoginTokenExpired),
		Status:    status,
		Nonce:     nonce,
	}

	b, e := json.Marshal(info)
	if e != nil {
		return "", e
	}

	b = crypt.PKCS5Padding(b, des.BlockSize)
	r, e := crypt.EnDes(b, LoginKey)
	if e != nil {
		return "", e
	}

	return base64.RawURLEncoding.EncodeToString(r), nil
}

func GenerateLoginToken(uid int64, nonce string) (string, error) {
	now := time.Now().Unix()
	info := &LoginToken{
		Uid:       uid,
		CreatedAt: now,
		ExpiredAt: now + int64(conf.LoginTokenExpired),
		Nonce:     nonce,
	}

	b, e := json.Marshal(info)
	if e != nil {
		return "", e
	}

	b = crypt.PKCS5Padding(b, des.BlockSize)
	r, e := crypt.EnDes(b, LoginKey)
	if e != nil {
		return "", e
	}

	return base64.RawURLEncoding.EncodeToString(r), nil
}

func ParseLoginToken(loginToken string) (*LoginToken, error) {
	b, e := base64.RawURLEncoding.DecodeString(loginToken)
	if e != nil {
		return nil, e
	}

	deced, e := crypt.DeDes(b, LoginKey)
	if e != nil {
		return nil, e
	}

	deced, e = crypt.PKCS5Unpadding(deced)
	if e != nil {
		return nil, e
	}

	info := new(LoginToken)

	e = json.Unmarshal(deced, info)
	if e != nil {
		return nil, e
	}

	now := time.Now().Unix()
	if now > info.ExpiredAt {
		return nil, conf.ErrLoginTokenExpired
	}

	return info, nil
}
