package core

import (
	"github.com/project-miko/miko/tools/timer"
)

type Crond interface {
	GetDurationMillisecond() uint32
	Init()
	Worker()
}

func RegisterCrond(c Crond) {

	c.Init()

	timer.DoTimer(c.GetDurationMillisecond(), c.Worker)
}
