<script setup lang="ts">
import { computed } from 'vue'
import type { ConsoleProfile, ConsoleStats } from '../api'
import { formatCost, formatMoney, formatNumber } from './format'
import { useCampaignMessages } from '../messages'

const props = defineProps<{
  stats?: ConsoleStats | null
  profile?: ConsoleProfile | null
  loading?: boolean
  error?: string
}>()

const t = useCampaignMessages()

const cards = computed(() => [
  {
    key: 'balance',
    label: t('console.balance'),
    value: formatMoney(props.profile?.balance),
    hint: props.profile?.email || '',
  },
  {
    key: 'todayRequests',
    label: t('console.todayRequests'),
    value: formatNumber(props.stats?.today_requests),
  },
  {
    key: 'todayCost',
    label: t('console.todayCost'),
    value: formatCost(props.stats?.today_actual_cost),
  },
  {
    key: 'todayTokens',
    label: t('console.todayTokens'),
    value: formatNumber(props.stats?.today_tokens),
  },
  {
    key: 'totalRequests',
    label: t('console.totalRequests'),
    value: formatNumber(props.stats?.total_requests),
  },
  {
    key: 'activeKeys',
    label: t('console.activeKeys'),
    value: formatNumber(props.stats?.active_api_keys),
  },
])
</script>

<template>
  <div>
    <p v-if="error" class="c-error mb-4">{{ error }}</p>

    <div class="c-stat-grid">
      <div v-for="card in cards" :key="card.key" class="c-card c-stat">
        <span class="c-stat-label">{{ card.label }}</span>
        <span v-if="loading" class="c-skeleton c-stat-skeleton"></span>
        <span v-else class="c-stat-value">{{ card.value }}</span>
        <span v-if="card.hint" class="c-stat-hint">{{ card.hint }}</span>
      </div>
    </div>
  </div>
</template>
