<template>
  <AppLayout>
    <div class="mx-auto flex max-w-7xl flex-col gap-6">
      <section
        v-if="!isDesktop"
        class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-200"
      >
        {{ t('desktopCodex.webOnlyNotice') }}
      </section>

      <section class="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
        <div class="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">
                {{ t('desktopCodex.connection') }}
              </p>
              <h2 class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                {{ t('desktopCodex.title') }}
              </h2>
              <p class="mt-2 max-w-2xl text-sm text-gray-500 dark:text-dark-400">
                {{ t('desktopCodex.description') }}
              </p>
            </div>
            <div class="flex flex-wrap gap-2">
              <button class="btn btn-secondary" :disabled="!isDesktop" type="button" @click="downloadCodex">
                <Icon name="download" size="sm" class="mr-2" />
                {{ t('desktopCodex.downloadCodex') }}
              </button>
              <button class="btn btn-secondary" :disabled="!isDesktop" type="button" @click="openConfigDir">
                <Icon name="document" size="sm" class="mr-2" />
                {{ t('desktopCodex.openConfigDir') }}
              </button>
            </div>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <label class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.apiBaseUrl') }}</span>
              <input
                v-model="apiBaseUrlDraft"
                class="input w-full"
                type="url"
                spellcheck="false"
                :placeholder="fallbackOrigin"
              />
            </label>
            <label class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.configMode') }}</span>
              <Select
                v-model="configMode"
                :options="configModeOptions"
              />
            </label>
          </div>
        </div>

        <div class="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                {{ t('desktopCodex.desktopServer') }}
              </p>
              <h3 class="mt-1 text-base font-semibold text-gray-900 dark:text-white">
                {{ t('desktopCodex.serverUrl') }}
              </h3>
            </div>
            <span
              class="rounded-full px-2.5 py-1 text-xs font-medium"
              :class="isDesktop ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'"
            >
              {{ isDesktop ? t('desktopCodex.desktopReady') : t('desktopCodex.webMode') }}
            </span>
          </div>
          <div class="mt-4 flex flex-col gap-3">
            <input
              v-model="serverUrlDraft"
              class="input w-full"
              type="url"
              spellcheck="false"
              :disabled="!isDesktop"
              placeholder="https://codex.apiz.ai"
            />
            <button class="btn btn-primary w-full" :disabled="!isDesktop || savingServer" type="button" @click="saveServerUrl">
              <Icon name="swap" size="sm" class="mr-2" />
              {{ savingServer ? t('common.saving') : t('desktopCodex.switchServer') }}
            </button>
          </div>
        </div>
      </section>

      <section class="grid gap-4 lg:grid-cols-3">
        <article class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.platform') }}</p>
          <strong class="mt-1 block text-lg text-gray-900 dark:text-white">{{ platformLabel }}</strong>
        </article>
        <article class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <p class="text-xs text-gray-500 dark:text-dark-400">config.toml</p>
          <code class="mt-1 block truncate text-sm text-gray-900 dark:text-white">{{ environment?.configPath || '-' }}</code>
        </article>
        <article class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <p class="text-xs text-gray-500 dark:text-dark-400">auth.json</p>
          <code class="mt-1 block truncate text-sm text-gray-900 dark:text-white">{{ environment?.authPath || '-' }}</code>
        </article>
      </section>

      <section class="grid gap-4 lg:grid-cols-[0.95fr_1.05fr]">
        <div class="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">
                {{ t('desktopCodex.keySelection') }}
              </p>
              <h3 class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('desktopCodex.selectOpenAIGroup') }}
              </h3>
            </div>
            <button class="btn btn-secondary" :disabled="loading" type="button" @click="loadRemoteData">
              <Icon name="refresh" size="sm" :class="loading ? 'mr-2 animate-spin' : 'mr-2'" />
              {{ t('common.refresh') }}
            </button>
          </div>

          <div class="mt-5 space-y-4">
            <label class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.openaiGroup') }}</span>
              <Select
                v-model="selectedGroupId"
                :options="groupOptions"
                :disabled="openaiGroups.length === 0"
                :placeholder="t('desktopCodex.selectGroupPlaceholder')"
              />
            </label>

            <label class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.apiKey') }}</span>
              <Select
                v-model="selectedKeyId"
                :options="keyOptions"
                :disabled="filteredKeys.length === 0"
                :placeholder="t('desktopCodex.selectKeyPlaceholder')"
              />
            </label>

            <button
              class="btn btn-secondary w-full"
              :disabled="!selectedGroupId || creatingKey"
              type="button"
              @click="createCodexKey"
            >
              <Icon name="plus" size="sm" class="mr-2" />
              {{ creatingKey ? t('desktopCodex.creatingKey') : t('desktopCodex.createKey') }}
            </button>

            <p v-if="openaiGroups.length === 0 && !loading" class="rounded-lg bg-amber-50 p-3 text-sm text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
              {{ t('desktopCodex.noOpenaiGroups') }}
            </p>
          </div>
        </div>

        <div class="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">
                {{ t('desktopCodex.writeAndVerify') }}
              </p>
              <h3 class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('desktopCodex.oneClickConfig') }}
              </h3>
            </div>
            <button
              class="btn btn-primary"
              :disabled="!canWrite || writing"
              type="button"
              @click="writeConfig"
            >
              <Icon name="terminal" size="sm" class="mr-2" />
              {{ writing ? t('desktopCodex.writing') : t('desktopCodex.writeConfig') }}
            </button>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.currentConfig') }}</p>
              <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                {{ environment?.hasConfig || environment?.hasAuth ? t('desktopCodex.detected') : t('desktopCodex.notWritten') }}
              </p>
              <p class="mt-1 truncate text-xs text-gray-500 dark:text-dark-400">
                {{ environment?.existingBaseUrl || t('desktopCodex.noBaseUrl') }}
              </p>
            </div>
            <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.selectedKey') }}</p>
              <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ selectedKey?.name || '-' }}</p>
              <p class="mt-1 truncate text-xs text-gray-500 dark:text-dark-400">{{ selectedKey ? maskKey(selectedKey.key) : '-' }}</p>
            </div>
          </div>

          <div v-if="writeResult" class="mt-4 rounded-lg border border-green-200 bg-green-50 p-3 text-sm text-green-800 dark:border-green-800 dark:bg-green-900/20 dark:text-green-200">
            <p>{{ t('desktopCodex.writeSuccess') }}</p>
            <p v-if="writeResult.configBackup" class="mt-1 break-all text-xs">
              config backup: {{ writeResult.configBackup }}
            </p>
            <p v-if="writeResult.authBackup" class="mt-1 break-all text-xs">
              auth backup: {{ writeResult.authBackup }}
            </p>
          </div>

          <div v-if="usageResult" class="mt-4 grid gap-3 md:grid-cols-3">
            <div class="rounded-lg border border-gray-100 p-3 dark:border-dark-700">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.available') }}</p>
              <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ formatMoney(usageResult.available, usageResult.unit) }}</strong>
            </div>
            <div class="rounded-lg border border-gray-100 p-3 dark:border-dark-700">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.status') }}</p>
              <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ usageResult.status || usageResult.mode || '-' }}</strong>
            </div>
            <div class="rounded-lg border border-gray-100 p-3 dark:border-dark-700">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.updatedAt') }}</p>
              <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ formatUpdatedAt(usageResult.updated_at) }}</strong>
            </div>
          </div>

          <p
            v-if="message"
            class="mt-4 rounded-lg p-3 text-sm"
            :class="messageKind === 'error' ? 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300' : 'bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300'"
          >
            {{ message }}
          </p>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { keysAPI, userGroupsAPI } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import { useAppStore } from '@/stores/app'
import type { ApiKey, Group, SelectOption } from '@/types'
import type { CodexConfigMode, CodexConfigWriteResult, CodexDesktopEnvironment } from '@/types/global'

const { t } = useI18n()
const appStore = useAppStore()

const isDesktop = computed(() => typeof window !== 'undefined' && !!window.sub2apiDesktop)
const environment = ref<CodexDesktopEnvironment | null>(null)
const writeResult = ref<CodexConfigWriteResult | null>(null)
const usageResult = ref<Record<string, any> | null>(null)
const loading = ref(false)
const creatingKey = ref(false)
const writing = ref(false)
const savingServer = ref(false)
const message = ref('')
const messageKind = ref<'info' | 'error'>('info')
const openaiGroups = ref<Group[]>([])
const apiKeys = ref<ApiKey[]>([])
const selectedGroupId = ref<number | null>(null)
const selectedKeyId = ref<number | null>(null)
const serverUrlDraft = ref('')
const apiBaseUrlDraft = ref('')
const configMode = ref<CodexConfigMode>('standard')

const fallbackOrigin = computed(() => serverUrlDraft.value || window.location.origin)
const apiBaseUrl = computed(() => normalizeUrl(apiBaseUrlDraft.value || appStore.apiBaseUrl || fallbackOrigin.value))
const filteredKeys = computed(() => apiKeys.value.filter((key) => key.group_id === selectedGroupId.value))
const selectedKey = computed(() => filteredKeys.value.find((key) => key.id === selectedKeyId.value) || null)
const canWrite = computed(() => isDesktop.value && selectedKey.value && apiBaseUrl.value)

const configModeOptions = computed<SelectOption[]>(() => [
  { value: 'standard', label: t('desktopCodex.standardMode') },
  { value: 'websocket', label: t('desktopCodex.websocketMode') },
])

const groupOptions = computed<SelectOption[]>(() =>
  openaiGroups.value.map((group) => ({
    value: group.id,
    label: group.name,
  }))
)

const keyOptions = computed<SelectOption[]>(() =>
  filteredKeys.value.map((key) => ({
    value: key.id,
    label: `${key.name} · ${maskKey(key.key)}`,
  }))
)

const platformLabel = computed(() => {
  const platform = environment.value?.platform
  if (platform === 'darwin') return 'macOS'
  if (platform === 'win32') return 'Windows'
  if (platform === 'linux') return 'Linux'
  return platform || '-'
})

function normalizeUrl(value: string): string {
  return String(value || '').trim().replace(/\/+$/, '')
}

function setMessage(text: string, kind: 'info' | 'error' = 'info'): void {
  message.value = text
  messageKind.value = kind
}

function maskKey(value: string): string {
  if (!value) return ''
  if (value.length <= 10) return '********'
  return `${value.slice(0, 6)}${'*'.repeat(Math.min(12, value.length - 10))}${value.slice(-4)}`
}

function formatMoney(value: unknown, unit = 'USD'): string {
  const amount = Number(value || 0)
  if (unit === 'USD') {
    return `$${amount.toLocaleString('en-US', { maximumFractionDigits: 4 })}`
  }
  return `${amount.toLocaleString('zh-CN', { maximumFractionDigits: 4 })} ${unit || ''}`.trim()
}

function formatUpdatedAt(value: unknown): string {
  if (!value || typeof value !== 'string') return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}

async function loadEnvironment(): Promise<void> {
  if (!window.sub2apiDesktop) return
  environment.value = await window.sub2apiDesktop.getEnvironment()
  if (environment.value.existingBaseUrl) {
    apiBaseUrlDraft.value = environment.value.existingBaseUrl
  }
  configMode.value = environment.value.existingMode || 'standard'
}

function applyAutoSelection(): void {
  if (!selectedGroupId.value && openaiGroups.value.length === 1) {
    selectedGroupId.value = openaiGroups.value[0].id
  }

  if (!selectedKeyId.value && filteredKeys.value.length === 1) {
    selectedKeyId.value = filteredKeys.value[0].id
  }
}

async function loadRemoteData(): Promise<void> {
  loading.value = true
  try {
    const [groups, keys] = await Promise.all([
      userGroupsAPI.getAvailable(),
      keysAPI.list(1, 100, { status: 'active' }),
    ])
    openaiGroups.value = groups.filter((group) => group.platform === 'openai' && group.status === 'active')
    const openaiGroupIds = new Set(openaiGroups.value.map((group) => group.id))
    apiKeys.value = keys.items
      .filter((key) => key.status === 'active' && key.group_id !== null && openaiGroupIds.has(key.group_id))
      .map((key) => ({
        ...key,
        group: key.group || openaiGroups.value.find((group) => group.id === key.group_id),
      }))
    applyAutoSelection()
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function createCodexKey(): Promise<void> {
  if (!selectedGroupId.value) return
  creatingKey.value = true
  try {
    const created = await keysAPI.create('Codex Desktop', selectedGroupId.value)
    apiKeys.value.unshift({
      ...created,
      group: created.group || openaiGroups.value.find((group) => group.id === selectedGroupId.value),
    })
    selectedKeyId.value = created.id
    appStore.showSuccess(t('desktopCodex.keyCreated'))
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.createKeyFailed'))
  } finally {
    creatingKey.value = false
  }
}

async function writeConfig(): Promise<void> {
  if (!window.sub2apiDesktop || !selectedKey.value) return
  writing.value = true
  writeResult.value = null
  usageResult.value = null
  setMessage('')
  try {
    writeResult.value = await window.sub2apiDesktop.writeCodexConfig({
      baseUrl: apiBaseUrl.value,
      apiKey: selectedKey.value.key,
      mode: configMode.value,
    })
    environment.value = await window.sub2apiDesktop.getEnvironment()
    usageResult.value = await window.sub2apiDesktop.queryGatewayUsage({
      baseUrl: apiBaseUrl.value,
      apiKey: selectedKey.value.key,
    })
    appStore.showSuccess(t('desktopCodex.writeSuccess'))
  } catch (error: any) {
    setMessage(error?.message || t('desktopCodex.writeFailed'), 'error')
  } finally {
    writing.value = false
  }
}

async function downloadCodex(): Promise<void> {
  if (!window.sub2apiDesktop) return
  try {
    await window.sub2apiDesktop.openCodexDownload()
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.downloadFailed'))
  }
}

async function openConfigDir(): Promise<void> {
  if (!window.sub2apiDesktop) return
  try {
    await window.sub2apiDesktop.openConfigDir()
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.openConfigDirFailed'))
  }
}

async function saveServerUrl(): Promise<void> {
  if (!window.sub2apiDesktop) return
  savingServer.value = true
  try {
    await window.sub2apiDesktop.setServerUrl(serverUrlDraft.value)
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.switchServerFailed'))
  } finally {
    savingServer.value = false
  }
}

watch(selectedGroupId, () => {
  if (!filteredKeys.value.some((key) => key.id === selectedKeyId.value)) {
    selectedKeyId.value = filteredKeys.value[0]?.id || null
  }
})

onMounted(async () => {
  const publicSettings = await appStore.fetchPublicSettings()
  apiBaseUrlDraft.value = publicSettings?.api_base_url || window.location.origin
  if (window.sub2apiDesktop) {
    serverUrlDraft.value = await window.sub2apiDesktop.getServerUrl()
    await loadEnvironment()
  }
  await loadRemoteData()
})
</script>
