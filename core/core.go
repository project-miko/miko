package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/project-miko/miko/conf"

	"github.com/gin-gonic/gin"
)

var (
	webEngine *gin.Engine
	LoginKey  []byte
	ParamKey  []byte
)

func init() {

	// web server init
	webEngine = gin.Default()
	webEngine.RedirectFixedPath = true
}

func GetEngine() *gin.Engine {
	return webEngine
}

func Run() {
	_loginKey := conf.GetConfigString("app", "login_key")
	if _loginKey == "" {
		panic(fmt.Errorf("login_key did not configured"))
	}
	LoginKey = []byte(_loginKey)

	_paramKey := conf.GetConfigString("app", "param_key")
	if _paramKey == "" {
		panic(fmt.Errorf("param_key did not configured"))
	}
	ParamKey = []byte(_paramKey)

	host := conf.GetConfigString("app", "host")
	port := conf.GetConfigString("app", "port")

	e := webEngine.Run(host + ":" + port)
	if e != nil {
		panic(e)
	}
}

func UseMiddleware(f func(*gin.Context)) {
	webEngine.Use(f)
}

func AutoGroupRoute(ctrl interface{}, group *gin.RouterGroup) {
	autoRoute(ctrl, group)
}

func AutoRoute(ctrl interface{}) {
	autoRoute(ctrl, nil)
}

func autoRoute(ctrl interface{}, group *gin.RouterGroup) {

	autoSupportMethods := []string{"GET", "POST"}

	cType := reflect.TypeOf(&gin.Context{})

	t := reflect.TypeOf(ctrl)
	tv := reflect.ValueOf(ctrl)

	_, ok := ctrl.(Controller)

	if !ok {
		panic(t.String() + " is not a Controller")
	}

	methodTotal := t.NumMethod()
	for i := 0; i < methodTotal; i++ {
		m := t.Method(i)

		if m.Type.NumIn() != 2 || m.Type.NumOut() != 0 {
			continue
		}

		if m.Type.In(1) != cType {
			continue
		}

		pkgName := strings.ReplaceAll(t.String(), "*controllers.", "")
		pkgName = strings.ReplaceAll(pkgName, "Controller", "")

		rStr := fmt.Sprintf("/%s/%s", pkgName, m.Name)
		rStr = strings.ToLower(rStr)

		msg := fmt.Sprintf("will add %s", rStr)
		fmt.Println(msg)

		if group != nil {
			for _, me := range autoSupportMethods {
				group.Handle(me, rStr, func(context *gin.Context) {
					tv.MethodByName(m.Name).Call([]reflect.Value{
						reflect.ValueOf(context),
					})
				})
			}
		} else {
			for _, me := range autoSupportMethods {
				webEngine.Handle(me, rStr, func(context *gin.Context) {
					tv.MethodByName(m.Name).Call([]reflect.Value{
						reflect.ValueOf(context),
					})
				})
			}
		}
	}
}
