<script setup lang="ts">
import { computed, ref } from 'vue'
import { copyText, fillSample, type SampleContext } from '../format'
import { useCampaignMessages } from '../messages'
import { AGENT_SAMPLES } from '../agents'

const props = defineProps<{
  context: SampleContext
}>()

const t = useCampaignMessages()

const activeId = ref<string>(AGENT_SAMPLES[0]?.id ?? '')
const copiedId = ref<string>('')
let copiedTimer: ReturnType<typeof setTimeout> | undefined

const activeAgent = computed(
  () => AGENT_SAMPLES.find((agent) => agent.id === activeId.value) ?? AGENT_SAMPLES[0],
)

const activeCode = computed(() =>
  activeAgent.value ? fillSample(activeAgent.value.code, props.context) : '',
)

async function handleCopy(): Promise<void> {
  if (!activeAgent.value || !activeCode.value) {
    return
  }
  const ok = await copyText(activeCode.value)
  if (!ok) {
    return
  }
  copiedId.value = activeAgent.value.id
  if (copiedTimer) {
    clearTimeout(copiedTimer)
  }
  copiedTimer = setTimeout(() => {
    copiedId.value = ''
  }, 1800)
}
</script>

<template>
  <div class="c-agents">
    <h3 class="c-section-title" style="font-size: 18px">{{ t('agents.title') }}</h3>
    <p class="c-section-desc">{{ t('agents.subtitle') }}</p>

    <div class="c-tabs mt-5" role="tablist">
      <button
        v-for="agent in AGENT_SAMPLES"
        :key="agent.id"
        class="c-tab"
        type="button"
        role="tab"
        :aria-selected="agent.id === activeId"
        @click="activeId = agent.id"
      >
        {{ agent.label }}
      </button>
    </div>

    <p v-if="activeAgent" class="c-agent-note">{{ t(activeAgent.noteKey) }}</p>

    <div class="c-code">
      <button class="c-copy" type="button" @click="handleCopy">
        {{ copiedId === activeId ? t('integration.copied') : t('integration.copy') }}
      </button>
      <pre><code>{{ activeCode }}</code></pre>
    </div>
  </div>
</template>
