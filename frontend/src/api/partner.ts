import { apiClient } from './client'

export type PartnerStatus = 'active' | 'disabled'
export type PartnerBoundSource = 'admin' | 'invite'

export interface PartnerOverview {
  is_partner: true
  min_rate_multiplier: number
  status: PartnerStatus
  downline_count: number
  total_commission: number
  month_commission: number
  settled_total: number
  pending_settlement: number
  aff_code: string
}

export interface PartnerUserGroupRate {
  group_id: number
  group_name: string
  rate_multiplier: number | null
  default_rate: number
}

export interface PartnerDownline {
  user_id: number
  email: string
  username: string
  bound_source: PartnerBoundSource
  created_at: string
  total_cost_30d: number
  commission_30d: number
  rates: PartnerUserGroupRate[]
}

export interface PartnerCommission {
  id: number
  created_at: string
  user_id: number
  user_email: string
  model: string
  total_cost: number
  base_rate: number
  rate_floor: number
  rate_multiplier: number
  commission: number
}

export interface PartnerCommissionQuery {
  user_id?: number
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

export interface PartnerCommissionResponse {
  items: PartnerCommission[]
  total: number
  total_commission: number
}

export interface PartnerRateUpdate {
  group_id: number
  rate_multiplier: number
}

export async function getOverview(): Promise<PartnerOverview> {
  const { data } = await apiClient.get<PartnerOverview>('/user/partner')
  return data
}

export async function listUsers(): Promise<{ items: PartnerDownline[] }> {
  const { data } = await apiClient.get<{ items: PartnerDownline[] }>('/user/partner/users')
  return data
}

export async function updateUserRates(userId: number, rates: PartnerRateUpdate[]): Promise<void> {
  await apiClient.put(`/user/partner/users/${userId}/rates`, { rates })
}

export async function listCommission(
  params: PartnerCommissionQuery = {},
): Promise<PartnerCommissionResponse> {
  const { data } = await apiClient.get<PartnerCommissionResponse>('/user/partner/commission', {
    params: {
      user_id: params.user_id || undefined,
      start_date: params.start_date || undefined,
      end_date: params.end_date || undefined,
      page: params.page ?? 1,
      page_size: params.page_size ?? 20,
    },
  })
  return data
}

export const partnerAPI = {
  getOverview,
  listUsers,
  updateUserRates,
  listCommission,
}

export default partnerAPI
