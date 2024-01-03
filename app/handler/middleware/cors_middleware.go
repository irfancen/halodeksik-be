package middleware

import (
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/env"
	"halodeksik-be/app/util"
)

var appClients string

func CORSMiddleware(ctx *gin.Context) {
	if util.IsEmptyString(appClients) {
		appClients = env.Get("APP_CLIENT")
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", appClients)
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
}
