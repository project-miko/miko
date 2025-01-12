package commands

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/project-miko/miko/core"
	"github.com/project-miko/miko/models"
	"github.com/project-miko/miko/tools"
	"github.com/project-miko/miko/tools/log"
	"github.com/urfave/cli"
)

func CreateAdminUser(c *cli.Context) {

	username := c.String("username")
	pwd := c.String("pwd")

	salt, pwdCrypted := core.HashPassword(pwd)

	now := time.Now()
	admin := &models.Admin{
		Username:    username,
		Userpwd:     pwdCrypted,
		Salt:        salt,
		Enable2FAGA: models.Enable2FAGA,
		CreatedAt:   tools.GetMillisecond(now),
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TwitterTools",
		AccountName: username,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		log.Error("", "totp.Generate error %s", err.Error())
		return
	}

	totpOpts := totp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}

	secret := key.Secret()
	passcode, err := totp.GenerateCodeCustom(secret, now, totpOpts)
	if err != nil {
		log.Error("", "totp.GenerateCodeCustom error %s", err.Error())
		return
	}

	valid, err := totp.ValidateCustom(passcode, secret, now, totpOpts)
	if err != nil {
		log.Error("", "totp.ValidateCustom error %s", err.Error())
		return
	}
	if !valid {
		log.Error("", "verify passcode failed")
		return
	}

	if err = admin.Save(); err != nil {
		log.Error("", "create admin user failed %s", err.Error())
		return
	}

	admin2FAGA := &models.Admin2FAGA{
		AdminId:           admin.Id,
		Secret:            secret,
		RecoverCode:       "",
		RecoverCodeStatus: 0,
		CreatedAt:         tools.GetMillisecond(now),
	}

	if err = admin2FAGA.Save(); err != nil {
		log.Error("", "admin2FAGA.Save() error %s", err.Error())
		return
	}

	log.Info("", "create admin done")
}
