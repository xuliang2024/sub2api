//go:build embed

package web

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	htmlpkg "html"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

const (
	// NonceHTMLPlaceholder is the placeholder for nonce in HTML script tags
	NonceHTMLPlaceholder = "__CSP_NONCE_VALUE__"
)

//go:embed all:dist
var frontendFS embed.FS

// PublicSettingsProvider is an interface to fetch public settings
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (any, error)
}

// CampaignHostResolver 判断某个域名是否绑定了已上线的广告落地页。
//
// 定义在 web 包而不是引用 service 包，避免 transport 层反向依赖业务层实现。
// 返回 false（含解析失败）时一律回落到主站 index.html——落地页分流绝不能
// 把主站打挂，因此这里 fail-open 到主站。
type CampaignHostResolver interface {
	IsActiveCampaignHost(ctx context.Context, host string) bool
}

// FrontendServer serves the embedded frontend with settings injection
type FrontendServer struct {
	distFS      fs.FS
	fileServer  http.Handler
	baseHTML    []byte
	cache       *HTMLCache
	settings    PublicSettingsProvider
	overrideDir string // local file override directory

	// campaignHTML 第二个前端入口（广告落地页）。为空表示该构建未包含 campaign.html，
	// 此时所有域名都回落到主站 index.html。
	campaignHTML  []byte
	campaignCache *HTMLCache
	campaignHosts CampaignHostResolver
}

// NewFrontendServer creates a new frontend server with settings injection
func NewFrontendServer(settingsProvider PublicSettingsProvider) (*FrontendServer, error) {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		return nil, err
	}

	// Read base HTML once
	file, err := distFS.Open("index.html")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	baseHTML, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cache := NewHTMLCache()
	cache.SetBaseHTML(baseHTML)

	server := &FrontendServer{
		distFS:      distFS,
		fileServer:  http.FileServer(http.FS(distFS)),
		baseHTML:    baseHTML,
		cache:       cache,
		settings:    settingsProvider,
		overrideDir: filepath.Join("data", "public"),
	}
	server.loadCampaignShell()
	return server, nil
}

// loadCampaignShell 读取可选的 campaign.html 第二入口。
//
// 构建未包含该入口（例如只产出了主站）时静默跳过：落地页分流是附加能力，
// 不能因为它让主站起不来。
func (s *FrontendServer) loadCampaignShell() {
	if s == nil || s.distFS == nil {
		return
	}
	file, err := s.distFS.Open("campaign.html")
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil || len(content) == 0 {
		return
	}
	campaignCache := NewHTMLCache()
	campaignCache.SetBaseHTML(content)
	s.campaignHTML = content
	s.campaignCache = campaignCache
}

// SetCampaignHostResolver 注入域名 -> 落地页的判定器。
//
// 必须在服务启动后、对外提供请求前调用；未注入时所有域名都走主站入口。
func (s *FrontendServer) SetCampaignHostResolver(resolver CampaignHostResolver) {
	if s == nil {
		return
	}
	s.campaignHosts = resolver
}

// HasCampaignShell 表示本次构建是否包含落地页入口。
func (s *FrontendServer) HasCampaignShell() bool {
	return s != nil && len(s.campaignHTML) > 0
}

// InvalidateCache invalidates the HTML cache (call when settings change)
func (s *FrontendServer) InvalidateCache() {
	if s != nil && s.cache != nil {
		s.cache.Invalidate()
	}
}

// Middleware returns the Gin middleware handler
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if shouldBypassEmbeddedFrontend(path) {
			c.Next()
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// For index.html / campaign.html / SPA routes, serve the shell with injected settings
		if cleanPath == "index.html" || cleanPath == "campaign.html" || !s.fileExists(cleanPath) {
			s.serveShellHTML(c, s.shellFor(c))
			return
		}

		// Try local override first
		if s.tryServeOverride(c, cleanPath) {
			return
		}

		// Serve static files normally (hashed assets get long-lived cache headers)
		applyStaticAssetCacheHeaders(c.Writer.Header(), cleanPath)
		s.fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func (s *FrontendServer) fileExists(path string) bool {
	file, err := s.distFS.Open(path)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

// shellFor 决定本次请求使用哪个前端入口。
//
// 判定全部在服务端完成（而不是让前端按 location.host 自己跳），原因：
//   - 落地页域名打开时首屏必须是落地页，不能先加载主站 SPA 再跳转；
//   - 主站与落地页是两个独立的 Vite entry，包体互不影响。
//
// 任何异常（未构建 campaign.html、解析器未注入、查询失败）都回落到 index.html。
func (s *FrontendServer) shellFor(c *gin.Context) string {
	if len(s.campaignHTML) == 0 || s.campaignHosts == nil {
		return "index.html"
	}
	host := requestHostForCampaign(c)
	if host == "" {
		return "index.html"
	}
	if !s.campaignHosts.IsActiveCampaignHost(c.Request.Context(), host) {
		return "index.html"
	}
	return "campaign.html"
}

// requestHostForCampaign 取请求域名（优先 X-Forwarded-Host，去掉端口）。
func requestHostForCampaign(c *gin.Context) string {
	raw := strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
	if raw != "" {
		if idx := strings.Index(raw, ","); idx >= 0 {
			raw = raw[:idx]
		}
	} else {
		raw = c.Request.Host
	}
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		raw = host
	}
	raw = strings.TrimSuffix(raw, ".")
	return strings.TrimPrefix(strings.TrimSuffix(raw, "]"), "[")
}

// tryServeOverride checks if a local override file exists and serves it.
// Files in overrideDir take precedence over embedded files.
func (s *FrontendServer) tryServeOverride(c *gin.Context, cleanPath string) bool {
	if s.overrideDir == "" {
		return false
	}
	filePath := filepath.Join(s.overrideDir, filepath.Clean("/"+cleanPath))
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	c.File(filePath)
	c.Abort()
	return true
}

// injectSettings 保留为单测入口：按主站 shell 注入公开配置。
//
// 业务请求走 injectSettingsInto（需要区分主站/落地页），此包装只用于既有测试。
func (s *FrontendServer) injectSettings(settingsJSON []byte) []byte {
	return injectSettingsInto(s.baseHTML, settingsJSON, true)
}

// serveIndexHTML 保留为单测入口：按主站 shell 渲染 index.html。
func (s *FrontendServer) serveIndexHTML(c *gin.Context) {
	s.serveShellHTML(c, "index.html")
}

func (s *FrontendServer) serveShellHTML(c *gin.Context, shell string) {
	baseHTML, cache := s.baseHTML, s.cache
	if shell == "campaign.html" && len(s.campaignHTML) > 0 && s.campaignCache != nil {
		baseHTML, cache = s.campaignHTML, s.campaignCache
	}

	// Get nonce from context (generated by SecurityHeaders middleware)
	nonce := middleware.GetNonceFromContext(c)

	// Check cache first
	cached := cache.Get()
	if cached != nil {
		// Check If-None-Match for 304 response
		if match := c.GetHeader("If-None-Match"); match == cached.ETag {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}

		// Replace nonce placeholder with actual nonce before serving
		content := replaceNoncePlaceholder(cached.Content, nonce)

		c.Header("ETag", cached.ETag)
		c.Header("Cache-Control", "no-cache") // Must revalidate
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
		c.Abort()
		return
	}

	// Cache miss - fetch settings and render
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	settings, err := s.settings.GetPublicSettingsForInjection(ctx)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", baseHTML)
		c.Abort()
		return
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", baseHTML)
		c.Abort()
		return
	}

	rendered := injectSettingsInto(baseHTML, settingsJSON, shell != "campaign.html")
	cache.Set(rendered, settingsJSON)

	// Replace nonce placeholder with actual nonce before serving
	content := replaceNoncePlaceholder(rendered, nonce)

	// 法务页给爬虫与审核一个像样的标题：页面正文是客户端渲染的，
	// 这里只改 <title>，不动缓存（缓存的是整壳，标题按路径后处理）。
	if title := legalPageTitle(c.Request.URL.Path); title != "" {
		content = replaceDocumentTitle(content, title)
	}

	cached = cache.Get()
	if cached != nil {
		c.Header("ETag", cached.ETag)
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

// injectSettingsInto 把公开配置与品牌信息注入任意前端 shell。
//
// 抽成包级函数是因为现在有两个入口（index.html / campaign.html），注入逻辑一致，
// 差别只在 base HTML。applySiteTitle 对落地页为 false：落地页的标题由活动自己的
// SEO 配置在客户端设置，服务端不该把它改成主站的「站点名 - AI API Gateway」。
func injectSettingsInto(baseHTML, settingsJSON []byte, applySiteTitle bool) []byte {
	// Create the script tag to inject with nonce placeholder
	// The placeholder will be replaced with actual nonce at request time
	script := []byte(`<script nonce="` + NonceHTMLPlaceholder + `">window.__APP_CONFIG__=` + string(settingsJSON) + `;</script>`)

	// Inject before </head>
	headClose := []byte("</head>")
	result := bytes.Replace(baseHTML, headClose, append(script, headClose...), 1)

	// Apply custom branding before the browser paints the static defaults.
	if applySiteTitle {
		result = injectSiteTitle(result, settingsJSON)
	}
	result = injectSiteFavicon(result, settingsJSON)

	return result
}

// legalPageTitle 返回法务页的服务端标题；非该路径返回空串。
func legalPageTitle(path string) string {
	switch strings.TrimRight(path, "/") {
	case "/privacy":
		return "Privacy Policy"
	case "/terms":
		return "Terms of Service"
	default:
		return ""
	}
}

// replaceDocumentTitle 把 HTML 里第一个 <title>…</title> 换成给定标题。
func replaceDocumentTitle(html []byte, title string) []byte {
	start := bytes.Index(html, []byte("<title>"))
	if start == -1 {
		return html
	}
	endOffset := bytes.Index(html[start:], []byte("</title>"))
	if endOffset == -1 {
		return html
	}
	end := start + endOffset + len("</title>")
	replacement := []byte("<title>" + htmlpkg.EscapeString(title) + "</title>")

	var buf bytes.Buffer
	buf.Grow(len(html) - (end - start) + len(replacement))
	buf.Write(html[:start])
	buf.Write(replacement)
	buf.Write(html[end:])
	return buf.Bytes()
}

// injectSiteFavicon replaces the static favicon with a configured, browser-safe image URL.
func injectSiteFavicon(html, settingsJSON []byte) []byte {
	var cfg struct {
		SiteLogo string `json:"site_logo"`
	}
	if err := json.Unmarshal(settingsJSON, &cfg); err != nil {
		return html
	}

	logoURL := safeImageURL(cfg.SiteLogo)
	if logoURL == "" {
		return html
	}

	linkStart := bytes.Index(html, []byte(`<link rel="icon"`))
	if linkStart == -1 {
		return html
	}
	linkEndOffset := bytes.IndexByte(html[linkStart:], '>')
	if linkEndOffset == -1 {
		return html
	}
	linkEnd := linkStart + linkEndOffset + 1
	replacement := []byte(`<link rel="icon" href="` + htmlpkg.EscapeString(logoURL) + `" />`)

	var buf bytes.Buffer
	buf.Write(html[:linkStart])
	buf.Write(replacement)
	buf.Write(html[linkEnd:])
	return buf.Bytes()
}

func safeImageURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "/") && !strings.HasPrefix(trimmed, "//") {
		return trimmed
	}
	if strings.HasPrefix(strings.ToLower(trimmed), "data:image/") {
		return trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return ""
	}
	return trimmed
}

// injectSiteTitle replaces the static <title> in HTML with the configured site name.
// This ensures the browser tab shows the correct title before JS executes.
func injectSiteTitle(html, settingsJSON []byte) []byte {
	var cfg struct {
		SiteName string `json:"site_name"`
	}
	if err := json.Unmarshal(settingsJSON, &cfg); err != nil || cfg.SiteName == "" {
		return html
	}

	// Find and replace the existing <title>...</title>
	titleStart := bytes.Index(html, []byte("<title>"))
	titleEnd := bytes.Index(html, []byte("</title>"))
	if titleStart == -1 || titleEnd == -1 || titleEnd <= titleStart {
		return html
	}

	newTitle := []byte("<title>" + htmlpkg.EscapeString(cfg.SiteName) + " - AI API Gateway</title>")
	var buf bytes.Buffer
	buf.Write(html[:titleStart])
	buf.Write(newTitle)
	buf.Write(html[titleEnd+len("</title>"):])
	return buf.Bytes()
}

// replaceNoncePlaceholder replaces the nonce placeholder with actual nonce value
func replaceNoncePlaceholder(html []byte, nonce string) []byte {
	return bytes.ReplaceAll(html, []byte(NonceHTMLPlaceholder), []byte(nonce))
}

// ServeEmbeddedFrontend returns a middleware for serving embedded frontend
// This is the legacy function for backward compatibility when no settings provider is available
func ServeEmbeddedFrontend() gin.HandlerFunc {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic("failed to get dist subdirectory: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))
	overrideDir := filepath.Join("data", "public")

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if shouldBypassEmbeddedFrontend(path) {
			c.Next()
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		if file, err := distFS.Open(cleanPath); err == nil {
			_ = file.Close()
			// Try local override first
			if tryServeOverrideFile(c, overrideDir, cleanPath) {
				return
			}
			applyStaticAssetCacheHeaders(c.Writer.Header(), cleanPath)
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		serveIndexHTML(c, distFS)
	}
}

// tryServeOverrideFile is a standalone version of tryServeOverride for legacy usage.
func tryServeOverrideFile(c *gin.Context, overrideDir, cleanPath string) bool {
	if overrideDir == "" {
		return false
	}
	filePath := filepath.Join(overrideDir, filepath.Clean("/"+cleanPath))
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	c.File(filePath)
	c.Abort()
	return true
}

func shouldBypassEmbeddedFrontend(path string) bool {
	trimmed := strings.TrimSpace(path)
	return strings.HasPrefix(trimmed, "/api/") ||
		strings.HasPrefix(trimmed, "/v1/") ||
		strings.HasPrefix(trimmed, "/v1beta/") ||
		strings.HasPrefix(trimmed, "/backend-api/") ||
		strings.HasPrefix(trimmed, "/antigravity/") ||
		strings.HasPrefix(trimmed, "/setup/") ||
		trimmed == "/health" ||
		trimmed == "/models" ||
		strings.HasPrefix(trimmed, "/models/") ||
		trimmed == "/responses" ||
		strings.HasPrefix(trimmed, "/responses/") ||
		trimmed == "/chat/completions" ||
		trimmed == "/embeddings" ||
		trimmed == "/messages/count_tokens" ||
		trimmed == "/alpha/search" ||
		trimmed == "/web_search" ||
		trimmed == "/x_search" ||
		trimmed == "/tts" ||
		trimmed == "/stt" ||
		trimmed == "/realtime" ||
		trimmed == "/custom-voices" ||
		strings.HasPrefix(trimmed, "/custom-voices/") ||
		trimmed == "/contents/generations/tasks" ||
		strings.HasPrefix(trimmed, "/contents/generations/tasks/") ||
		trimmed == "/v3/contents/generations/tasks" ||
		strings.HasPrefix(trimmed, "/v3/contents/generations/tasks/") ||
		strings.HasPrefix(trimmed, "/images/") ||
		trimmed == "/videos" ||
		strings.HasPrefix(trimmed, "/videos/")
}

func serveIndexHTML(c *gin.Context, fsys fs.FS) {
	file, err := fsys.Open("index.html")
	if err != nil {
		c.String(http.StatusNotFound, "Frontend not found")
		c.Abort()
		return
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read index.html")
		c.Abort()
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

func HasEmbeddedFrontend() bool {
	_, err := frontendFS.ReadFile("dist/index.html")
	return err == nil
}
