package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

const headerName = "X-Trace-Id"

func AccessLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		h := ctx.ClientIP() // the IP address of the client (remote host)
		// u - the userID that requested the information
		u := "-"
		t := time.Now().String()                               // the time that the request was received
		r := req.Method + " " + req.URL.Path + " " + req.Proto // the client request line, ex: "GET /image.png HTTP/1.0"

		traceID := ctx.Request.Header.Get(headerName)
		if traceID == "" {
			b := make([]byte, 12)
			_, _ = rand.Read(b)
			traceID = hex.EncodeToString(b)
		}
		ctx.Request.Header.Set(headerName, traceID)
		ctx.Writer.Header().Set(headerName, traceID)

		ctx.Next()

		s := ctx.Writer.Status()              // the response status code
		b := ctx.Writer.Size()                // the size of the object returned to the client
		ti := ctx.Request.Header.Get(traceID) // ti - the request trace id

		zap.S().Infof("%s %s %s %s %d %d %s", h, u, t, r, s, b, ti)
	}
}
