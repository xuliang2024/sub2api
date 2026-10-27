import { apiClient } from '../client'
import type { PartnerBoundSource, PartnerStatus } from '../partner'

export interface AdminPartner {
  user_id: number
  email: string
  username: string
  min_rate_multiplier: number
  status: PartnerStatus
  downline_count: number
  total_commission: number
  settled_total: number
  pending_settlement: number
  created_at: string
  aff_code?: string
}

export interface AdminPartnerListResponse {
  items: AdminPartner[]
  total: number
}

export interface AdminPartnerListQuery {
  query?: string
  status?: PartnerStatus | ''
  page?: number
  page_size?: number
}

export interface UpdatePartnerRequest {
  min_rate_multiplier: number
  status?: PartnerStatus
}

export interface AdminPartnerRelationRate {
  group_id: number
  group_name: string
  rate_multiplier: number | null
}

export interface AdminPartnerRelation {
  user_id: number
  email: string
  username: string
  bound_source: PartnerBoundSource
  created_at: string
  rates: AdminPartnerRelationRate[]
}

export interface AdminPartnerSettlement {
  id: number
  amount: number
  method: 'offline'
  note?: string
  operator_id: number | null
  created_at: string
}

export interface CreateSettlementRequest {
  amount: number
  method: 'offline'
  note?: string
}

export async function list(params: AdminPartnerListQuery = {}): Promise<AdminPartnerListResponse> {
  const { data } = await apiClient.get<AdminPartnerListResponse>('/admin/partners', {
    params: {
      query: params.query?.trim() || undefined,
      status: params.status || undefined,
      page: params.page ?? 1,
      page_size: params.page_size ?? 20,
    },
  })
  return data
}

export async function upsert(userId: number, payload: UpdatePartnerRequest): Promise<void> {
  await apiClient.put(`/admin/partners/${userId}`, payload)
}

export async function remove(userId: number): Promise<void> {
  await apiClient.delete(`/admin/partners/${userId}`)
}

export async function listRelations(userId: number): Promise<{ items: AdminPartnerRelation[] }> {
  const { data } = await apiClient.get<{ items: AdminPartnerRelation[] }>(
    `/admin/partners/${userId}/relations`,
  )
  return data
}

export async function bindRelation(partnerUserId: number, userId: number): Promise<void> {
  await apiClient.post(`/admin/partners/${partnerUserId}/relations`, { user_id: userId })
}

export async function unbindRelation(partnerUserId: number, targetUserId: number): Promise<void> {
  await apiClient.delete(`/admin/partners/${partnerUserId}/relations/${targetUserId}`)
}

export async function createSettlement(
  userId: number,
  payload: CreateSettlementRequest,
): Promise<void> {
  await apiClient.post(`/admin/partners/${userId}/settlements`, payload)
}

export async function listSettlements(
  userId: number,
): Promise<{ items: AdminPartnerSettlement[] }> {
  const { data } = await apiClient.get<{ items: AdminPartnerSettlement[] }>(
    `/admin/partners/${userId}/settlements`,
  )
  return data
}

export const partnersAPI = {
  list,
  upsert,
  remove,
  listRelations,
  bindRelation,
  unbindRelation,
  createSettlement,
  listSettlements,
}

export default partnersAPI
