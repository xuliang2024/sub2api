<!--
  落地页的独立法务页（隐私政策 / 服务条款）。

  落地页没有路由，这里按 pathname 单独挂载：/privacy 与 /terms。
  文案与主站共用 src/shared/legal.ts，避免两处表述不一致。
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import { legalDocuments, normalizeLegalLocale, type LegalLocale } from '../src/shared/legal'

type Doc = 'privacy' | 'terms'

function resolveDoc(): Doc {
  return window.location.pathname.replace(/\/+$/, '').endsWith('/terms') ? 'terms' : 'privacy'
}

function resolveLocale(): LegalLocale {
  const fromQuery = new URLSearchParams(window.location.search).get('lang')
  const stored = window.localStorage.getItem('campaign_locale')
  return normalizeLegalLocale(fromQuery || stored || navigator.language)
}

const doc = ref<Doc>(resolveDoc())
const locale = ref<LegalLocale>(resolveLocale())

const content = computed(() => legalDocuments[doc.value][locale.value])
const isZh = computed(() => locale.value === 'zh')

function switchLocale(next: LegalLocale): void {
  if (next === locale.value) {
    return
  }
  locale.value = next
  const url = new URL(window.location.href)
  url.searchParams.set('lang', next)
  window.history.replaceState(null, '', url.toString())
  window.localStorage.setItem('campaign_locale', next)
  document.documentElement.lang = next
}

document.title = content.value.title
</script>

<template>
  <div class="c-legal">
    <header class="c-legal-head">
      <a class="c-legal-brand" href="/">apiz.ai</a>
      <nav class="c-legal-nav">
        <a href="/privacy" :class="{ 'is-active': doc === 'privacy' }">
          {{ isZh ? '隐私政策' : 'Privacy' }}
        </a>
        <a href="/terms" :class="{ 'is-active': doc === 'terms' }">
          {{ isZh ? '服务条款' : 'Terms' }}
        </a>
        <span class="c-legal-lang">
          <button type="button" :class="{ 'is-active': locale === 'en' }" @click="switchLocale('en')">
            EN
          </button>
          <button type="button" :class="{ 'is-active': locale === 'zh' }" @click="switchLocale('zh')">
            中文
          </button>
        </span>
      </nav>
    </header>

    <main class="c-legal-body">
      <h1 class="c-legal-title">{{ content.title }}</h1>
      <p class="c-legal-updated">
        {{ isZh ? '最后更新：' : 'Last updated: ' }}{{ content.updated }}
      </p>
      <p class="c-legal-intro">{{ content.intro }}</p>

      <section v-for="section in content.sections" :key="section.heading" class="c-legal-section">
        <h2>{{ section.heading }}</h2>
        <p v-for="(para, index) in section.paras ?? []" :key="`p-${index}`">{{ para }}</p>
        <ul v-if="section.bullets?.length">
          <li v-for="(item, index) in section.bullets" :key="`b-${index}`">{{ item }}</li>
        </ul>
      </section>

      <footer class="c-legal-foot">
        <a href="/">{{ isZh ? '返回首页' : 'Back to home' }}</a>
      </footer>
    </main>
  </div>
</template>
