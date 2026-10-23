<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { describeError, fetchConsoleLogs, type ConsoleUsageLog } from '../api'
import { formatCost, formatDateTime, formatNumber } from './format'
import { useCampaignMessages } from '../messages'

const t = useCampaignMessages()

const items = ref<ConsoleUsageLog[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const pages = ref(1)

async function load(target = page.value): Promise<void> {
  loading.value = true
  error.value = ''
  try {
    const data = await fetchConsoleLogs(target, 20)
    items.value = data.items ?? []
    page.value = data.page ?? target
    pages.value = data.pages || 1
  } catch (err) {
    error.value = describeError(err, t('console.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => load(1))

function tokens(log: ConsoleUsageLog): string {
  const total = log.total_tokens ?? log.input_tokens + log.output_tokens + log.cache_read_tokens
  return formatNumber(total)
}
</script>

<template>
  <div>
    <p v-if="error" class="c-error mb-3">{{ error }}</p>

    <div v-if="loading" class="space-y-3">
      <div class="c-skeleton h-10 w-full"></div>
      <div class="c-skeleton h-10 w-full"></div>
    </div>

    <p v-else-if="!items.length" class="c-note">{{ t('console.emptyLogs') }}</p>

    <div v-else class="c-table-wrap">
      <table class="c-table c-table--plain">
        <thead>
          <tr>
            <th scope="col">{{ t('console.time') }}</th>
            <th scope="col">{{ t('console.model') }}</th>
            <th scope="col" class="c-num">{{ t('console.tokens') }}</th>
            <th scope="col" class="c-num">{{ t('console.cost') }}</th>
            <th scope="col" class="c-num">{{ t('console.key') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in items" :key="log.id">
            <td class="c-table-muted">{{ formatDateTime(log.created_at) }}</td>
            <td>
              <div class="c-table-name">{{ log.model }}</div>
              <div v-if="log.duration_ms" class="c-table-id">{{ log.duration_ms }} ms</div>
            </td>
            <td class="c-table-price c-num">
              <div class="c-table-ours">{{ tokens(log) }}</div>
              <div class="c-table-id">
                {{ formatNumber(log.input_tokens) }} / {{ formatNumber(log.output_tokens) }}
              </div>
            </td>
            <td class="c-table-price c-num">
              <span class="c-table-ours">{{ formatCost(log.actual_cost) }}</span>
            </td>
            <td class="c-table-muted c-num">#{{ log.api_key_id }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="pages > 1" class="c-pager">
      <button class="c-btn c-btn-ghost c-btn-sm" type="button" :disabled="page <= 1" @click="load(page - 1)">
        {{ t('console.prev') }}
      </button>
      <span class="c-table-muted">{{ t('console.page', { page, pages }) }}</span>
      <button
        class="c-btn c-btn-ghost c-btn-sm"
        type="button"
        :disabled="page >= pages"
        @click="load(page + 1)"
      >
        {{ t('console.next') }}
      </button>
    </div>
  </div>
</template>
