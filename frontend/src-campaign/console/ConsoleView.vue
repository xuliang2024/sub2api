<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  describeError,
  fetchConsoleProfile,
  fetchConsoleStats,
  type ConsoleProfile,
  type ConsoleStats,
} from '../api'
import StatsPanel from './StatsPanel.vue'
import KeysPanel from './KeysPanel.vue'
import LogsPanel from './LogsPanel.vue'
import { useCampaignMessages } from '../messages'

const t = useCampaignMessages()

type ConsoleTab = 'overview' | 'keys' | 'logs'

const active = ref<ConsoleTab>('overview')
const stats = ref<ConsoleStats | null>(null)
const profile = ref<ConsoleProfile | null>(null)
const loading = ref(true)
const error = ref('')

const tabs = computed(() => [
  { key: 'overview' as const, label: t('console.overview') },
  { key: 'keys' as const, label: t('console.keys') },
  { key: 'logs' as const, label: t('console.logs') },
])

/**
 * 概览数据一次拉两份：余额在 profile 里，用量在 dashboard/stats。
 * 任一失败都不应让整个控制台空白，所以分开 catch。
 */
async function loadSummary(): Promise<void> {
  loading.value = true
  error.value = ''
  const [statsResult, profileResult] = await Promise.allSettled([
    fetchConsoleStats(),
    fetchConsoleProfile(),
  ])
  if (statsResult.status === 'fulfilled') {
    stats.value = statsResult.value
  }
  if (profileResult.status === 'fulfilled') {
    profile.value = profileResult.value
  }
  if (statsResult.status === 'rejected' && profileResult.status === 'rejected') {
    error.value = describeError(statsResult.reason, t('console.loadFailed'))
  }
  loading.value = false
}

onMounted(loadSummary)
</script>

<template>
  <section class="c-section">
    <div class="campaign-container">
      <div class="c-tabs" role="tablist">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="c-tab"
          type="button"
          role="tab"
          :aria-selected="active === tab.key"
          @click="active = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <div class="mt-6">
        <StatsPanel
          v-if="active === 'overview'"
          :stats="stats"
          :profile="profile"
          :loading="loading"
          :error="error"
        />
        <KeysPanel v-else-if="active === 'keys'" @changed="loadSummary" />
        <LogsPanel v-else />
      </div>
    </div>
  </section>
</template>
