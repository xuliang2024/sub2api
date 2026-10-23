<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6">
      <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 class="text-xl font-semibold text-gray-900 dark:text-white">广告投放落地页</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
            一个活动一套落地页，绑定域名后访问该域名即直接展示。绑定分组决定可用模型与价格。
          </p>
        </div>
        <div class="flex gap-2">
          <button class="btn btn-secondary btn-md" :disabled="loading" @click="loadAll">刷新</button>
          <button class="btn btn-primary btn-md" @click="openCreate">新建活动</button>
        </div>
      </div>

      <p
        v-if="banner"
        class="mb-4 rounded-lg border border-amber-300 bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:border-amber-700 dark:bg-amber-900/30 dark:text-amber-200"
      >
        {{ banner }}
      </p>

      <div
        class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
      >
        <table class="w-full text-left text-sm">
          <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-900 dark:text-dark-400">
            <tr>
              <th class="px-4 py-3 font-medium">活动</th>
              <th class="px-4 py-3 font-medium">域名</th>
              <th class="px-4 py-3 font-medium">状态</th>
              <th class="px-4 py-3 font-medium">分组</th>
              <th class="px-4 py-3 font-medium">更新时间</th>
              <th class="px-4 py-3 text-right font-medium">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
            <tr v-if="loading">
              <td colspan="6" class="px-4 py-10 text-center text-gray-500">加载中…</td>
            </tr>
            <tr v-else-if="!items.length">
              <td colspan="6" class="px-4 py-10 text-center text-gray-500">
                还没有活动。点「新建活动」开始配置第一个投放页。
              </td>
            </tr>
            <tr v-for="item in items" :key="item.id" class="align-top">
              <td class="px-4 py-3">
                <div class="font-medium text-gray-900 dark:text-white">{{ item.name }}</div>
                <div class="text-xs text-gray-500">{{ item.slug }}</div>
              </td>
              <td class="px-4 py-3">
                <span v-if="item.host" class="text-gray-700 dark:text-dark-200">{{ item.host }}</span>
                <span v-else class="text-gray-400">未绑定</span>
              </td>
              <td class="px-4 py-3">
                <span :class="statusClass(item.status)">{{ statusLabel(item.status) }}</span>
              </td>
              <td class="px-4 py-3 text-gray-700 dark:text-dark-200">
                {{ groupLabel(item.group_id) }}
              </td>
              <td class="px-4 py-3 text-gray-500">{{ formatTime(item.updated_at) }}</td>
              <td class="px-4 py-3">
                <div class="flex justify-end gap-2">
                  <a
                    class="btn btn-ghost btn-sm"
                    :href="`/campaign.html?c=${encodeURIComponent(item.slug)}`"
                    target="_blank"
                    rel="noopener"
                  >
                    预览
                  </a>
                  <button class="btn btn-secondary btn-sm" @click="openEdit(item)">编辑</button>
                  <button class="btn btn-danger btn-sm" @click="handleDelete(item)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="editorOpen" class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/50 p-3 sm:p-6">
      <div class="w-full max-w-3xl rounded-xl bg-white shadow-xl dark:bg-dark-800">
        <div class="flex items-center justify-between border-b border-gray-200 px-5 py-4 dark:border-dark-700">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">
            {{ form.id ? '编辑活动' : '新建活动' }}
          </h2>
          <button class="text-sm text-gray-500 hover:text-gray-800" @click="editorOpen = false">
            关闭
          </button>
        </div>

        <div class="flex flex-wrap gap-1 border-b border-gray-200 px-4 pt-3 dark:border-dark-700">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            class="rounded-t-lg px-3 py-2 text-sm"
            :class="
              activeTab === tab.key
                ? 'bg-gray-100 font-medium text-gray-900 dark:bg-dark-700 dark:text-white'
                : 'text-gray-500 hover:text-gray-800 dark:hover:text-dark-200'
            "
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </button>
        </div>

        <div class="max-h-[70vh] overflow-y-auto px-5 py-4">
          <!-- 基础 -->
          <div v-show="activeTab === 'basic'" class="space-y-4">
            <div class="grid gap-4 sm:grid-cols-2">
              <label class="block">
                <span class="input-label">活动名称</span>
                <input v-model="form.name" class="input" placeholder="例如：DeepSeek V4 投放" />
              </label>
              <label class="block">
                <span class="input-label">slug（唯一，小写字母/数字/中划线）</span>
                <input v-model="form.slug" class="input" placeholder="deepseek-v4" />
              </label>
              <label class="block">
                <span class="input-label">绑定域名（不含 https:// 与端口）</span>
                <input v-model="form.host" class="input" placeholder="try.example.com" />
              </label>
              <label class="block">
                <span class="input-label">状态</span>
                <select v-model="form.status" class="input">
                  <option value="draft">草稿（仅预览）</option>
                  <option value="active">已上线</option>
                  <option value="archived">已下线</option>
                </select>
              </label>
              <label class="block">
                <span class="input-label">界面语言</span>
                <select v-model="form.locale" class="input">
                  <option value="en">English（海外投放，默认）</option>
                  <option value="zh">简体中文</option>
                </select>
              </label>
              <label class="block sm:col-span-2">
                <span class="input-label">隐藏模型（每行一个，别名/退役模型不展示）</span>
                <textarea
                  v-model="hiddenModelsText"
                  class="input min-h-20 font-mono text-xs"
                  placeholder="deepseek-flash&#10;deepseek-v4-flash"
                />
                <span class="mt-1 block text-xs text-gray-500">
                  渠道定价里的历史别名会计入模型清单，填在这里可在卡片与模型广场中隐藏。
                </span>
              </label>
              <label class="block sm:col-span-2">
                <span class="input-label">绑定分组（决定模型清单与价格）</span>
                <select v-model="groupIdProxy" class="input">
                  <option :value="0">不绑定（纯展示，不发放 Key）</option>
                  <option v-for="group in groups" :key="group.id" :value="group.id">
                    {{ group.name }}（{{ group.platform }} · {{ group.rate_multiplier }}x）
                  </option>
                </select>
              </label>
            </div>

            <div class="grid gap-4 sm:grid-cols-3">
              <label class="block">
                <span class="input-label">站点名</span>
                <input v-model="form.brand.site_name" class="input" placeholder="留空用活动名" />
              </label>
              <label class="block">
                <span class="input-label">主色</span>
                <input v-model="form.brand.primary" class="input" placeholder="#4f46e5" />
              </label>
              <label class="block">
                <span class="input-label">强调色</span>
                <input v-model="form.brand.accent" class="input" placeholder="#06b6d4" />
              </label>
              <label class="block">
                <span class="input-label">Logo URL</span>
                <input v-model="form.brand.logo" class="input" placeholder="/logo.svg" />
              </label>
              <label class="block sm:col-span-2">
                <span class="input-label">控制台地址（导航「控制台」入口）</span>
                <input v-model="form.brand.console_url" class="input" placeholder="https://codex.apiz.ai/dashboard" />
              </label>
            </div>
          </div>

          <!-- 首屏 -->
          <div v-show="activeTab === 'hero'" class="space-y-4">
            <label class="block">
              <span class="input-label">小标题（eyebrow）</span>
              <input v-model="form.hero.eyebrow" class="input" placeholder="新模型上线" />
            </label>
            <label class="block">
              <span class="input-label">主标题</span>
              <input v-model="form.hero.title" class="input" placeholder="DeepSeek V4：更低成本的长上下文" />
            </label>
            <label class="block">
              <span class="input-label">副标题</span>
              <textarea v-model="form.hero.subtitle" class="input min-h-20" />
            </label>
            <label class="block">
              <span class="input-label">卖点（每行一条）</span>
              <textarea v-model="highlightsText" class="input min-h-24" />
            </label>
            <label class="block">
              <span class="input-label">标签（每行一条，显示为圆角 chip）</span>
              <textarea v-model="badgesText" class="input min-h-20" />
            </label>
            <label class="block">
              <span class="input-label">主视觉图 URL</span>
              <input v-model="form.hero.image" class="input" />
            </label>
          </div>

          <!-- 使用说明 -->
          <div v-show="activeTab === 'content'" class="space-y-3">
            <p class="text-xs text-gray-500">
              支持 Markdown（标题、列表、表格、代码块）。落地页会用 marked + DOMPurify 渲染。
            </p>
            <textarea v-model="form.content_md" class="input min-h-[46vh] font-mono text-xs" />
          </div>

          <!-- 翻译 -->
          <div v-show="activeTab === 'translations'" class="space-y-4">
            <p class="text-xs text-gray-500">
              界面上的按钮、表头、提示语已内置英/中/日/韩/泰五套；这里覆盖标题、副标题与使用说明。
              留空表示沿用主文案。
            </p>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="option in TRANSLATION_LOCALES"
                :key="option.code"
                class="rounded-lg border px-3 py-1.5 text-sm"
                :class="
                  translationLocale === option.code
                    ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-200'
                    : 'border-gray-200 text-gray-600 dark:border-dark-700 dark:text-dark-300'
                "
                @click="translationLocale = option.code"
              >
                {{ option.label }}
                <span v-if="hasTranslation(option.code)" class="ml-1 text-xs text-green-600">●</span>
              </button>
            </div>

            <label class="block">
              <span class="input-label">标题（{{ localeLabel(translationLocale) }}）</span>
              <input
                v-model="translationTitle"
                class="input"
                :placeholder="form.hero.title || '沿用主文案'"
              />
            </label>
            <label class="block">
              <span class="input-label">副标题</span>
              <textarea
                v-model="translationSubtitle"
                class="input min-h-20"
                :placeholder="form.hero.subtitle || '沿用主文案'"
              />
            </label>
            <label class="block">
              <span class="input-label">主按钮文案</span>
              <input
                v-model="translationCta"
                class="input"
                :placeholder="form.cta.label || '沿用主文案'"
              />
            </label>
            <label class="block">
              <span class="input-label">使用说明（Markdown）</span>
              <textarea
                v-model="translationContent"
                class="input min-h-[34vh] font-mono text-xs"
                placeholder="沿用主文案"
              />
            </label>
            <button class="btn btn-secondary btn-sm" @click="clearTranslation">
              清空该语言（回落主文案）
            </button>
          </div>

          <!-- 模型 -->
          <div v-show="activeTab === 'models'" class="space-y-3">
            <p class="text-xs text-gray-500">
              留空表示展示绑定分组下的全部模型。填写后只展示这些模型，并可覆盖展示名与卖点；
              分组里没有的模型会标注「待开通」。
            </p>
            <div
              v-for="(entry, index) in form.model_entries"
              :key="index"
              class="grid gap-2 rounded-lg border border-gray-200 p-3 sm:grid-cols-12 dark:border-dark-700"
            >
              <input v-model="entry.model" class="input sm:col-span-4" placeholder="模型名 deepseek-v4-pro" />
              <input v-model="entry.display_name" class="input sm:col-span-3" placeholder="展示名" />
              <input v-model="entry.badge" class="input sm:col-span-2" placeholder="角标" />
              <input v-model.number="entry.order" type="number" class="input sm:col-span-1" placeholder="排序" />
              <label class="flex items-center gap-2 text-xs sm:col-span-2">
                <input v-model="entry.featured" type="checkbox" />
                主推
              </label>
              <input v-model="entry.tagline" class="input sm:col-span-11" placeholder="一句话卖点" />
              <button class="btn btn-ghost btn-sm sm:col-span-1" @click="form.model_entries.splice(index, 1)">
                删除
              </button>
            </div>
            <button class="btn btn-secondary btn-sm" @click="addModelEntry">+ 添加模型</button>
          </div>

          <!-- 代码示例 -->
          <div v-show="activeTab === 'samples'" class="space-y-3">
            <p class="text-xs text-gray-500">
              代码里的
              <code v-for="placeholder in placeholders" :key="placeholder" class="mr-1">{{ placeholder }}</code>
              会被自动替换为真实值（未登录时 API_KEY 显示占位符）。
            </p>
            <div
              v-for="(sample, index) in form.code_samples"
              :key="index"
              class="space-y-2 rounded-lg border border-gray-200 p-3 dark:border-dark-700"
            >
              <div class="grid gap-2 sm:grid-cols-3">
                <input v-model="sample.label" class="input" placeholder="标签 curl / Python" />
                <input v-model="sample.id" class="input" placeholder="id（唯一）" />
                <input v-model="sample.lang" class="input" placeholder="lang（可选）" />
              </div>
              <textarea v-model="sample.code" class="input min-h-32 font-mono text-xs" />
              <button class="btn btn-ghost btn-sm" @click="form.code_samples.splice(index, 1)">删除</button>
            </div>
            <button class="btn btn-secondary btn-sm" @click="addSample">+ 添加示例</button>
          </div>

          <!-- 转化 -->
          <div v-show="activeTab === 'cta'" class="space-y-4">
            <div class="grid gap-4 sm:grid-cols-2">
              <label class="block">
                <span class="input-label">主按钮文案</span>
                <input v-model="form.cta.label" class="input" placeholder="立即领取 API Key" />
              </label>
              <label class="block">
                <span class="input-label">按钮行为</span>
                <select v-model="form.cta.mode" class="input">
                  <option value="key">登录后一键领取 Key</option>
                  <option value="register">只引导注册</option>
                  <option value="topup">跳转充值</option>
                </select>
              </label>
              <label class="block sm:col-span-2">
                <span class="input-label">按钮副文案</span>
                <input v-model="form.cta.sublabel" class="input" />
              </label>
              <label class="block sm:col-span-2">
                <span class="input-label">充值页地址（可选）</span>
                <input v-model="form.cta.topup_url" class="input" placeholder="/purchase" />
              </label>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <label class="block">
                <span class="input-label">注册赠送额度（USD，0 = 不送）</span>
                <input v-model.number="form.trial.balance" type="number" step="0.01" class="input" />
              </label>
              <label class="block">
                <span class="input-label">赠送文案</span>
                <input v-model="form.trial.note" class="input" placeholder="新用户注册即送 $1 试用额度" />
              </label>
            </div>

            <p class="text-xs text-gray-500">
              注意：赠送额度会计入用户余额，需要站点「注册」开关处于开启状态才会生效。
            </p>
          </div>

          <!-- SEO -->
          <div v-show="activeTab === 'seo'" class="space-y-4">
            <label class="block">
              <span class="input-label">页面标题</span>
              <input v-model="form.seo.title" class="input" />
            </label>
            <label class="block">
              <span class="input-label">页面描述</span>
              <textarea v-model="form.seo.description" class="input min-h-20" />
            </label>
            <label class="block">
              <span class="input-label">分享图 URL</span>
              <input v-model="form.seo.og_image" class="input" />
            </label>
            <label class="block">
              <span class="input-label">关键词（每行一条）</span>
              <textarea v-model="keywordsText" class="input min-h-20" />
            </label>
          </div>

          <p v-if="saveError" class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700 dark:bg-red-900/30 dark:text-red-200">
            {{ saveError }}
          </p>
        </div>

        <div class="flex items-center justify-end gap-2 border-t border-gray-200 px-5 py-4 dark:border-dark-700">
          <button class="btn btn-secondary btn-md" @click="editorOpen = false">取消</button>
          <button class="btn btn-primary btn-md" :disabled="saving" @click="handleSave">
            {{ saving ? '保存中…' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import * as campaignsAPI from '@/api/admin/campaigns'
import type {
  AdminCampaign,
  CampaignPayload,
} from '@/api/admin/campaigns'
import * as groupsAPI from '@/api/admin/groups'

interface GroupOption {
  id: number
  name: string
  platform: string
  rate_multiplier: number
}

const tabs = [
  { key: 'basic', label: '基础' },
  { key: 'hero', label: '首屏' },
  { key: 'content', label: '使用说明' },
  { key: 'translations', label: '翻译' },
  { key: 'models', label: '模型' },
  { key: 'samples', label: '示例' },
  { key: 'cta', label: '转化' },
  { key: 'seo', label: 'SEO' },
] as const

/** 可覆盖正文的语言（en 是主文案，不需要覆盖）。 */
const TRANSLATION_LOCALES = [
  { code: 'zh', label: '中文' },
  { code: 'ja', label: '日本語' },
  { code: 'ko', label: '한국어' },
  { code: 'th', label: 'ไทย' },
]

const translationLocale = ref('ja')

function localeLabel(code: string): string {
  return TRANSLATION_LOCALES.find((item) => item.code === code)?.label ?? code
}

function hasTranslation(code: string): boolean {
  const entry = form.translations?.[code]
  return Boolean(
    entry && (entry.hero_title || entry.hero_subtitle || entry.content_md || entry.cta_label),
  )
}

type TranslationField = 'hero_title' | 'hero_subtitle' | 'content_md' | 'cta_label'

function setTranslationField(field: TranslationField, value: string): void {
  if (!form.translations) {
    form.translations = {}
  }
  const current = form.translations[translationLocale.value] ?? {}
  form.translations[translationLocale.value] = { ...current, [field]: value }
}

function currentTranslationField(field: TranslationField): string {
  return form.translations?.[translationLocale.value]?.[field] ?? ''
}

const translationTitle = computed({
  get: () => currentTranslationField('hero_title'),
  set: (value: string) => setTranslationField('hero_title', value),
})

const translationSubtitle = computed({
  get: () => currentTranslationField('hero_subtitle'),
  set: (value: string) => setTranslationField('hero_subtitle', value),
})

const translationContent = computed({
  get: () => currentTranslationField('content_md'),
  set: (value: string) => setTranslationField('content_md', value),
})

const translationCta = computed({
  get: () => currentTranslationField('cta_label'),
  set: (value: string) => setTranslationField('cta_label', value),
})

function clearTranslation(): void {
  if (!form.translations) {
    return
  }
  delete form.translations[translationLocale.value]
}

/** 示例代码支持的占位符，仅用于文案展示。 */
const placeholders = ['{{BASE_URL}}', '{{API_KEY}}', '{{MODEL}}']

const items = ref<AdminCampaign[]>([])
const groups = ref<GroupOption[]>([])
const loading = ref(false)
const saving = ref(false)
const banner = ref('')
const saveError = ref('')
const editorOpen = ref(false)
const activeTab = ref<(typeof tabs)[number]['key']>('basic')

function emptyPayload(): CampaignPayload {
  return {
    slug: '',
    name: '',
    status: 'draft',
    host: '',
    locale: 'en',
    group_id: null,
    brand: { primary: '#4f46e5', accent: '#06b6d4' },
    hero: { highlights: [], badges: [] },
    content_md: '',
    model_entries: [],
    code_samples: [],
    cta: { mode: 'key', label: '立即领取 API Key' },
    trial: { balance: 0 },
    seo: { keywords: [] },
    translations: {},
    hidden_models: [],
  }
}

const form = reactive<CampaignPayload & { id?: number }>(emptyPayload())

const groupIdProxy = computed<number>({
  get: () => form.group_id ?? 0,
  set: (value) => {
    form.group_id = value > 0 ? value : null
  },
})

const highlightsText = computed<string>({
  get: () => (form.hero.highlights ?? []).join('\n'),
  set: (value) => {
    form.hero.highlights = splitLines(value)
  },
})

const badgesText = computed<string>({
  get: () => (form.hero.badges ?? []).join('\n'),
  set: (value) => {
    form.hero.badges = splitLines(value)
  },
})

const hiddenModelsText = computed<string>({
  get: () => (form.hidden_models ?? []).join('\n'),
  set: (value) => {
    form.hidden_models = splitLines(value)
  },
})

const keywordsText = computed<string>({
  get: () => (form.seo.keywords ?? []).join('\n'),
  set: (value) => {
    form.seo.keywords = splitLines(value)
  },
})

function splitLines(value: string): string[] {
  return value
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
}

function statusLabel(status: string): string {
  if (status === 'active') return '已上线'
  if (status === 'archived') return '已下线'
  return '草稿'
}

function statusClass(status: string): string {
  const base = 'inline-flex rounded-full px-2 py-0.5 text-xs font-medium '
  if (status === 'active') return base + 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-200'
  if (status === 'archived') return base + 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'
  return base + 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-200'
}

function groupLabel(groupId: number | null): string {
  if (!groupId) return '未绑定'
  const group = groups.value.find((item) => item.id === groupId)
  return group ? group.name : `#${groupId}`
}

function formatTime(value: string): string {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

async function loadAll(): Promise<void> {
  loading.value = true
  banner.value = ''
  try {
    const [campaigns, groupPage] = await Promise.all([
      campaignsAPI.list(),
      groupsAPI.list(1, 200).catch(() => null),
    ])
    items.value = campaigns
    groups.value = (groupPage?.items ?? []) as unknown as GroupOption[]
  } catch (error) {
    banner.value = describe(error)
  } finally {
    loading.value = false
  }
}

function describe(error: unknown): string {
  if (typeof error === 'object' && error && 'message' in error) {
    return String((error as { message?: string }).message)
  }
  return '操作失败，请稍后重试'
}

function openCreate(): void {
  Object.assign(form, emptyPayload())
  form.id = undefined
  activeTab.value = 'basic'
  saveError.value = ''
  editorOpen.value = true
}

function openEdit(item: AdminCampaign): void {
  Object.assign(form, {
    ...emptyPayload(),
    ...item,
    brand: { ...item.brand },
    hero: { ...item.hero },
    cta: { ...item.cta },
    trial: { ...item.trial },
    seo: { ...item.seo },
    translations: { ...(item.translations ?? {}) },
    hidden_models: [...(item.hidden_models ?? [])],
    model_entries: item.model_entries.map((entry) => ({ ...entry })),
    code_samples: item.code_samples.map((sample) => ({ ...sample })),
  })
  activeTab.value = 'basic'
  saveError.value = ''
  editorOpen.value = true
}

function addModelEntry(): void {
  form.model_entries.push({ model: '', featured: false, order: form.model_entries.length })
}

function addSample(): void {
  form.code_samples.push({ id: `sample-${form.code_samples.length + 1}`, label: 'curl', code: '' })
}

async function handleSave(): Promise<void> {
  if (saving.value) return
  saveError.value = ''
  if (!form.name.trim() || !form.slug.trim()) {
    saveError.value = '活动名称与 slug 不能为空'
    return
  }
  saving.value = true
  try {
    const payload: CampaignPayload = {
      slug: form.slug.trim(),
      name: form.name.trim(),
      status: form.status,
      host: form.host.trim(),
      locale: form.locale,
      group_id: form.group_id,
      brand: form.brand,
      hero: form.hero,
      content_md: form.content_md,
      model_entries: form.model_entries.filter((entry) => entry.model.trim()),
      code_samples: form.code_samples.filter((sample) => sample.id.trim() && sample.code.trim()),
      cta: form.cta,
      trial: form.trial,
      seo: form.seo,
      translations: form.translations ?? {},
      hidden_models: form.hidden_models ?? [],
    }
    if (form.id) {
      await campaignsAPI.update(form.id, payload)
    } else {
      await campaignsAPI.create(payload)
    }
    editorOpen.value = false
    await loadAll()
  } catch (error) {
    saveError.value = describe(error)
  } finally {
    saving.value = false
  }
}

async function handleDelete(item: AdminCampaign): Promise<void> {
  if (!window.confirm(`确定删除活动「${item.name}」？该活动已发放的 Key 记录会一并清理。`)) {
    return
  }
  try {
    await campaignsAPI.remove(item.id)
    await loadAll()
  } catch (error) {
    banner.value = describe(error)
  }
}

onMounted(loadAll)
</script>
