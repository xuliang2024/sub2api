<script setup lang="ts">
import { computed, onMounted, provide, ref, watch } from 'vue'
import type {
  CampaignAppConfig,
  CampaignHero,
  CampaignKeyGrant,
  CampaignUser,
  CampaignView,
} from './types'
import { claimKey, describeError, fetchCampaign } from './api'
import { clearSession, getStoredUser, getToken } from './session'
import { primaryModelName } from './format'
import {
  CampaignMessagesKey,
  SUPPORTED_LOCALES,
  createMessages,
  isSupportedLocale,
  normalizeCampaignLocale,
  type CampaignLocale,
  type Translate,
} from './messages'
import SiteNav from './components/SiteNav.vue'
import HeroSection from './components/HeroSection.vue'
import ModelGrid from './components/ModelGrid.vue'
import ModelCatalog from './components/ModelCatalog.vue'
import ConsoleView from './console/ConsoleView.vue'
import UsageDoc from './components/UsageDoc.vue'
import CodeSampleTabs from './components/CodeSampleTabs.vue'
import AuthModal from './components/AuthModal.vue'
import KeyPanel from './components/KeyPanel.vue'

/** 由后端注入到 campaign.html 的站点公开配置（注册开关等）。 */
const appConfig: CampaignAppConfig =
  (window as unknown as { __APP_CONFIG__?: CampaignAppConfig }).__APP_CONFIG__ ?? {}

/** ?c=<slug> 允许同域名预览不同活动；线上正式投放不传，由 Host 解析。 */
const slugOverride = new URLSearchParams(window.location.search).get('c')?.trim() || ''

const campaign = ref<CampaignView | null>(null)
const loading = ref(true)
const loadError = ref('')

const LOCALE_STORAGE_KEY = 'campaign_locale'

/**
 * 显式选择的语言。
 *
 * 优先级：?lang=xx（投放可给不同语言各做一个广告链接）> localStorage（用户上次选择）
 * > 活动配置的 locale（默认 en）。都没有时由 campaign.locale 决定。
 */
function readExplicitLocale(): CampaignLocale | null {
  const fromQuery = new URLSearchParams(window.location.search).get('lang')?.trim().toLowerCase()
  if (fromQuery && isSupportedLocale(fromQuery)) {
    return fromQuery
  }
  try {
    const stored = localStorage.getItem(LOCALE_STORAGE_KEY)
    if (stored && isSupportedLocale(stored)) {
      return stored
    }
  } catch {
    // 隐私模式下读不到 localStorage，用活动默认值即可。
  }
  return null
}

type CampaignTab = 'overview' | 'catalog' | 'console'

const CAMPAIGN_TABS: CampaignTab[] = ['overview', 'catalog', 'console']

/** 页签来自 URL（?tab=catalog / ?tab=console），方便单独投放某个视图。 */
function readInitialTab(): CampaignTab {
  const tab = new URLSearchParams(window.location.search).get('tab')?.trim().toLowerCase()
  return CAMPAIGN_TABS.includes(tab as CampaignTab) ? (tab as CampaignTab) : 'overview'
}

const activeTab = ref<CampaignTab>(readInitialTab())

const explicitLocale = ref<CampaignLocale | null>(readExplicitLocale())
const activeLocale = computed<CampaignLocale>(
  () => explicitLocale.value ?? normalizeCampaignLocale(campaign.value?.locale),
)

const translate = computed(() => createMessages(activeLocale.value))
provide(CampaignMessagesKey, (key: string, vars?: Record<string, string | number>) =>
  translate.value(key, vars),
)
// 注意：同一组件里 inject 不到自己 provide 的值，根组件必须直接用 translate。
const t: Translate = (key, vars) => translate.value(key, vars)

function setLocale(next: CampaignLocale): void {
  explicitLocale.value = next
  try {
    localStorage.setItem(LOCALE_STORAGE_KEY, next)
  } catch {
    // 不阻断：仅本次会话生效。
  }
  // 同步到 URL，方便用户把某个语言的页面直接复制分享。
  const url = new URL(window.location.href)
  url.searchParams.set('lang', next)
  window.history.replaceState({}, '', url.toString())
}

/** 当前语言的正文覆盖（标题/副标题/使用说明），没有则回落主文案。 */
const localeContent = computed(() => campaign.value?.translations?.[activeLocale.value])

const heroView = computed<CampaignHero>(() => {
  const base = campaign.value?.hero ?? {}
  const override = localeContent.value
  return {
    ...base,
    title: override?.hero_title || base.title || campaign.value?.name || '',
    subtitle: override?.hero_subtitle || base.subtitle,
  }
})

const contentMarkdown = computed(
  () => localeContent.value?.content_md || campaign.value?.content_md || '',
)

const user = ref<CampaignUser | null>(getStoredUser())

// localStorage 不是响应式的：登录/登出时自增这个计数器，让依赖它的 computed 重算。
const sessionTick = ref(0)
const authed = computed(() => {
  sessionTick.value
  return Boolean(getToken())
})

const authOpen = ref(false)
const authInitialMode = ref<'login' | 'register'>('register')
const authTrialNote = ref('')

const keyPanelOpen = ref(false)
const keyGrant = ref<CampaignKeyGrant | undefined>()
const keyLoading = ref(false)
const keyError = ref('')

const ctaLabel = computed(
  () => localeContent.value?.cta_label || campaign.value?.cta.label || t('cta.defaultLabel'),
)

const trialNote = computed(() => {
  const trial = campaign.value?.trial
  if (!trial?.balance || trial.balance <= 0) {
    return trial?.note || ''
  }
  const amount = `$${trial.balance.toFixed(2)}`
  return trial.note || amount
})

/**
 * 顶部页签。
 *
 * 概览放主推与说明，模型广场是独立一页（全量价目表）——分组没有模型时不出现第二个页签。
 */
const tabs = computed(() => {
  const items: { key: CampaignTab; label: string }[] = [
    { key: 'overview', label: t('tabs.overview') },
  ]
  if ((campaign.value?.catalog?.length ?? 0) > 0) {
    items.push({ key: 'catalog', label: t('tabs.catalog') })
  }
  // 控制台需要登录态：未登录时前端还没法区分「没有 Key」和「没权限」
  if (authed.value) {
    items.push({ key: 'console', label: t('tabs.console') })
  }
  return items
})

function setTab(tab: string): void {
  activeTab.value = CAMPAIGN_TABS.includes(tab as CampaignTab) ? (tab as CampaignTab) : 'overview'
  const url = new URL(window.location.href)
  url.searchParams.set('tab', activeTab.value)
  window.history.replaceState({}, '', url.toString())
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 活动没配模型时 URL 里的 ?tab=catalog 不应停在空页面
watch(tabs, (items) => {
  if (!items.some((item) => item.key === activeTab.value)) {
    activeTab.value = 'overview'
  }
})

const sampleContext = computed(() => ({
  baseUrl: campaign.value?.api_base_url || window.location.origin,
  apiKey: keyGrant.value?.api_key,
  model: campaign.value ? primaryModelName(campaign.value.models) : '',
}))

onMounted(loadCampaign)

// 语言变化时同步 <html lang> 与页面标题：切到日语后，浏览器标签页与无障碍朗读
// 都应该跟着变，否则分享出去的链接在别人机器上读屏是错的。
watch(
  activeLocale,
  (locale) => {
    const option = SUPPORTED_LOCALES.find((item) => item.code === locale)
    if (option) {
      document.documentElement.lang = option.htmlLang
    }
    applyDocumentMeta()
  },
  { immediate: true },
)

/**
 * 应用标题/描述等文档元信息。
 *
 * 先按活动配置写一遍（保证 description / og:image / keywords 一定存在），
 * 再用当前语言的 hero 标题覆盖 title——否则切到日语，浏览器标签页还是英文。
 */
function applyDocumentMeta(): void {
  const view = campaign.value
  if (!view) {
    return
  }
  applySEO(view)
  const overrideTitle = localeContent.value?.hero_title
  if (overrideTitle) {
    document.title = overrideTitle
    setMeta('property', 'og:title', overrideTitle)
  }
}

async function loadCampaign(): Promise<void> {
  loading.value = true
  loadError.value = ''
  try {
    const data = await fetchCampaign(slugOverride)
    campaign.value = data
    applyTheme(data)
    applyDocumentMeta()
  } catch (error) {
    loadError.value = describeError(error, t('error.loadFailed'))
  } finally {
    loading.value = false
  }
}

function applyTheme(view: CampaignView): void {
  const root = document.documentElement
  const primary = view.brand?.primary?.trim()
  const accent = view.brand?.accent?.trim()
  if (primary) {
    root.style.setProperty('--c-primary', primary)
    root.style.setProperty('--c-primary-soft', hexToRgba(primary, 0.1))
  }
  if (accent) {
    root.style.setProperty('--c-accent', accent)
  }
}

function applySEO(view: CampaignView): void {
  const title = view.seo?.title || view.hero?.title || view.name
  if (title) {
    document.title = title
  }
  const description = view.seo?.description || view.hero?.subtitle
  if (description) {
    setMeta('name', 'description', description)
    setMeta('property', 'og:description', description)
  }
  setMeta('property', 'og:title', document.title)
  if (view.seo?.og_image) {
    setMeta('property', 'og:image', view.seo.og_image)
  }
  if (view.seo?.keywords?.length) {
    setMeta('name', 'keywords', view.seo.keywords.join(','))
  }
}

function setMeta(attr: 'name' | 'property', key: string, content: string): void {
  let tag = document.head.querySelector<HTMLMetaElement>(`meta[${attr}="${key}"]`)
  if (!tag) {
    tag = document.createElement('meta')
    tag.setAttribute(attr, key)
    document.head.appendChild(tag)
  }
  tag.setAttribute('content', content)
}

function hexToRgba(hex: string, alpha: number): string {
  const normalized = hex.replace('#', '')
  if (normalized.length !== 6) {
    return hex
  }
  const value = Number.parseInt(normalized, 16)
  const r = (value >> 16) & 255
  const g = (value >> 8) & 255
  const b = value & 255
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

function openAuth(mode: 'login' | 'register'): void {
  authInitialMode.value = mode
  authTrialNote.value = mode === 'register' ? trialNote.value : ''
  authOpen.value = true
}

/** 主 CTA：未登录先注册/登录，已登录直接领 Key。 */
async function handleCTA(): Promise<void> {
  const view = campaign.value
  if (!view) {
    return
  }
  if (view.cta.mode === 'topup') {
    if (view.cta.topup_url) {
      window.location.href = view.cta.topup_url
    } else {
      openAuth('login')
    }
    return
  }
  if (!authed.value) {
    openAuth('register')
    return
  }
  if (!view.group_id) {
    // 未绑定分组：只做展示，没有可发的 Key
    if (view.cta.topup_url) {
      window.location.href = view.cta.topup_url
      return
    }
    openAuth('login')
    return
  }
  await requestKey()
}

async function requestKey(): Promise<void> {
  const view = campaign.value
  if (!view) {
    return
  }
  keyPanelOpen.value = true
  keyLoading.value = true
  keyError.value = ''
  try {
    keyGrant.value = await claimKey(view.slug)
  } catch (error) {
    keyError.value = describeError(error, t('error.claimFailed'))
  } finally {
    keyLoading.value = false
  }
}

function handleAuthSuccess(payload: { user?: CampaignUser; isNew: boolean }): void {
  user.value = payload.user ?? getStoredUser()
  sessionTick.value += 1
  authOpen.value = false
  void requestKey()
}

/**
 * 退出登录。
 *
 * 只清当前 origin 的会话：token/user 都在 localStorage，清掉后 authed 立刻变 false，
 * 控制台页签随之消失；如果正停在控制台，退回概览，避免停在空白页。
 */
function handleLogout(): void {
  clearSession()
  user.value = null
  keyGrant.value = undefined
  keyPanelOpen.value = false
  sessionTick.value += 1
  if (activeTab.value === 'console') {
    setTab('overview')
  }
}
</script>

<template>
  <div class="campaign-shell">
    <template v-if="loading">
      <div class="campaign-container" style="padding-top: 120px; padding-bottom: 120px">
        <div class="c-skeleton" style="width: 180px"></div>
        <div class="c-skeleton mt-6" style="width: 70%; height: 36px"></div>
        <div class="c-skeleton mt-4" style="width: 52%"></div>
        <div class="c-skeleton mt-10" style="width: 100%; height: 180px"></div>
      </div>
    </template>

    <template v-else-if="loadError">
      <div class="campaign-container" style="padding-top: 120px; padding-bottom: 120px">
        <h1 class="c-hero-title" style="font-size: 26px">{{ t('error.title') }}</h1>
        <p class="c-hero-sub">{{ loadError }}</p>
        <button class="c-btn c-btn-ghost mt-6" type="button" @click="loadCampaign">{{ t('error.reload') }}</button>
      </div>
    </template>

    <template v-else-if="campaign">
      <SiteNav
        :brand="campaign.brand"
        :site-name="campaign.brand.site_name || campaign.name"
        :authed="authed"
        :user-label="user?.email || user?.username"
        :console-url="campaign.brand.console_url"
        :cta-label="ctaLabel"
        :locale="activeLocale"
        :tabs="tabs"
        :active-tab="activeTab"
        @login="openAuth('login')"
        @logout="handleLogout"
        @cta="handleCTA"
        @update:locale="setLocale"
        @update:tab="setTab"
      />

      <template v-if="activeTab === 'overview'">
        <HeroSection
          :hero="heroView"
          :cta-label="ctaLabel"
          :trial-note="trialNote"
          @cta="handleCTA"
        />

        <ModelGrid
          :models="campaign.models"
          :group="campaign.group"
          :group-label="campaign.brand.group_label"
          :catalog-count="campaign.catalog?.length ?? 0"
          @open-catalog="setTab('catalog')"
        />

        <UsageDoc :markdown="contentMarkdown" />

        <CodeSampleTabs :samples="campaign.code_samples" :context="sampleContext" />
      </template>

      <ModelCatalog v-else-if="activeTab === 'catalog'" :models="campaign.catalog ?? []" />

      <ConsoleView v-else />

      <!-- 控制台里不再出现拉新 CTA：用户已经登录并在这里管 Key 了 -->
      <section v-if="activeTab !== 'console'" class="c-section c-cta">
        <div class="campaign-container">
          <div class="c-cta-card">
            <div>
              <h2>{{ campaign.cta.label || t('cta.defaultLabel') }}</h2>
              <p v-if="campaign.cta.sublabel">{{ campaign.cta.sublabel }}</p>
            </div>
            <button class="c-btn" type="button" @click="handleCTA">{{ ctaLabel }}</button>
          </div>
        </div>
      </section>

      <footer class="c-footer">
        <div class="campaign-container">
          {{ campaign.brand.footer || campaign.brand.site_name || campaign.name }}
        </div>
      </footer>

      <AuthModal
        :open="authOpen"
        :config="appConfig"
        :campaign-slug="campaign.slug"
        :initial-mode="authInitialMode"
        :trial-note="authTrialNote"
        @close="authOpen = false"
        @success="handleAuthSuccess"
      />

      <KeyPanel
        :open="keyPanelOpen"
        :grant="keyGrant"
        :model="sampleContext.model"
        :topup-url="campaign.cta.topup_url"
        :loading="keyLoading"
        :error="keyError"
        @close="keyPanelOpen = false"
        @retry="requestKey"
      />
    </template>
  </div>
</template>
