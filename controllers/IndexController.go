package controllers

import (
	"time"

	"github.com/project-miko/miko/tools"

	"github.com/gin-gonic/gin"
	"github.com/project-miko/miko/core"
)

type IndexController struct {
	core.BaseController
}

func (ctrl *IndexController) ServerStatus(c *gin.Context) {
	ctrl.JsonSuccess(c, map[string]interface{}{
		"millisecond": tools.GetMillisecond(time.Now()),
	})
}
