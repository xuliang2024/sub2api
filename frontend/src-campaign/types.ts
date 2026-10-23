/**
 * 广告落地页（campaign）的类型定义。
 *
 * 与后端 internal/service/campaign.go 的 JSON 结构一一对应；
 * 价格字段单位是「美元 / token」（与 /api/v1/model-plaza 同口径），
 * 展示时统一乘 1e6 换算成「美元 / 百万 token」。
 */

export interface CampaignBrand {
  site_name?: string
  logo?: string
  primary?: string
  accent?: string
  footer?: string
  hide_auth?: boolean
  /** 主站控制台地址；配置后导航会显示「控制台」入口。 */
  console_url?: string
  /** 分组对外展示名；中文分组名不该直接给海外用户看，留空则不在标题里提分组。 */
  group_label?: string
}

export interface CampaignHero {
  eyebrow?: string
  title?: string
  subtitle?: string
  highlights?: string[]
  badges?: string[]
  image?: string
}

export interface CampaignTimePricingPeriod {
  start_time: string
  end_time: string
  multiplier: number
}

export interface CampaignModel {
  model: string
  display_name: string
  tagline?: string
  badge?: string
  featured?: boolean
  platform?: string
  available: boolean
  input_price?: number
  output_price?: number
  cache_read_price?: number
  cache_write_price?: number
  per_request_price?: number
  /** 上游官方价（USD/token），用于展示折扣力度；缺失则不展示对比。 */
  official_input_price?: number
  official_output_price?: number
  time_pricing?: CampaignTimePricingPeriod[]
  time_pricing_timezone?: string
}

export interface CampaignCodeSample {
  id: string
  label: string
  lang?: string
  code: string
}

export interface CampaignCTA {
  mode?: 'key' | 'register' | 'topup'
  label?: string
  sublabel?: string
  topup_url?: string
}

export interface CampaignTrial {
  balance?: number
  note?: string
}

export interface CampaignGroup {
  id: number
  name: string
  platform: string
  rate_multiplier: number
  peak_rate_multiplier?: number
  peak_start?: string
  peak_end?: string
}

export interface CampaignSEO {
  title?: string
  description?: string
  og_image?: string
  keywords?: string[]
}

/** 某一语言的正文覆盖；未配置的字段回落活动主文案。 */
export interface CampaignLocaleContent {
  hero_title?: string
  hero_subtitle?: string
  content_md?: string
  cta_label?: string
}

export interface CampaignView {
  slug: string
  name: string
  /** 界面语言：en（默认）/ zh / ja / ko / th */
  locale?: string
  brand: CampaignBrand
  hero: CampaignHero
  content_md: string
  models: CampaignModel[]
  /** 绑定分组的全量模型（含未被运营选中的），用于模型广场价目表。 */
  catalog?: CampaignModel[]
  code_samples: CampaignCodeSample[]
  cta: CampaignCTA
  trial: CampaignTrial
  seo: CampaignSEO
  /** 按语言的正文覆盖，key 为语言代码；缺失则回落 hero / content_md。 */
  translations?: Record<string, CampaignLocaleContent>
  group?: CampaignGroup
  group_id?: number
  api_base_url: string
}

export interface CampaignKeyGrant {
  api_key_id: number
  api_key: string
  name: string
  base_url: string
  group_id: number
  group_name?: string
  created_at: string
}

export interface CampaignUser {
  id: number
  email?: string
  username?: string
  role?: string
  balance?: number
}

/** 注入到 campaign.html 的站点公开配置（由后端 SSR 注入，只用到其中几项）。 */
export interface CampaignAppConfig {
  registration_enabled?: boolean
  email_verify_enabled?: boolean
  google_oauth_enabled?: boolean
  promo_code_enabled?: boolean
  turnstile_enabled?: boolean
  turnstile_site_key?: string
  site_name?: string
  [key: string]: unknown
}
