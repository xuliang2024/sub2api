<script setup lang="ts">
import { SUPPORTED_LOCALES, type CampaignLocale } from '../messages'

defineProps<{ modelValue: CampaignLocale }>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: CampaignLocale): void
}>()

function handleChange(event: Event): void {
  const value = (event.target as HTMLSelectElement).value as CampaignLocale
  emit('update:modelValue', value)
}
</script>

<template>
  <label class="c-lang">
    <span class="sr-only">Language</span>
    <svg
      class="c-lang-icon"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.6"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="9" />
      <path d="M3 12h18M12 3c2.5 2.6 2.5 15.4 0 18M12 3c-2.5 2.6-2.5 15.4 0 18" stroke-linecap="round" />
    </svg>
    <select :value="modelValue" @change="handleChange">
      <option v-for="locale in SUPPORTED_LOCALES" :key="locale.code" :value="locale.code">
        {{ locale.label }}
      </option>
    </select>
  </label>
</template>
