package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// simple per-user rate limiter in-memory (dev/prototype only)
type limiterEntry struct {
	limiter *rate.Limiter
	last    time.Time
}

var (
	userLimiters   = make(map[string]*limiterEntry)
	userLimiterMux sync.Mutex
)

// UserRateLimit limits requests per user_id (from JWT). If no user, skip.
func UserRateLimit(r rate.Limit, burst int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, _ := c.Get("user_id").(string)
			if uid == "" {
				return next(c)
			}
			userLimiterMux.Lock()
			le, ok := userLimiters[uid]
			if !ok {
				le = &limiterEntry{limiter: rate.NewLimiter(r, burst), last: time.Now()}
				userLimiters[uid] = le
			}
			le.last = time.Now()
			userLimiterMux.Unlock()

			if !le.limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "user rate limit exceeded")
			}
			return next(c)
		}
	}
}

