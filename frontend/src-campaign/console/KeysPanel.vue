<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  createConsoleKey,
  deleteConsoleKey,
  describeError,
  fetchConsoleKeys,
  updateConsoleKeyStatus,
  type ConsoleKey,
} from '../api'
import { copyText } from '../format'
import { formatDateTime, maskKey } from './format'
import { useCampaignMessages } from '../messages'

const emit = defineEmits<{
  (event: 'changed'): void
}>()

const t = useCampaignMessages()

const items = ref<ConsoleKey[]>([])
const loading = ref(true)
const error = ref('')
const busyId = ref(0)
const creating = ref(false)
const name = ref('')
const copiedId = ref(0)
const page = ref(1)
const pages = ref(1)

async function load(target = page.value): Promise<void> {
  loading.value = true
  error.value = ''
  try {
    const data = await fetchConsoleKeys(target, 20)
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

async function handleCreate(): Promise<void> {
  if (creating.value) {
    return
  }
  creating.value = true
  error.value = ''
  try {
    await createConsoleKey(name.value.trim() || 'console')
    name.value = ''
    await load(1)
    emit('changed')
  } catch (err) {
    error.value = describeError(err, t('console.loadFailed'))
  } finally {
    creating.value = false
  }
}

async function handleCopy(item: ConsoleKey): Promise<void> {
  const ok = await copyText(item.key)
  if (ok) {
    copiedId.value = item.id
    setTimeout(() => {
      if (copiedId.value === item.id) {
        copiedId.value = 0
      }
    }, 1800)
  }
}

async function toggleStatus(item: ConsoleKey): Promise<void> {
  if (busyId.value) {
    return
  }
  busyId.value = item.id
  try {
    await updateConsoleKeyStatus(item.id, item.status === 'active' ? 'inactive' : 'active')
    await load()
    emit('changed')
  } catch (err) {
    error.value = describeError(err, t('console.loadFailed'))
  } finally {
    busyId.value = 0
  }
}

async function handleDelete(item: ConsoleKey): Promise<void> {
  if (busyId.value || !window.confirm(t('console.confirmRemove'))) {
    return
  }
  busyId.value = item.id
  try {
    await deleteConsoleKey(item.id)
    await load()
    emit('changed')
  } catch (err) {
    error.value = describeError(err, t('console.loadFailed'))
  } finally {
    busyId.value = 0
  }
}

function statusLabel(status: ConsoleKey['status']): string {
  if (status === 'active') {
    return t('console.statusActive')
  }
  if (status === 'inactive') {
    return t('console.statusInactive')
  }
  return t('console.statusOther')
}
</script>

<template>
  <div>
    <div class="c-console-toolbar">
      <input
        v-model="name"
        class="c-input c-console-name"
        type="text"
        :placeholder="t('console.keyName')"
        @keyup.enter="handleCreate"
      />
      <button class="c-btn c-btn-primary c-btn-sm" type="button" :disabled="creating" @click="handleCreate">
        {{ creating ? t('console.creating') : t('console.newKey') }}
      </button>
    </div>

    <p v-if="error" class="c-error mb-3">{{ error }}</p>

    <div v-if="loading" class="space-y-3">
      <div class="c-skeleton h-10 w-full"></div>
      <div class="c-skeleton h-10 w-full"></div>
    </div>

    <p v-else-if="!items.length" class="c-note">{{ t('console.emptyKeys') }}</p>

    <div v-else class="c-table-wrap">
      <table class="c-table c-table--plain">
        <thead>
          <tr>
            <th scope="col">{{ t('console.keyName') }}</th>
            <th scope="col">{{ t('console.key') }}</th>
            <th scope="col">{{ t('console.createdAt') }}</th>
            <th scope="col">{{ t('console.lastUsed') }}</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id">
            <td>
              <div class="c-table-name">{{ item.name }}</div>
              <span class="c-key-status" :class="`is-${item.status}`">{{ statusLabel(item.status) }}</span>
            </td>
            <td>
              <code class="c-key-code">{{ maskKey(item.key) }}</code>
            </td>
            <td class="c-table-muted">{{ formatDateTime(item.created_at) }}</td>
            <td class="c-table-muted">
              {{ item.last_used_at ? formatDateTime(item.last_used_at) : t('console.never') }}
            </td>
            <td class="c-key-actions">
              <button class="c-btn c-btn-ghost c-btn-sm" type="button" @click="handleCopy(item)">
                {{ copiedId === item.id ? t('console.copied') : t('console.copy') }}
              </button>
              <button
                class="c-btn c-btn-ghost c-btn-sm"
                type="button"
                :disabled="busyId === item.id"
                @click="toggleStatus(item)"
              >
                {{ item.status === 'active' ? t('console.disable') : t('console.enable') }}
              </button>
              <button
                class="c-btn c-btn-ghost c-btn-sm c-danger"
                type="button"
                :disabled="busyId === item.id"
                @click="handleDelete(item)"
              >
                {{ t('console.remove') }}
              </button>
            </td>
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
