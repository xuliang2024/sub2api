<template>
  <AppLayout>
    <div class="space-y-6">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <EmptyState
        v-else-if="!overview"
        :title="t('partner.notPartnerTitle')"
        :description="t('partner.notPartnerDescription')"
      />

      <template v-else>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
          <StatCard :title="t('partner.stats.minRate')" :value="`×${overview.min_rate_multiplier}`" />
          <StatCard :title="t('partner.stats.downlines')" :value="overview.downline_count" />
          <StatCard :title="t('partner.stats.totalCommission')" :value="formatCurrency(overview.total_commission)" icon-variant="success" />
          <StatCard :title="t('partner.stats.monthCommission')" :value="formatCurrency(overview.month_commission)" icon-variant="success" />
          <StatCard :title="t('partner.stats.pendingSettlement')" :value="formatCurrency(overview.pending_settlement)" icon-variant="warning" />
        </div>

        <div class="card p-6">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('partner.invite.title') }}</h3>
                <span :class="['badge', overview.status === 'active' ? 'badge-success' : 'badge-gray']">
                  {{ overview.status === 'active' ? t('common.active') : t('partner.disabled') }}
                </span>
              </div>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('partner.invite.description') }}</p>
              <code class="mt-3 block truncate text-sm text-gray-700 dark:text-gray-300">{{ inviteLink }}</code>
            </div>
            <button class="btn btn-secondary shrink-0" @click="copyInviteLink">
              <Icon name="copy" size="sm" />{{ t('partner.invite.copy') }}
            </button>
          </div>
        </div>

        <div class="card overflow-hidden">
          <div class="flex border-b border-gray-200 dark:border-dark-700" role="tablist">
            <button
              type="button"
              :class="['px-6 py-4 text-sm font-medium transition-colors', activeTab === 'users' ? 'border-b-2 border-primary-500 text-primary-600 dark:text-primary-400' : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-gray-200']"
              @click="activeTab = 'users'"
            >
              {{ t('partner.tabs.users') }}
            </button>
            <button
              type="button"
              :class="['px-6 py-4 text-sm font-medium transition-colors', activeTab === 'commission' ? 'border-b-2 border-primary-500 text-primary-600 dark:text-primary-400' : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-gray-200']"
              @click="activeTab = 'commission'"
            >
              {{ t('partner.tabs.commission') }}
            </button>
          </div>

          <div v-if="activeTab === 'users'" class="p-6">
            <div class="mb-5 flex items-start justify-between gap-4">
              <div>
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('partner.users.title') }}</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('partner.users.description', { min: overview.min_rate_multiplier }) }}</p>
              </div>
              <button class="btn btn-secondary btn-sm shrink-0" :disabled="usersLoading" @click="loadUsers">
                <Icon name="refresh" size="sm" :class="usersLoading ? 'animate-spin' : ''" />
              </button>
            </div>

            <div v-if="usersLoading" class="flex justify-center py-10">
              <div class="h-7 w-7 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <div v-else-if="downlines.length === 0" class="rounded-xl border border-dashed border-gray-300 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
              {{ t('partner.users.empty') }}
            </div>
            <div v-else class="space-y-3">
              <div v-for="user in downlines" :key="user.user_id" class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
                <button type="button" class="flex w-full flex-col gap-3 p-4 text-left hover:bg-gray-50 dark:hover:bg-dark-800 lg:flex-row lg:items-center" @click="toggleUser(user.user_id)">
                  <div class="min-w-0 flex-1">
                    <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ user.email || '-' }}</div>
                    <div class="mt-0.5 text-xs text-gray-500 dark:text-dark-400">#{{ user.user_id }} · {{ user.username || '-' }}</div>
                  </div>
                  <div class="grid flex-1 grid-cols-2 gap-3 text-sm lg:max-w-md lg:grid-cols-3">
                    <div>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('partner.users.source') }}</p>
                      <p class="mt-1 text-gray-700 dark:text-gray-300">{{ boundSourceLabel(user.bound_source) }}</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('partner.users.cost30d') }}</p>
                      <p class="mt-1 text-gray-700 dark:text-gray-300">{{ formatCurrency(user.total_cost_30d) }}</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('partner.users.commission30d') }}</p>
                      <p class="mt-1 font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(user.commission_30d) }}</p>
                    </div>
                  </div>
                  <Icon name="chevronDown" size="sm" :class="['shrink-0 text-gray-400 transition-transform', expandedUserId === user.user_id ? 'rotate-180' : '']" />
                </button>

                <div v-if="expandedUserId === user.user_id" class="border-t border-gray-200 bg-gray-50/60 p-4 dark:border-dark-700 dark:bg-dark-800/50">
                  <div v-if="user.rates.length === 0" class="py-4 text-center text-sm text-gray-500 dark:text-dark-400">{{ t('partner.users.noGroups') }}</div>
                  <div v-else class="space-y-3">
                    <div v-for="rate in user.rates" :key="rate.group_id" class="rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
                      <div class="flex flex-col gap-3 lg:flex-row lg:items-center">
                        <div class="min-w-0 flex-1">
                          <label class="text-sm font-medium text-gray-900 dark:text-white">{{ rate.group_name }}</label>
                          <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('partner.users.defaultRate', { rate: rate.default_rate }) }}</p>
                        </div>
                        <div class="w-full lg:w-36">
                          <input
                            v-model="rateInputs[user.user_id][rate.group_id]"
                            type="number"
                            :min="overview.min_rate_multiplier"
                            max="1"
                            step="0.01"
                            class="input"
                            :placeholder="t('partner.users.followDefault')"
                            :aria-label="t('partner.users.rateForGroup', { group: rate.group_name })"
                          />
                        </div>
                      </div>
                      <p class="mt-3 text-xs text-gray-500 dark:text-dark-400">
                        {{ t('partner.users.rateRange', { min: overview.min_rate_multiplier }) }}
                      </p>
                      <p class="mt-1 text-xs text-primary-600 dark:text-primary-400">
                        {{ t('partner.users.simulation', simulationFor(user.user_id, rate)) }}
                      </p>
                    </div>
                    <div class="flex items-center justify-between gap-3">
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('partner.users.blankHint') }}</p>
                      <button class="btn btn-primary btn-sm shrink-0" :disabled="savingUserId === user.user_id" @click="saveRates(user)">
                        {{ savingUserId === user.user_id ? t('common.saving') : t('common.save') }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="space-y-5 p-6">
            <div class="flex flex-col gap-3 xl:flex-row xl:items-end">
              <div class="w-full xl:w-64">
                <label class="input-label">{{ t('partner.commission.user') }}</label>
                <Select v-model="commissionFilters.userId" :options="commissionUserOptions" />
              </div>
              <div>
                <label class="input-label">{{ t('partner.commission.startDate') }}</label>
                <input v-model="commissionFilters.startDate" type="date" class="input w-full sm:w-44" />
              </div>
              <div>
                <label class="input-label">{{ t('partner.commission.endDate') }}</label>
                <input v-model="commissionFilters.endDate" type="date" class="input w-full sm:w-44" />
              </div>
              <button class="btn btn-primary" @click="applyCommissionFilters">{{ t('partner.commission.filter') }}</button>
              <div class="rounded-xl bg-emerald-50 px-4 py-2.5 text-sm xl:ml-auto dark:bg-emerald-900/20">
                <span class="text-emerald-700 dark:text-emerald-300">{{ t('partner.commission.filteredTotal') }}</span>
                <span class="ml-2 font-semibold text-emerald-700 dark:text-emerald-300">{{ formatCurrency(commissionTotal) }}</span>
              </div>
            </div>

            <div v-if="commissionLoading" class="flex justify-center py-10">
              <div class="h-7 w-7 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <div v-else-if="commissions.length === 0" class="rounded-xl border border-dashed border-gray-300 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
              {{ t('partner.commission.empty') }}
            </div>
            <div v-else class="overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-700">
              <table class="w-full min-w-[1040px] text-left text-sm">
                <thead class="bg-gray-50 text-gray-500 dark:bg-dark-800 dark:text-dark-400">
                  <tr>
                    <th class="px-3 py-3 font-medium">{{ t('partner.commission.columns.time') }}</th>
                    <th class="px-3 py-3 font-medium">{{ t('partner.commission.columns.user') }}</th>
                    <th class="px-3 py-3 font-medium">{{ t('partner.commission.columns.model') }}</th>
                    <th class="px-3 py-3 text-right font-medium">{{ t('partner.commission.columns.totalCost') }}</th>
                    <th class="px-3 py-3 text-right font-medium">{{ t('partner.commission.columns.baseRate') }}</th>
                    <th class="px-3 py-3 text-right font-medium">{{ t('partner.commission.columns.rateFloor') }}</th>
                    <th class="px-3 py-3 text-right font-medium">{{ t('partner.commission.columns.difference') }}</th>
                    <th class="px-3 py-3 text-right font-medium">{{ t('partner.commission.columns.commission') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                  <tr v-for="item in commissions" :key="item.id">
                    <td class="whitespace-nowrap px-3 py-3 text-gray-700 dark:text-gray-300">{{ formatDateTime(item.created_at) }}</td>
                    <td class="px-3 py-3">
                      <div class="max-w-48 truncate text-gray-900 dark:text-white">{{ item.user_email || '-' }}</div>
                      <div class="font-mono text-xs text-gray-500">#{{ item.user_id }}</div>
                    </td>
                    <td class="max-w-48 truncate px-3 py-3 font-mono text-xs text-gray-700 dark:text-gray-300" :title="item.model">{{ item.model }}</td>
                    <td class="px-3 py-3 text-right text-gray-900 dark:text-white">{{ formatCurrency(item.total_cost) }}</td>
                    <td class="px-3 py-3 text-right text-gray-700 dark:text-gray-300">×{{ item.base_rate }}</td>
                    <td class="px-3 py-3 text-right text-gray-700 dark:text-gray-300">×{{ item.rate_floor }}</td>
                    <td class="px-3 py-3 text-right text-gray-700 dark:text-gray-300">×{{ formatRateDifference(item) }}</td>
                    <td class="px-3 py-3 text-right font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.commission) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('partner.commission.note') }}</p>
            <Pagination
              v-if="commissionPagination.total > 0"
              :page="commissionPagination.page"
              :total="commissionPagination.total"
              :page-size="commissionPagination.page_size"
              @update:page="handleCommissionPage"
              @update:pageSize="handleCommissionPageSize"
            />
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import StatCard from '@/components/common/StatCard.vue'
import Icon from '@/components/icons/Icon.vue'
import partnerAPI from '@/api/partner'
import type { PartnerBoundSource, PartnerCommission, PartnerDownline, PartnerOverview, PartnerUserGroupRate } from '@/api/partner'
import { useAppStore } from '@/stores/app'
import { useClipboard } from '@/composables/useClipboard'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const appStore = useAppStore()
const { copyToClipboard } = useClipboard()

const loading = ref(true)
const overview = ref<PartnerOverview | null>(null)
const activeTab = ref<'users' | 'commission'>('users')
const downlines = ref<PartnerDownline[]>([])
const usersLoading = ref(false)
const expandedUserId = ref<number | null>(null)
const savingUserId = ref<number | null>(null)
const rateInputs = reactive<Record<number, Record<number, string>>>({})

const commissions = ref<PartnerCommission[]>([])
const commissionLoading = ref(false)
const commissionLoaded = ref(false)
const commissionTotal = ref(0)
const commissionFilters = reactive({ userId: '' as string | number, startDate: '', endDate: '' })
const commissionPagination = reactive({ page: 1, page_size: 20, total: 0 })

const inviteLink = computed(() => {
  if (!overview.value?.aff_code) return ''
  const origin = typeof window === 'undefined' ? '' : window.location.origin
  return `${origin}/register?aff=${encodeURIComponent(overview.value.aff_code)}`
})
const commissionUserOptions = computed(() => [
  { value: '', label: t('partner.commission.allUsers') },
  ...downlines.value.map((user) => ({ value: user.user_id, label: user.email || `#${user.user_id}` })),
])

async function loadOverview() {
  loading.value = true
  try {
    overview.value = await partnerAPI.getOverview()
    await loadUsers()
  } catch (error) {
    overview.value = null
    if ((error as { status?: number })?.status !== 403) {
      appStore.showError(extractApiErrorMessage(error, t('partner.loadFailed')))
    }
  } finally {
    loading.value = false
  }
}

async function loadUsers() {
  usersLoading.value = true
  try {
    const res = await partnerAPI.listUsers()
    downlines.value = res.items || []
    for (const user of downlines.value) {
      rateInputs[user.user_id] = {}
      for (const rate of user.rates) {
        rateInputs[user.user_id][rate.group_id] = rate.rate_multiplier == null ? '' : String(rate.rate_multiplier)
      }
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('partner.users.loadFailed')))
  } finally {
    usersLoading.value = false
  }
}

async function loadCommission() {
  commissionLoading.value = true
  try {
    const res = await partnerAPI.listCommission({
      user_id: commissionFilters.userId === '' ? undefined : Number(commissionFilters.userId),
      start_date: commissionFilters.startDate || undefined,
      end_date: commissionFilters.endDate || undefined,
      page: commissionPagination.page,
      page_size: commissionPagination.page_size,
    })
    commissions.value = res.items || []
    commissionPagination.total = res.total || 0
    commissionTotal.value = res.total_commission || 0
    commissionLoaded.value = true
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('partner.commission.loadFailed')))
  } finally {
    commissionLoading.value = false
  }
}

function toggleUser(userId: number) {
  expandedUserId.value = expandedUserId.value === userId ? null : userId
}
function boundSourceLabel(source: PartnerBoundSource): string {
  return t(`partner.users.sources.${source}`)
}
function effectiveRate(userId: number, rate: PartnerUserGroupRate): number {
  const raw = rateInputs[userId]?.[rate.group_id]?.trim()
  if (!raw) return rate.default_rate
  const parsed = Number(raw)
  return Number.isFinite(parsed) ? parsed : rate.default_rate
}
function simulationFor(userId: number, rate: PartnerUserGroupRate) {
  const min = overview.value?.min_rate_multiplier ?? 0
  const effective = effectiveRate(userId, rate)
  const commission = Math.max(0, effective - min)
  return { userPay: effective.toFixed(2), commission: commission.toFixed(2), min: String(min) }
}
async function saveRates(user: PartnerDownline) {
  if (!overview.value) return
  const min = overview.value.min_rate_multiplier
  const rates: Array<{ group_id: number; rate_multiplier: number }> = []
  for (const rate of user.rates) {
    const raw = rateInputs[user.user_id]?.[rate.group_id]?.trim() ?? ''
    if (!raw) continue
    const value = Number(raw)
    if (!Number.isFinite(value) || value < min || value > 1) {
      appStore.showError(t('partner.users.invalidRate', { group: rate.group_name, min }))
      return
    }
    rates.push({ group_id: rate.group_id, rate_multiplier: value })
  }
  if (rates.length === 0) {
    appStore.showError(t('partner.users.noRatesToSave'))
    return
  }
  savingUserId.value = user.user_id
  try {
    await partnerAPI.updateUserRates(user.user_id, rates)
    appStore.showSuccess(t('partner.users.saved'))
    await loadUsers()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('partner.users.saveFailed')))
  } finally {
    savingUserId.value = null
  }
}
async function copyInviteLink() {
  if (inviteLink.value) await copyToClipboard(inviteLink.value, t('partner.invite.copied'))
}
function applyCommissionFilters() {
  if (commissionFilters.startDate && commissionFilters.endDate && commissionFilters.startDate > commissionFilters.endDate) {
    appStore.showError(t('partner.commission.invalidDateRange'))
    return
  }
  commissionPagination.page = 1
  void loadCommission()
}
function handleCommissionPage(page: number) {
  commissionPagination.page = page
  void loadCommission()
}
function handleCommissionPageSize(size: number) {
  commissionPagination.page_size = size
  commissionPagination.page = 1
  void loadCommission()
}
function formatRateDifference(item: PartnerCommission): string {
  const difference = Math.max(0, item.base_rate - item.rate_floor)
  return String(Number(difference.toFixed(4)))
}

watch(activeTab, (tab) => {
  if (tab === 'commission' && !commissionLoaded.value) void loadCommission()
})
onMounted(() => {
  void loadOverview()
})
</script>
