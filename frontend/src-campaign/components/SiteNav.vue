<script setup lang="ts">
import type { CampaignBrand } from '../types'
import { useCampaignMessages, type CampaignLocale } from '../messages'
import LocaleSwitcher from './LocaleSwitcher.vue'

defineProps<{
  brand: CampaignBrand
  siteName: string
  authed: boolean
  userLabel?: string
  consoleUrl?: string
  ctaLabel: string
  locale: CampaignLocale
  /** 顶部页签（概览 / 模型广场…） */
  tabs: { key: string; label: string }[]
  activeTab: string
}>()

const emit = defineEmits<{
  (event: 'login'): void
  (event: 'cta'): void
  (event: 'update:locale', value: CampaignLocale): void
  (event: 'update:tab', value: string): void
  (event: 'logout'): void
}>()

const t = useCampaignMessages()
</script>

<template>
  <header class="c-nav">
    <div class="campaign-container c-nav-inner">
      <a class="c-brand" href="#top">
        <img v-if="brand.logo" :src="brand.logo" :alt="siteName" />
        <span>{{ siteName }}</span>
      </a>

      <nav v-if="tabs.length > 1" class="c-nav-tabs" role="tablist">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="c-nav-tab"
          :class="{ 'is-active': tab.key === activeTab }"
          type="button"
          role="tab"
          :aria-selected="tab.key === activeTab"
          @click="emit('update:tab', tab.key)"
        >
          {{ tab.label }}
        </button>
      </nav>

      <div class="flex flex-wrap items-center justify-end gap-2 sm:gap-3">
        <LocaleSwitcher
          :model-value="locale"
          @update:model-value="emit('update:locale', $event)"
        />
        <template v-if="authed">
          <span v-if="userLabel" class="hidden text-sm sm:inline" style="color: var(--c-text-soft)">
            {{ userLabel }}
          </span>
          <a v-if="consoleUrl" class="c-btn c-btn-ghost c-btn-sm" :href="consoleUrl">{{ t('nav.console') }}</a>
          <button v-else class="c-btn c-btn-ghost c-btn-sm" type="button" @click="emit('cta')">
            {{ ctaLabel }}
          </button>
          <button
            class="c-btn c-btn-ghost c-btn-sm c-nav-logout"
            type="button"
            @click="emit('logout')"
          >
            {{ t('nav.logout') }}
          </button>
        </template>
        <template v-else>
          <button
            v-if="!brand.hide_auth"
            class="c-btn c-btn-ghost c-btn-sm"
            type="button"
            @click="emit('login')"
          >
            {{ t('nav.login') }}
          </button>
          <button class="c-btn c-btn-primary c-btn-sm" type="button" @click="emit('cta')">
            {{ ctaLabel }}
          </button>
        </template>
      </div>
    </div>
  </header>
</template>
