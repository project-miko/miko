package router

import (
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/controllers"
	"github.com/project-miko/miko/core"
	"github.com/project-miko/miko/middlewares"
)

func Router() {
	core.GetEngine().Use(middlewares.Cors)

	templatesPath := conf.GetConfigString("app", "templates")
	core.GetEngine().LoadHTMLGlob(templatesPath + "/*.tmpl")

	// register routes
	core.AutoRoute(&controllers.IndexController{})
	middlewareInst := new(middlewares.Middleware)

	// /security/**
	securityRouterGroup := core.GetEngine().Group("/security")
	securityRouterGroup.Use(middlewareInst.AdminToken)
}
