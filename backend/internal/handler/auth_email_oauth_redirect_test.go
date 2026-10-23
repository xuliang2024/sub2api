package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResolveEmailOAuthRedirectURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name       string
		host       string
		forwarded  string
		configured string
		want       string
	}{
		{
			name:       "same host keeps configured",
			host:       "codex.apiz.ai",
			configured: "https://codex.apiz.ai/api/v1/auth/oauth/google/callback",
			want:       "https://codex.apiz.ai/api/v1/auth/oauth/google/callback",
		},
		{
			name:       "different host derives per request",
			host:       "cn.apiz.ai",
			configured: "https://codex.apiz.ai/api/v1/auth/oauth/google/callback",
			want:       "https://cn.apiz.ai/api/v1/auth/oauth/google/callback",
		},
		{
			name:       "empty configured derives",
			host:       "cn.apiz.ai",
			configured: "",
			want:       "https://cn.apiz.ai/api/v1/auth/oauth/google/callback",
		},
		{
			name:      "forwarded host wins",
			host:      "127.0.0.1:8080",
			forwarded: "cn.apiz.ai",
			want:      "https://cn.apiz.ai/api/v1/auth/oauth/google/callback",
		},
		{
			name:       "http when request is not https",
			host:       "local.test",
			configured: "https://codex.apiz.ai/api/v1/auth/oauth/google/callback",
			want:       "http://local.test/api/v1/auth/oauth/google/callback",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+tc.host+"/api/v1/auth/oauth/google/start", nil)
			req.Host = tc.host
			if tc.host != "local.test" {
				req.Header.Set("X-Forwarded-Proto", "https")
			}
			if tc.forwarded != "" {
				req.Header.Set("X-Forwarded-Host", tc.forwarded)
			}
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			if got := resolveEmailOAuthRedirectURL(c, "google", tc.configured); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
