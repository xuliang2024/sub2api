package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

const officialPricingCatalogJSON = `{
  "deepseek-flash": {
    "input_cost_per_token": 3e-07,
    "output_cost_per_token": 1.2e-06,
    "cache_creation_input_token_cost": 3e-07,
    "cache_read_input_token_cost": 6e-09,
    "litellm_provider": "deepseek"
  },
  "claude-sonnet-4.5": {
    "input_cost_per_token": 3e-06,
    "output_cost_per_token": 1.5e-05,
    "cache_creation_input_token_cost": 3.75e-06,
    "cache_creation_input_token_cost_above_1hr": 6e-06,
    "cache_read_input_token_cost": 3e-07,
    "litellm_provider": "anthropic"
  },
  "glm-5.3": {
    "input_cost_per_token": 6e-07,
    "output_cost_per_token": 2.2e-06,
    "litellm_provider": "zhipu"
  },
  "gpt-image-1": {
    "output_cost_per_image": 0.04,
    "litellm_provider": "openai"
  }
}`

func newOfficialPricingTestService(t *testing.T) *BillingService {
	t.Helper()
	return NewBillingService(&config.Config{}, newStubPricingServiceFromJSON(t, officialPricingCatalogJSON))
}

func TestGetOfficialPricing_ExactCatalogHit(t *testing.T) {
	svc := newOfficialPricingTestService(t)

	p := svc.GetOfficialPricing("claude-sonnet-4.5")
	require.NotNil(t, p)
	require.InDelta(t, 3e-06, *p.InputPrice, 1e-15)
	require.InDelta(t, 1.5e-05, *p.OutputPrice, 1e-15)
	require.InDelta(t, 3.75e-06, *p.CacheWritePrice, 1e-15)
	require.InDelta(t, 3e-07, *p.CacheReadPrice, 1e-15)
	// 1h 价只在高于 5m 价时给出。
	require.NotNil(t, p.CacheWrite1hPrice)
	require.InDelta(t, 6e-06, *p.CacheWrite1hPrice, 1e-15)
}

func TestGetOfficialPricing_NoCacheWrite1hWhenNotHigher(t *testing.T) {
	svc := newOfficialPricingTestService(t)

	p := svc.GetOfficialPricing("deepseek-flash")
	require.NotNil(t, p)
	require.Nil(t, p.CacheWrite1hPrice)
}

func TestGetOfficialPricing_AliasResolvesRename(t *testing.T) {
	svc := newOfficialPricingTestService(t)

	for _, name := range []string{
		"deepseek-v4.1-flash",
		"deepseek-v4.1-flash-20260910",
		"deepseek-v4-1-flash-260910",
	} {
		p := svc.GetOfficialPricing(name)
		require.NotNilf(t, p, "别名 %s 应命中 deepseek-flash", name)
		require.InDelta(t, 3e-07, *p.InputPrice, 1e-15)
		require.InDelta(t, 1.2e-06, *p.OutputPrice, 1e-15)
	}
}

func TestGetOfficialPricing_UnknownModelStaysNil(t *testing.T) {
	svc := newOfficialPricingTestService(t)

	// 目录只有 glm-5.3，没有 5.5-flash；禁止跨代映射，必须返回 nil。
	require.Nil(t, svc.GetOfficialPricing("glm-5.5-flash"))
	// 目录只有 gpt-5-nano；不许按 "gpt-" 前缀兜底出价格。
	require.Nil(t, svc.GetOfficialPricing("gpt-5.6-nano"))
}

func TestGetOfficialPricing_ImageOnlyEntry(t *testing.T) {
	svc := newOfficialPricingTestService(t)

	p := svc.GetOfficialPricing("gpt-image-1")
	require.NotNil(t, p)
	require.Nil(t, p.InputPrice)
	require.NotNil(t, p.ImageOutputPrice)
	require.InDelta(t, 0.04, *p.ImageOutputPrice, 1e-12)
}

func TestGetOfficialPricing_NilPricingServiceReturnsNil(t *testing.T) {
	svc := NewBillingService(&config.Config{}, nil)
	require.Nil(t, svc.GetOfficialPricing("deepseek-flash"))

	var typedNil *BillingService
	require.Nil(t, typedNil.GetOfficialPricing("deepseek-flash"))
}
