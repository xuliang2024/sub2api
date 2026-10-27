package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	gocache "github.com/patrickmn/go-cache"
)

const defaultPartnerUsageCacheTTL = 30 * time.Second

// PartnerUsageResolver 计费高频路径使用的合伙人快照解析器（进程内缓存 + 失败降级）。
type PartnerUsageResolver struct {
	repo  PartnerRepository
	cache *gocache.Cache
	ttl   time.Duration
}

func NewPartnerUsageResolver(repo PartnerRepository, ttl time.Duration) *PartnerUsageResolver {
	if ttl <= 0 {
		ttl = defaultPartnerUsageCacheTTL
	}
	return &PartnerUsageResolver{
		repo:  repo,
		cache: gocache.New(ttl, time.Minute),
		ttl:   ttl,
	}
}

type partnerUsageShareCacheEntry struct {
	share *PartnerUsageShare
}

// Resolve 返回该用户×分组命中的合伙人定价快照；未命中/出错均返回 nil（不影响计费主流程）。
func (r *PartnerUsageResolver) Resolve(ctx context.Context, userID, groupID int64) *PartnerUsageShare {
	if r == nil || r.repo == nil || userID <= 0 || groupID <= 0 {
		return nil
	}
	key := fmt.Sprintf("%d:%d", userID, groupID)
	if cached, ok := r.cache.Get(key); ok {
		if entry, castOK := cached.(partnerUsageShareCacheEntry); castOK {
			return entry.share
		}
	}
	share, err := r.repo.ResolveUsageShare(ctx, userID, groupID)
	if err != nil {
		// 降级：本次不计收益，绝不阻塞扣费。
		logger.LegacyPrintf("service.gateway", "resolve partner usage share failed, skip commission: user=%d group=%d err=%v", userID, groupID, err)
		return nil
	}
	r.cache.Set(key, partnerUsageShareCacheEntry{share: share}, r.ttl)
	return share
}

// ApplyPartnerCommission 把合伙人收益快照写入 usageLog。
// 命中即记录 partner_id/base_rate/rate_floor（commission 可为 0，用于下线消费统计）；
// 收益 = 折前费用 × (基础倍率 − 最低倍率)，高峰因子不参与，负值 clamp 为 0。
func ApplyPartnerCommission(log *UsageLog, share *PartnerUsageShare, totalCost float64) {
	if log == nil || share == nil || share.PartnerUserID <= 0 {
		return
	}
	partnerID := share.PartnerUserID
	baseRate := share.BaseRate
	floor := share.MinRate
	log.PartnerID = &partnerID
	log.PartnerBaseRate = &baseRate
	log.PartnerRateFloor = &floor
	if totalCost > 0 && baseRate > floor {
		commission := roundTo(totalCost*(baseRate-floor), 8)
		if commission > 0 {
			log.PartnerCommission = commission
		}
	}
}

// partnerIndependentMediaBilling 判断该笔是否"图片/视频独立倍率"计费：
// 独立倍率分组下这些消费不按用户级倍率收费，因此不产生合伙人差价收益。
func partnerIndependentMediaBilling(imageCount, videoCount int, group *Group) bool {
	if group == nil {
		return false
	}
	if imageCount > 0 && group.ImageRateIndependent {
		return true
	}
	if videoCount > 0 && group.VideoRateIndependent {
		return true
	}
	return false
}
