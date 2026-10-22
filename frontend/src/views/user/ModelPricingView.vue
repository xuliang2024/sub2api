<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- 筛选与操作 -->
      <div class="flex flex-col justify-between gap-3 xl:flex-row xl:items-center">
        <div class="flex flex-1 flex-wrap items-center gap-3">
          <div class="relative w-full sm:w-80">
            <Icon
              name="search"
              size="md"
              class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-dark-500"
            />
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('modelPricing.searchPlaceholder')"
              class="input pl-10"
            />
          </div>
          <Select v-model="selectedGroup" :options="groupOptions" class="w-full sm:w-48" />
          <Select v-model="sortKey" :options="sortOptions" class="w-full sm:w-56" />
        </div>

        <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 xl:w-auto">
          <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
            <Toggle v-model="onlyPriced" />
            <span class="cursor-pointer select-none" @click="onlyPriced = !onlyPriced">
              {{ t('modelPricing.onlyPriced') }}
            </span>
          </div>
          <button
            @click="load"
            :disabled="loading"
            class="btn btn-secondary"
            :title="t('common.refresh', 'Refresh')"
          >
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </div>

      <!-- 价格表 -->
      <div class="card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full min-w-[930px] table-fixed border-collapse text-[13px]">
            <colgroup>
              <col class="w-[148px]" />
              <col class="w-[82px]" />
              <col class="w-[116px]" />
              <col class="w-[66px]" />
              <col class="w-[60px]" />
              <col class="w-[60px]" />
              <col class="w-[60px]" />
              <col class="w-[60px]" />
              <col class="w-[54px]" />
              <col class="w-[100px]" />
              <col class="w-[124px]" />
            </colgroup>
            <thead>
              <tr
                class="border-b border-gray-100 bg-gray-50/60 text-[11px] uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-400"
              >
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.model') }}
                </th>
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.platform') }}
                </th>
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.channel') }}
                </th>
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.billingMode') }}
                </th>
                <th
                  colspan="4"
                  class="border-b border-gray-100 px-2 pt-2 text-center align-bottom font-medium normal-case dark:border-dark-700"
                >
                  {{ t('modelPricing.columns.tokenGroup') }}
                  <span class="mt-0.5 block text-[10px] font-normal tracking-normal text-gray-400 dark:text-dark-500">
                    {{ t('modelPricing.units.perMillion') }}
                  </span>
                </th>
                <th rowspan="2" class="px-2 py-2 text-right align-middle font-medium">
                  {{ t('modelPricing.columns.perRequest') }}
                  <span class="mt-0.5 block text-[10px] font-normal tracking-normal text-gray-400 dark:text-dark-500">
                    {{ t('modelPricing.units.perRequest') }}
                  </span>
                </th>
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.groups') }}
                </th>
                <th
                  rowspan="2"
                  class="px-2 py-2 text-right align-middle font-medium text-primary-700 dark:text-primary-400"
                >
                  {{ t('modelPricing.columns.paid') }}
                </th>
              </tr>
              <tr
                class="border-b border-gray-100 text-[11px] font-medium text-gray-500 dark:border-dark-700 dark:text-dark-400"
              >
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.input') }}</th>
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.output') }}</th>
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.cacheRead') }}</th>
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.cacheWrite') }}</th>
              </tr>
            </thead>

            <tbody v-if="loading">
              <tr>
                <td colspan="11" class="py-16 text-center">
                  <Icon name="refresh" size="lg" class="inline-block animate-spin text-gray-400" />
                </td>
              </tr>
            </tbody>
            <tbody v-else-if="filteredRows.length === 0">
              <tr>
                <td colspan="11" class="py-16 text-center">
                  <Icon name="inbox" size="xl" class="mx-auto mb-3 h-10 w-10 text-gray-300 dark:text-dark-500" />
                  <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('modelPricing.empty') }}</p>
                </td>
              </tr>
            </tbody>
            <tbody v-else>
              <tr
                v-for="row in filteredRows"
                :key="row.key"
                class="border-b border-gray-100 transition-colors last:border-b-0 hover:bg-gray-50/50 dark:border-dark-800 dark:hover:bg-dark-800/40"
              >
                <td class="px-2 py-2.5 align-middle">
                  <div class="flex items-center gap-1.5">
                    <span
                      class="truncate font-medium text-gray-900 dark:text-white"
                      :title="row.model"
                    >{{ row.model }}</span>
                    <span
                      v-if="row.tiered"
                      class="flex-shrink-0 rounded bg-gray-100 px-1 py-0.5 text-[10px] text-gray-500 dark:bg-dark-700 dark:text-dark-300"
                      :title="t('modelPricing.tieredHint')"
                    >{{ t('modelPricing.tiered') }}</span>
                  </div>
                </td>
                <td class="px-2 py-2.5 align-middle">
                  <span
                    :class="[
                      'inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium uppercase',
                      platformBadgeClass(row.platform),
                    ]"
                  >
                    <PlatformIcon :platform="row.platform as GroupPlatform" size="xs" />
                    {{ row.platform }}
                  </span>
                </td>
                <td class="px-2 py-2.5 align-middle text-xs text-gray-500 dark:text-dark-400">
                  <span class="block truncate" :title="row.channel">{{ row.channel }}</span>
                </td>
                <td class="px-2 py-2.5 align-middle text-xs text-gray-500 dark:text-dark-400">
                  {{ billingLabel(row.billingMode) }}
                </td>
                <td class="px-2 py-2.5 text-right tabular-nums text-gray-800 dark:text-gray-200">
                  {{ money(row.input, PER_MILLION) }}
                </td>
                <td class="px-2 py-2.5 text-right tabular-nums text-gray-800 dark:text-gray-200">
                  {{ money(row.output, PER_MILLION) }}
                </td>
                <td class="px-2 py-2.5 text-right tabular-nums text-gray-800 dark:text-gray-200">
                  {{ money(row.cacheRead, PER_MILLION) }}
                </td>
                <td class="px-2 py-2.5 text-right tabular-nums text-gray-800 dark:text-gray-200">
                  {{ money(row.cacheWrite, PER_MILLION) }}
                </td>
                <td class="px-2 py-2.5 text-right tabular-nums text-gray-800 dark:text-gray-200">
                  {{ row.perRequest == null ? '-' : money(row.perRequest, 1) }}
                </td>
                <td class="px-2 py-2.5 align-middle">
                  <div v-if="row.groups.length" class="flex flex-wrap items-center gap-1">
                    <span
                      v-for="g in row.groups"
                      :key="g.id"
                      class="inline-flex items-center gap-1 text-[11px] text-gray-600 dark:text-dark-300"
                    >
                      <span class="max-w-[80px] truncate" :title="g.name">{{ g.name }}</span>
                      <span
                        class="rounded bg-primary-50 px-1 text-[10px] text-primary-700 dark:bg-primary-900/30 dark:text-primary-300"
                      >×{{ formatRate(g.rate) }}</span>
                    </span>
                  </div>
                  <span v-else class="text-xs text-gray-400">-</span>
                </td>
                <td class="px-2 py-2.5 text-right align-middle">
                  <span
                    class="whitespace-nowrap font-medium text-primary-700 dark:text-primary-400"
                    :title="row.rateVaries ? t('modelPricing.rateVariesHint') : undefined"
                  >
                    <template v-if="row.billingMode === BILLING_MODE_TOKEN">
                      {{ money(row.paidInput, PER_MILLION) }} / {{ money(row.paidOutput, PER_MILLION) }}
                    </template>
                    <template v-else>
                      {{ row.paidPerRequest == null ? '-' : money(row.paidPerRequest, 1) }}
                      <span class="font-normal text-gray-400 dark:text-dark-500">{{ t('modelPricing.units.perRequestShort') }}</span>
                    </template>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div
          v-if="hiddenUnpriced > 0"
          class="flex items-center gap-1.5 border-t border-gray-100 bg-gray-50/40 px-4 py-2 text-xs text-gray-400 dark:border-dark-700 dark:bg-dark-800/40 dark:text-dark-500"
        >
          <Icon name="infoCircle" size="xs" class="h-3.5 w-3.5 flex-shrink-0" />
          {{ t('modelPricing.hiddenUnpriced', { count: hiddenUnpriced }) }}
        </div>
      </div>

      <p class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-400 dark:text-dark-500">
        <span>{{ t('modelPricing.notes.base') }}</span>
        <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
        <span>{{ t('modelPricing.notes.paid') }}</span>
        <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
        <span>{{ t('modelPricing.notes.peak') }}</span>
        <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
        <span>{{ t('modelPricing.notes.unit') }}</span>
      </p>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import Select, { type SelectOption } from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import userChannelsAPI, {
  type UserAvailableChannel,
  type UserChannelPlatformSection,
  type UserSupportedModel,
} from '@/api/channels'
import userGroupsAPI from '@/api/groups'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { formatScaled } from '@/utils/pricing'
import { platformBadgeClass } from '@/utils/platformColors'
import {
  BILLING_MODE_TOKEN,
  BILLING_MODE_PER_REQUEST,
  BILLING_MODE_IMAGE,
  BILLING_MODE_VIDEO,
  type BillingMode,
} from '@/constants/channel'
import type { GroupPlatform } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

/** 表格按 1M token 展示单价；按次/按图片按单次展示。 */
const PER_MILLION = 1_000_000

type SortKey = 'input_asc' | 'input_desc' | 'output_asc' | 'output_desc' | 'cache_asc' | 'name'

interface RowGroup {
  id: number
  name: string
  /** 生效倍率：用户专属倍率优先，否则用分组默认倍率。 */
  rate: number
}

interface PricingRow {
  key: string
  model: string
  platform: string
  channel: string
  billingMode: BillingMode
  input: number | null
  output: number | null
  cacheRead: number | null
  cacheWrite: number | null
  perRequest: number | null
  /** true = 仅配置了阶梯价，表格展示首档价格。 */
  tiered: boolean
  groups: RowGroup[]
  /** 实付计算所用倍率（多个分组时取最低倍率）。 */
  rate: number
  /** 同一模型在不同分组下倍率不一致时为 true（实付按最低倍率，"起"价）。 */
  rateVaries: boolean
  paidInput: number | null
  paidOutput: number | null
  paidPerRequest: number | null
  hasPricing: boolean
}

const channels = ref<UserAvailableChannel[]>([])
const userGroupRates = ref<Record<number, number>>({})
const loading = ref(false)
const searchQuery = ref('')
const selectedGroup = ref<string | number>('all')
const sortKey = ref<string>('input_asc')
const onlyPriced = ref(true)

const sortOptions = computed<SelectOption[]>(() => [
  { value: 'input_asc', label: t('modelPricing.sort.inputAsc') },
  { value: 'input_desc', label: t('modelPricing.sort.inputDesc') },
  { value: 'output_asc', label: t('modelPricing.sort.outputAsc') },
  { value: 'output_desc', label: t('modelPricing.sort.outputDesc') },
  { value: 'cache_asc', label: t('modelPricing.sort.cacheAsc') },
  { value: 'name', label: t('modelPricing.sort.name') },
])

/** 分组下拉：只列出当前用户可见分组，值用分组 ID。 */
const groupOptions = computed<SelectOption[]>(() => {
  const seen = new Map<number, string>()
  for (const ch of channels.value) {
    for (const section of ch.platforms) {
      for (const g of section.groups) {
        if (!seen.has(g.id)) seen.set(g.id, g.name)
      }
    }
  }
  const options: SelectOption[] = [{ value: 'all', label: t('modelPricing.allGroups') }]
  for (const [id, name] of [...seen.entries()].sort((a, b) => a[0] - b[0])) {
    options.push({ value: id, label: name })
  }
  return options
})

function resolvePrices(model: UserSupportedModel): {
  mode: BillingMode
  input: number | null
  output: number | null
  cacheRead: number | null
  cacheWrite: number | null
  perRequest: number | null
  tiered: boolean
} {
  const p = model.pricing
  const mode = (p?.billing_mode ?? BILLING_MODE_TOKEN) as BillingMode
  let input = p?.input_price ?? null
  let output = p?.output_price ?? null
  let cacheRead = p?.cache_read_price ?? null
  let cacheWrite = p?.cache_write_price ?? null
  let perRequest: number | null = null
  let tiered = false

  if (mode === BILLING_MODE_TOKEN) {
    // 只配置阶梯价时用首档做展示，完整阶梯仍在「可用渠道」的悬停气泡里。
    if (input == null && output == null && cacheRead == null && cacheWrite == null) {
      const first = p?.intervals?.[0]
      if (first) {
        input = first.input_price ?? null
        output = first.output_price ?? null
        cacheRead = first.cache_read_price ?? null
        cacheWrite = first.cache_write_price ?? null
        tiered = true
      }
    }
  } else if (mode === BILLING_MODE_PER_REQUEST) {
    perRequest = p?.per_request_price ?? null
  } else if (mode === BILLING_MODE_IMAGE) {
    perRequest = p?.image_output_price ?? p?.per_request_price ?? null
  } else if (mode === BILLING_MODE_VIDEO) {
    perRequest = p?.per_request_price ?? p?.image_output_price ?? null
  }

  return { mode, input, output, cacheRead, cacheWrite, perRequest, tiered }
}

function buildGroups(
  section: UserChannelPlatformSection,
  rates: Record<number, number>,
): RowGroup[] {
  return section.groups.map((g) => ({
    id: g.id,
    name: g.name,
    rate: rates[g.id] ?? g.rate_multiplier,
  }))
}

const allRows = computed<PricingRow[]>(() => {
  const out: PricingRow[] = []
  for (const ch of channels.value) {
    for (const section of ch.platforms) {
      const groups = buildGroups(section, userGroupRates.value)
      const rate = groups.length ? Math.min(...groups.map((g) => g.rate)) : 1
      const rateVaries = new Set(groups.map((g) => g.rate)).size > 1
      for (const model of section.supported_models) {
        const prices = resolvePrices(model)
        const hasPricing =
          prices.input != null ||
          prices.output != null ||
          prices.cacheRead != null ||
          prices.cacheWrite != null ||
          prices.perRequest != null
        out.push({
          key: `${ch.name}|${section.platform}|${model.name}`,
          model: model.name,
          platform: section.platform || model.platform,
          channel: ch.name,
          billingMode: prices.mode,
          input: prices.input,
          output: prices.output,
          cacheRead: prices.cacheRead,
          cacheWrite: prices.cacheWrite,
          perRequest: prices.perRequest,
          tiered: prices.tiered,
          groups,
          rate,
          rateVaries,
          paidInput: prices.input == null ? null : prices.input * rate,
          paidOutput: prices.output == null ? null : prices.output * rate,
          paidPerRequest: prices.perRequest == null ? null : prices.perRequest * rate,
          hasPricing,
        })
      }
    }
  }
  return out
})

function matchesSearch(row: PricingRow, q: string): boolean {
  if (!q) return true
  return (
    row.model.toLowerCase().includes(q) ||
    row.channel.toLowerCase().includes(q) ||
    row.groups.some((g) => g.name.toLowerCase().includes(q))
  )
}

function matchesGroup(row: PricingRow, group: string | number): boolean {
  if (group === 'all') return true
  return row.groups.some((g) => g.id === group)
}

function priceForSort(row: PricingRow, key: SortKey): number | null {
  if (key.startsWith('input')) return row.input
  if (key.startsWith('output')) return row.output
  if (key.startsWith('cache')) return row.cacheRead
  return null
}

const filteredRows = computed<PricingRow[]>(() => {
  const q = searchQuery.value.trim().toLowerCase()
  let rows = allRows.value
  if (onlyPriced.value) rows = rows.filter((r) => r.hasPricing)
  rows = rows.filter((r) => matchesGroup(r, selectedGroup.value) && matchesSearch(r, q))

  const sorted = [...rows]
  if (sortKey.value === 'name') {
    sorted.sort((a, b) => a.model.localeCompare(b.model))
    return sorted
  }
  const desc = sortKey.value.endsWith('_desc')
  sorted.sort((a, b) => {
    const av = priceForSort(a, sortKey.value as SortKey)
    const bv = priceForSort(b, sortKey.value as SortKey)
    // 无价的排最后，不与 0 混排。
    if (av == null && bv == null) return a.model.localeCompare(b.model)
    if (av == null) return 1
    if (bv == null) return -1
    if (av === bv) return a.model.localeCompare(b.model)
    return desc ? bv - av : av - bv
  })
  return sorted
})

/** 命中筛选但被「仅看有定价」隐藏的模型数量。 */
const hiddenUnpriced = computed(() => {
  if (!onlyPriced.value) return 0
  const q = searchQuery.value.trim().toLowerCase()
  return allRows.value.filter(
    (r) => !r.hasPricing && matchesGroup(r, selectedGroup.value) && matchesSearch(r, q),
  ).length
})

/**
 * 表格内不带 "$"（表头已标注 USD，11 列时省下约 8px/格）；
 * formatScaled 负责精度与 null 处理。
 */
function money(value: number | null, scale: number): string {
  const formatted = formatScaled(value, scale, 2)
  return formatted === '-' ? '-' : formatted.replace(/^\$/, '')
}

function formatRate(rate: number): string {
  if (!Number.isFinite(rate)) return '1.0'
  return String(Number(rate.toFixed(4)))
}

function billingLabel(mode: BillingMode): string {
  switch (mode) {
    case BILLING_MODE_PER_REQUEST:
      return t('modelPricing.billing.perRequest')
    case BILLING_MODE_IMAGE:
      return t('modelPricing.billing.image')
    case BILLING_MODE_VIDEO:
      return t('modelPricing.billing.video')
    default:
      return t('modelPricing.billing.token')
  }
}

async function load() {
  loading.value = true
  try {
    // 专属倍率失败不阻塞价格表：降级为分组默认倍率。
    const [list, rates] = await Promise.all([
      userChannelsAPI.getAvailable(),
      userGroupsAPI.getUserGroupRates().catch((err: unknown) => {
        console.error('Failed to load user group rates:', err)
        return {} as Record<number, number>
      }),
    ])
    channels.value = list
    userGroupRates.value = rates
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
