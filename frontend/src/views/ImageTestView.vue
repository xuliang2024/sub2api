<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200 bg-white/90 px-5 py-4 backdrop-blur dark:border-dark-800 dark:bg-dark-950/90">
      <nav class="mx-auto flex max-w-7xl items-center justify-between gap-4">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-gray-900 text-white shadow-sm dark:bg-white dark:text-gray-950">
            <Icon name="beaker" size="md" />
          </div>
          <div>
            <p class="text-sm font-semibold leading-5">Image API Test</p>
            <p class="text-xs text-gray-500 dark:text-dark-400">gpt-image-2 / OpenAI compatible</p>
          </div>
        </router-link>

        <div class="flex items-center gap-2">
          <span
            class="hidden rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-xs font-medium text-emerald-700 dark:border-emerald-900/60 dark:bg-emerald-950/30 dark:text-emerald-300 sm:inline-flex"
          >
            {{ localProxyEnabled ? 'Local proxy' : 'Direct URL' }}
          </span>
          <router-link
            to="/key-usage"
            class="inline-flex h-9 items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 transition-colors hover:border-gray-300 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-dark-800"
          >
            <Icon name="key" size="sm" />
            Key Usage
          </router-link>
        </div>
      </nav>
    </header>

    <main class="mx-auto grid max-w-7xl gap-5 px-5 py-5 lg:grid-cols-[420px_minmax(0,1fr)]">
      <section class="space-y-5">
        <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="mb-4 flex items-center justify-between gap-3">
            <div>
              <h1 class="text-lg font-semibold">本地图像接口测试</h1>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
                文生图、图生图、尺寸和分辨率矩阵
              </p>
            </div>
            <Icon name="sparkles" size="lg" class="text-amber-500" />
          </div>

          <div class="grid grid-cols-2 gap-2 rounded-lg bg-gray-100 p-1 dark:bg-dark-800">
            <button
              type="button"
              class="rounded-md px-3 py-2 text-sm font-medium transition-colors"
              :class="mode === 'text' ? 'bg-white text-gray-950 shadow-sm dark:bg-dark-700 dark:text-white' : 'text-gray-500 hover:text-gray-800 dark:text-dark-300 dark:hover:text-white'"
              @click="mode = 'text'"
            >
              文生图
            </button>
            <button
              type="button"
              class="rounded-md px-3 py-2 text-sm font-medium transition-colors"
              :class="mode === 'edit' ? 'bg-white text-gray-950 shadow-sm dark:bg-dark-700 dark:text-white' : 'text-gray-500 hover:text-gray-800 dark:text-dark-300 dark:hover:text-white'"
              @click="mode = 'edit'"
            >
              图生图
            </button>
          </div>

          <div class="mt-4 space-y-4">
            <Input
              v-model="gatewayBaseUrl"
              label="Gateway URL"
              placeholder="https://codex.apiz.ai"
              autocomplete="off"
            />

            <div>
              <label class="input-label mb-1.5 block">API Key</label>
              <div class="relative">
                <input
                  v-model="apiKey"
                  :type="showKey ? 'text' : 'password'"
                  class="input w-full pr-11"
                  placeholder="sk-..."
                  autocomplete="off"
                  @keydown.enter="runSingle"
                />
                <button
                  type="button"
                  class="absolute inset-y-0 right-0 flex w-10 items-center justify-center text-gray-400 hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
                  :title="showKey ? '隐藏 key' : '显示 key'"
                  @click="showKey = !showKey"
                >
                  <Icon :name="showKey ? 'eyeOff' : 'eye'" size="sm" />
                </button>
              </div>
            </div>

            <Input v-model="model" label="Model" placeholder="gpt-image-2" autocomplete="off" />

            <div>
              <label class="input-label mb-1.5 block">Prompt</label>
              <textarea
                v-model="prompt"
                rows="5"
                class="input min-h-[126px] w-full resize-y"
                placeholder="一张清晰的产品海报，白底，高级商业摄影..."
              ></textarea>
            </div>
          </div>
        </div>

        <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-300">尺寸和输出</h2>
            <span class="rounded-md bg-cyan-50 px-2 py-1 text-xs font-medium text-cyan-700 dark:bg-cyan-950/40 dark:text-cyan-300">
              {{ activeSize }}
            </span>
          </div>

          <div class="space-y-4">
            <div>
              <label class="input-label mb-1.5 block">单次尺寸</label>
              <Select v-model="sizePreset" :options="sizeOptions" />
            </div>

            <Input
              v-if="sizePreset === 'custom'"
              v-model="customSize"
              label="自定义尺寸"
              placeholder="2048x1152"
              autocomplete="off"
            />

            <div class="grid grid-cols-2 gap-3">
              <label>
                <span class="input-label mb-1.5 block">Quality</span>
                <Select v-model="quality" :options="qualityOptions" />
              </label>
              <label>
                <span class="input-label mb-1.5 block">Response</span>
                <Select v-model="responseFormat" :options="responseFormatOptions" />
              </label>
              <label>
                <span class="input-label mb-1.5 block">Output</span>
                <Select v-model="outputFormat" :options="outputFormatOptions" />
              </label>
              <label>
                <span class="input-label mb-1.5 block">Background</span>
                <Select v-model="background" :options="backgroundOptions" />
              </label>
            </div>

            <Input
              v-model="imageCount"
              type="number"
              label="Images"
              placeholder="1"
              autocomplete="off"
            />
          </div>
        </div>

        <div
          v-if="mode === 'edit'"
          class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900"
        >
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-300">图生图输入</h2>
            <Select v-model="inputFidelity" class="w-32" :options="inputFidelityOptions" />
          </div>

          <div class="space-y-4">
            <Input v-model="sourceImageUrl" label="源图 URL" placeholder="https://.../source.png" autocomplete="off" />
            <div class="rounded-lg border border-dashed border-gray-300 p-3 dark:border-dark-700">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-medium">上传源图</p>
                  <p class="text-xs text-gray-500 dark:text-dark-400">{{ sourceImageName || 'PNG / JPG / WebP，自动转成 data URL' }}</p>
                </div>
                <label class="inline-flex h-9 cursor-pointer items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-100">
                  <Icon name="upload" size="sm" />
                  选择
                  <input class="sr-only" type="file" accept="image/*" @change="handleImageFile($event, 'source')" />
                </label>
              </div>
              <img v-if="sourceImagePreview" :src="sourceImagePreview" alt="Source preview" class="mt-3 aspect-video w-full rounded-lg object-contain bg-gray-100 dark:bg-dark-950" />
            </div>

            <Input v-model="maskImageUrl" label="Mask URL（可选）" placeholder="https://.../mask.png" autocomplete="off" />
            <div class="rounded-lg border border-dashed border-gray-300 p-3 dark:border-dark-700">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-medium">上传 Mask（可选）</p>
                  <p class="text-xs text-gray-500 dark:text-dark-400">{{ maskImageName || '透明区域会被重绘' }}</p>
                </div>
                <label class="inline-flex h-9 cursor-pointer items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-100">
                  <Icon name="upload" size="sm" />
                  选择
                  <input class="sr-only" type="file" accept="image/*" @change="handleImageFile($event, 'mask')" />
                </label>
              </div>
              <img v-if="maskImagePreview" :src="maskImagePreview" alt="Mask preview" class="mt-3 aspect-video w-full rounded-lg object-contain bg-gray-100 dark:bg-dark-950" />
            </div>
          </div>
        </div>

        <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-300">批量尺寸</h2>
            <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-dark-400">
              <Toggle v-model="localProxyEnabled" />
              Proxy
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2">
            <button
              v-for="preset in batchSizeOptions"
              :key="String(preset.value)"
              type="button"
              class="rounded-lg border px-3 py-2 text-left text-xs transition-colors"
              :class="selectedBatchSizes.includes(String(preset.value))
                ? 'border-cyan-500 bg-cyan-50 text-cyan-800 dark:border-cyan-700 dark:bg-cyan-950/40 dark:text-cyan-200'
                : 'border-gray-200 bg-white text-gray-600 hover:border-gray-300 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:border-dark-600'"
              @click="toggleBatchSize(String(preset.value))"
            >
              <span class="block font-semibold">{{ preset.value }}</span>
              <span class="text-gray-500 dark:text-dark-400">{{ preset.label }}</span>
            </button>
          </div>

          <div class="mt-4 grid gap-2 sm:grid-cols-2">
            <button
              type="button"
              class="inline-flex h-11 items-center justify-center gap-2 rounded-lg bg-gray-900 px-4 text-sm font-semibold text-white transition-colors hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-white dark:text-gray-950 dark:hover:bg-gray-100"
              :disabled="isRunning || !canRun"
              @click="runSingle"
            >
              <Icon :name="isRunning ? 'refresh' : 'play'" size="sm" :class="isRunning ? 'animate-spin' : ''" />
              运行当前测试
            </button>
            <button
              type="button"
              class="inline-flex h-11 items-center justify-center gap-2 rounded-lg border border-gray-200 bg-white px-4 text-sm font-semibold text-gray-800 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-100 dark:hover:bg-dark-700"
              :disabled="isRunning || !canRun || selectedBatchSizes.length === 0"
              @click="runBatch"
            >
              <Icon name="grid" size="sm" />
              跑选中尺寸
            </button>
          </div>

          <p v-if="validationMessage" class="mt-3 rounded-lg bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:bg-amber-950/40 dark:text-amber-200">
            {{ validationMessage }}
          </p>
        </div>
      </section>

      <section class="space-y-5">
        <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="mb-3 flex items-center justify-between gap-3">
            <div>
              <h2 class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-300">请求预览</h2>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ requestEndpoint }}</p>
            </div>
            <button
              type="button"
              class="inline-flex h-8 items-center gap-2 rounded-lg border border-gray-200 px-3 text-xs font-medium text-gray-600 hover:bg-gray-50 dark:border-dark-700 dark:text-dark-300 dark:hover:bg-dark-800"
              @click="copyText(JSON.stringify(redactedCurrentPayload, null, 2))"
            >
              <Icon name="copy" size="xs" />
              Copy
            </button>
          </div>
          <pre class="max-h-72 overflow-auto rounded-lg bg-gray-950 p-4 text-xs leading-5 text-gray-100">{{ JSON.stringify(redactedCurrentPayload, null, 2) }}</pre>
        </div>

        <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-800 dark:bg-dark-900">
          <div class="mb-4 flex items-center justify-between">
            <div>
              <h2 class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-300">测试结果</h2>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ results.length ? `${results.length} requests` : '还没有请求记录' }}</p>
            </div>
            <button
              type="button"
              class="inline-flex h-8 items-center gap-2 rounded-lg border border-gray-200 px-3 text-xs font-medium text-gray-600 hover:bg-gray-50 dark:border-dark-700 dark:text-dark-300 dark:hover:bg-dark-800"
              :disabled="results.length === 0"
              @click="results = []"
            >
              <Icon name="trash" size="xs" />
              Clear
            </button>
          </div>

          <div v-if="results.length === 0" class="flex min-h-[420px] items-center justify-center rounded-lg border border-dashed border-gray-300 bg-gray-50 dark:border-dark-700 dark:bg-dark-950">
            <div class="max-w-sm px-6 text-center">
              <Icon name="beaker" size="xl" class="mx-auto mb-3 text-gray-400" />
              <p class="text-sm font-medium text-gray-700 dark:text-dark-200">准备好后运行一次测试</p>
              <p class="mt-1 text-xs leading-5 text-gray-500 dark:text-dark-400">
                结果会按请求记录展示图片、延迟、尺寸、错误和响应摘要。
              </p>
            </div>
          </div>

          <div v-else class="grid gap-4 xl:grid-cols-2">
            <article
              v-for="result in results"
              :key="result.id"
              class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-950"
            >
              <div class="flex items-start justify-between gap-3 border-b border-gray-200 p-3 dark:border-dark-800">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span
                      class="inline-flex rounded-md px-2 py-1 text-xs font-semibold"
                      :class="result.status === 'success'
                        ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300'
                        : result.status === 'error'
                          ? 'bg-rose-50 text-rose-700 dark:bg-rose-950/40 dark:text-rose-300'
                          : 'bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-300'"
                    >
                      {{ result.status }}
                    </span>
                    <span class="text-xs text-gray-500 dark:text-dark-400">{{ result.mode === 'text' ? 'text-to-image' : 'image-to-image' }}</span>
                    <span class="text-xs text-gray-500 dark:text-dark-400">{{ result.size }}</span>
                  </div>
                  <p class="mt-2 text-xs text-gray-500 dark:text-dark-400">
                    {{ result.endpoint }} · {{ result.latencyMs ? `${result.latencyMs}ms` : 'running' }}
                  </p>
                </div>
                <button
                  type="button"
                  class="rounded-lg p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dark-800 dark:hover:text-white"
                  title="Copy response"
                  @click="copyText(JSON.stringify(result.responsePreview || result.error || {}, null, 2))"
                >
                  <Icon name="copy" size="sm" />
                </button>
              </div>

              <div class="p-3">
                <div v-if="result.status === 'pending'" class="flex aspect-square items-center justify-center rounded-lg bg-gray-100 dark:bg-dark-900">
                  <Icon name="refresh" size="lg" class="animate-spin text-gray-400" />
                </div>

                <div v-else-if="result.status === 'error'" class="rounded-lg bg-rose-50 p-3 text-sm text-rose-800 dark:bg-rose-950/40 dark:text-rose-200">
                  {{ result.error }}
                </div>

                <div v-else class="space-y-3">
                  <div v-for="image in result.images" :key="image.id" class="overflow-hidden rounded-lg bg-gray-100 dark:bg-dark-900">
                    <img :src="image.src" :alt="image.alt" class="aspect-square w-full object-contain" />
                    <div class="flex items-center justify-between gap-2 border-t border-gray-200 bg-white p-2 dark:border-dark-800 dark:bg-dark-950">
                      <span class="truncate text-xs text-gray-500 dark:text-dark-400">{{ image.label }}</span>
                      <div class="flex items-center gap-1">
                        <button
                          type="button"
                          class="rounded-md p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
                          title="Open image"
                          @click="openImage(image.src)"
                        >
                          <Icon name="externalLink" size="xs" />
                        </button>
                        <button
                          type="button"
                          class="rounded-md p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
                          title="Download image"
                          @click="downloadImage(image.src, image.filename)"
                        >
                          <Icon name="download" size="xs" />
                        </button>
                      </div>
                    </div>
                  </div>

                  <p v-if="result.revisedPrompt" class="rounded-lg bg-gray-50 p-3 text-xs leading-5 text-gray-600 dark:bg-dark-900 dark:text-dark-300">
                    {{ result.revisedPrompt }}
                  </p>
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import Input from '@/components/common/Input.vue'
import Select, { type SelectOption } from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'

type ImageMode = 'text' | 'edit'
type ResultStatus = 'pending' | 'success' | 'error'
type FileSlot = 'source' | 'mask'

interface ResultImage {
  id: string
  src: string
  alt: string
  label: string
  filename: string
}

interface ImageTestResult {
  id: string
  mode: ImageMode
  size: string
  endpoint: string
  status: ResultStatus
  latencyMs: number | null
  images: ResultImage[]
  revisedPrompt?: string
  responsePreview?: unknown
  error?: string
}

const gatewayBaseUrl = ref('https://codex.apiz.ai')
const apiKey = ref('')
const showKey = ref(false)
const model = ref('gpt-image-2')
const prompt = ref('A clean product poster of a matte black wireless speaker on a white background, studio lighting, crisp shadow')
const mode = ref<ImageMode>('text')
const sizePreset = ref('1024x1024')
const customSize = ref('2048x1152')
const quality = ref('auto')
const responseFormat = ref('b64_json')
const outputFormat = ref('png')
const background = ref('auto')
const imageCount = ref(1)
const inputFidelity = ref('auto')
const localProxyEnabled = ref(true)

const sourceImageUrl = ref('')
const sourceImageDataUrl = ref('')
const sourceImageName = ref('')
const maskImageUrl = ref('')
const maskImageDataUrl = ref('')
const maskImageName = ref('')

const selectedBatchSizes = ref<string[]>(['1024x1024', '1536x1024', '1024x1536'])
const results = ref<ImageTestResult[]>([])
const isRunning = ref(false)

const sizeOptions: SelectOption[] = [
  { value: 'auto', label: 'auto' },
  { value: '1024x1024', label: '1024x1024 · square 1K' },
  { value: '1536x1024', label: '1536x1024 · landscape 2K' },
  { value: '1024x1536', label: '1024x1536 · portrait 2K' },
  { value: '1792x1024', label: '1792x1024 · wide' },
  { value: '1024x1792', label: '1024x1792 · tall' },
  { value: '2048x2048', label: '2048x2048 · square 2K' },
  { value: '2560x1440', label: '2560x1440 · 16:9 QHD' },
  { value: '3840x2160', label: '3840x2160 · 4K wide' },
  { value: '2160x3840', label: '2160x3840 · 4K tall' },
  { value: 'custom', label: 'custom' },
]

const batchSizeOptions = sizeOptions.filter((option) => !['auto', 'custom'].includes(String(option.value)))

const qualityOptions: SelectOption[] = [
  { value: 'auto', label: 'auto' },
  { value: 'low', label: 'low' },
  { value: 'medium', label: 'medium' },
  { value: 'high', label: 'high' },
]

const responseFormatOptions: SelectOption[] = [
  { value: 'b64_json', label: 'b64_json' },
  { value: 'url', label: 'url' },
]

const outputFormatOptions: SelectOption[] = [
  { value: 'png', label: 'png' },
  { value: 'jpeg', label: 'jpeg' },
  { value: 'webp', label: 'webp' },
]

const backgroundOptions: SelectOption[] = [
  { value: 'auto', label: 'auto' },
  { value: 'opaque', label: 'opaque' },
  { value: 'transparent', label: 'transparent' },
]

const inputFidelityOptions: SelectOption[] = [
  { value: 'auto', label: 'auto' },
  { value: 'low', label: 'low' },
  { value: 'high', label: 'high' },
]

const activeSize = computed(() => {
  if (sizePreset.value === 'custom') {
    return customSize.value.trim()
  }
  return sizePreset.value
})

const sourceImageValue = computed(() => sourceImageUrl.value.trim() || sourceImageDataUrl.value)
const sourceImagePreview = computed(() => sourceImageDataUrl.value || sourceImageUrl.value.trim())
const maskImageValue = computed(() => maskImageUrl.value.trim() || maskImageDataUrl.value)
const maskImagePreview = computed(() => maskImageDataUrl.value || maskImageUrl.value.trim())

const requestEndpoint = computed(() => mode.value === 'text' ? '/v1/images/generations' : '/v1/images/edits')

const validationMessage = computed(() => {
  if (!apiKey.value.trim()) return '需要 API Key 才能请求。'
  if (!prompt.value.trim()) return '需要填写 prompt。'
  if (!activeSize.value) return '需要选择或填写尺寸。'
  if (mode.value === 'edit' && !sourceImageValue.value) return '图生图需要源图 URL 或上传源图。'
  return ''
})

const canRun = computed(() => validationMessage.value === '')

const redactedCurrentPayload = computed(() => redactPayload(buildPayload(activeSize.value)))

onMounted(() => {
  const storedGateway = window.localStorage.getItem('sub2api:image-test:gateway')
  const storedProxy = window.localStorage.getItem('sub2api:image-test:proxy')
  const sessionKey = window.sessionStorage.getItem('sub2api:image-test:key')
  if (storedGateway) gatewayBaseUrl.value = storedGateway
  if (storedProxy) localProxyEnabled.value = storedProxy === 'true'
  if (sessionKey) apiKey.value = sessionKey
})

watch(gatewayBaseUrl, (value) => {
  window.localStorage.setItem('sub2api:image-test:gateway', value)
})

watch(localProxyEnabled, (value) => {
  window.localStorage.setItem('sub2api:image-test:proxy', String(value))
})

watch(apiKey, (value) => {
  if (value.trim()) {
    window.sessionStorage.setItem('sub2api:image-test:key', value)
  } else {
    window.sessionStorage.removeItem('sub2api:image-test:key')
  }
})

function buildPayload(size: string) {
  const payload: Record<string, unknown> = {
    model: model.value.trim() || 'gpt-image-2',
    prompt: prompt.value.trim(),
    n: clampNumber(Number(imageCount.value), 1, 4),
    size,
    response_format: responseFormat.value,
  }

  if (quality.value !== 'auto') payload.quality = quality.value
  if (background.value !== 'auto') payload.background = background.value
  if (outputFormat.value) payload.output_format = outputFormat.value

  if (mode.value === 'edit') {
    payload.images = [{ image_url: sourceImageValue.value }]
    if (maskImageValue.value) payload.mask = { image_url: maskImageValue.value }
    if (inputFidelity.value !== 'auto') payload.input_fidelity = inputFidelity.value
  }

  return payload
}

async function runSingle() {
  await runSizes([activeSize.value])
}

async function runBatch() {
  await runSizes([...selectedBatchSizes.value])
}

async function runSizes(sizes: string[]) {
  if (!canRun.value || isRunning.value) return

  isRunning.value = true
  try {
    for (const size of sizes) {
      await runRequest(size)
    }
  } finally {
    isRunning.value = false
  }
}

async function runRequest(size: string) {
  const endpoint = requestEndpoint.value
  const result: ImageTestResult = {
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    mode: mode.value,
    size,
    endpoint,
    status: 'pending',
    latencyMs: null,
    images: [],
  }

  results.value = [result, ...results.value]

  const startedAt = performance.now()
  const payload = buildPayload(size)

  try {
    const response = await fetch(resolveRequestUrl(endpoint), {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${apiKey.value.trim()}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })

    const responseText = await response.text()
    const data = parseJson(responseText)
    result.latencyMs = Math.round(performance.now() - startedAt)

    if (!response.ok) {
      throw new Error(extractErrorMessage(data, responseText, response.status))
    }

    result.images = normalizeImages(data, result)
    result.revisedPrompt = extractRevisedPrompt(data)
    result.responsePreview = redactPayload(data)
    result.status = 'success'

    if (result.images.length === 0) {
      result.status = 'error'
      result.error = '响应成功，但没有找到 data[].url 或 data[].b64_json。'
    }
  } catch (error) {
    result.latencyMs = Math.round(performance.now() - startedAt)
    result.status = 'error'
    result.error = error instanceof Error ? error.message : String(error)
  }
}

function resolveRequestUrl(endpoint: string) {
  if (localProxyEnabled.value) return endpoint
  const base = gatewayBaseUrl.value.trim().replace(/\/+$/, '')
  return `${base}${endpoint}`
}

function normalizeImages(data: unknown, result: ImageTestResult): ResultImage[] {
  if (!data || typeof data !== 'object') return []
  const dataItems = Array.isArray((data as { data?: unknown }).data)
    ? (data as { data: Array<Record<string, unknown>> }).data
    : []

  return dataItems.flatMap((item, index) => {
    const url = typeof item.url === 'string' ? item.url : ''
    const b64 = typeof item.b64_json === 'string' ? item.b64_json : ''
    const src = b64 ? `data:image/${outputFormat.value || 'png'};base64,${b64}` : url
    if (!src) return []

    return [{
      id: `${result.id}-image-${index}`,
      src,
      alt: `${result.mode} ${result.size} result ${index + 1}`,
      label: `${result.size} · ${outputFormat.value}`,
      filename: `image-test-${result.mode}-${result.size}-${index + 1}.${outputFormat.value === 'jpeg' ? 'jpg' : outputFormat.value}`,
    }]
  })
}

function extractRevisedPrompt(data: unknown) {
  if (!data || typeof data !== 'object') return ''
  const items = (data as { data?: unknown }).data
  if (!Array.isArray(items)) return ''
  const first = items.find((item) => item && typeof item === 'object' && typeof (item as { revised_prompt?: unknown }).revised_prompt === 'string')
  return first ? String((first as { revised_prompt: string }).revised_prompt) : ''
}

function extractErrorMessage(data: unknown, text: string, status: number) {
  if (data && typeof data === 'object') {
    const error = (data as { error?: unknown }).error
    if (error && typeof error === 'object' && typeof (error as { message?: unknown }).message === 'string') {
      return `HTTP ${status}: ${(error as { message: string }).message}`
    }
    if (typeof error === 'string') return `HTTP ${status}: ${error}`
    if (typeof (data as { message?: unknown }).message === 'string') {
      return `HTTP ${status}: ${(data as { message: string }).message}`
    }
  }
  return `HTTP ${status}: ${text.slice(0, 500) || 'empty response'}`
}

function parseJson(text: string) {
  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}

function redactPayload(value: unknown): unknown {
  if (typeof value === 'string') {
    if (value.startsWith('data:image/')) return `${value.slice(0, 38)}...[${Math.round(value.length / 1024)}KB]`
    if (value.length > 600) return `${value.slice(0, 600)}...`
    return value
  }

  if (Array.isArray(value)) {
    return value.map((item) => redactPayload(item))
  }

  if (value && typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>).map(([key, item]) => [key, redactPayload(item)])
    )
  }

  return value
}

async function handleImageFile(event: Event, slot: FileSlot) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const dataUrl = await readFileAsDataUrl(file)
  if (slot === 'source') {
    sourceImageDataUrl.value = dataUrl
    sourceImageName.value = `${file.name} · ${formatBytes(file.size)}`
  } else {
    maskImageDataUrl.value = dataUrl
    maskImageName.value = `${file.name} · ${formatBytes(file.size)}`
  }
  input.value = ''
}

function readFileAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(reader.error || new Error('读取图片失败'))
    reader.readAsDataURL(file)
  })
}

function toggleBatchSize(size: string) {
  if (selectedBatchSizes.value.includes(size)) {
    selectedBatchSizes.value = selectedBatchSizes.value.filter((item) => item !== size)
  } else {
    selectedBatchSizes.value = [...selectedBatchSizes.value, size]
  }
}

function clampNumber(value: number, min: number, max: number) {
  if (Number.isNaN(value)) return min
  return Math.min(max, Math.max(min, Math.round(value)))
}

async function copyText(text: string) {
  await navigator.clipboard?.writeText(text)
}

function openImage(src: string) {
  window.open(src, '_blank', 'noopener,noreferrer')
}

function downloadImage(src: string, filename: string) {
  const link = document.createElement('a')
  link.href = src
  link.download = filename
  link.rel = 'noopener noreferrer'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes}B`
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)}KB`
  return `${(bytes / 1024 / 1024).toFixed(1)}MB`
}
</script>
