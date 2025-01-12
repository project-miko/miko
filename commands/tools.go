package commands

import (
	"github.com/project-miko/miko/core"
	"github.com/project-miko/miko/tools/log"

	"github.com/urfave/cli"
)

func ParseUserToken(c *cli.Context) {
	loginToken := c.String("login-token")

	tokenInfo, e := core.ParseLoginToken(loginToken)
	if e != nil {
		log.Error("", "parse login token failed %v", tokenInfo)
		return
	}

	log.Info("", "token info %v", tokenInfo)
}
