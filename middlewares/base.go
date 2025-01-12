package middlewares

import (
	"net/http"

	"github.com/project-miko/miko/conf"

	"github.com/gin-gonic/gin"
)

var (
	allowOrigin = map[string]bool{
		// local develope
		"http://localhost:3000": true,

		// production
	}

	allowOriginUri = map[string]bool{
		// "/index/xxx": true,
	}
)

func Cors(c *gin.Context) {

	if conf.Env == conf.EnvTest {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*") // replace * with your domain
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
		return
	}

	requri := c.Request.RequestURI
	//fmt.Println("req uri", c.Request.RequestURI)
	if _, ok := allowOriginUri[requri]; len(requri) > 0 && ok && allowOriginUri[requri] {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
		return
	}

	origin := c.Request.Header.Get("origin")
	if _, ok := allowOrigin[origin]; len(origin) > 0 && ok && allowOrigin[origin] {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
