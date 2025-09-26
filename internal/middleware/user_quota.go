package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// inMemoryMonthlyQuota tracks per-identity usage per month.
type inMemoryMonthlyQuota struct {
	mu sync.Mutex
	// key: identity (user_id or anon_id), value: map[YYYY-MM]int
	counts map[string]map[string]int
}

var monthlyQuota = &inMemoryMonthlyQuota{counts: make(map[string]map[string]int)}

// UserQuotaMiddleware enforces per-user monthly quotas for generate endpoints without requiring API keys.
// Identity resolution order: JWT user_id (if set by OptionalJWT) -> anon cookie. If anon cookie missing, it will be created.
// Limits:
// - Free: 3/month. If cookie "shared=yes" then 11/month.
// - Pro plan (cookie plan=pro): 50/month.
// - Lifetime (cookie plan=lifetime): effectively unlimited.
func UserQuotaMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only enforce on relevant endpoints (generate actions under enhanced API)
			path := c.Path()
			if !strings.HasPrefix(path, "/api/v1/enhanced/") || !strings.HasSuffix(path, "/generate") || c.Request().Method != http.MethodPost {
				return next(c)
			}

			// Resolve identity
			identity := ""
			if uid, ok := c.Get("user_id").(string); ok && uid != "" {
				identity = "user:" + uid
			} else {
				// Ensure anon cookie exists
				cookie, err := c.Cookie("anon_id")
				if err != nil || cookie == nil || cookie.Value == "" {
					anon := uuid.New().String()
					newC := &http.Cookie{Name: "anon_id", Value: anon, Path: "/", Expires: time.Now().Add(365 * 24 * time.Hour), HttpOnly: true}
					c.SetCookie(newC)
					identity = "anon:" + anon
				} else {
					identity = "anon:" + cookie.Value
				}
			}

			// Determine plan
			plan := "free"
			if pc, err := c.Cookie("plan"); err == nil && pc != nil && pc.Value != "" {
				plan = strings.ToLower(pc.Value)
			}

			// Determine monthly limit
			limit := 3
			if plan == "pro" {
				limit = 50
			} else if plan == "lifetime" {
				limit = 1_000_000 // effectively unlimited
			} else {
				// free plan bonus if shared promise cookie present
				if sc, err := c.Cookie("shared"); err == nil && sc != nil && strings.ToLower(sc.Value) == "yes" {
					limit = 11
				}
			}

			// Count usage
			now := time.Now().UTC()
			monthKey := now.Format("2006-01")

			monthlyQuota.mu.Lock()
			if _, ok := monthlyQuota.counts[identity]; !ok {
				monthlyQuota.counts[identity] = make(map[string]int)
			}
			used := monthlyQuota.counts[identity][monthKey]
			if used >= limit {
				monthlyQuota.mu.Unlock()
				return echo.NewHTTPError(http.StatusTooManyRequests, "monthly quota exceeded")
			}
			monthlyQuota.counts[identity][monthKey] = used + 1
			monthlyQuota.mu.Unlock()

			return next(c)
		}
	}
}
