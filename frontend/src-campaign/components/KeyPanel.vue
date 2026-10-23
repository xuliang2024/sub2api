<script setup lang="ts">
import { computed, ref } from 'vue'
import type { CampaignKeyGrant } from '../types'
import { copyText } from '../format'
import { useCampaignMessages } from '../messages'

const props = defineProps<{
  open: boolean
  grant?: CampaignKeyGrant
  model?: string
  topupUrl?: string
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'retry'): void
}>()

const t = useCampaignMessages()
const copied = ref('')

const curlSample = computed(() => {
  if (!props.grant) {
    return ''
  }
  const base = props.grant.base_url.replace(/\/+$/, '')
  const model = props.model || 'your-model'
  return [
    `curl ${base}/v1/chat/completions \\`,
    `  -H "Authorization: Bearer ${props.grant.api_key}" \\`,
    `  -H "Content-Type: application/json" \\`,
    `  -d '{"model":"${model}","messages":[{"role":"user","content":"Hello"}]}'`,
  ].join('\n')
})

async function copy(label: string, text: string): Promise<void> {
  const ok = await copyText(text)
  if (ok) {
    copied.value = label
    setTimeout(() => {
      if (copied.value === label) {
        copied.value = ''
      }
    }, 1800)
  }
}
</script>

<template>
  <div v-if="open" class="c-modal-mask" @click.self="emit('close')">
    <div class="c-modal c-modal-wide" role="dialog" aria-modal="true">
      <div class="flex items-start justify-between gap-4">
        <h3 class="text-lg font-semibold">{{ t('key.title') }}</h3>
        <button
          type="button"
          class="text-sm"
          style="color: var(--c-text-soft)"
          :aria-label="t('auth.close')"
          @click="emit('close')"
        >
          {{ t('auth.close') }}
        </button>
      </div>

      <div v-if="loading" class="mt-6 space-y-3">
        <div class="c-skeleton w-2/3"></div>
        <div class="c-skeleton w-full"></div>
      </div>

      <p v-else-if="error" class="c-error mt-4">
        {{ error }}
        <button type="button" class="ml-2 underline" @click="emit('retry')">
          {{ t('key.retry') }}
        </button>
      </p>

      <template v-else-if="grant">
        <p class="mt-3 text-sm" style="color: var(--c-text-soft)">
          {{ t('key.desc') }}<template v-if="grant.group_name">
            {{ t('key.group', { name: grant.group_name }) }}</template
          >
        </p>

        <div class="mt-5">
          <div class="flex items-center justify-between">
            <span class="c-field-label mb-0">{{ t('key.apiKey') }}</span>
            <button
              class="c-btn c-btn-ghost c-btn-sm"
              type="button"
              @click="copy('key', grant.api_key)"
            >
              {{ copied === 'key' ? t('key.copied') : t('key.copy') }}
            </button>
          </div>
          <div class="c-key-box">{{ grant.api_key }}</div>
        </div>

        <div class="mt-4">
          <div class="flex items-center justify-between">
            <span class="c-field-label mb-0">{{ t('key.baseUrl') }}</span>
            <button
              class="c-btn c-btn-ghost c-btn-sm"
              type="button"
              @click="copy('base', grant.base_url)"
            >
              {{ copied === 'base' ? t('key.copied') : t('key.copy') }}
            </button>
          </div>
          <div class="c-key-box">{{ grant.base_url }}</div>
        </div>

        <div class="c-code mt-4">
          <button class="c-copy" type="button" @click="copy('curl', curlSample)">
            {{ copied === 'curl' ? t('key.copied') : t('key.copy') }}
          </button>
          <pre><code>{{ curlSample }}</code></pre>
        </div>

        <p class="c-note mt-3">{{ t('key.hint') }}</p>

        <div class="mt-5 flex flex-wrap gap-3">
          <a v-if="topupUrl" class="c-btn c-btn-primary" :href="topupUrl">{{ t('key.topup') }}</a>
          <button class="c-btn c-btn-ghost" type="button" @click="emit('close')">
            {{ t('key.continue') }}
          </button>
        </div>
      </template>
    </div>
  </div>
</template>
