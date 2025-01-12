package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project-miko/miko/conf"
)

type ApiRet struct {
	Code int32                  `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

type Controller interface {
	JsonSuccess(c *gin.Context, data map[string]interface{})
	JsonError(c *gin.Context, code int32, msg string)
}

type BaseController struct {
}

func (ctrl *BaseController) JsonSuccessMsg(c *gin.Context) {
	ret := &ApiRet{
		Code: conf.ApiCodeSuccess,
		Msg:  "",
		Data: map[string]interface{}{
			"msg": "success",
		},
	}

	c.JSON(http.StatusOK, ret)
}

func (ctrl *BaseController) JsonSuccess(c *gin.Context, data map[string]interface{}) {
	ret := &ApiRet{
		Code: conf.ApiCodeSuccess,
		Msg:  "",
		Data: data,
	}

	c.JSON(http.StatusOK, ret)
}

func (ctrl *BaseController) JsonError(c *gin.Context, code int32, msg string) {
	ret := &ApiRet{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, ret)
}
