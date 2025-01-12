package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/core"
	"github.com/project-miko/miko/tools/log"
	"github.com/project-miko/miko/tools/rpchelper"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	core.BaseController
}

func (ctrl *Middleware) AdminToken(c *gin.Context) {
	// limit to POST
	_m := strings.ToLower(c.Request.Method)
	if _m != "post" {
		ctrl.JsonError(c, conf.ApiCodeMethodNotAllowed, "method not allowed")
		c.Abort()
		return
	}

	token := rpchelper.RequestParameterString(c, "login_token") //c.PostForm("login_token")
	if token == "" {
		lt := new(core.LoginTokenJsonReq)

		err := c.ShouldBindBodyWith(&lt, binding.JSON)
		if err != nil || len(lt.LoginToken) == 0 {
			ctrl.JsonError(c, conf.ApiCodeNoAuth, "wrong token")
			c.Abort()
			return
		}

		token = lt.LoginToken
	}

	tokenInfo, e := core.ParseAdminToken(token)
	if e == conf.ErrLoginTokenExpired {
		ctrl.JsonError(c, conf.ApiCodeTokenExpired, "token expired")
		c.Abort()
		return
	}
	if e != nil {
		log.Error("", "parse token failed %s %s", token, e.Error())
		ctrl.JsonError(c, conf.ApiCodeNoAuth, "wrong token")
		c.Abort()
		return
	}
	isVerify2FA := strings.Contains(c.Request.RequestURI, "verify2faga")
	if tokenInfo.Status != core.TokenVerified && !isVerify2FA {
		ctrl.JsonError(c, conf.ApiCodeTokenNotVerified, "token not verified")
		c.Abort()
		return
	}

	c.Set("admin_token", tokenInfo)
}

func (ctrl *Middleware) LoginToken(c *gin.Context) {
	// limit to POST
	_m := strings.ToLower(c.Request.Method)
	if _m != "post" {
		ctrl.JsonError(c, conf.ApiCodeMethodNotAllowed, "method not allowed")
		c.Abort()
		return
	}

	token := rpchelper.RequestParameterString(c, "login_token") //c.PostForm("login_token")
	if token == "" {
		ctrl.JsonError(c, conf.ApiCodeNoAuth, "wrong token")
		c.Abort()
		return
	}

	uid, ok := rpchelper.RequestParameterInt(c, "uid")
	if !ok {
		ctrl.JsonError(c, conf.ApiCodeAuthWrongUid, "wrong uid")
		c.Abort()
		return
	}

	tokenInfo, e := core.ParseLoginToken(token)
	if e == conf.ErrLoginTokenExpired {
		ctrl.JsonError(c, conf.ApiCodeTokenExpired, "token expired")
		c.Abort()
		return
	}
	if e != nil {
		log.Error("", "parse token failed %s %s", token, e.Error())
		ctrl.JsonError(c, conf.ApiCodeNoAuth, "wrong token")
		c.Abort()
		return
	}

	if tokenInfo.Uid != uid {
		ctrl.JsonError(c, conf.ApiCodeAddressNotEqual, "wrong uid")
		c.Abort()
		return
	}

	c.Set("login_token", tokenInfo)
}
