package middlewares

import (
	"time"

	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func RateLimiter(rate time.Duration, limit int) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate: rate,
		Limit: uint(limit),
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(ctx *gin.Context, i ratelimit.Info) {
			utils.Error(ctx, 429, "Too many request", nil)
		},
		KeyFunc: func(ctx *gin.Context) string {
			return ctx.ClientIP()
		},
	})
}

