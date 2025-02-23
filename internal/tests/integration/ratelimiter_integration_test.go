package middleware_test

import (
	"net/http"
	"net/http/httptest"
	config "server/config"
	util "server/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRateLimiterWithinLimit(t *testing.T) {
	configValue := config.LoadConfig()
	limit := configValue.Server.RequestLimitPerMin
	window := configValue.Server.RequestWindowSec

	w, r, rateLimiter := util.SetupRateLimiter(limit, window)

	for i := int64(0); i < limit; i++ {
		rateLimiter.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimiterExceedsLimit(t *testing.T) {
	configValue := config.LoadConfig()
	limit := configValue.Server.RequestLimitPerMin
	window := configValue.Server.RequestWindowSec

	w, r, rateLimiter := util.SetupRateLimiter(limit, window)

	for i := int64(0); i < limit; i++ {
		rateLimiter.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	rateLimiter.ServeHTTP(w, r)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
	require.Contains(t, w.Body.String(), "429 - Too Many Requests")
}

func TestRateLimitHeadersSet(t *testing.T) {
	configValue := config.LoadConfig()
	limit := configValue.Server.RequestLimitPerMin
	window := configValue.Server.RequestWindowSec

	w, r, rateLimiter := util.SetupRateLimiter(limit, window)

	rateLimiter.ServeHTTP(w, r)

	require.NotEmpty(t, w.Header().Get("x-ratelimit-limit"), "x-ratelimit-limit should be set")
	require.NotEmpty(t, w.Header().Get("x-ratelimit-remaining"), "x-ratelimit-remaining should be set")
	require.NotEmpty(t, w.Header().Get("x-ratelimit-used"), "x-ratelimit-used should be set")
	require.NotEmpty(t, w.Header().Get("x-ratelimit-reset"), "x-ratelimit-reset should be set")
}
