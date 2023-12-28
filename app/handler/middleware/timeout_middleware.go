package middleware

import (
	"context"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/env"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var duration time.Duration

func initDuration() {
	seconds, err := strconv.Atoi(env.Get("REQUEST_TIMEOUT"))
	if err != nil {
		seconds = appconstant.DefaultRequestTimeout
	}
	duration = time.Duration(seconds * 1e9)
}

func TimeoutHandler(c *gin.Context) {
	if duration == time.Duration(0) {
		initDuration()
	}
	ctx := c.Request.Context()
	ctx2, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	c.Request = c.Request.WithContext(ctx2)
	c.Next()
}
