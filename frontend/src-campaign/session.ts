/**
 * 落地页的登录态存储。
 *
 * 刻意与主站 SPA 使用同一组 localStorage key（auth_token / refresh_token /
 * auth_user / token_expires_at）：用户在落地页注册后直接去主站，无需二次登录。
 */

import type { CampaignUser } from './types'

const TOKEN_KEY = 'auth_token'
const REFRESH_KEY = 'refresh_token'
const USER_KEY = 'auth_user'
const EXPIRES_KEY = 'token_expires_at'

export interface TokenPairPayload {
  access_token: string
  refresh_token?: string
  expires_in?: number
  user?: CampaignUser
}

export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

export function getStoredUser(): CampaignUser | null {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? (JSON.parse(raw) as CampaignUser) : null
  } catch {
    return null
  }
}

export function saveSession(payload: TokenPairPayload): void {
  if (!payload?.access_token) {
    return
  }
  try {
    localStorage.setItem(TOKEN_KEY, payload.access_token)
    if (payload.refresh_token) {
      localStorage.setItem(REFRESH_KEY, payload.refresh_token)
    }
    if (typeof payload.expires_in === 'number' && payload.expires_in > 0) {
      localStorage.setItem(EXPIRES_KEY, String(Date.now() + payload.expires_in * 1000))
    }
    if (payload.user) {
      localStorage.setItem(USER_KEY, JSON.stringify(payload.user))
    }
  } catch {
    // 隐私模式下 localStorage 可能不可写：登录态仅本次会话有效，不阻断流程。
  }
}

export function clearSession(): void {
  try {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(EXPIRES_KEY)
  } catch {
    // ignore
  }
}
