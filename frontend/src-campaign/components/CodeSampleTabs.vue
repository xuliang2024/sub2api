<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { CampaignCodeSample } from '../types'
import { copyText, fillSample, type SampleContext } from '../format'
import { useCampaignMessages } from '../messages'
import AgentSetup from './AgentSetup.vue'

const props = defineProps<{
  samples: CampaignCodeSample[]
  context: SampleContext
}>()

const t = useCampaignMessages()

const activeId = ref<string>('')
const copiedId = ref<string>('')
let copiedTimer: ReturnType<typeof setTimeout> | undefined

/**
 * 网关真实支持的协议面（见 backend/internal/server/routes/gateway.go）。
 * 这张表是「一次性对接、兼容接入」的实证：不需要新客户端，换 path 即可。
 */
const protocols = [
  { key: 'compat.openai', method: 'POST', path: '/v1/chat/completions' },
  { key: 'compat.anthropic', method: 'POST', path: '/v1/messages' },
  { key: 'compat.responses', method: 'POST', path: '/v1/responses' },
  { key: 'compat.gemini', method: 'POST', path: '/v1beta/models/{model}:generateContent' },
]

watch(
  () => props.samples,
  (samples) => {
    if (!samples.some((sample) => sample.id === activeId.value)) {
      activeId.value = samples[0]?.id ?? ''
    }
  },
  { immediate: true },
)

const activeSample = computed(
  () => props.samples.find((sample) => sample.id === activeId.value) ?? props.samples[0],
)

const activeCode = computed(() =>
  activeSample.value ? fillSample(activeSample.value.code, props.context) : '',
)

async function handleCopy(): Promise<void> {
  if (!activeSample.value || !activeCode.value) {
    return
  }
  const ok = await copyText(activeCode.value)
  if (!ok) {
    return
  }
  copiedId.value = activeSample.value.id
  if (copiedTimer) {
    clearTimeout(copiedTimer)
  }
  copiedTimer = setTimeout(() => {
    copiedId.value = ''
  }, 1800)
}
</script>

<template>
  <section class="c-section" id="integration">
    <div class="campaign-container">
      <div class="c-section-head">
        <h2 class="c-section-title">{{ t('compat.title') }}</h2>
        <p class="c-section-desc">{{ t('compat.note') }}</p>
      </div>

      <div class="c-compat">
        <div v-for="protocol in protocols" :key="protocol.path" class="c-compat-row">
          <span class="c-compat-name">{{ t(protocol.key) }}</span>
          <span class="c-compat-method">{{ protocol.method }}</span>
          <code class="c-compat-path">{{ protocol.path }}</code>
        </div>
      </div>

      <div v-if="samples.length" class="mt-10">
        <h3 class="c-section-title" style="font-size: 18px">{{ t('integration.title') }}</h3>
        <p class="c-section-desc">{{ t('integration.subtitle') }}</p>

        <div class="c-tabs mt-5" role="tablist">
          <button
            v-for="sample in samples"
            :key="sample.id"
            class="c-tab"
            type="button"
            role="tab"
            :aria-selected="sample.id === activeId"
            @click="activeId = sample.id"
          >
            {{ sample.label }}
          </button>
        </div>

        <div class="c-code">
          <button class="c-copy" type="button" @click="handleCopy">
            {{ copiedId === activeId ? t('integration.copied') : t('integration.copy') }}
          </button>
          <pre><code>{{ activeCode }}</code></pre>
        </div>
      </div>

      <AgentSetup :context="context" />
    </div>
  </section>
</template>
