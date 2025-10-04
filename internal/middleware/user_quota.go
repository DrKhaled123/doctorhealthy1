package middleware

import (
	"fmt"
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

// getNextMonthReset returns the timestamp for when the quota resets (next month)
func getNextMonthReset() string {
	now := time.Now().UTC()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	return nextMonth.Format(time.RFC3339)
}

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

			// Resolve identity - fix identity resolution order
			identity := ""
			if uid, ok := c.Get("user_id").(string); ok && uid != "" && strings.TrimSpace(uid) != "" {
				identity = "user:" + strings.TrimSpace(uid)
			} else {
				// Ensure anon cookie exists
				cookie, err := c.Cookie("anon_id")
				if err != nil || cookie == nil || cookie.Value == "" || strings.TrimSpace(cookie.Value) == "" {
					anon := uuid.New().String()
					newC := &http.Cookie{
						Name:     "anon_id",
						Value:    anon,
						Path:     "/",
						Expires:  time.Now().Add(365 * 24 * time.Hour),
						HttpOnly: true,
						Secure:   true,
						SameSite: http.SameSiteLaxMode,
					}
					c.SetCookie(newC)
					identity = "anon:" + anon
				} else {
					identity = "anon:" + strings.TrimSpace(cookie.Value)
				}
			}

			// Determine plan - fix case sensitivity and validation
			plan := "free"
			if pc, err := c.Cookie("plan"); err == nil && pc != nil && strings.TrimSpace(pc.Value) != "" {
				planValue := strings.ToLower(strings.TrimSpace(pc.Value))
				// Validate plan values
				switch planValue {
				case "pro", "lifetime", "free":
					plan = planValue
				default:
					plan = "free" // Default to free for invalid plans
				}
			}

			// Determine monthly limit - fix shared bonus calculation
			limit := 3
			switch plan {
			case "pro":
				limit = 50
			case "lifetime":
				limit = 1_000_000 // effectively unlimited
			default: // free plan
				limit = 3
				// Fix shared bonus calculation - only apply to free plan
				if sc, err := c.Cookie("shared"); err == nil && sc != nil {
					sharedValue := strings.ToLower(strings.TrimSpace(sc.Value))
					if sharedValue == "yes" || sharedValue == "true" {
						limit = 11
					}
				}
			}

			// Count usage with better month key generation
			now := time.Now().UTC()
			monthKey := now.Format("2006-01")

			monthlyQuota.mu.Lock()
			defer monthlyQuota.mu.Unlock()

			if monthlyQuota.counts[identity] == nil {
				monthlyQuota.counts[identity] = make(map[string]int)
			}

			used := monthlyQuota.counts[identity][monthKey]
			if used >= limit {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]interface{}{
					"error":  "monthly quota exceeded",
					"used":   used,
					"limit":  limit,
					"plan":   plan,
					"resets": getNextMonthReset(),
				})
			}

			monthlyQuota.counts[identity][monthKey] = used + 1

			// Add quota info to response headers for debugging
			c.Response().Header().Set("X-Quota-Used", fmt.Sprintf("%d", used+1))
			c.Response().Header().Set("X-Quota-Limit", fmt.Sprintf("%d", limit))
			c.Response().Header().Set("X-Quota-Plan", plan)

			return next(c)
		}
	}
}
