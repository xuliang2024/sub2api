<template>
  <AppLayout>
    <div class="mx-auto flex w-full max-w-5xl flex-col gap-6">
      <div class="card">
        <div class="flex flex-wrap items-end gap-4 px-6 py-4">
          <div>
            <label class="input-label">{{ t('dashboard.timeRange') }}</label>
            <DateRangePicker
              v-model:start-date="startDate"
              v-model:end-date="endDate"
              @change="onDateRangeChange"
            />
          </div>
          <button class="btn btn-secondary ml-auto" :disabled="loading" @click="loadRanking">
            {{ t('common.refresh') }}
          </button>
        </div>
      </div>

      <UserDashboardSpendingRanking
        :items="rankingItems"
        :total-actual-cost="rankingTotalActualCost"
        :loading="loading"
        :range-label="rangeLabel"
      />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import UserDashboardSpendingRanking from '@/components/user/dashboard/UserDashboardSpendingRanking.vue'
import { usageAPI, type UserSpendingRankingItem } from '@/api/usage'

const { t } = useI18n()

const formatLocalDate = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const endDate = ref(formatLocalDate(new Date()))
const startDate = ref(formatLocalDate(new Date(Date.now() - 6 * 86400000)))
const loading = ref(false)
const rankingItems = ref<UserSpendingRankingItem[]>([])
const rankingTotalActualCost = ref(0)

const rangeLabel = computed(() => `${startDate.value} - ${endDate.value}`)

const loadRanking = async () => {
  loading.value = true
  try {
    const res = await usageAPI.getDashboardUsersRanking({
      start_date: startDate.value,
      end_date: endDate.value,
      limit: 10,
    })
    rankingItems.value = res.ranking || []
    rankingTotalActualCost.value = res.total_actual_cost || 0
  } catch (error) {
    console.error('Failed to load spending ranking:', error)
    rankingItems.value = []
    rankingTotalActualCost.value = 0
  } finally {
    loading.value = false
  }
}

const onDateRangeChange = (range: { startDate: string; endDate: string }) => {
  startDate.value = range.startDate
  endDate.value = range.endDate
  loadRanking()
}

onMounted(() => {
  loadRanking()
})
</script>
