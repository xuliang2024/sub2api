/**
 * 静态法务页（隐私政策 / 服务条款）。
 *
 * 主站与落地页共用同一份文案（src/shared/legal.ts），这里只负责主站的排版与语言选择。
 * 路由通过 meta.doc 指明渲染哪一篇。
 */
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { legalDocuments, normalizeLegalLocale } from '@/shared/legal'

const route = useRoute()
const { locale } = useI18n()

const doc = computed(() => {
  const key = route.meta.doc === 'terms' ? 'terms' : 'privacy'
  return legalDocuments[key][normalizeLegalLocale(locale.value)]
})

const isZh = computed(() => normalizeLegalLocale(locale.value) === 'zh')
</script>

<template>
  <main class="mx-auto w-full max-w-3xl px-5 py-12 sm:px-8">
    <RouterLink
      :to="'/'"
      class="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
    >
      ← {{ isZh ? '返回首页' : 'Back to home' }}
    </RouterLink>

    <header class="mt-6 border-b border-gray-200 pb-6 dark:border-dark-700">
      <h1 class="text-3xl font-semibold tracking-tight text-gray-900 dark:text-white">
        {{ doc.title }}
      </h1>
      <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
        {{ isZh ? '最后更新：' : 'Last updated: ' }}{{ doc.updated }}
      </p>
      <p class="mt-4 text-[15px] leading-7 text-gray-600 dark:text-dark-300">{{ doc.intro }}</p>
    </header>

    <article class="mt-8 space-y-8">
      <section v-for="section in doc.sections" :key="section.heading">
        <h2 class="text-lg font-medium text-gray-900 dark:text-white">{{ section.heading }}</h2>
        <p
          v-for="(para, index) in section.paras ?? []"
          :key="`p-${index}`"
          class="mt-3 text-[15px] leading-7 text-gray-600 dark:text-dark-300"
        >
          {{ para }}
        </p>
        <ul
          v-if="section.bullets?.length"
          class="mt-3 list-disc space-y-2 pl-6 text-[15px] leading-7 text-gray-600 dark:text-dark-300"
        >
          <li v-for="(item, index) in section.bullets" :key="`b-${index}`">{{ item }}</li>
        </ul>
      </section>
    </article>
  </main>
</template>
