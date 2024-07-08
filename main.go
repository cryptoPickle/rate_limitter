package main

import (
	"net/http"
	"time"

	ratelimitter "github.com/cryptoPickle/rate_limitter/pkg/rate_limitter"
	"github.com/cryptoPickle/rate_limitter/pkg/util"
	"github.com/cryptoPickle/rate_limitter/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*gin.Engine
	RateLimit types.IRateLimitter
}

func NewServer(r types.IRateLimitter) *Server {
	s := gin.Default()
	return &Server{
		Engine:    s,
		RateLimit: r,
	}
}

func (s *Server) Router() {
	s.GET("/test", RateLimitterMiddleWare(s.RateLimit),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})
}

func RateLimitterMiddleWare(rl types.IRateLimitter) func(*gin.Context) {
	return func(ctx *gin.Context) {
		stop := util.Took("RateLimitterMiddleware")
		defer stop()
		clientIp := ctx.ClientIP()
		ctx.Set("ip", clientIp)
		if err := rl.Start(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			return
		}
	}
}

func main() {
	rl := ratelimitter.New(10, time.Second*1)
	s := NewServer(rl)
	s.Router()
	if err := s.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
}
