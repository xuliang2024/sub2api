<template>
  <AppLayout>
    <div class="mx-auto flex max-w-6xl flex-col gap-5">
      <section
        v-if="!isDesktop"
        class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-200"
      >
        {{ t('desktopCodex.webOnlyNotice') }}
      </section>

      <section class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="grid lg:grid-cols-[1.08fr_0.92fr]">
          <div class="p-6 lg:p-7">
            <div class="max-w-2xl">
              <p class="text-xs font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">
                {{ t('desktopCodex.connection') }}
              </p>
              <h2 class="mt-2 text-2xl font-semibold text-gray-950 dark:text-white">
                {{ t('desktopCodex.guidedTitle') }}
              </h2>
              <p class="mt-2 text-sm leading-6 text-gray-500 dark:text-dark-400">
                {{ t('desktopCodex.guidedDescription') }}
              </p>
            </div>

            <div class="mt-6 grid gap-3 sm:grid-cols-2">
              <button
                class="group flex min-h-[104px] items-center gap-4 rounded-lg border border-gray-200 bg-white px-5 py-4 text-left shadow-sm transition hover:-translate-y-0.5 hover:border-primary-200 hover:shadow-md disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0 dark:border-dark-700 dark:bg-dark-800 dark:hover:border-primary-800"
                :disabled="!isDesktop"
                type="button"
                @click="downloadCodex"
              >
                <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-600 transition group-hover:bg-primary-100 dark:bg-primary-900/30 dark:text-primary-300">
                  <Icon name="download" size="lg" />
                </span>
                <span class="min-w-0">
                  <span class="block text-base font-semibold text-gray-950 dark:text-white">
                    {{ t('desktopCodex.downloadCodex') }}
                  </span>
                  <span class="mt-1 block text-sm leading-5 text-gray-500 dark:text-dark-400">
                    {{ t('desktopCodex.downloadHint') }}
                  </span>
                </span>
              </button>

              <button
                class="flex min-h-[104px] items-center gap-4 rounded-lg bg-gradient-to-r from-primary-500 to-primary-700 px-5 py-4 text-left text-white shadow-xl shadow-primary-500/30 transition hover:-translate-y-0.5 hover:shadow-2xl hover:shadow-primary-500/35 disabled:cursor-not-allowed disabled:opacity-60 disabled:shadow-none disabled:hover:translate-y-0 dark:shadow-primary-950/40"
                :disabled="!canWrite || writing || loading"
                type="button"
                @click="writeConfig"
              >
                <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-white/15 text-white ring-1 ring-white/25">
                  <Icon name="terminal" size="lg" />
                </span>
                <span class="min-w-0">
                  <span class="block text-lg font-semibold">
                    {{ writing ? t('desktopCodex.writing') : t('desktopCodex.writeConfig') }}
                  </span>
                  <span class="mt-1 block text-sm leading-5 text-primary-50">
                    {{ writeButtonHint }}
                  </span>
                </span>
              </button>
            </div>

            <p
              v-if="message"
              class="mt-5 rounded-lg p-3 text-sm"
              :class="messageKind === 'error' ? 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300' : 'bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300'"
            >
              {{ message }}
            </p>

            <div
              v-if="writeResult"
              class="mt-5 rounded-lg border border-green-200 bg-green-50 p-4 text-sm text-green-800 dark:border-green-800 dark:bg-green-900/20 dark:text-green-200"
            >
              <div class="flex items-start gap-3">
                <Icon name="checkCircle" size="md" class="mt-0.5 shrink-0" />
                <div class="min-w-0">
                  <p class="font-medium">{{ t('desktopCodex.writeSuccess') }}</p>
                  <p v-if="writeResult.configBackup" class="mt-1 break-all text-xs">
                    config backup: {{ writeResult.configBackup }}
                  </p>
                  <p v-if="writeResult.authBackup" class="mt-1 break-all text-xs">
                    auth backup: {{ writeResult.authBackup }}
                  </p>
                </div>
              </div>
            </div>

            <div v-if="usageResult" class="mt-5 grid gap-3 md:grid-cols-3">
              <div class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800">
                <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.available') }}</p>
                <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ formatMoney(usageResult.available, usageResult.unit) }}</strong>
              </div>
              <div class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800">
                <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.status') }}</p>
                <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ usageResult.status || usageResult.mode || '-' }}</strong>
              </div>
              <div class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800">
                <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.updatedAt') }}</p>
                <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ formatUpdatedAt(usageResult.updated_at) }}</strong>
              </div>
            </div>
          </div>

          <div class="border-t border-gray-100 bg-gray-50/80 p-6 dark:border-dark-800 dark:bg-dark-950/40 lg:border-l lg:border-t-0 lg:p-7">
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                  {{ t('desktopCodex.autoSetup') }}
                </p>
                <h3 class="mt-1 text-lg font-semibold text-gray-950 dark:text-white">
                  {{ t('desktopCodex.readyToWrite') }}
                </h3>
              </div>
              <span
                class="shrink-0 rounded-full px-2.5 py-1 text-xs font-medium"
                :class="isDesktop ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'"
              >
                {{ isDesktop ? t('desktopCodex.desktopReady') : t('desktopCodex.webMode') }}
              </span>
            </div>

            <div class="mt-5 space-y-3">
              <div class="flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
                <span class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
                  <Icon name="server" size="sm" />
                </span>
                <div class="min-w-0">
                  <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('desktopCodex.serviceEndpoint') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-gray-950 dark:text-white">{{ apiBaseUrl || '-' }}</p>
                </div>
              </div>

              <div class="flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
                <span class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
                  <Icon name="users" size="sm" />
                </span>
                <div class="min-w-0">
                  <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('desktopCodex.autoSelectedGroup') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-gray-950 dark:text-white">{{ selectedGroupName }}</p>
                </div>
              </div>

              <div class="flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-900">
                <span class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
                  <Icon name="key" size="sm" />
                </span>
                <div class="min-w-0">
                  <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('desktopCodex.autoSelectedKey') }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-gray-950 dark:text-white">{{ selectedKeyName }}</p>
                  <p class="mt-1 truncate text-xs text-gray-500 dark:text-dark-400">{{ selectedKeySecret }}</p>
                </div>
              </div>
            </div>

            <p
              v-if="openaiGroups.length === 0 && !loading"
              class="mt-4 rounded-lg bg-amber-50 p-3 text-sm text-amber-700 dark:bg-amber-900/20 dark:text-amber-300"
            >
              {{ t('desktopCodex.noOpenaiGroups') }}
            </p>

            <button class="btn btn-secondary mt-4 w-full" :disabled="loading" type="button" @click="loadRemoteData">
              <Icon name="refresh" size="sm" :class="loading ? 'mr-2 animate-spin' : 'mr-2'" />
              {{ t('desktopCodex.refreshSetup') }}
            </button>
          </div>
        </div>
      </section>

      <details class="group rounded-lg border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <summary class="flex cursor-pointer list-none items-center justify-between gap-4 px-5 py-4">
          <div>
            <p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
              {{ t('desktopCodex.advancedSettings') }}
            </p>
            <h3 class="mt-1 text-base font-semibold text-gray-900 dark:text-white">
              {{ t('desktopCodex.advancedSettingsHint') }}
            </h3>
          </div>
          <Icon name="chevronDown" size="md" class="shrink-0 text-gray-400 transition group-open:rotate-180" />
        </summary>

        <div class="border-t border-gray-100 p-5 dark:border-dark-800">
          <div class="grid gap-4 lg:grid-cols-3">
            <label class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.serverUrl') }}</span>
              <input
                v-model="serverUrlDraft"
                class="input w-full"
                type="url"
                spellcheck="false"
                :disabled="!isDesktop"
                placeholder="https://codex.apiz.ai"
              />
            </label>

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

            <div class="space-y-1">
              <span class="input-label">{{ t('desktopCodex.desktopServer') }}</span>
              <button class="btn btn-primary w-full" :disabled="!isDesktop || savingServer" type="button" @click="saveServerUrl">
                <Icon name="swap" size="sm" class="mr-2" />
                {{ savingServer ? t('common.saving') : t('desktopCodex.switchServer') }}
              </button>
            </div>
          </div>

          <div class="mt-4 grid gap-3 lg:grid-cols-3">
            <article class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800">
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('desktopCodex.platform') }}</p>
              <strong class="mt-1 block text-base text-gray-900 dark:text-white">{{ platformLabel }}</strong>
            </article>
            <article class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800">
              <p class="text-xs text-gray-500 dark:text-dark-400">config.toml</p>
              <code class="mt-1 block truncate text-sm text-gray-900 dark:text-white">{{ environment?.configPath || '-' }}</code>
            </article>
            <article class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800">
              <p class="text-xs text-gray-500 dark:text-dark-400">auth.json</p>
              <code class="mt-1 block truncate text-sm text-gray-900 dark:text-white">{{ environment?.authPath || '-' }}</code>
            </article>
          </div>

          <div class="mt-4 flex flex-wrap gap-3">
            <button
              class="btn btn-secondary"
              :disabled="!selectedGroupId || creatingKey"
              type="button"
              @click="createCodexKey"
            >
              <Icon name="plus" size="sm" class="mr-2" />
              {{ creatingKey ? t('desktopCodex.creatingKey') : t('desktopCodex.createKey') }}
            </button>
            <button class="btn btn-secondary" :disabled="!isDesktop" type="button" @click="openConfigDir">
              <Icon name="document" size="sm" class="mr-2" />
              {{ t('desktopCodex.openConfigDir') }}
            </button>
          </div>
        </div>
      </details>
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
const selectedGroup = computed(() => openaiGroups.value.find((group) => group.id === selectedGroupId.value) || null)
const selectedKey = computed(() => filteredKeys.value.find((key) => key.id === selectedKeyId.value) || null)
const canWrite = computed(() => isDesktop.value && !!selectedGroup.value && !!apiBaseUrl.value)

const selectedGroupName = computed(() => {
  if (loading.value) return t('common.loading')
  return selectedGroup.value?.name || t('desktopCodex.notReady')
})

const selectedKeyName = computed(() => {
  if (loading.value) return t('common.loading')
  if (selectedKey.value) return selectedKey.value.name
  if (selectedGroup.value) return t('desktopCodex.willCreateKey')
  return t('desktopCodex.notReady')
})

const selectedKeySecret = computed(() => {
  if (selectedKey.value) return maskKey(selectedKey.value.key)
  if (selectedGroup.value) return t('desktopCodex.autoCreateKeyShort')
  return '-'
})

const writeButtonHint = computed(() => {
  if (!isDesktop.value) return t('desktopCodex.needDesktop')
  if (loading.value) return t('common.loading')
  if (!selectedGroup.value) return t('desktopCodex.needOpenaiGroup')
  if (!selectedKey.value) return t('desktopCodex.autoCreateKeyHint')
  return t('desktopCodex.writeButtonHint')
})

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

function normalizeName(value: unknown): string {
  return String(value || '').trim().toLowerCase()
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

function pickCodexGroup(groups: Group[]): Group | null {
  return (
    groups.find((group) => normalizeName(group.name) === 'codex') ||
    groups.find((group) => normalizeName(group.name).includes('codex')) ||
    groups[0] ||
    null
  )
}

function pickCodexKey(keys: ApiKey[]): ApiKey | null {
  return (
    keys.find((key) => normalizeName(key.name) === 'codex desktop') ||
    keys.find((key) => normalizeName(key.name).includes('codex')) ||
    keys[0] ||
    null
  )
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
  if (selectedGroupId.value && !openaiGroups.value.some((group) => group.id === selectedGroupId.value)) {
    selectedGroupId.value = null
  }

  if (!selectedGroupId.value) {
    selectedGroupId.value = pickCodexGroup(openaiGroups.value)?.id || null
  }

  if (selectedKeyId.value && !filteredKeys.value.some((key) => key.id === selectedKeyId.value)) {
    selectedKeyId.value = null
  }

  if (!selectedKeyId.value) {
    selectedKeyId.value = pickCodexKey(filteredKeys.value)?.id || null
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

async function createKeyForSelectedGroup(showToast = true): Promise<ApiKey | null> {
  if (!selectedGroupId.value) return null
  creatingKey.value = true
  try {
    const created = await keysAPI.create('Codex Desktop', selectedGroupId.value)
    const codexKey = {
      ...created,
      group: created.group || openaiGroups.value.find((group) => group.id === selectedGroupId.value),
    }
    apiKeys.value.unshift(codexKey)
    selectedKeyId.value = created.id
    if (showToast) {
      appStore.showSuccess(t('desktopCodex.keyCreated'))
    }
    return codexKey
  } catch (error: any) {
    appStore.showError(error?.message || t('desktopCodex.createKeyFailed'))
    return null
  } finally {
    creatingKey.value = false
  }
}

async function createCodexKey(): Promise<void> {
  await createKeyForSelectedGroup(true)
}

async function ensureSelectedKey(): Promise<ApiKey | null> {
  if (selectedKey.value) return selectedKey.value
  return createKeyForSelectedGroup(false)
}

async function writeConfig(): Promise<void> {
  if (!window.sub2apiDesktop || !selectedGroup.value) return
  writing.value = true
  writeResult.value = null
  usageResult.value = null
  setMessage('')
  try {
    const key = await ensureSelectedKey()
    if (!key) {
      setMessage(t('desktopCodex.createKeyFailed'), 'error')
      return
    }
    writeResult.value = await window.sub2apiDesktop.writeCodexConfig({
      baseUrl: apiBaseUrl.value,
      apiKey: key.key,
      mode: configMode.value,
    })
    environment.value = await window.sub2apiDesktop.getEnvironment()
    usageResult.value = await window.sub2apiDesktop.queryGatewayUsage({
      baseUrl: apiBaseUrl.value,
      apiKey: key.key,
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
  selectedKeyId.value = pickCodexKey(filteredKeys.value)?.id || null
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
