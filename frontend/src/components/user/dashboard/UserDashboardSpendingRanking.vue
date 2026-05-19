<template>
  <div class="card">
    <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <div>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.spendingRanking') }}</h2>
        <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('dashboard.top10Users') }}</p>
      </div>
      <span class="badge badge-gray">{{ rangeLabel || t('dashboard.last7Days') }}</span>
    </div>

    <div class="p-6">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
      <div v-else-if="items.length === 0" class="py-8">
        <EmptyState :title="t('dashboard.noRankingData')" :description="t('dashboard.rankingDataHint')" />
      </div>
      <div v-else class="space-y-3">
        <div
          v-for="item in items"
          :key="item.rank"
          class="flex items-center justify-between rounded-xl p-3 transition-colors"
          :class="item.is_current_user
            ? 'bg-primary-50 ring-1 ring-primary-100 dark:bg-primary-900/20 dark:ring-primary-800/40'
            : 'bg-gray-50 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800'"
        >
          <div class="flex min-w-0 items-center gap-3">
            <div
              class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg text-xs font-bold"
              :class="rankClass(item.rank)"
            >
              {{ item.rank }}
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <p class="truncate text-sm font-medium text-gray-900 dark:text-white" :title="item.display_name">
                  {{ item.display_name }}
                </p>
                <span v-if="item.is_current_user" class="badge badge-primary">{{ t('dashboard.you') }}</span>
              </div>
              <p class="text-xs text-gray-500 dark:text-dark-400">
                {{ formatNumber(item.requests) }} {{ t('dashboard.requests') }} · {{ formatTokens(item.tokens) }} {{ t('dashboard.tokens') }}
              </p>
            </div>
          </div>

          <div class="ml-3 text-right">
            <p class="text-sm font-semibold text-green-600 dark:text-green-400">${{ formatCost(item.actual_cost) }}</p>
            <p class="text-xs text-gray-500 dark:text-dark-400">{{ percentage(item.actual_cost) }}%</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { formatCostFixed as formatCost, formatNumberLocaleString as formatNumber, formatTokensK as formatTokens } from '@/utils/format'
import type { UserSpendingRankingItem } from '@/api/usage'

const props = defineProps<{
  items: UserSpendingRankingItem[]
  totalActualCost: number
  loading: boolean
  rangeLabel?: string
}>()

const { t } = useI18n()

const rankClass = (rank: number) => {
  if (rank === 1) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  if (rank === 2) return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
  if (rank === 3) return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300'
  return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'
}

const percentage = (actualCost: number) => {
  if (!props.totalActualCost) return '0.0'
  return ((actualCost / props.totalActualCost) * 100).toFixed(1)
}
</script>
