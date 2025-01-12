package commands

import (
	"github.com/urfave/cli"
)

var (
	ToolCommands = []cli.Command{
		{
			Name:        "admin-createuser",
			Description: "create admin user",
			Action:      CreateAdminUser,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username",
					Usage: "the username",
				},
				cli.StringFlag{
					Name:  "pwd",
					Usage: "the password",
				},
			},
		},
		{
			Name:        "parse-user-token",
			Description: "parse user token",
			Action:      ParseUserToken,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "login-token",
					Usage: "user login token",
				},
			},
		}}
)
