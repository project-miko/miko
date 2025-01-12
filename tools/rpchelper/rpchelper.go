package rpchelper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequestParameterExists(c *gin.Context, key string) (string, bool) {
	v, ok := c.GetQuery(key)

	if !ok {
		v, ok = c.GetPostForm(key)
	}

	if !ok {
		return "", false
	}

	return v, true
}

// get string value
func RequestParameterString(c *gin.Context, key string) string {
	v := c.Query(key)

	if v == "" {
		v = c.PostForm(key)
	}

	return v
}

// get integer value
func RequestParameterInt(c *gin.Context, key string) (int64, bool) {
	v := RequestParameterString(c, key)

	if v == "" {
		return 0, false
	}

	i, e := strconv.ParseInt(v, 10, 64)
	if e != nil {
		return 0, false
	}

	return i, true
}

// get float value
func RequestParameterFloat(c *gin.Context, key string) (float64, bool) {
	v := RequestParameterString(c, key)

	if v == "" {
		return 0, false
	}

	i, e := strconv.ParseFloat(v, 64)
	if e != nil {
		return 0, false
	}

	return i, true
}
