//go:build !embed

// Package web provides embedded web assets for the application.
package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PublicSettingsProvider is an interface to fetch public settings
// This stub is needed for compilation when frontend is not embedded
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (any, error)
}

// CampaignHostResolver 是 embed 构建里同名接口的占位实现。
//
// 非 embed 构建没有前端入口，也就不存在按域名分流，方法全部退化为 no-op，
// 以保证 internal/server 在两个构建标签下都能编译。
type CampaignHostResolver interface {
	IsActiveCampaignHost(ctx context.Context, host string) bool
}

// FrontendServer is a stub for non-embed builds
type FrontendServer struct{}

// SetCampaignHostResolver is a no-op for non-embed builds
func (s *FrontendServer) SetCampaignHostResolver(resolver CampaignHostResolver) {}

// HasCampaignShell 非 embed 构建没有落地页入口
func (s *FrontendServer) HasCampaignShell() bool { return false }

// NewFrontendServer returns an error when frontend is not embedded
func NewFrontendServer(settingsProvider PublicSettingsProvider) (*FrontendServer, error) {
	return nil, errors.New("frontend not embedded")
}

// InvalidateCache is a no-op for non-embed builds
func (s *FrontendServer) InvalidateCache() {}

// Middleware returns a handler that returns 404 for non-embed builds
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
	}
}

func ServeEmbeddedFrontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
	}
}

func HasEmbeddedFrontend() bool {
	return false
}
