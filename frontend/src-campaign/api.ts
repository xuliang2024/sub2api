/**
 * 落地页专用 HTTP 客户端。
 *
 * 为什么不复用 src/api/client.ts：那个客户端带 i18n、token 自动刷新、管理面
 * header 等一堆依赖，会把主站的 i18n 文案包拖进落地页 bundle。落地页只需要
 * 「带 Bearer + 解包 {code,message,data} + 401 清登录态」这三件事。
 */

import axios, { AxiosError } from 'axios'
import type { CampaignKeyGrant, CampaignUser, CampaignView } from './types'
import { clearSession, getToken, saveSession } from './session'

const DEFAULT_BASE = '/api/v1'

function resolveBaseURL(): string {
  const raw = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim()
  if (!raw) {
    return DEFAULT_BASE
  }
  return raw.replace(/\/+$/, '')
}

export const API_BASE_URL = resolveBaseURL()

const http = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      clearSession()
    }
    return Promise.reject(error)
  },
)

interface Envelope<T> {
  code: number
  message: string
  data: T
}

async function unwrap<T>(request: Promise<{ data: Envelope<T> | T }>): Promise<T> {
  const response = await request
  const body = response.data as Envelope<T>
  if (body && typeof body === 'object' && 'code' in body) {
    if (body.code !== 0) {
      throw new Error(body.message || '请求失败')
    }
    return body.data
  }
  return response.data as T
}

/** 把 axios 错误翻译成可直接展示的中文文案。 */
export function describeError(error: unknown, fallback = '请求失败，请稍后重试'): string {
  if (axios.isAxiosError(error)) {
    const payload = error.response?.data as { message?: string; error?: string } | undefined
    const message = payload?.message || payload?.error
    if (message) {
      return message
    }
    if (error.code === 'ECONNABORTED') {
      return '请求超时，请稍后重试'
    }
    if (!error.response) {
      return '网络异常，请检查网络后重试'
    }
  }
  if (error instanceof Error && error.message) {
    return error.message
  }
  return fallback
}

// ==================== 落地页数据 ====================

/**
 * 读取落地页配置。
 * slug 为空时由后端按 Host 解析——这也是线上唯一需要的方式。
 */
export function fetchCampaign(slug?: string): Promise<CampaignView> {
  const params = slug ? { c: slug } : undefined
  return unwrap<CampaignView>(http.get('/public/campaign', { params }))
}

// ==================== 认证 ====================

export interface LoginPayload {
  access_token: string
  refresh_token?: string
  expires_in?: number
  user?: CampaignUser
}

export async function login(email: string, password: string): Promise<LoginPayload> {
  const data = await unwrap<LoginPayload>(http.post('/auth/login', { email, password }))
  saveSession(data)
  return data
}

export interface RegisterPayload {
  email: string
  password: string
  verify_code?: string
  promo_code?: string
  campaign?: string
  turnstile_token?: string
}

export async function register(payload: RegisterPayload): Promise<LoginPayload> {
  const data = await unwrap<LoginPayload>(http.post('/auth/register', payload))
  saveSession(data)
  return data
}

export async function sendVerifyCode(email: string, turnstileToken?: string): Promise<void> {
  await unwrap(http.post('/auth/send-verify-code', {
    email,
    turnstile_token: turnstileToken || '',
  }))
}

// ==================== 转化动作 ====================

/** 一键建 Key（后端幂等：同一活动同一用户永远返回同一把）。 */
export function claimKey(slug: string): Promise<CampaignKeyGrant> {
  return unwrap<CampaignKeyGrant>(http.post(`/campaigns/${encodeURIComponent(slug)}/key`))
}

// ==================== 用户控制台 ====================
//
// 复用主站已实现的用户端接口（/user/profile、/keys、/usage、/usage/dashboard/stats），
// 落地页只做展示，不重复造数据口径。

export interface ConsoleProfile {
  id: number
  email?: string
  username?: string
  role?: string
  balance?: number
}

export interface ConsoleKey {
  id: number
  name: string
  key: string
  status: 'active' | 'inactive' | 'quota_exhausted' | 'expired'
  created_at: string
  last_used_at?: string | null
  quota?: number
  quota_used?: number
  expires_at?: string | null
  group?: { id: number; name: string } | null
}

export interface ConsoleUsageLog {
  id: number
  created_at: string
  model: string
  api_key_id: number
  input_tokens: number
  output_tokens: number
  cache_read_tokens: number
  total_tokens?: number
  actual_cost: number
  duration_ms?: number
  status_code?: number
}

export interface ConsoleStats {
  total_requests: number
  today_requests: number
  total_tokens: number
  today_tokens: number
  total_actual_cost: number
  today_actual_cost: number
  total_api_keys: number
  active_api_keys: number
  rpm: number
  tpm: number
  average_duration_ms: number
}

export interface ConsolePage<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}

export function fetchConsoleProfile(): Promise<ConsoleProfile> {
  return unwrap<ConsoleProfile>(http.get('/user/profile'))
}

export function fetchConsoleStats(): Promise<ConsoleStats> {
  return unwrap<ConsoleStats>(http.get('/usage/dashboard/stats'))
}

export function fetchConsoleKeys(page = 1, pageSize = 20): Promise<ConsolePage<ConsoleKey>> {
  return unwrap<ConsolePage<ConsoleKey>>(http.get('/keys', { params: { page, page_size: pageSize } }))
}

export function createConsoleKey(name: string): Promise<ConsoleKey> {
  return unwrap<ConsoleKey>(http.post('/keys', { name }))
}

export function updateConsoleKeyStatus(
  id: number,
  status: 'active' | 'inactive',
): Promise<ConsoleKey> {
  return unwrap<ConsoleKey>(http.put(`/keys/${id}`, { status }))
}

export function deleteConsoleKey(id: number): Promise<void> {
  return unwrap<void>(http.delete(`/keys/${id}`))
}

export function fetchConsoleLogs(
  page = 1,
  pageSize = 20,
  apiKeyId?: number,
): Promise<ConsolePage<ConsoleUsageLog>> {
  const params: Record<string, number> = { page, page_size: pageSize }
  if (apiKeyId) {
    params.api_key_id = apiKeyId
  }
  return unwrap<ConsolePage<ConsoleUsageLog>>(http.get('/usage', { params }))
}
