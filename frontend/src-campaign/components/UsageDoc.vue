<script setup lang="ts">
import { computed } from 'vue'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { useCampaignMessages } from '../messages'

const props = defineProps<{ markdown: string }>()

const t = useCampaignMessages()

const rendered = computed(() => {
  const source = props.markdown?.trim()
  if (!source) {
    return ''
  }
  // 使用说明由管理员配置，仍按不可信输入处理：markdown -> HTML -> DOMPurify。
  // 允许 iframe 是为了嵌入视频演示，与主站 CustomPageView 的策略保持一致。
  const html = marked.parse(source, { async: false })
  return DOMPurify.sanitize(html, {
    ADD_TAGS: ['iframe'],
    ADD_ATTR: ['allow', 'allowfullscreen', 'frameborder', 'scrolling', 'target'],
  })
})
</script>

<template>
  <section v-if="rendered" class="c-section" id="usage">
    <div class="campaign-container">
      <div class="c-section-head">
        <h2 class="c-section-title">{{ t('usage.title') }}</h2>
      </div>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div class="markdown-body" v-html="rendered"></div>
    </div>
  </section>
</template>
