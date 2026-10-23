//go:build embed

package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubCampaignHosts map[string]bool

func (s stubCampaignHosts) IsActiveCampaignHost(_ context.Context, host string) bool {
	return s[host]
}

func newShellTestContext(host string, forwardedHost string) *gin.Context {
	c, _ := newShellTestRecorder(host, forwardedHost)
	return c
}

func newShellTestRecorder(host string, forwardedHost string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Host = host
	if forwardedHost != "" {
		c.Request.Header.Set("X-Forwarded-Host", forwardedHost)
	}
	return c, recorder
}

// 域名命中时返回落地页入口，未命中回主站——分流失败绝不能把主站打挂。
func TestShellForCampaignHost(t *testing.T) {
	server, err := NewFrontendServer(&mockSettingsProvider{settings: map[string]string{}})
	if err != nil {
		t.Fatalf("NewFrontendServer: %v", err)
	}
	if !server.HasCampaignShell() {
		t.Fatal("campaign.html 未嵌入：请确认构建带 campaign 入口且 Go 构建带 -tags embed")
	}
	server.SetCampaignHostResolver(stubCampaignHosts{"try.example.com": true})

	cases := []struct {
		name          string
		host          string
		forwardedHost string
		want          string
	}{
		{name: "命中落地页域名", host: "try.example.com", want: "campaign.html"},
		{name: "保留端口", host: "try.example.com:8443", want: "campaign.html"},
		{name: "X-Forwarded-Host 优先", host: "127.0.0.1:8081", forwardedHost: "try.example.com", want: "campaign.html"},
		{name: "未命中回落主站", host: "codex.apiz.ai", want: "index.html"},
		{name: "空 Host 回落主站", host: "", want: "index.html"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := server.shellFor(newShellTestContext(tc.host, tc.forwardedHost)); got != tc.want {
				t.Fatalf("shellFor() = %q, want %q", got, tc.want)
			}
		})
	}
}

// 未注入判定器时一律主站：保证落地页功能异常不影响现有站点。
func TestShellForWithoutResolver(t *testing.T) {
	server, err := NewFrontendServer(&mockSettingsProvider{settings: map[string]string{}})
	if err != nil {
		t.Fatalf("NewFrontendServer: %v", err)
	}
	if got := server.shellFor(newShellTestContext("try.example.com", "")); got != "index.html" {
		t.Fatalf("shellFor() without resolver = %q, want index.html", got)
	}
}

// 两个入口各自独立缓存：落地页请求不会污染主站注入结果。
func TestServeShellHTMLSeparatesCaches(t *testing.T) {
	provider := &mockSettingsProvider{settings: map[string]string{"site_name": "Sub2API"}}
	server, err := NewFrontendServer(provider)
	if err != nil {
		t.Fatalf("NewFrontendServer: %v", err)
	}
	if !server.HasCampaignShell() {
		t.Skip("campaign.html 未嵌入")
	}

	mainCtx, mainRecorder := newShellTestRecorder("codex.apiz.ai", "")
	server.serveShellHTML(mainCtx, "index.html")

	campaignCtx, campaignRecorder := newShellTestRecorder("try.example.com", "")
	server.serveShellHTML(campaignCtx, "campaign.html")

	mainHTML := mainRecorder.Body.String()
	campaignHTML := campaignRecorder.Body.String()
	if mainHTML == "" || campaignHTML == "" {
		t.Fatalf("empty shell body: main=%d bytes campaign=%d bytes", len(mainHTML), len(campaignHTML))
	}
	if mainHTML == campaignHTML {
		t.Fatal("两个入口返回了同一份 HTML，缓存可能被共用")
	}
	if !strings.Contains(campaignHTML, "window.__APP_CONFIG__") {
		t.Errorf("落地页 shell 缺少配置注入")
	}
	if !strings.Contains(campaignHTML, "/assets/campaign-") {
		t.Errorf("落地页 shell 未引用 campaign 入口脚本")
	}
	if provider.called == 0 {
		t.Error("公开配置未被读取")
	}
}
