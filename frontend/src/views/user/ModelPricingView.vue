<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- 筛选与操作 -->
      <div class="flex flex-col justify-between gap-3 xl:flex-row xl:items-center">
        <div class="flex flex-1 flex-wrap items-center gap-3">
          <div class="relative w-full sm:w-72">
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
          <div
            v-if="platformOptions.length > 1"
            class="flex items-center gap-1 rounded-lg bg-gray-100 p-0.5 dark:bg-dark-800"
          >
            <button
              v-for="p in platformOptions"
              :key="p.value"
              type="button"
              class="rounded-md px-2.5 py-1 text-xs transition-colors"
              :class="
                selectedPlatform === p.value
                  ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'
              "
              @click="selectedPlatform = p.value"
            >
              {{ p.label }}
            </button>
          </div>
          <Select v-model="selectedGroup" :options="groupOptions" class="w-full sm:w-40" />
          <Select v-model="sortKey" :options="sortOptions" class="w-full sm:w-52" />
        </div>

        <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 xl:w-auto">
          <div
            v-if="SHOW_OFFICIAL_DISCOUNT"
            class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300"
          >
            <Toggle v-model="onlyDiscounted" />
            <span class="cursor-pointer select-none" @click="onlyDiscounted = !onlyDiscounted">
              {{ t('modelPricing.onlyDiscounted') }}
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

      <!-- 汇总 -->
      <div
        v-if="!loading && rows.length"
        class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-dark-400"
      >
        <span>{{ t('modelPricing.stats.summary', { models: rows.length, offers: totalOffers }) }}</span>
        <template v-if="SHOW_OFFICIAL_DISCOUNT">
          <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
          <span>
            {{ t('modelPricing.stats.avgDiscount') }}
            <span class="font-medium text-gray-700 dark:text-dark-200">{{ avgDiscountText }}</span>
          </span>
        </template>
        <template v-if="SHOW_OFFICIAL_DISCOUNT && aboveOfficialCount > 0">
          <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
          <span class="text-amber-600 dark:text-amber-400">
            {{ t('modelPricing.stats.aboveOfficial', { count: aboveOfficialCount }) }}
          </span>
        </template>
      </div>

      <!-- 价格表：按模型聚合，展开看分组 -->
      <div class="card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full min-w-[1150px] table-fixed border-collapse text-[13px]">
            <colgroup>
              <col class="w-[212px]" />
              <col class="w-[86px]" />
              <col class="w-[86px]" />
              <col class="w-[118px]" />
              <col class="w-[94px]" />
              <col class="w-[94px]" />
              <col v-if="SHOW_OFFICIAL_DISCOUNT" class="w-[122px]" />
              <col class="w-[112px]" />
              <col class="w-[140px]" />
              <col class="w-[36px]" />
            </colgroup>
            <thead>
              <tr
                class="border-b border-gray-100 bg-gray-50/60 text-[11px] uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-400"
              >
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.model') }}
                </th>
                <th
                  colspan="3"
                  class="border-b border-gray-100 bg-gray-50/40 px-2 pt-2 text-center align-bottom font-medium normal-case dark:border-dark-700 dark:bg-dark-800/40"
                >
                  {{ t('modelPricing.blocks.official') }}
                  <span class="mt-0.5 block text-[10px] font-normal tracking-normal text-gray-400 dark:text-dark-500">
                    {{ t('modelPricing.blocks.unit') }}
                  </span>
                </th>
                <th
                  colspan="3"
                  class="border-b border-l border-gray-100 px-2 pt-2 text-center align-bottom font-medium normal-case dark:border-dark-700"
                >
                  {{ t('modelPricing.blocks.ours') }}
                  <span class="mt-0.5 block text-[10px] font-normal tracking-normal text-gray-400 dark:text-dark-500">
                    {{ t('modelPricing.blocks.unitPaid') }}
                  </span>
                </th>
                <th
                  v-if="SHOW_OFFICIAL_DISCOUNT"
                  rowspan="2"
                  class="px-2 py-2 text-right align-middle font-medium"
                >
                  {{ t('modelPricing.columns.discount') }}
                </th>
                <th rowspan="2" class="px-2 py-2 text-left align-middle font-medium">
                  {{ t('modelPricing.columns.groups') }}
                </th>
                <th rowspan="2" class="px-2 py-2"></th>
              </tr>
              <tr
                class="border-b border-gray-100 text-[11px] font-medium text-gray-500 dark:border-dark-700 dark:text-dark-400"
              >
                <th class="bg-gray-50/40 px-2 py-1.5 text-right font-medium dark:bg-dark-800/40">
                  {{ t('modelPricing.columns.input') }}
                </th>
                <th class="bg-gray-50/40 px-2 py-1.5 text-right font-medium dark:bg-dark-800/40">
                  {{ t('modelPricing.columns.output') }}
                </th>
                <th class="bg-gray-50/40 px-2 py-1.5 text-right font-medium dark:bg-dark-800/40">
                  {{ t('modelPricing.columns.cache') }}
                </th>
                <th class="border-l border-gray-100 px-2 py-1.5 text-right font-medium dark:border-dark-700">
                  {{ t('modelPricing.columns.input') }}
                </th>
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.output') }}</th>
                <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.cache') }}</th>
              </tr>
            </thead>

            <tbody v-if="loading">
              <tr>
                <td :colspan="tableColspan" class="py-16 text-center">
                  <Icon name="refresh" size="lg" class="inline-block animate-spin text-gray-400" />
                </td>
              </tr>
            </tbody>
            <tbody v-else-if="filteredRows.length === 0">
              <tr>
                <td :colspan="tableColspan" class="py-16 text-center">
                  <Icon name="inbox" size="xl" class="mx-auto mb-3 h-10 w-10 text-gray-300 dark:text-dark-500" />
                  <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('modelPricing.empty') }}</p>
                </td>
              </tr>
            </tbody>
            <tbody v-else>
              <template v-for="row in filteredRows" :key="row.name">
                <tr
                  class="cursor-pointer border-b border-gray-100 transition-colors last:border-b-0 hover:bg-gray-50/50 dark:border-dark-800 dark:hover:bg-dark-800/40"
                  :class="isOpen(row.name) && 'bg-gray-50/40 dark:bg-dark-800/30'"
                  tabindex="0"
                  role="button"
                  :aria-expanded="isOpen(row.name)"
                  @click="toggleRow(row.name)"
                  @keydown.enter.prevent="toggleRow(row.name)"
                  @keydown.space.prevent="toggleRow(row.name)"
                >
                  <!-- 模型 -->
                  <td class="px-2 py-2.5 align-middle">
                    <div class="flex items-center gap-1.5">
                      <span class="truncate font-medium text-gray-900 dark:text-white" :title="row.name">
                        {{ row.name }}
                      </span>
                      <span
                        v-if="row.tiered"
                        class="flex-shrink-0 rounded bg-violet-50 px-1 py-0.5 text-[10px] text-violet-700 dark:bg-violet-950/40 dark:text-violet-300"
                        :title="t('modelPricing.tieredHint')"
                      >
                        {{ t('modelPricing.tiered') }}
                      </span>
                    </div>
                    <div class="mt-1 flex items-center gap-1.5 text-[11px] text-gray-500 dark:text-dark-400">
                      <span
                        :class="[
                          'inline-flex items-center gap-0.5 rounded px-1 py-0.5 text-[9px] font-semibold uppercase',
                          platformBadgeClass(row.platform),
                        ]"
                      >
                        <PlatformIcon :platform="row.platform as GroupPlatform" size="xs" />
                        {{ row.platform }}
                      </span>
                      <span class="truncate" :title="row.best?.channel">{{ row.best?.channel }}</span>
                    </div>
                  </td>

                  <!-- 官方价 -->
                  <td class="bg-gray-50/40 px-2 py-2.5 text-right tabular-nums text-gray-500 dark:bg-dark-800/40 dark:text-dark-400">
                    {{ officialInput(row) }}
                  </td>
                  <td class="bg-gray-50/40 px-2 py-2.5 text-right tabular-nums text-gray-500 dark:bg-dark-800/40 dark:text-dark-400">
                    {{ officialOutput(row) }}
                  </td>
                  <td class="bg-gray-50/40 px-2 py-2.5 text-right tabular-nums text-gray-500 dark:bg-dark-800/40 dark:text-dark-400">
                    <div
                      v-for="c in officialCache(row)"
                      :key="c.label"
                      class="flex items-baseline justify-end gap-1 leading-5"
                    >
                      <span class="text-[9px] text-gray-400 dark:text-dark-500">{{ c.label }}</span>
                      <span>{{ c.value }}</span>
                    </div>
                  </td>

                  <!-- 本站价 -->
                  <td class="border-l border-gray-100 px-2 py-2.5 text-right font-medium tabular-nums text-gray-900 dark:border-dark-700 dark:text-white">
                    {{ paidInput(row) }}
                  </td>
                  <td class="px-2 py-2.5 text-right font-medium tabular-nums text-gray-900 dark:text-white">
                    {{ paidOutput(row) }}
                  </td>
                  <td class="px-2 py-2.5 text-right tabular-nums">
                    <div
                      v-for="c in paidCache(row)"
                      :key="c.label"
                      class="flex items-baseline justify-end gap-1 leading-5"
                      :class="c.above && 'text-amber-600 dark:text-amber-400'"
                      :title="c.above ? t('modelPricing.cacheAboveOfficial') : undefined"
                    >
                      <span class="text-[9px] text-gray-400 dark:text-dark-500">{{ c.label }}</span>
                      <span class="text-gray-900 dark:text-white">{{ c.value }}</span>
                    </div>
                  </td>

                  <!-- 折扣 -->
                  <td v-if="SHOW_OFFICIAL_DISCOUNT" class="px-2 py-2.5 text-right align-middle">
                    <div class="flex flex-col items-end gap-0.5">
                      <span
                        class="inline-flex items-center whitespace-nowrap rounded-md border px-1.5 py-0.5 text-xs font-semibold"
                        :class="discountClass(row.bestRatioMax ?? row.bestRatio)"
                        :title="t('modelPricing.discount.inputHint')"
                      >
                        {{ discountText(row.bestRatio, row.bestRatioMax) }}
                      </span>
                      <span
                        class="whitespace-nowrap text-[10px]"
                        :class="subDiscount(row)?.above ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400 dark:text-dark-500'"
                      >
                        {{ subDiscount(row)?.text }}
                      </span>
                    </div>
                  </td>

                  <!-- 分组 -->
                  <td class="px-2 py-2.5 align-middle">
                    <div class="flex flex-wrap items-center gap-1">
                      <span
                        v-if="row.best"
                        class="inline-flex items-center gap-0.5 rounded border border-primary-200 bg-primary-50 px-1 py-0.5 text-[10px] text-primary-700 dark:border-primary-900/60 dark:bg-primary-950/40 dark:text-primary-300"
                      >
                        <span class="max-w-[72px] truncate" :title="row.best.groupName">{{ row.best.groupName }}</span>
                        <b class="font-semibold">×{{ formatRate(row.best.rate) }}</b>
                      </span>
                      <span
                        v-if="row.groupNames.length > 1"
                        class="rounded border border-dashed border-gray-300 px-1 py-0.5 text-[10px] text-gray-400 dark:border-dark-600 dark:text-dark-500"
                      >
                        +{{ row.groupNames.length - 1 }}
                      </span>
                    </div>
                  </td>

                  <td class="px-2 py-2.5 text-center align-middle">
                    <Icon
                      name="chevronRight"
                      size="sm"
                      aria-hidden="true"
                      class="inline-block text-gray-400 transition-transform dark:text-dark-500"
                      :class="isOpen(row.name) && 'rotate-90'"
                    />
                  </td>
                </tr>

                <!-- 展开：分组 × 线路报价 -->
                <tr v-if="isOpen(row.name)" class="border-b border-gray-100 dark:border-dark-800">
                  <td :colspan="tableColspan" class="bg-gray-50/40 px-4 pb-4 pt-2 dark:bg-dark-800/30">
                    <div class="mb-2 flex flex-wrap items-center gap-2">
                      <span
                        class="rounded bg-gray-200/70 px-1.5 py-0.5 text-[10px] font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-300"
                      >
                        {{ t('modelPricing.expand.officialRef') }}
                      </span>
                      <span class="text-[11px] text-gray-600 dark:text-dark-300">{{ officialSummary(row) }}</span>
                      <span class="ml-auto text-[11px] text-gray-400 dark:text-dark-500">
                        {{ t('modelPricing.expand.summary', { groups: row.groupNames.length, offers: row.offers.length }) }}
                      </span>
                    </div>
                    <table class="w-full table-fixed border-collapse overflow-hidden rounded-lg border border-gray-200 text-[12px] dark:border-dark-700">
                      <colgroup>
                        <col class="w-[150px]" />
                        <col class="w-[170px]" />
                        <col class="w-[62px]" />
                        <col class="w-[90px]" />
                        <col class="w-[90px]" />
                        <col class="w-[96px]" />
                        <col class="w-[96px]" />
                        <col v-if="SHOW_OFFICIAL_DISCOUNT" class="w-[190px]" />
                        <col />
                      </colgroup>
                      <thead>
                        <tr class="bg-gray-100/70 text-[10px] text-gray-500 dark:bg-dark-700/60 dark:text-dark-400">
                          <th class="px-2 py-1.5 text-left font-medium">{{ t('modelPricing.columns.group') }}</th>
                          <th class="px-2 py-1.5 text-left font-medium">{{ t('modelPricing.columns.channel') }}</th>
                          <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.rate') }}</th>
                          <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.input') }}</th>
                          <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.output') }}</th>
                          <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.cacheRead') }}</th>
                          <th class="px-2 py-1.5 text-right font-medium">{{ t('modelPricing.columns.cacheWrite') }}</th>
                          <th v-if="SHOW_OFFICIAL_DISCOUNT" class="px-2 py-1.5 text-right font-medium">
                            {{ t('modelPricing.expand.vsOfficial') }}
                            <span class="font-normal text-gray-400 dark:text-dark-500">
                              ({{ t('modelPricing.columns.input') }})
                            </span>
                          </th>
                          <th></th>
                        </tr>
                      </thead>
                      <tbody class="bg-white dark:bg-dark-900">
                        <tr
                          v-for="offer in sortedOffers(row)"
                          :key="offer.key"
                          class="border-t border-gray-100 dark:border-dark-800"
                          :class="SHOW_OFFICIAL_DISCOUNT && offer.above && 'bg-amber-50/40 dark:bg-amber-950/10'"
                        >
                          <td class="px-2 py-1.5 font-medium text-gray-800 dark:text-dark-100">
                            <span class="inline-flex items-center gap-1">
                              <span
                                v-if="SHOW_OFFICIAL_DISCOUNT && offer.above"
                                class="h-3 w-[3px] rounded-sm bg-amber-400"
                                :aria-label="t('modelPricing.discount.aboveLabel')"
                              ></span>
                              {{ offer.groupName }}
                            </span>
                          </td>
                          <td class="truncate px-2 py-1.5 text-gray-500 dark:text-dark-400">
                            <span class="mr-1 text-[9px] uppercase text-gray-400 dark:text-dark-500">
                              {{ platformLabel(offer.platform) }}
                            </span>
                            {{ offer.channel }}
                          </td>
                          <td class="px-2 py-1.5 text-right text-gray-500 dark:text-dark-400">
                            ×{{ formatRate(offer.rate) }}
                          </td>
                          <td class="px-2 py-1.5 text-right tabular-nums text-gray-800 dark:text-dark-100">
                            {{ money(offer.input) }}
                          </td>
                          <td class="px-2 py-1.5 text-right tabular-nums text-gray-800 dark:text-dark-100">
                            {{ money(offer.output) }}
                          </td>
                          <td class="px-2 py-1.5 text-right tabular-nums text-gray-800 dark:text-dark-100">
                            {{ money(offer.cacheRead) }}
                          </td>
                          <td class="px-2 py-1.5 text-right tabular-nums text-gray-800 dark:text-dark-100">
                            {{ money(offer.cacheWrite) }}
                          </td>
                          <td v-if="SHOW_OFFICIAL_DISCOUNT" class="px-2 py-1.5">
                            <div v-if="offer.ratio != null" class="flex flex-col items-end gap-1">
                              <div class="flex items-center justify-end gap-1.5">
                                <span
                                  v-if="offer.ratio > 1"
                                  class="whitespace-nowrap text-[10px] font-medium text-amber-600 dark:text-amber-400"
                                >
                                  {{ t('modelPricing.discount.above', { percent: Math.round((offer.ratio - 1) * 100) }) }}
                                </span>
                                <span v-else class="whitespace-nowrap text-[10px] text-gray-400 dark:text-dark-500">
                                  {{ t('modelPricing.discount.save', { percent: Math.round((1 - offer.ratio) * 100) }) }}
                                </span>
                                <span
                                  class="inline-flex items-center whitespace-nowrap rounded border px-1.5 py-0.5 text-[11px] font-semibold"
                                  :class="discountClass(offer.ratioMax ?? offer.ratio)"
                                >
                                  {{ discountText(offer.ratio, offer.ratioMax) }}
                                </span>
                              </div>
                              <div class="h-[5px] w-[120px] overflow-hidden rounded-full bg-gray-200/80 dark:bg-dark-700">
                                <div
                                  class="h-full rounded-full"
                                  :class="barClass(offer.ratioMax ?? offer.ratio)"
                                  :style="{ width: barWidth(offer.ratioMax ?? offer.ratio) }"
                                ></div>
                              </div>
                            </div>
                            <span
                              v-else
                              class="flex items-center justify-end gap-1 text-[11px] text-gray-400 dark:text-dark-500"
                            >
                              <span class="rounded border border-gray-200 px-1.5 py-0.5 dark:border-dark-700">
                                {{ t('modelPricing.discount.notListed') }}
                              </span>
                            </span>
                          </td>
                          <td></td>
                        </tr>
                      </tbody>
                    </table>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>

        <div
          v-if="hiddenByDiscount > 0"
          class="flex items-center gap-1.5 border-t border-gray-100 bg-gray-50/40 px-4 py-2 text-xs text-gray-400 dark:border-dark-700 dark:bg-dark-800/40 dark:text-dark-500"
        >
          <Icon name="infoCircle" size="xs" class="h-3.5 w-3.5 flex-shrink-0" />
          {{ t('modelPricing.hiddenUndiscounted', { count: hiddenByDiscount }) }}
        </div>
      </div>

      <p class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-400 dark:text-dark-500">
        <span>{{ t('modelPricing.notes.paid') }}</span>
        <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
        <span>{{ t('modelPricing.notes.official') }}</span>
        <template v-if="SHOW_OFFICIAL_DISCOUNT">
          <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
          <span>{{ t('modelPricing.notes.discount') }}</span>
        </template>
        <span class="h-3 w-px bg-gray-200 dark:bg-dark-700"></span>
        <span>{{ t('modelPricing.notes.peak') }}</span>
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
  type UserOfficialPricing,
  type UserSupportedModel,
} from '@/api/channels'
import userGroupsAPI from '@/api/groups'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { platformBadgeClass } from '@/utils/platformColors'
import { SHOW_OFFICIAL_DISCOUNT } from '@/constants/pricingUi'
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

/** 表格按 1M token 展示单价；按次/按图片按单次展示（接口给的是 USD / token）。 */
const PER_MILLION = 1_000_000

/**
 * 单条报价 = 一个「分组 × 渠道」。同一模型在不同分组倍率不同、在不同渠道基础价
 * 也可能不同，所以价格按报价逐条保存；折叠时只展示其中最便宜的一条。
 */
interface Offer {
  key: string
  channel: string
  platform: string
  groupId: number
  groupName: string
  rate: number
  billingMode: BillingMode
  /** 实付单价（已乘分组倍率），USD / token 或 USD / 次。 */
  input: number | null
  output: number | null
  cacheRead: number | null
  cacheWrite: number | null
  perRequest: number | null
  /** 分时定价的最高档倍率；无分时或不分档为 1。 */
  maxMult: number
  /** 相对官方的输入价折扣（当前价 ÷ 官方）；官方未收录时为 null。 */
  ratio: number | null
  /** 标准时段折扣（当前价 × 分时倍率 ÷ 官方），无分时时等于 ratio。 */
  ratioMax: number | null
  /** 高于官方价（按标准时段价判断），展开表里标黄。 */
  above: boolean
}

interface ModelRow {
  name: string
  platform: string
  platforms: string[]
  billingMode: BillingMode
  tiered: boolean
  official: UserOfficialPricing | null
  offers: Offer[]
  best: Offer | null
  bestRatio: number | null
  bestRatioMax: number | null
  bestMaxMult: number
  bestOutputRatio: number | null
  groupNames: string[]
}

const channels = ref<UserAvailableChannel[]>([])
const userGroupRates = ref<Record<number, number>>({})
const loading = ref(false)
const searchQuery = ref('')
const selectedGroup = ref<string | number>('all')
const selectedPlatform = ref<string>('all')
/** 折扣体系关闭时默认按实付输入价升序，避免出现「看不出依据」的排序。 */
const sortKey = ref<string>(SHOW_OFFICIAL_DISCOUNT ? 'discount' : 'input')
const onlyDiscounted = ref(false)
const expanded = ref<Set<string>>(new Set())

/** 表格总列数：模型 + 官方 3 + 本站 3 + 折扣 + 分组 + 展开。 */
const tableColspan = SHOW_OFFICIAL_DISCOUNT ? 10 : 9

const sortOptions = computed<SelectOption[]>(() => {
  const options: SelectOption[] = []
  if (SHOW_OFFICIAL_DISCOUNT) {
    options.push({ value: 'discount', label: t('modelPricing.sort.discountAsc') })
  }
  options.push(
    { value: 'input', label: t('modelPricing.sort.inputAsc') },
    { value: 'output', label: t('modelPricing.sort.outputAsc') },
    { value: 'name', label: t('modelPricing.sort.name') },
  )
  return options
})

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

interface ResolvedPricing {
  mode: BillingMode
  input: number | null
  output: number | null
  cacheRead: number | null
  cacheWrite: number | null
  perRequest: number | null
  tiered: boolean
}

/** 渠道定价 → 基础单价（未乘分组倍率）。阶梯价只取首档展示。 */
function resolvePricing(model: UserSupportedModel): ResolvedPricing {
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

function scaled(v: number | null, rate: number): number | null {
  return v == null ? null : v * rate
}

function ratioOf(paid: number | null, official: number | null | undefined): number | null {
  if (paid == null || official == null || official <= 0) return null
  return paid / official
}

/**
 * 官方价对比区间：min = 当前价 ÷ 官价，max = 当前价 × 分时最高倍率 ÷ 官价。
 * 分时模型（如 deepseek 谷/峰）两个数才有区分，其余模型两者相同。
 */
function offerRatios(
  offer: Offer,
  official: UserOfficialPricing | null,
): { min: number | null; max: number | null } {
  if (!official) return { min: null, max: null }
  if (offer.perRequest != null) {
    const ratio = ratioOf(offer.perRequest, official.image_output_price ?? null)
    return { min: ratio, max: ratio }
  }
  const min = ratioOf(offer.input, official.input_price)
  return { min, max: min == null ? null : min * offer.maxMult }
}

function outputRatio(offer: Offer | null, official: UserOfficialPricing | null): number | null {
  if (!offer || !official || offer.perRequest != null) return null
  return ratioOf(offer.output, official.output_price)
}

/** 报价的主价：按次计费用单次价，token 计费用输入价。 */
function primaryPrice(o: Offer | null): number | null {
  if (!o) return null
  return o.perRequest != null ? o.perRequest : o.input
}

function comparePrice(a: Offer, b: Offer): number {
  const av = primaryPrice(a)
  const bv = primaryPrice(b)
  if (av == null && bv == null) return a.groupName.localeCompare(b.groupName)
  if (av == null) return 1
  if (bv == null) return -1
  if (av === bv) return a.groupName.localeCompare(b.groupName)
  return av - bv
}

/**
 * 把「渠道 × 平台 × 模型 × 分组」摊平成报价，再按模型名合并。
 *
 * 合并键只用模型名：同一模型挂在多个渠道/平台时不应重复出行（原页面按
 * 「渠道 × 平台 × 模型」平铺，导致同一个模型出现多行）。跨平台报价在展开表里
 * 各自带平台标记。
 */
const rows = computed<ModelRow[]>(() => {
  const byName = new Map<string, ModelRow>()

  for (const ch of channels.value) {
    for (const section of ch.platforms) {
      for (const model of section.supported_models) {
        // "*" 是渠道的兜底通配条目，不是真实模型，也没有官方价可比。
        if (!model.name || model.name === '*') continue
        const base = resolvePricing(model)
        const priced =
          base.input != null ||
          base.output != null ||
          base.cacheRead != null ||
          base.cacheWrite != null ||
          base.perRequest != null
        // 完全没配价格的模型不出行，避免整行只剩「-」。
        if (!priced) continue

        const platform = section.platform || model.platform
        let row = byName.get(model.name)
        if (!row) {
          row = {
            name: model.name,
            platform,
            platforms: [],
            billingMode: base.mode,
            tiered: base.tiered,
            official: model.official_pricing ?? null,
            offers: [],
            best: null,
            bestRatio: null,
            bestRatioMax: null,
            bestMaxMult: 1,
            bestOutputRatio: null,
            groupNames: [],
          }
          byName.set(model.name, row)
        }
        // 同一模型只要有一个渠道带官价就够用；没带的渠道不覆盖已有值。
        if (!row.official && model.official_pricing) row.official = model.official_pricing
        row.tiered = row.tiered || base.tiered
        if (!row.platforms.includes(platform)) row.platforms.push(platform)

        for (const g of section.groups) {
          const rate = userGroupRates.value[g.id] ?? g.rate_multiplier
          const offer: Offer = {
            key: `${ch.name}|${platform}|${model.name}|${g.id}`,
            channel: ch.name,
            platform,
            groupId: g.id,
            groupName: g.name,
            rate,
            billingMode: base.mode,
            input: scaled(base.input, rate),
            output: scaled(base.output, rate),
            cacheRead: scaled(base.cacheRead, rate),
            cacheWrite: scaled(base.cacheWrite, rate),
            perRequest: scaled(base.perRequest, rate),
            maxMult: model.pricing?.max_time_multiplier ?? 1,
            ratio: null,
            ratioMax: null,
            above: false,
          }
          const ratios = offerRatios(offer, row.official)
          offer.ratio = ratios.min
          offer.ratioMax = ratios.max
          offer.above = ratios.max != null && ratios.max > 1
          row.offers.push(offer)
          if (!row.groupNames.includes(g.name)) row.groupNames.push(g.name)
        }
      }
    }
  }

  const out = [...byName.values()].filter((r) => r.offers.length > 0)
  for (const r of out) {
    r.offers.sort(comparePrice)
    r.best = r.offers[0] ?? null
    r.bestRatio = r.best ? r.best.ratio : null
    r.bestRatioMax = r.best ? r.best.ratioMax : null
    r.bestMaxMult = r.best ? r.best.maxMult : 1
    r.bestOutputRatio = outputRatio(r.best, r.official)
    if (r.best) r.platform = r.best.platform
  }
  return out
})

/** 平台切换按钮（不是 Select，值是纯 string，便于和 selectedPlatform 比较）。 */
interface PlatformOption {
  value: string
  label: string
}

const platformOptions = computed<PlatformOption[]>(() => {
  const seen = new Set<string>()
  for (const r of rows.value) for (const p of r.platforms) seen.add(p)
  const options: PlatformOption[] = [{ value: 'all', label: t('modelPricing.allPlatforms') }]
  for (const p of [...seen].sort()) options.push({ value: p, label: platformLabel(p) })
  return options
})

function matchesFilters(row: ModelRow, query: string): boolean {
  if (selectedGroup.value !== 'all' && !row.offers.some((o) => o.groupId === selectedGroup.value)) {
    return false
  }
  if (selectedPlatform.value !== 'all' && !row.platforms.includes(selectedPlatform.value)) {
    return false
  }
  if (!query) return true
  return (
    row.name.toLowerCase().includes(query) ||
    row.offers.some(
      (o) => o.channel.toLowerCase().includes(query) || o.groupName.toLowerCase().includes(query),
    )
  )
}

/** 分组筛选时只保留该分组的报价，折扣随之重算。 */
function visibleOffers(row: ModelRow): Offer[] {
  if (selectedGroup.value === 'all') return row.offers
  return row.offers.filter((o) => o.groupId === selectedGroup.value)
}

/** 命中搜索/平台/分组筛选的模型，报价按所选分组收敛、折扣重算。 */
const scopedRows = computed<ModelRow[]>(() => {
  const query = searchQuery.value.trim().toLowerCase()
  const out: ModelRow[] = []
  for (const row of rows.value) {
    if (!matchesFilters(row, query)) continue
    const offers = visibleOffers(row)
    if (!offers.length) continue
    const best = offers[0]
    out.push({
      ...row,
      offers,
      best,
      platform: best.platform,
      bestRatio: best.ratio,
      bestRatioMax: best.ratioMax,
      bestMaxMult: best.maxMult,
      bestOutputRatio: outputRatio(best, row.official),
      groupNames: [...new Set(offers.map((o) => o.groupName))],
    })
  }
  return out
})

function hasDiscount(row: ModelRow): boolean {
  return row.bestRatio != null && row.bestRatio < 0.999
}

const filteredRows = computed<ModelRow[]>(() => {
  const out = onlyDiscounted.value ? scopedRows.value.filter(hasDiscount) : scopedRows.value

  const sorted = [...out]
  switch (sortKey.value) {
    case 'name':
      sorted.sort((a, b) => a.name.localeCompare(b.name))
      break
    case 'input':
      sorted.sort(
        (a, b) =>
          compareNullable(primaryPrice(a.best), primaryPrice(b.best)) || a.name.localeCompare(b.name),
      )
      break
    case 'output':
      sorted.sort(
        (a, b) =>
          compareNullable(a.best?.output ?? null, b.best?.output ?? null) || a.name.localeCompare(b.name),
      )
      break
    default:
      // 折扣升序：最便宜的排前面；没有官价（算不出折扣）的排最后。
      sorted.sort((a, b) => compareNullable(a.bestRatio, b.bestRatio) || a.name.localeCompare(b.name))
  }
  return sorted
})

function compareNullable(a: number | null, b: number | null): number {
  if (a == null && b == null) return 0
  if (a == null) return 1
  if (b == null) return -1
  return a - b
}

const totalOffers = computed(() => rows.value.reduce((n, r) => n + r.offers.length, 0))

/** 平均折扣只统计有官价的模型，避免被「未收录」拉偏。 */
const avgDiscountText = computed(() => {
  const ratio = filteredRows.value.map((r) => r.bestRatio).filter((r): r is number => r != null)
  if (!ratio.length) return t('modelPricing.discount.notListed')
  const avg = ratio.reduce((s, r) => s + r, 0) / ratio.length
  return t('modelPricing.discount.value', { value: (avg * 10).toFixed(1) })
})

const aboveOfficialCount = computed(
  () => rows.value.reduce((n, r) => n + r.offers.filter((o) => o.above).length, 0),
)

/** 当前筛选下被「仅看有折扣」隐藏的模型数。 */
const hiddenByDiscount = computed(() => {
  if (!onlyDiscounted.value) return 0
  return scopedRows.value.length - filteredRows.value.length
})

function isOpen(name: string): boolean {
  return expanded.value.has(name)
}

function toggleRow(name: string): void {
  const next = new Set(expanded.value)
  if (next.has(name)) next.delete(name)
  else next.add(name)
  expanded.value = next
}

function sortedOffers(row: ModelRow): Offer[] {
  return [...row.offers].sort(comparePrice)
}

/* ---------------- 展示格式化 ---------------- */

/**
 * 表格内不带 "$"（表头已标注 USD），统一保留两位小数：0.117994 / 1.474926 这类长尾
 * 对读价没有意义，反而让整列难以扫读。
 *
 * 例外：换算后不足 1 分（典型是缓存读价，如 $0.0015 / 1M）会保留 4 位小数并去掉末尾
 * 的 0，否则会被四舍五入成 0.00，等于把价格抹掉。
 */
function money(value: number | null, scale = PER_MILLION): string {
  if (value == null) return '-'
  const scaled = value * scale
  if (scaled === 0) return '0.00'
  if (Math.abs(scaled) < 0.01) {
    return scaled.toFixed(4).replace(/0+$/, '').replace(/\.$/, '')
  }
  return scaled.toFixed(2)
}

function officialInput(row: ModelRow): string {
  if (!row.official) return '-'
  if (row.billingMode === BILLING_MODE_IMAGE || row.billingMode === BILLING_MODE_VIDEO) {
    return row.official.image_output_price != null
      ? `${money(row.official.image_output_price, 1)}/${t('modelPricing.units.perImage')}`
      : '-'
  }
  return money(row.official.input_price)
}

function officialOutput(row: ModelRow): string {
  if (!row.official || row.billingMode !== BILLING_MODE_TOKEN) return '-'
  return money(row.official.output_price)
}

function officialCache(row: ModelRow): { label: string; value: string }[] {
  if (!row.official || row.billingMode !== BILLING_MODE_TOKEN) return []
  return [
    { label: t('modelPricing.cacheLabels.read'), value: money(row.official.cache_read_price) },
    { label: t('modelPricing.cacheLabels.write'), value: money(row.official.cache_write_price) },
  ]
}

function paidInput(row: ModelRow): string {
  const offer = row.best
  if (!offer) return '-'
  if (offer.perRequest != null) {
    const unit =
      offer.billingMode === BILLING_MODE_IMAGE
        ? t('modelPricing.units.perImage')
        : t('modelPricing.units.perRequestShort')
    return `${money(offer.perRequest, 1)}/${unit}`
  }
  return money(offer.input)
}

function paidOutput(row: ModelRow): string {
  if (!row.best || row.best.perRequest != null) return '-'
  return money(row.best.output)
}

function paidCache(row: ModelRow): { label: string; value: string; above: boolean }[] {
  const offer = row.best
  if (!offer || offer.perRequest != null) return []
  const official = row.official
  const above = (paid: number | null, off: number | null | undefined) =>
    paid != null && off != null && off > 0 && paid > off * 1.001
  return [
    {
      label: t('modelPricing.cacheLabels.read'),
      value: money(offer.cacheRead),
      above: above(offer.cacheRead, official?.cache_read_price),
    },
    {
      label: t('modelPricing.cacheLabels.write'),
      value: money(offer.cacheWrite),
      above: above(offer.cacheWrite, official?.cache_write_price),
    },
  ]
}

function officialSummary(row: ModelRow): string {
  if (!row.official) return t('modelPricing.discount.notListedHint')
  if (row.billingMode === BILLING_MODE_IMAGE || row.billingMode === BILLING_MODE_VIDEO) {
    return row.official.image_output_price != null
      ? `${money(row.official.image_output_price, 1)} ${t('modelPricing.units.perImage')}`
      : t('modelPricing.discount.notListedHint')
  }
  const parts = [
    `${t('modelPricing.columns.input')} ${money(row.official.input_price)}`,
    `${t('modelPricing.columns.output')} ${money(row.official.output_price)}`,
    `${t('modelPricing.columns.cacheRead')} ${money(row.official.cache_read_price)}`,
    `${t('modelPricing.columns.cacheWrite')} ${money(row.official.cache_write_price)}`,
  ]
  return `${parts.join(' · ')} ${t('modelPricing.blocks.unit')}`
}

function discountText(min: number | null, max: number | null = null): string {
  if (min == null) return t('modelPricing.discount.notListed')
  // 分时模型给出区间（谷时 / 峰时），非分时两者相同只显示一个数。
  if (max != null && max - min > 0.005) {
    return t('modelPricing.discount.range', { min: (min * 10).toFixed(1), max: (max * 10).toFixed(1) })
  }
  return t('modelPricing.discount.value', { value: (min * 10).toFixed(1) })
}

function discountClass(ratio: number | null): string {
  if (ratio == null) {
    return 'border-gray-200 bg-gray-50 text-gray-400 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-500'
  }
  if (ratio <= 0.55) {
    return 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/60 dark:bg-emerald-950/40 dark:text-emerald-300'
  }
  if (ratio <= 0.75) {
    return 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-900/60 dark:bg-blue-950/40 dark:text-blue-300'
  }
  if (ratio <= 1) {
    return 'border-gray-200 bg-gray-50 text-gray-600 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300'
  }
  return 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-900/60 dark:bg-amber-950/40 dark:text-amber-300'
}

function barClass(ratio: number | null): string {
  if (ratio == null) return 'bg-gray-300 dark:bg-dark-600'
  if (ratio <= 0.55) return 'bg-emerald-500'
  if (ratio <= 0.75) return 'bg-blue-500'
  if (ratio <= 1) return 'bg-gray-400'
  return 'bg-amber-500'
}

function barWidth(ratio: number | null): string {
  const pct = ratio == null ? 0 : Math.min(100, Math.max(2, ratio * 100))
  return `${pct}%`
}

/**
 * 折扣副标注：输入/输出折扣不一致时写明输出折扣，避免只看到一个数字被误导；
 * 一致时显示省了百分之多少。
 */
function subDiscount(row: ModelRow): { text: string; above: boolean } | null {
  if (!row.best || row.bestRatio == null) return null
  // 分时模型：折扣是区间，标明倍率比"省 X%"更有用。
  if (row.bestRatioMax != null && row.bestRatioMax - row.bestRatio > 0.005) {
    return {
      text: t('modelPricing.discount.timePricing', { multiplier: formatRate(row.bestMaxMult) }),
      above: false,
    }
  }
  // 高于官价是最需要被看到的信息，优先于"输出折扣不一致"。
  if (row.bestRatio > 1) {
    return {
      text: t('modelPricing.discount.above', { percent: Math.round((row.bestRatio - 1) * 100) }),
      above: true,
    }
  }
  const out = row.bestOutputRatio
  if (out != null && Math.abs(out - row.bestRatio) > 0.03) {
    return {
      text: t('modelPricing.discount.outputDiscount', { value: (out * 10).toFixed(1) }),
      above: out > 1,
    }
  }
  return { text: t('modelPricing.discount.save', { percent: Math.round((1 - row.bestRatio) * 100) }), above: false }
}

/** 平台展示名：接口给的是小写标识，界面上按厂商习惯写法展示。 */
const PLATFORM_LABELS: Record<string, string> = {
  openai: 'OpenAI',
  anthropic: 'Anthropic',
  gemini: 'Gemini',
  vertex: 'Vertex',
  bedrock: 'Bedrock',
  azure: 'Azure',
  xai: 'xAI',
  deepseek: 'DeepSeek',
  qwen: 'Qwen',
  glm: 'GLM',
  kimi: 'Kimi',
  minimax: 'MiniMax',
  bytedance: 'ByteDance',
  composite: 'Composite',
}

function platformLabel(platform: string): string {
  if (!platform) return ''
  return PLATFORM_LABELS[platform.toLowerCase()] ?? platform
}

function formatRate(rate: number): string {
  if (!Number.isFinite(rate)) return '1.0'
  return Number.isInteger(rate) ? rate.toFixed(1) : String(Number(rate.toFixed(4)))
}

async function load(): Promise<void> {
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
