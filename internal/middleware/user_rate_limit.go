package middleware

import (
	"net/http"
	"strings"
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
			uid = strings.TrimSpace(uid)
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

			// Clean up old limiters periodically (basic memory management)
			if len(userLimiters) > 10000 { // arbitrary limit to prevent memory leaks
				for id, entry := range userLimiters {
					if time.Since(entry.last) > time.Hour { // remove limiters unused for 1 hour
						delete(userLimiters, id)
					}
				}
			}
			userLimiterMux.Unlock()

			if !le.limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "user rate limit exceeded",
					"user_id":     uid,
					"retry_after": "1s", // suggest retry time
				})
			}
			return next(c)
		}
	}
}
