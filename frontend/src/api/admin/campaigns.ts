/**
 * Admin Campaigns API endpoints（广告投放落地页配置）
 */

import { apiClient } from '../client'

export interface CampaignBrandInput {
  site_name?: string
  logo?: string
  primary?: string
  accent?: string
  footer?: string
  hide_auth?: boolean
  console_url?: string
}

export interface CampaignHeroInput {
  eyebrow?: string
  title?: string
  subtitle?: string
  highlights?: string[]
  badges?: string[]
  image?: string
}

export interface CampaignModelEntryInput {
  model: string
  display_name?: string
  tagline?: string
  badge?: string
  featured?: boolean
  order?: number
}

export interface CampaignCodeSampleInput {
  id: string
  label: string
  lang?: string
  code: string
}

export interface CampaignCTAInput {
  mode?: 'key' | 'register' | 'topup'
  label?: string
  sublabel?: string
  topup_url?: string
}

export interface CampaignTrialInput {
  balance?: number
  note?: string
}

export interface CampaignSEOInput {
  title?: string
  description?: string
  og_image?: string
  keywords?: string[]
}

/** 某一语言的正文覆盖（未配置的字段回落主文案）。 */
export interface CampaignLocaleContentInput {
  hero_title?: string
  hero_subtitle?: string
  content_md?: string
  cta_label?: string
}

export interface CampaignPayload {
  slug: string
  name: string
  status: string
  host: string
  /** 界面语言：en（默认，海外投放）/ zh */
  locale: string
  group_id: number | null
  brand: CampaignBrandInput
  hero: CampaignHeroInput
  content_md: string
  model_entries: CampaignModelEntryInput[]
  code_samples: CampaignCodeSampleInput[]
  cta: CampaignCTAInput
  trial: CampaignTrialInput
  seo: CampaignSEOInput
  /** 按语言的正文覆盖，key 为 en/zh/ja/ko/th */
  translations: Record<string, CampaignLocaleContentInput>
  /** 不在落地页展示的模型名（别名/退役模型） */
  hidden_models: string[]
}

export interface AdminCampaign extends CampaignPayload {
  id: number
  created_at: string
  updated_at: string
}

/** 列表（配置数量很少，不分页）。 */
export async function list(): Promise<AdminCampaign[]> {
  const { data } = await apiClient.get<{ items: AdminCampaign[] }>('/admin/campaigns')
  return data.items ?? []
}

export async function create(payload: CampaignPayload): Promise<AdminCampaign> {
  const { data } = await apiClient.post<AdminCampaign>('/admin/campaigns', payload)
  return data
}

export async function update(id: number, payload: CampaignPayload): Promise<AdminCampaign> {
  const { data } = await apiClient.put<AdminCampaign>(`/admin/campaigns/${id}`, payload)
  return data
}

export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/campaigns/${id}`)
}
