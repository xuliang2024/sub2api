package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 广告投放落地页（campaign）相关错误。
var (
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrCampaignNotActive   = errors.New("campaign is not active")
	ErrCampaignNoGroup     = errors.New("campaign has no bound group")
	ErrCampaignSlugInvalid = errors.New("campaign slug is invalid")
	ErrCampaignNameEmpty   = errors.New("campaign name is required")
	ErrCampaignStatusBad   = errors.New("campaign status is invalid")
	ErrCampaignHostTaken   = errors.New("campaign host is already bound by another active campaign")
	ErrCampaignSlugTaken   = errors.New("campaign slug is already in use")
)

// CampaignStatus 落地页状态。
//
//	draft    仅管理员可见（预览需带管理员身份）
//	active   对外可访问，并占用绑定的域名
//	archived 下线但保留配置
type CampaignStatus string

const (
	CampaignStatusDraft    CampaignStatus = "draft"
	CampaignStatusActive   CampaignStatus = "active"
	CampaignStatusArchived CampaignStatus = "archived"
)

const (
	// campaignGroupModelsTTL 分组模型与价格的进程内缓存时长。
	// 落地页是投放入口，QPS 远高于后台；而 ListGroups 会遍历全部渠道并为每个
	// 模型解析计费阶梯，成本不低，因此按分组维度缓存。
	campaignGroupModelsTTL = 60 * time.Second
	// campaignMaxModelEntries 主推模型条数上限，防止配置把首屏撑爆。
	campaignMaxModelEntries = 24
	// campaignMaxCodeSamples 接入示例条数上限。
	campaignMaxCodeSamples = 12
	// campaignMaxContentLen 使用说明长度上限（字节），与自定义页面 1MB 上限保持一致。
	campaignMaxContentLen = 1 << 20
)

// CampaignBrand 落地页外观配置。
type CampaignBrand struct {
	SiteName string `json:"site_name,omitempty"`
	Logo     string `json:"logo,omitempty"`
	// Primary/Accent 为主色与强调色（#RRGGBB）。为空时前端用内置默认值。
	Primary  string `json:"primary,omitempty"`
	Accent   string `json:"accent,omitempty"`
	Footer   string `json:"footer,omitempty"`
	HideAuth bool   `json:"hide_auth,omitempty"`
	// ConsoleURL 主站控制台地址；配置后落地页导航展示「控制台」入口。
	ConsoleURL string `json:"console_url,omitempty"`
	// GroupLabel 分组对外展示名。分组名多为内部命名（如「国产全模型」），
	// 留空时落地页不展示分组，避免中文内部标识出现在海外页面上。
	GroupLabel string `json:"group_label,omitempty"`
}

// CampaignHero 首屏配置。
type CampaignHero struct {
	Eyebrow    string   `json:"eyebrow,omitempty"`
	Title      string   `json:"title,omitempty"`
	Subtitle   string   `json:"subtitle,omitempty"`
	Highlights []string `json:"highlights,omitempty"`
	Badges     []string `json:"badges,omitempty"`
	Image      string   `json:"image,omitempty"`
}

// CampaignModelEntry 主推模型条目；为空清单时展示绑定分组下的全部模型。
type CampaignModelEntry struct {
	Model       string `json:"model"`
	DisplayName string `json:"display_name,omitempty"`
	Tagline     string `json:"tagline,omitempty"`
	Badge       string `json:"badge,omitempty"`
	Featured    bool   `json:"featured,omitempty"`
	Order       int    `json:"order,omitempty"`
	// OfficialInputPrice/OfficialOutputPrice 覆盖官方价（USD/token）。
	//
	// 平台展示的官方价默认取「计费口径」的刊例价，但个别模型两者不等价：
	// 例如 deepseek-v4-pro 自 2026-09-14 起被上游按 Flash 价路由计费，
	// 计费口径会返回 Flash 的价，比我们售价还低，页面就不显示折扣了。
	// 需要按厂商公布价展示时在这里显式覆盖。
	OfficialInputPrice  *float64 `json:"official_input_price,omitempty"`
	OfficialOutputPrice *float64 `json:"official_output_price,omitempty"`
}

// CampaignCodeSample 接入示例；Code 中支持 {{BASE_URL}} / {{API_KEY}} / {{MODEL}} 占位。
type CampaignCodeSample struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Lang  string `json:"lang,omitempty"`
	Code  string `json:"code"`
}

// CampaignCTA 转化动作配置。
type CampaignCTA struct {
	// Mode: key（一键建 Key，默认）| register（只注册）| topup（引导充值）
	Mode     string `json:"mode,omitempty"`
	Label    string `json:"label,omitempty"`
	Sublabel string `json:"sublabel,omitempty"`
	TopupURL string `json:"topup_url,omitempty"`
}

// CampaignTrial 注册赠送策略。
type CampaignTrial struct {
	// Balance 注册时赠送的余额（USD）。<=0 表示不赠送。
	Balance float64 `json:"balance,omitempty"`
	Note    string  `json:"note,omitempty"`
}

// CampaignSEO 页面元信息。
type CampaignSEO struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	OGImage     string   `json:"og_image,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

// Campaign 落地页配置实体。
// CampaignLocaleContent 某一语言的正文覆盖。
//
// 只覆盖运营内容（标题/副标题/使用说明）；模型名与价格是数据，不在这里翻译。
// 未配置的字段回落活动主文案。
type CampaignLocaleContent struct {
	HeroTitle    string `json:"hero_title,omitempty"`
	HeroSubtitle string `json:"hero_subtitle,omitempty"`
	ContentMD    string `json:"content_md,omitempty"`
	// CTALabel 主按钮文案（营销动作里最关键的一句，必须能跟着语言走）。
	CTALabel string `json:"cta_label,omitempty"`
}

// campaignSupportedLocales 前端内置词表支持的语言。
var campaignSupportedLocales = map[string]struct{}{
	"en": {}, "zh": {}, "ja": {}, "ko": {}, "th": {},
}

type Campaign struct {
	ID      int64          `json:"id"`
	Slug    string         `json:"slug"`
	Name    string         `json:"name"`
	Status  CampaignStatus `json:"status"`
	Host    string         `json:"host"`
	GroupID *int64         `json:"group_id"`
	// Locale 界面语言：en（默认，海外投放）/ zh / ja / ko / th。
	Locale       string               `json:"locale"`
	Brand        CampaignBrand        `json:"brand"`
	Hero         CampaignHero         `json:"hero"`
	ContentMD    string               `json:"content_md"`
	ModelEntries []CampaignModelEntry `json:"model_entries"`
	CodeSamples  []CampaignCodeSample `json:"code_samples"`
	CTA          CampaignCTA          `json:"cta"`
	Trial        CampaignTrial        `json:"trial"`
	SEO          CampaignSEO          `json:"seo"`
	// HiddenModels 不在落地页展示的模型名（别名/退役模型）。
	HiddenModels []string `json:"hidden_models"`
	// Translations 按语言的正文覆盖，key 为语言代码。
	Translations map[string]CampaignLocaleContent `json:"translations"`
	CreatedBy    *int64                           `json:"created_by"`
	CreatedAt    time.Time                        `json:"created_at"`
	UpdatedAt    time.Time                        `json:"updated_at"`
}

// CampaignGroupView 落地页展示的分组摘要。
type CampaignGroupView struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Platform       string  `json:"platform"`
	RateMultiplier float64 `json:"rate_multiplier"`
	// PeakRateMultiplier 高峰时段倍率（>1 时生效）；用于页面提示。
	PeakRateMultiplier float64 `json:"peak_rate_multiplier,omitempty"`
	PeakStart          string  `json:"peak_start,omitempty"`
	PeakEnd            string  `json:"peak_end,omitempty"`
}

// CampaignTimePricingPeriod 分时倍率时段，供落地页展示「峰时 ×N」。
type CampaignTimePricingPeriod struct {
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	Multiplier float64 `json:"multiplier"`
}

// CampaignModelView 落地页模型条目（合并配置与真实分组定价）。
type CampaignModelView struct {
	Model       string `json:"model"`
	DisplayName string `json:"display_name"`
	Tagline     string `json:"tagline,omitempty"`
	Badge       string `json:"badge,omitempty"`
	Featured    bool   `json:"featured,omitempty"`
	Platform    string `json:"platform,omitempty"`
	// Available = 该模型在绑定分组的渠道定价里存在（不代表上游账号此刻可用）。
	Available bool `json:"available"`
	// 以下价格单位为「美元 / token」，与 /api/v1/model-plaza 口径一致，前端自行 ×1e6 展示。
	InputPrice      *float64 `json:"input_price,omitempty"`
	OutputPrice     *float64 `json:"output_price,omitempty"`
	CacheReadPrice  *float64 `json:"cache_read_price,omitempty"`
	CacheWritePrice *float64 `json:"cache_write_price,omitempty"`
	PerRequestPrice *float64 `json:"per_request_price,omitempty"`
	// OfficialInputPrice/OfficialOutputPrice 上游官方价（USD/token），用于展示折扣力度。
	// 官方价缺失时为空，前端不展示对比。
	OfficialInputPrice  *float64 `json:"official_input_price,omitempty"`
	OfficialOutputPrice *float64 `json:"official_output_price,omitempty"`
	// TimePricing 命中的分时倍率时段（仅倍率 ≠ 1 时给出）。
	TimePricing []CampaignTimePricingPeriod `json:"time_pricing,omitempty"`
	// TimePricingTimezone 分时倍率时区。
	TimePricingTimezone string `json:"time_pricing_timezone,omitempty"`
}

// CampaignView 对外公开的落地页数据。
type CampaignView struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	// Locale 界面语言（en/zh/ja/ko/th），前端据此选择内置文案并作为默认语言。
	Locale    string              `json:"locale"`
	Brand     CampaignBrand       `json:"brand"`
	Hero      CampaignHero        `json:"hero"`
	ContentMD string              `json:"content_md"`
	Models    []CampaignModelView `json:"models"`
	// Catalog 绑定分组下的**全部**模型（含未被 ModelEntries 选中的）。
	// 卡片区只展示运营挑选的主推模型，广场用这份全量价目表。
	Catalog     []CampaignModelView  `json:"catalog,omitempty"`
	CodeSamples []CampaignCodeSample `json:"code_samples"`
	CTA         CampaignCTA          `json:"cta"`
	Trial       CampaignTrial        `json:"trial"`
	SEO         CampaignSEO          `json:"seo"`
	// Translations 按语言的正文覆盖；缺失的语言回落 Hero / ContentMD。
	Translations map[string]CampaignLocaleContent `json:"translations,omitempty"`
	Group        *CampaignGroupView               `json:"group,omitempty"`
	// APIBaseURL 用户调用网关时应填写的 base URL。
	APIBaseURL string `json:"api_base_url"`
	// GroupID 与 CTA 组合使用：前端据此判断能否一键建 Key。
	GroupID int64 `json:"group_id,omitempty"`
}

// CampaignKeyView 一键建 Key 的返回体。
type CampaignKeyView struct {
	APIKeyID  int64     `json:"api_key_id"`
	APIKey    string    `json:"api_key"`
	Name      string    `json:"name"`
	BaseURL   string    `json:"base_url"`
	GroupID   int64     `json:"group_id"`
	GroupName string    `json:"group_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CampaignRepository 落地页持久化接口（原生 SQL 实现，见 repository/campaign_repo.go）。
type CampaignRepository interface {
	GetByID(ctx context.Context, id int64) (*Campaign, error)
	GetBySlug(ctx context.Context, slug string) (*Campaign, error)
	GetActiveByHost(ctx context.Context, host string) (*Campaign, error)
	List(ctx context.Context) ([]Campaign, error)
	Create(ctx context.Context, campaign *Campaign) (*Campaign, error)
	Update(ctx context.Context, campaign *Campaign) (*Campaign, error)
	Delete(ctx context.Context, id int64) error

	// GetKeyGrantAPIKeyID 返回该活动已发放给该用户的主键 ID；found=false 表示未发放。
	GetKeyGrantAPIKeyID(ctx context.Context, campaignID, userID int64) (apiKeyID int64, found bool, err error)
	// CreateKeyGrant 记录发放关系；唯一约束冲突由实现层翻译为 ErrCampaignHostTaken 之外的
	// 通用冲突错误，由服务层回读已存在记录。
	CreateKeyGrant(ctx context.Context, campaignID, userID, apiKeyID int64) error
	// AttributeSignup 写入注册归因（仅首次写入生效）。
	AttributeSignup(ctx context.Context, campaignID, userID int64) error
}

type campaignModelCacheEntry struct {
	expires time.Time
	group   *CampaignGroupView
	models  []CampaignModelView
}

// campaignHostCacheEntry 域名 -> 是否绑定 active 活动的进程内快照。
//
// 该判定在每次 HTML 请求前都会发生（决定返回哪个前端入口），必须避免一请求一查库；
// 活动配置变更频率极低，30s 内的延迟不可感知。
type campaignHostCacheEntry struct {
	expires time.Time
	hosts   map[string]struct{}
}

const campaignHostCacheTTL = 30 * time.Second

// CampaignService 组装落地页数据并处理转化动作。
type CampaignService struct {
	repo      CampaignRepository
	groupRepo GroupRepository
	plaza     *ModelPlazaService
	apiKeys   *APIKeyService
	settings  *SettingService
	users     *UserService

	modelCache sync.Map     // map[int64]*campaignModelCacheEntry
	hostCache  atomic.Value // *campaignHostCacheEntry
}

// NewCampaignService 创建落地页服务。
func NewCampaignService(
	repo CampaignRepository,
	groupRepo GroupRepository,
	plaza *ModelPlazaService,
	apiKeys *APIKeyService,
	settings *SettingService,
	users *UserService,
) *CampaignService {
	return &CampaignService{
		repo:      repo,
		groupRepo: groupRepo,
		plaza:     plaza,
		apiKeys:   apiKeys,
		settings:  settings,
		users:     users,
	}
}

// NormalizeCampaignHost 归一化域名：去掉协议、端口、路径，转小写，去掉末尾点。
func NormalizeCampaignHost(raw string) string {
	host := strings.TrimSpace(strings.ToLower(raw))
	if host == "" {
		return ""
	}
	if idx := strings.Index(host, "://"); idx >= 0 {
		host = host[idx+3:]
	}
	if idx := strings.IndexAny(host, "/?#"); idx >= 0 {
		host = host[:idx]
	}
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		host = parsed
	}
	host = strings.TrimSuffix(host, ".")
	// IPv6 字面量形如 [::1]，去掉方括号便于比较。
	host = strings.TrimPrefix(strings.TrimSuffix(host, "]"), "[")
	return host
}

// ValidateCampaignSlug 校验 slug：[a-z0-9]，可含中划线，长度 2-64。
func ValidateCampaignSlug(slug string) error {
	s := strings.TrimSpace(slug)
	if len(s) < 1 || len(s) > 64 {
		return ErrCampaignSlugInvalid
	}
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-':
			if i == 0 || i == len(s)-1 {
				return ErrCampaignSlugInvalid
			}
		default:
			return ErrCampaignSlugInvalid
		}
	}
	return nil
}

// normalizeCampaignLocale 归一化界面语言。
//
// 未知值一律回落 en：这套页面的默认场景是国产模型出海，英文优先。
func normalizeCampaignLocale(locale string) string {
	switch strings.ToLower(strings.TrimSpace(locale)) {
	case "zh", "zh-cn", "zh_cn", "cn":
		return "zh"
	default:
		return "en"
	}
}

// sanitizeCampaignTranslations 丢弃不支持的语言与空条目，并截断超长正文，
// 避免脏配置把公开接口的响应撑大。
func sanitizeCampaignTranslations(input map[string]CampaignLocaleContent) map[string]CampaignLocaleContent {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]CampaignLocaleContent, len(input))
	for key, value := range input {
		locale := strings.ToLower(strings.TrimSpace(key))
		if _, ok := campaignSupportedLocales[locale]; !ok {
			continue
		}
		if len(value.ContentMD) > campaignMaxContentLen {
			value.ContentMD = value.ContentMD[:campaignMaxContentLen]
		}
		if value.HeroTitle == "" && value.HeroSubtitle == "" && value.ContentMD == "" && value.CTALabel == "" {
			continue
		}
		out[locale] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ValidateCampaign 校验并归一化配置。返回的错误均为可直接回给管理端的业务错误。
func ValidateCampaign(campaign *Campaign) error {
	if campaign == nil {
		return ErrCampaignNotFound
	}
	if err := ValidateCampaignSlug(campaign.Slug); err != nil {
		return err
	}
	campaign.Slug = strings.TrimSpace(campaign.Slug)
	campaign.Name = strings.TrimSpace(campaign.Name)
	if campaign.Name == "" {
		return ErrCampaignNameEmpty
	}
	if len(campaign.Name) > 128 {
		campaign.Name = campaign.Name[:128]
	}
	switch campaign.Status {
	case "":
		campaign.Status = CampaignStatusDraft
	case CampaignStatusDraft, CampaignStatusActive, CampaignStatusArchived:
	default:
		return ErrCampaignStatusBad
	}
	campaign.Host = NormalizeCampaignHost(campaign.Host)
	if len(campaign.Host) > 255 {
		return fmt.Errorf("campaign host too long")
	}
	campaign.Locale = normalizeCampaignLocale(campaign.Locale)
	campaign.Translations = sanitizeCampaignTranslations(campaign.Translations)
	campaign.HiddenModels = sanitizeHiddenModels(campaign.HiddenModels)
	if len(campaign.ContentMD) > campaignMaxContentLen {
		return fmt.Errorf("campaign content too large")
	}
	if len(campaign.ModelEntries) > campaignMaxModelEntries {
		campaign.ModelEntries = campaign.ModelEntries[:campaignMaxModelEntries]
	}
	if len(campaign.CodeSamples) > campaignMaxCodeSamples {
		campaign.CodeSamples = campaign.CodeSamples[:campaignMaxCodeSamples]
	}
	if campaign.GroupID != nil && *campaign.GroupID <= 0 {
		campaign.GroupID = nil
	}
	switch campaign.CTA.Mode {
	case "", "key", "register", "topup":
	default:
		return fmt.Errorf("campaign cta mode is invalid")
	}
	if campaign.Trial.Balance < 0 {
		campaign.Trial.Balance = 0
	}
	return nil
}

// CampaignResolveOptions 解析落地页的入参。
type CampaignResolveOptions struct {
	// Host 请求域名（优先取 X-Forwarded-Host，由 handler 决定）。
	Host string
	// Slug 显式指定活动，优先于 Host（用于预览与同域名多活动）。
	Slug string
	// AllowInactive 允许 draft/archived（仅管理员预览时置真）。
	AllowInactive bool
	// FallbackBaseURL 当 api_base_url 未配置时使用的网关地址（通常取请求 origin）。
	FallbackBaseURL string
}

// Resolve 按 slug 或域名解析落地页数据。
func (s *CampaignService) Resolve(ctx context.Context, opts CampaignResolveOptions) (*CampaignView, error) {
	var (
		campaign *Campaign
		err      error
	)

	if slug := strings.TrimSpace(opts.Slug); slug != "" {
		campaign, err = s.repo.GetBySlug(ctx, slug)
	} else if host := NormalizeCampaignHost(opts.Host); host != "" {
		campaign, err = s.repo.GetActiveByHost(ctx, host)
	}
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, ErrCampaignNotFound
	}
	if campaign.Status != CampaignStatusActive && !opts.AllowInactive {
		return nil, ErrCampaignNotFound
	}

	view := &CampaignView{
		Slug:         campaign.Slug,
		Name:         campaign.Name,
		Locale:       normalizeCampaignLocale(campaign.Locale),
		Brand:        campaign.Brand,
		Hero:         campaign.Hero,
		ContentMD:    campaign.ContentMD,
		CodeSamples:  campaign.CodeSamples,
		CTA:          campaign.CTA,
		Trial:        campaign.Trial,
		SEO:          campaign.SEO,
		Translations: campaign.Translations,
		APIBaseURL:   s.apiBaseURL(ctx, opts.FallbackBaseURL),
	}

	if campaign.GroupID != nil && *campaign.GroupID > 0 {
		group, models, err := s.groupModels(ctx, *campaign.GroupID)
		if err != nil {
			return nil, err
		}
		view.Group = group
		view.GroupID = *campaign.GroupID
		hidden := hiddenModelSet(campaign.HiddenModels)
		view.Models = filterHiddenModels(mergeCampaignModels(campaign.ModelEntries, models), hidden)
		view.Catalog = filterHiddenModels(catalogModels(models, campaign.ModelEntries), hidden)
	} else {
		// 未绑定分组：仅展示配置里的模型名，标注不可用（纯展示型页面）。
		view.Models = filterHiddenModels(
			mergeCampaignModels(campaign.ModelEntries, nil),
			hiddenModelSet(campaign.HiddenModels),
		)
	}

	if view.CTA.Mode == "" {
		view.CTA.Mode = "key"
		if view.GroupID == 0 {
			view.CTA.Mode = "register"
		}
	}
	return view, nil
}

// CampaignSignupHooks 由 AuthService 在注册成功后回调，用于落地页归因与试用赠送。
//
// 回调失败绝不能让注册失败——广告归因是附加动作，注册可用性优先。
type CampaignSignupHooks interface {
	OnCampaignSignup(ctx context.Context, userID int64, slug string) error
}

// OnCampaignSignup 记录来源活动，并按活动配置发放试用额度。
//
// 余额走 UserService.UpdateBalance（原子自增 + 缓存失效），与后台改余额同一条路径；
// 副作用是 trial 金额会计入 total_recharged，对投放赠送而言口径略有偏差，
// 但保证了缓存一致性，后续如需精确区分应改为独立的赠送流水。
func (s *CampaignService) OnCampaignSignup(ctx context.Context, userID int64, slug string) error {
	slug = strings.TrimSpace(slug)
	if userID <= 0 || slug == "" {
		return nil
	}
	campaign, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if campaign == nil || campaign.Status != CampaignStatusActive {
		return nil
	}
	if err := s.repo.AttributeSignup(ctx, campaign.ID, userID); err != nil {
		return err
	}
	if campaign.Trial.Balance <= 0 || s.users == nil {
		return nil
	}
	if err := s.users.UpdateBalance(ctx, userID, campaign.Trial.Balance); err != nil {
		return fmt.Errorf("grant campaign trial balance: %w", err)
	}
	return nil
}

// List 返回全部落地页（管理端）。
func (s *CampaignService) List(ctx context.Context) ([]Campaign, error) {
	return s.repo.List(ctx)
}

// IsActiveCampaignHost 判断域名是否绑定了已上线活动（供前端入口分流使用）。
//
// 任何查询失败都返回 false，让请求回落到主站入口——分流是附加能力，
// 不能因为它把主站打成 500 或空白页。
func (s *CampaignService) IsActiveCampaignHost(ctx context.Context, host string) bool {
	normalized := NormalizeCampaignHost(host)
	if normalized == "" {
		return false
	}
	if cached, ok := s.hostCache.Load().(*campaignHostCacheEntry); ok && time.Now().Before(cached.expires) {
		_, hit := cached.hosts[normalized]
		return hit
	}

	campaigns, err := s.repo.List(ctx)
	if err != nil {
		return false
	}
	hosts := make(map[string]struct{}, len(campaigns))
	for i := range campaigns {
		if campaigns[i].Status != CampaignStatusActive {
			continue
		}
		if h := NormalizeCampaignHost(campaigns[i].Host); h != "" {
			hosts[h] = struct{}{}
		}
	}
	s.hostCache.Store(&campaignHostCacheEntry{
		expires: time.Now().Add(campaignHostCacheTTL),
		hosts:   hosts,
	})
	_, hit := hosts[normalized]
	return hit
}

// Get 按 ID 返回落地页（管理端）。
func (s *CampaignService) Get(ctx context.Context, id int64) (*Campaign, error) {
	return s.repo.GetByID(ctx, id)
}

// Create 新建落地页（管理端）。
func (s *CampaignService) Create(ctx context.Context, campaign *Campaign, operatorID int64) (*Campaign, error) {
	if err := s.validateWithGroup(ctx, campaign); err != nil {
		return nil, err
	}
	if operatorID > 0 {
		campaign.CreatedBy = &operatorID
	}
	created, err := s.repo.Create(ctx, campaign)
	if err != nil {
		return nil, err
	}
	s.invalidateHostCache()
	return created, nil
}

// Update 更新落地页（管理端）。
func (s *CampaignService) Update(ctx context.Context, campaign *Campaign) (*Campaign, error) {
	if campaign == nil || campaign.ID <= 0 {
		return nil, ErrCampaignNotFound
	}
	if err := s.validateWithGroup(ctx, campaign); err != nil {
		return nil, err
	}
	s.invalidateModelCache(func(groupID int64) bool { return campaign.GroupID != nil && groupID == *campaign.GroupID })
	updated, err := s.repo.Update(ctx, campaign)
	if err != nil {
		return nil, err
	}
	s.invalidateHostCache()
	return updated, nil
}

// Delete 删除落地页（管理端）。
func (s *CampaignService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateHostCache()
	return nil
}

// InvalidateCache 清空分组模型缓存（分组或渠道定价变更后调用）。
func (s *CampaignService) InvalidateCache() {
	s.invalidateModelCache(func(int64) bool { return true })
}

func (s *CampaignService) invalidateModelCache(match func(groupID int64) bool) {
	s.modelCache.Range(func(key, _ any) bool {
		if groupID, ok := key.(int64); ok && match(groupID) {
			s.modelCache.Delete(key)
		}
		return true
	})
}

func (s *CampaignService) invalidateHostCache() {
	s.hostCache.Store(&campaignHostCacheEntry{})
}

func (s *CampaignService) validateWithGroup(ctx context.Context, campaign *Campaign) error {
	if err := ValidateCampaign(campaign); err != nil {
		return err
	}
	if campaign.GroupID == nil {
		return nil
	}
	if s.groupRepo == nil {
		return nil
	}
	group, err := s.groupRepo.GetByID(ctx, *campaign.GroupID)
	if err != nil {
		return fmt.Errorf("load campaign group: %w", err)
	}
	if group == nil {
		return ErrCampaignNoGroup
	}
	return nil
}

func (s *CampaignService) apiBaseURL(ctx context.Context, fallback string) string {
	if s.settings != nil {
		if value, err := s.settings.settingRepo.GetValue(ctx, SettingKeyAPIBaseURL); err == nil {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				return trimmed
			}
		}
	}
	return strings.TrimSpace(fallback)
}

// groupModels 返回绑定分组的展示信息与模型清单（带 60s 进程内缓存）。
func (s *CampaignService) groupModels(ctx context.Context, groupID int64) (*CampaignGroupView, []CampaignModelView, error) {
	if cached, ok := s.modelCache.Load(groupID); ok {
		if entry, ok := cached.(*campaignModelCacheEntry); ok && time.Now().Before(entry.expires) {
			return entry.group, entry.models, nil
		}
		s.modelCache.Delete(groupID)
	}

	if s.plaza == nil {
		return nil, nil, nil
	}
	groups, err := s.plaza.ListGroups(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list plaza groups: %w", err)
	}

	var (
		groupView *CampaignGroupView
		models    []CampaignModelView
	)
	for i := range groups {
		if groups[i].ID != groupID {
			continue
		}
		groupView = &CampaignGroupView{
			ID:                 groups[i].ID,
			Name:               groups[i].Name,
			Platform:           groups[i].Platform,
			RateMultiplier:     groups[i].RateMultiplier,
			PeakRateMultiplier: groups[i].PeakRateMultiplier,
			PeakStart:          groups[i].PeakStart,
			PeakEnd:            groups[i].PeakEnd,
		}
		models = make([]CampaignModelView, 0, len(groups[i].Models))
		for j := range groups[i].Models {
			models = append(models, plazaModelToCampaignView(&groups[i].Models[j]))
		}
		break
	}

	// 分组可能尚未配置渠道定价（models 为空）；此处仍缓存空结果，避免每次回源。
	s.modelCache.Store(groupID, &campaignModelCacheEntry{
		expires: time.Now().Add(campaignGroupModelsTTL),
		group:   groupView,
		models:  models,
	})
	return groupView, models, nil
}

func plazaModelToCampaignView(model *PlazaModel) CampaignModelView {
	view := CampaignModelView{
		Model:       model.Name,
		DisplayName: model.Name,
		Platform:    model.Platform,
		Available:   true,
	}
	if model.Pricing != nil {
		view.InputPrice = model.Pricing.InputPrice
		view.OutputPrice = model.Pricing.OutputPrice
		view.CacheReadPrice = model.Pricing.CacheReadPrice
		view.CacheWritePrice = model.Pricing.CacheWritePrice
		view.PerRequestPrice = model.Pricing.PerRequestPrice
	}
	if model.OfficialPricing != nil {
		view.OfficialInputPrice = model.OfficialPricing.InputPrice
		view.OfficialOutputPrice = model.OfficialPricing.OutputPrice
	}
	if model.TimePricing != nil {
		view.TimePricingTimezone = model.TimePricing.Timezone
		for _, period := range model.TimePricing.Periods {
			if period.Multiplier == 0 || period.Multiplier == 1 {
				continue
			}
			view.TimePricing = append(view.TimePricing, CampaignTimePricingPeriod{
				StartTime:  period.StartTime,
				EndTime:    period.EndTime,
				Multiplier: period.Multiplier,
			})
		}
	}
	return view
}

// sanitizeHiddenModels 归一化隐藏清单：去空、去重、限长。
func sanitizeHiddenModels(input []string) []string {
	if len(input) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, raw := range input {
		name := strings.TrimSpace(raw)
		if name == "" || len(name) > 128 {
			continue
		}
		key := strings.ToLower(name)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, name)
		if len(out) >= 200 {
			break
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// hiddenModelSet 构造大小写无关的隐藏集合。
func hiddenModelSet(models []string) map[string]struct{} {
	if len(models) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(models))
	for _, name := range models {
		out[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}
	return out
}

// filterHiddenModels 过滤掉别名/退役模型；hidden 为空时原样返回，避免额外分配。
func filterHiddenModels(models []CampaignModelView, hidden map[string]struct{}) []CampaignModelView {
	if len(hidden) == 0 || len(models) == 0 {
		return models
	}
	out := make([]CampaignModelView, 0, len(models))
	for _, model := range models {
		if _, skip := hidden[strings.ToLower(model.Model)]; skip {
			continue
		}
		out = append(out, model)
	}
	return out
}

// catalogModels 返回分组全量模型（按名称排序），用于「模型广场」价目表。
//
// 与 mergeCampaignModels 的区别：这里不做运营筛选、不覆盖展示名，只保证顺序稳定，
// 让用户能扫完整个分组的模型；空清单返回 nil，前端据此隐藏区块。
// catalogModels 返回广场全量模型，并按 entries 应用官方价覆盖。
//
// 官方价覆盖必须在广场这一层也生效：主推卡片走 mergeCampaignModels，
// 而广场是独立的全量列表，只覆盖前者会出现「卡片有折扣、广场没有」。
func catalogModels(models []CampaignModelView, entries []CampaignModelEntry) []CampaignModelView {
	if len(models) == 0 {
		return nil
	}
	overrides := make(map[string]CampaignModelEntry, len(entries))
	for _, entry := range entries {
		if entry.Model == "" {
			continue
		}
		if entry.OfficialInputPrice != nil || entry.OfficialOutputPrice != nil {
			overrides[strings.ToLower(entry.Model)] = entry
		}
	}

	out := make([]CampaignModelView, len(models))
	copy(out, models)
	for i := range out {
		entry, ok := overrides[strings.ToLower(out[i].Model)]
		if !ok {
			continue
		}
		if entry.OfficialInputPrice != nil {
			out[i].OfficialInputPrice = entry.OfficialInputPrice
		}
		if entry.OfficialOutputPrice != nil {
			out[i].OfficialOutputPrice = entry.OfficialOutputPrice
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Model != out[j].Model {
			return out[i].Model < out[j].Model
		}
		return out[i].Platform < out[j].Platform
	})
	return out
}

// mergeCampaignModels 把运营配置的主推顺序与真实分组模型合并。
//
// 规则：
//   - entries 为空 → 直接返回分组全部模型（按名称排序）；
//   - entries 非空 → 只输出 entries 里的模型：命中则合并定价并可用，未命中则保留
//     配置文案但标记 available=false（分组未配置该模型 / 渠道下线）；
//   - 排序：featured 优先，其次 order 升序，最后保持配置原序。
func mergeCampaignModels(entries []CampaignModelEntry, available []CampaignModelView) []CampaignModelView {
	if len(available) == 0 && len(entries) == 0 {
		return nil
	}
	if len(entries) == 0 {
		out := make([]CampaignModelView, len(available))
		copy(out, available)
		sort.SliceStable(out, func(i, j int) bool { return out[i].Model < out[j].Model })
		return out
	}

	byName := make(map[string]CampaignModelView, len(available))
	for _, model := range available {
		byName[model.Model] = model
	}

	type ordered struct {
		view  CampaignModelView
		order int
		index int
	}
	items := make([]ordered, 0, len(entries))
	for i, entry := range entries {
		view, ok := byName[entry.Model]
		if !ok {
			view = CampaignModelView{Model: entry.Model, DisplayName: entry.Model}
		}
		view.Available = ok
		if entry.DisplayName != "" {
			view.DisplayName = entry.DisplayName
		} else if view.DisplayName == "" {
			view.DisplayName = entry.Model
		}
		if entry.Tagline != "" {
			view.Tagline = entry.Tagline
		}
		if entry.Badge != "" {
			view.Badge = entry.Badge
		}
		if entry.OfficialInputPrice != nil {
			view.OfficialInputPrice = entry.OfficialInputPrice
		}
		if entry.OfficialOutputPrice != nil {
			view.OfficialOutputPrice = entry.OfficialOutputPrice
		}
		view.Featured = entry.Featured
		items = append(items, ordered{view: view, order: entry.Order, index: i})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].view.Featured != items[j].view.Featured {
			return items[i].view.Featured
		}
		if items[i].order != items[j].order {
			return items[i].order < items[j].order
		}
		return items[i].index < items[j].index
	})

	out := make([]CampaignModelView, 0, len(items))
	for _, item := range items {
		out = append(out, item.view)
	}
	return out
}

// ClaimKey 为当前用户在该活动绑定的分组下发放一把 Key（幂等）。
func (s *CampaignService) ClaimKey(ctx context.Context, userID int64, slug string, fallbackBaseURL string) (*CampaignKeyView, error) {
	if userID <= 0 {
		return nil, ErrCampaignNotFound
	}
	campaign, err := s.repo.GetBySlug(ctx, strings.TrimSpace(slug))
	if err != nil {
		return nil, err
	}
	if campaign == nil || campaign.Status != CampaignStatusActive {
		return nil, ErrCampaignNotFound
	}
	if campaign.GroupID == nil || *campaign.GroupID <= 0 {
		return nil, ErrCampaignNoGroup
	}
	baseURL := s.apiBaseURL(ctx, fallbackBaseURL)

	if keyID, found, err := s.repo.GetKeyGrantAPIKeyID(ctx, campaign.ID, userID); err != nil {
		return nil, err
	} else if found {
		if view, ok := s.loadGrantedKey(ctx, campaign, keyID, userID, baseURL); ok {
			return view, nil
		}
		// 授权记录指向的 Key 已被删除或转手：允许重新发放。
	}

	groupID := *campaign.GroupID
	created, err := s.apiKeys.Create(ctx, userID, CreateAPIKeyRequest{
		Name:    campaignKeyName(campaign),
		GroupID: &groupID,
	})
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateKeyGrant(ctx, campaign.ID, userID, created.ID); err != nil {
		// 并发点击：唯一约束挡住第二次写入，以已存在记录为准。
		if keyID, found, readErr := s.repo.GetKeyGrantAPIKeyID(ctx, campaign.ID, userID); readErr == nil && found {
			if view, ok := s.loadGrantedKey(ctx, campaign, keyID, userID, baseURL); ok {
				return view, nil
			}
		}
		return nil, err
	}
	return s.buildKeyView(ctx, campaign, created.Key, created.ID, created.Name, baseURL), nil
}

func (s *CampaignService) loadGrantedKey(ctx context.Context, campaign *Campaign, keyID, userID int64, baseURL string) (*CampaignKeyView, bool) {
	key, err := s.apiKeys.GetByID(ctx, keyID)
	if err != nil || key == nil || key.UserID != userID {
		return nil, false
	}
	if key.Status != "" && key.Status != StatusActive {
		return nil, false
	}
	return s.buildKeyView(ctx, campaign, key.Key, key.ID, key.Name, baseURL), true
}

func (s *CampaignService) buildKeyView(_ context.Context, campaign *Campaign, apiKey string, keyID int64, name string, baseURL string) *CampaignKeyView {
	view := &CampaignKeyView{
		APIKeyID: keyID,
		APIKey:   apiKey,
		Name:     name,
		BaseURL:  baseURL,
	}
	if campaign.GroupID != nil {
		view.GroupID = *campaign.GroupID
	}
	if group, _, err := s.groupModels(context.Background(), view.GroupID); err == nil && group != nil {
		view.GroupName = group.Name
	}
	return view
}

func campaignKeyName(campaign *Campaign) string {
	name := fmt.Sprintf("promo-%s", campaign.Slug)
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}
