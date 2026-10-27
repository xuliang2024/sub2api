<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full md:w-80">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model="filters.query"
              type="text"
              class="input pl-10"
              :placeholder="t('admin.partners.searchPlaceholder')"
              @input="debounceLoad"
            />
          </div>
          <div class="w-full sm:w-40">
            <Select v-model="filters.status" :options="statusFilterOptions" @change="reloadFromFirstPage" />
          </div>
          <button class="btn btn-secondary px-2 md:px-3" :disabled="loading" :title="t('common.refresh')" @click="loadPartners">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
          <button class="btn btn-primary ml-auto" @click="openCreate">
            <Icon name="plus" size="sm" />
            {{ t('admin.partners.add') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="partners" :loading="loading" row-key="user_id" :actions-count="4">
          <template #cell-user="{ row }">
            <div class="space-y-0.5">
              <div class="font-mono text-xs text-gray-500 dark:text-dark-400">#{{ row.user_id }}</div>
              <div class="max-w-56 truncate text-sm font-medium text-gray-900 dark:text-white">{{ row.email || '-' }}</div>
              <div class="max-w-56 truncate text-xs text-gray-500 dark:text-dark-400">{{ row.username || '-' }}</div>
            </div>
          </template>
          <template #cell-min_rate_multiplier="{ row }">
            <span class="font-medium text-gray-900 dark:text-white">×{{ row.min_rate_multiplier }}</span>
          </template>
          <template #cell-status="{ row }">
            <span :class="['badge', row.status === 'active' ? 'badge-success' : 'badge-gray']">
              {{ row.status === 'active' ? t('common.active') : t('admin.partners.disabled') }}
            </span>
          </template>
          <template #cell-total_commission="{ row }">
            <span class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(row.total_commission) }}</span>
          </template>
          <template #cell-pending_settlement="{ row }">
            <span class="font-medium text-amber-600 dark:text-amber-400">{{ formatCurrency(row.pending_settlement) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-1">
              <button class="btn btn-ghost btn-sm" @click="openEdit(row)">{{ t('common.edit') }}</button>
              <button class="btn btn-ghost btn-sm" @click="openRelations(row)">{{ t('admin.partners.relations.action') }}</button>
              <button class="btn btn-ghost btn-sm" @click="openSettlements(row)">{{ t('admin.partners.settlements.action') }}</button>
              <button class="btn btn-ghost btn-sm text-red-600 dark:text-red-400" @click="requestRemove(row)">
                {{ t('admin.partners.remove') }}
              </button>
            </div>
          </template>
          <template #empty>
            <EmptyState :title="t('admin.partners.empty')" :description="t('admin.partners.emptyDescription')" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <BaseDialog
      :show="editorOpen"
      :title="editingPartner ? t('admin.partners.editor.editTitle') : t('admin.partners.editor.addTitle')"
      width="normal"
      @close="editorOpen = false"
    >
      <form id="partner-editor-form" class="space-y-5" @submit.prevent="savePartner">
        <div>
          <label class="input-label">{{ t('admin.partners.editor.user') }}</label>
          <input
            v-model="editorForm.identifier"
            type="text"
            class="input"
            :disabled="!!editingPartner"
            :placeholder="t('admin.partners.editor.userPlaceholder')"
          />
          <p class="input-hint">{{ t('admin.partners.editor.userHint') }}</p>
        </div>
        <div v-if="editingPartner" class="rounded-lg bg-gray-50 p-3 text-sm text-gray-600 dark:bg-dark-800 dark:text-gray-300">
          #{{ editingPartner.user_id }} · {{ editingPartner.email }}
        </div>
        <div>
          <label class="input-label">{{ t('admin.partners.editor.minRate') }}</label>
          <input v-model="editorForm.minRate" type="number" min="0.0001" max="1" step="0.01" class="input" />
          <p class="input-hint">{{ t('admin.partners.editor.minRateHint') }}</p>
        </div>
        <div>
          <label class="input-label">{{ t('admin.partners.editor.status') }}</label>
          <Select v-model="editorForm.status" :options="partnerStatusOptions" />
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-secondary" @click="editorOpen = false">{{ t('common.cancel') }}</button>
          <button type="submit" form="partner-editor-form" class="btn btn-primary" :disabled="editorSaving">
            {{ editorSaving ? t('common.saving') : t('common.save') }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <BaseDialog
      :show="relationsOpen"
      :title="t('admin.partners.relations.title', { email: selectedPartner?.email || '' })"
      width="extra-wide"
      @close="relationsOpen = false"
    >
      <div class="space-y-5">
        <div v-if="inviteLink" class="rounded-xl border border-primary-200 bg-primary-50 p-4 dark:border-primary-900/40 dark:bg-primary-900/20">
          <p class="text-sm font-medium text-primary-800 dark:text-primary-200">{{ t('admin.partners.relations.inviteLink') }}</p>
          <div class="mt-2 flex flex-col gap-2 sm:flex-row sm:items-center">
            <code class="min-w-0 flex-1 break-all text-sm text-primary-700 dark:text-primary-300">{{ inviteLink }}</code>
            <button class="btn btn-secondary btn-sm shrink-0" @click="copyInviteLink">
              <Icon name="copy" size="sm" />{{ t('admin.partners.relations.copyLink') }}
            </button>
          </div>
        </div>

        <form class="flex flex-col gap-2 sm:flex-row sm:items-end" @submit.prevent="bindUser">
          <div class="flex-1">
            <label class="input-label">{{ t('admin.partners.relations.bindUserId') }}</label>
            <input v-model="bindUserId" type="number" min="1" step="1" class="input" :placeholder="t('admin.partners.relations.bindUserPlaceholder')" />
          </div>
          <button type="submit" class="btn btn-primary" :disabled="binding">
            {{ binding ? t('admin.partners.relations.binding') : t('admin.partners.relations.bind') }}
          </button>
        </form>
        <p class="-mt-3 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.partners.relations.rebindHint') }}</p>

        <div v-if="relationsLoading" class="flex justify-center py-8">
          <div class="h-7 w-7 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>
        <div v-else-if="relations.length === 0" class="rounded-xl border border-dashed border-gray-300 py-8 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
          {{ t('admin.partners.relations.empty') }}
        </div>
        <div v-else class="overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-700">
          <table class="w-full min-w-[760px] text-left text-sm">
            <thead class="bg-gray-50 text-gray-500 dark:bg-dark-800 dark:text-dark-400">
              <tr>
                <th class="px-4 py-3 font-medium">{{ t('admin.partners.relations.columns.user') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('admin.partners.relations.columns.source') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('admin.partners.relations.columns.boundAt') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('admin.partners.relations.columns.rates') }}</th>
                <th class="px-4 py-3 text-right font-medium">{{ t('admin.partners.relations.columns.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="relation in relations" :key="relation.user_id">
                <td class="px-4 py-3">
                  <div class="font-mono text-xs text-gray-500">#{{ relation.user_id }}</div>
                  <div class="font-medium text-gray-900 dark:text-white">{{ relation.email || '-' }}</div>
                  <div class="text-xs text-gray-500 dark:text-dark-400">{{ relation.username || '-' }}</div>
                </td>
                <td class="px-4 py-3 text-gray-700 dark:text-gray-300">{{ boundSourceLabel(relation.bound_source) }}</td>
                <td class="px-4 py-3 text-gray-700 dark:text-gray-300">{{ formatDateTime(relation.created_at) }}</td>
                <td class="px-4 py-3">
                  <div v-if="relation.rates.length" class="flex max-w-md flex-wrap gap-1.5">
                    <span v-for="rate in relation.rates" :key="rate.group_id" class="badge badge-gray whitespace-nowrap">
                      {{ rate.group_name }}: {{ rate.rate_multiplier == null ? t('admin.partners.relations.defaultRate') : `×${rate.rate_multiplier}` }}
                    </span>
                  </div>
                  <span v-else class="text-gray-400">-</span>
                </td>
                <td class="px-4 py-3 text-right">
                  <button class="btn btn-ghost btn-sm text-red-600 dark:text-red-400" @click="requestUnbind(relation)">
                    {{ t('admin.partners.relations.unbind') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </BaseDialog>

    <BaseDialog
      :show="settlementsOpen"
      :title="t('admin.partners.settlements.title', { email: selectedPartner?.email || '' })"
      width="wide"
      @close="settlementsOpen = false"
    >
      <div v-if="selectedPartner" class="space-y-5">
        <div class="grid gap-3 sm:grid-cols-3">
          <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-800">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.partners.settlements.totalCommission') }}</p>
            <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">{{ formatCurrency(selectedPartner.total_commission) }}</p>
          </div>
          <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-800">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.partners.settlements.settledTotal') }}</p>
            <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">{{ formatCurrency(selectedPartner.settled_total) }}</p>
          </div>
          <div class="rounded-xl bg-amber-50 p-4 dark:bg-amber-900/20">
            <p class="text-sm text-amber-700 dark:text-amber-300">{{ t('admin.partners.settlements.pending') }}</p>
            <p class="mt-1 text-xl font-semibold text-amber-700 dark:text-amber-300">{{ formatCurrency(selectedPartner.pending_settlement) }}</p>
          </div>
        </div>

        <form class="rounded-xl border border-gray-200 p-4 dark:border-dark-700" @submit.prevent="registerSettlement">
          <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.partners.settlements.register') }}</h4>
          <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.partners.settlements.offlineHint') }}</p>
          <div class="mt-4 grid gap-4 sm:grid-cols-2">
            <div>
              <label class="input-label">{{ t('admin.partners.settlements.amount') }}</label>
              <input v-model="settlementForm.amount" type="number" min="0.01" step="0.01" class="input" />
            </div>
            <div>
              <label class="input-label">{{ t('admin.partners.settlements.note') }}</label>
              <input v-model="settlementForm.note" type="text" class="input" :placeholder="t('admin.partners.settlements.notePlaceholder')" />
            </div>
          </div>
          <div class="mt-4 flex justify-end">
            <button type="submit" class="btn btn-primary" :disabled="settlementSaving || selectedPartner.pending_settlement <= 0">
              {{ settlementSaving ? t('common.saving') : t('admin.partners.settlements.register') }}
            </button>
          </div>
        </form>

        <div>
          <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.partners.settlements.history') }}</h4>
          <div v-if="settlementsLoading" class="flex justify-center py-6">
            <div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
          </div>
          <div v-else-if="settlements.length === 0" class="rounded-xl border border-dashed border-gray-300 py-6 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
            {{ t('admin.partners.settlements.empty') }}
          </div>
          <div v-else class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
            <table class="w-full text-left text-sm">
              <thead class="bg-gray-50 text-gray-500 dark:bg-dark-800 dark:text-dark-400">
                <tr>
                  <th class="px-4 py-3 font-medium">{{ t('admin.partners.settlements.columns.time') }}</th>
                  <th class="px-4 py-3 font-medium">{{ t('admin.partners.settlements.columns.amount') }}</th>
                  <th class="px-4 py-3 font-medium">{{ t('admin.partners.settlements.columns.method') }}</th>
                  <th class="px-4 py-3 font-medium">{{ t('admin.partners.settlements.columns.operator') }}</th>
                  <th class="px-4 py-3 font-medium">{{ t('admin.partners.settlements.columns.note') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                <tr v-for="item in settlements" :key="item.id">
                  <td class="px-4 py-3 text-gray-700 dark:text-gray-300">{{ formatDateTime(item.created_at) }}</td>
                  <td class="px-4 py-3 font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.amount) }}</td>
                  <td class="px-4 py-3 text-gray-700 dark:text-gray-300">{{ t('admin.partners.settlements.offline') }}</td>
                  <td class="px-4 py-3 font-mono text-gray-700 dark:text-gray-300">{{ item.operator_id == null ? '-' : `#${item.operator_id}` }}</td>
                  <td class="max-w-48 truncate px-4 py-3 text-gray-700 dark:text-gray-300" :title="item.note">{{ item.note || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </BaseDialog>

    <ConfirmDialog
      :show="confirmState !== null"
      :title="confirmTitle"
      :message="confirmMessage"
      :confirm-text="confirmState === 'remove' ? t('admin.partners.remove') : t('admin.partners.relations.unbind')"
      danger
      @confirm="confirmDangerAction"
      @cancel="cancelDangerAction"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import type { Column } from '@/components/common/types'
import { adminAPI } from '@/api/admin'
import type { AdminPartner, AdminPartnerRelation, AdminPartnerSettlement } from '@/api/admin/partners'
import type { PartnerBoundSource, PartnerStatus } from '@/api/partner'
import { useAppStore } from '@/stores/app'
import { useClipboard } from '@/composables/useClipboard'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const appStore = useAppStore()
const { copyToClipboard } = useClipboard()

const partners = ref<AdminPartner[]>([])
const loading = ref(false)
const filters = reactive({ query: '', status: '' as PartnerStatus | '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const columns = computed<Column[]>(() => [
  { key: 'user', label: t('admin.partners.columns.user') },
  { key: 'min_rate_multiplier', label: t('admin.partners.columns.minRate') },
  { key: 'status', label: t('admin.partners.columns.status') },
  { key: 'downline_count', label: t('admin.partners.columns.downlines') },
  { key: 'total_commission', label: t('admin.partners.columns.totalCommission') },
  { key: 'pending_settlement', label: t('admin.partners.columns.pendingSettlement') },
  { key: 'actions', label: t('admin.partners.columns.actions') },
])
const statusFilterOptions = computed(() => [
  { value: '', label: t('admin.partners.allStatus') },
  { value: 'active', label: t('common.active') },
  { value: 'disabled', label: t('admin.partners.disabled') },
])
const partnerStatusOptions = computed(() => statusFilterOptions.value.slice(1))

const editorOpen = ref(false)
const editorSaving = ref(false)
const editingPartner = ref<AdminPartner | null>(null)
const editorForm = reactive({ identifier: '', minRate: '0.5', status: 'active' as PartnerStatus })

const selectedPartner = ref<AdminPartner | null>(null)
const relationsOpen = ref(false)
const relationsLoading = ref(false)
const relations = ref<AdminPartnerRelation[]>([])
const bindUserId = ref('')
const binding = ref(false)

const settlementsOpen = ref(false)
const settlementsLoading = ref(false)
const settlements = ref<AdminPartnerSettlement[]>([])
const settlementSaving = ref(false)
const settlementForm = reactive({ amount: '', note: '' })

const confirmState = ref<'remove' | 'unbind' | null>(null)
const relationToUnbind = ref<AdminPartnerRelation | null>(null)

const inviteLink = computed(() => {
  const affCode = selectedPartner.value?.aff_code
  if (!affCode || typeof window === 'undefined') return ''
  return `${window.location.origin}/register?aff=${encodeURIComponent(affCode)}`
})
const confirmTitle = computed(() =>
  confirmState.value === 'remove' ? t('admin.partners.confirmRemoveTitle') : t('admin.partners.relations.confirmUnbindTitle'),
)
const confirmMessage = computed(() => {
  if (confirmState.value === 'remove') {
    return t('admin.partners.confirmRemoveMessage', { email: selectedPartner.value?.email || '' })
  }
  return t('admin.partners.relations.confirmUnbindMessage', { email: relationToUnbind.value?.email || '' })
})

async function loadPartners() {
  loading.value = true
  try {
    const res = await adminAPI.partners.list({
      query: filters.query,
      status: filters.status,
      page: pagination.page,
      page_size: pagination.page_size,
    })
    partners.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.loadFailed')))
  } finally {
    loading.value = false
  }
}

function debounceLoad() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(reloadFromFirstPage, 300)
}
function reloadFromFirstPage() {
  pagination.page = 1
  void loadPartners()
}
function handlePageChange(page: number) {
  pagination.page = page
  void loadPartners()
}
function handlePageSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  void loadPartners()
}

function openCreate() {
  editingPartner.value = null
  Object.assign(editorForm, { identifier: '', minRate: '0.5', status: 'active' })
  editorOpen.value = true
}
function openEdit(partner: AdminPartner) {
  editingPartner.value = partner
  Object.assign(editorForm, {
    identifier: String(partner.user_id),
    minRate: String(partner.min_rate_multiplier),
    status: partner.status,
  })
  editorOpen.value = true
}

async function resolveUserId(identifier: string): Promise<number | null> {
  const trimmed = identifier.trim()
  if (/^\d+$/.test(trimmed)) {
    const id = Number(trimmed)
    return Number.isSafeInteger(id) && id > 0 ? id : null
  }
  const res = await adminAPI.users.list(1, 20, { search: trimmed })
  const exact = res.items.find((user) => user.email.toLowerCase() === trimmed.toLowerCase())
  return exact?.id ?? null
}

async function savePartner() {
  const minRate = Number(editorForm.minRate)
  if (!Number.isFinite(minRate) || minRate <= 0 || minRate > 1) {
    appStore.showError(t('admin.partners.editor.invalidMinRate'))
    return
  }
  editorSaving.value = true
  try {
    const userId = editingPartner.value?.user_id ?? await resolveUserId(editorForm.identifier)
    if (!userId) {
      appStore.showError(t('admin.partners.editor.userNotFound'))
      return
    }
    await adminAPI.partners.upsert(userId, {
      min_rate_multiplier: minRate,
      status: editorForm.status,
    })
    appStore.showSuccess(t('admin.partners.editor.saved'))
    editorOpen.value = false
    await loadPartners()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.editor.saveFailed')))
  } finally {
    editorSaving.value = false
  }
}

async function loadRelations() {
  if (!selectedPartner.value) return
  relationsLoading.value = true
  try {
    const res = await adminAPI.partners.listRelations(selectedPartner.value.user_id)
    relations.value = res.items || []
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.relations.loadFailed')))
  } finally {
    relationsLoading.value = false
  }
}
function openRelations(partner: AdminPartner) {
  selectedPartner.value = partner
  relations.value = []
  bindUserId.value = ''
  relationsOpen.value = true
  void loadRelations()
}
async function bindUser() {
  if (!selectedPartner.value) return
  const userId = Number(bindUserId.value)
  if (!Number.isSafeInteger(userId) || userId <= 0) {
    appStore.showError(t('admin.partners.relations.invalidUserId'))
    return
  }
  binding.value = true
  try {
    await adminAPI.partners.bindRelation(selectedPartner.value.user_id, userId)
    appStore.showSuccess(t('admin.partners.relations.bound'))
    bindUserId.value = ''
    await Promise.all([loadRelations(), loadPartners()])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.relations.bindFailed')))
  } finally {
    binding.value = false
  }
}
async function copyInviteLink() {
  if (inviteLink.value) {
    await copyToClipboard(inviteLink.value, t('admin.partners.relations.linkCopied'))
  }
}
function boundSourceLabel(source: PartnerBoundSource): string {
  return t(`admin.partners.relations.sources.${source}`)
}

function openSettlements(partner: AdminPartner) {
  selectedPartner.value = partner
  settlementForm.amount = ''
  settlementForm.note = ''
  settlements.value = []
  settlementsOpen.value = true
  void loadSettlements()
}
async function loadSettlements() {
  if (!selectedPartner.value) return
  settlementsLoading.value = true
  try {
    const res = await adminAPI.partners.listSettlements(selectedPartner.value.user_id)
    settlements.value = res.items || []
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.settlements.loadFailed')))
  } finally {
    settlementsLoading.value = false
  }
}
async function registerSettlement() {
  if (!selectedPartner.value) return
  const amount = Number(settlementForm.amount)
  if (!Number.isFinite(amount) || amount <= 0) {
    appStore.showError(t('admin.partners.settlements.invalidAmount'))
    return
  }
  if (amount > selectedPartner.value.pending_settlement) {
    appStore.showError(t('admin.partners.settlements.amountExceedsPending'))
    return
  }
  settlementSaving.value = true
  try {
    await adminAPI.partners.createSettlement(selectedPartner.value.user_id, {
      amount,
      method: 'offline',
      note: settlementForm.note.trim() || undefined,
    })
    appStore.showSuccess(t('admin.partners.settlements.registered'))
    settlementForm.amount = ''
    settlementForm.note = ''
    await Promise.all([loadSettlements(), loadPartners()])
    const updated = partners.value.find((item) => item.user_id === selectedPartner.value?.user_id)
    if (updated) selectedPartner.value = updated
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.settlements.registerFailed')))
  } finally {
    settlementSaving.value = false
  }
}

function requestRemove(partner: AdminPartner) {
  selectedPartner.value = partner
  relationToUnbind.value = null
  confirmState.value = 'remove'
}
function requestUnbind(relation: AdminPartnerRelation) {
  relationToUnbind.value = relation
  confirmState.value = 'unbind'
}
function cancelDangerAction() {
  confirmState.value = null
  relationToUnbind.value = null
}
async function confirmDangerAction() {
  const action = confirmState.value
  try {
    if (action === 'remove' && selectedPartner.value) {
      await adminAPI.partners.remove(selectedPartner.value.user_id)
      appStore.showSuccess(t('admin.partners.removed'))
      relationsOpen.value = false
      settlementsOpen.value = false
      await loadPartners()
    } else if (action === 'unbind' && selectedPartner.value && relationToUnbind.value) {
      await adminAPI.partners.unbindRelation(selectedPartner.value.user_id, relationToUnbind.value.user_id)
      appStore.showSuccess(t('admin.partners.relations.unbound'))
      await Promise.all([loadRelations(), loadPartners()])
    }
    cancelDangerAction()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.partners.actionFailed')))
  }
}

onMounted(() => {
  void loadPartners()
})
onUnmounted(() => {
  if (debounceTimer) clearTimeout(debounceTimer)
})
</script>
