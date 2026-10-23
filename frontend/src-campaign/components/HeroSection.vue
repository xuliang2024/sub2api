<script setup lang="ts">
import type { CampaignHero } from '../types'
import { useCampaignMessages } from '../messages'

defineProps<{
  hero: CampaignHero
  ctaLabel: string
  trialNote?: string
}>()

const emit = defineEmits<{
  (event: 'cta'): void
}>()

const t = useCampaignMessages()
</script>

<template>
  <section class="c-hero" id="top">
    <div class="campaign-container">
      <span v-if="hero.eyebrow" class="c-eyebrow">{{ hero.eyebrow }}</span>

      <h1 class="c-hero-title">{{ hero.title }}</h1>
      <p v-if="hero.subtitle" class="c-hero-sub">{{ hero.subtitle }}</p>

      <div class="c-hero-actions">
        <button class="c-btn c-btn-primary" type="button" @click="emit('cta')">
          {{ ctaLabel }}
        </button>
        <a class="c-btn c-btn-ghost" href="#models">{{ t('hero.secondary') }}</a>
      </div>

      <p v-if="trialNote" class="mt-4 text-sm" style="color: var(--c-text-soft)">
        {{ trialNote }}
      </p>

      <div v-if="hero.badges?.length" class="c-badges">
        <span v-for="badge in hero.badges" :key="badge" class="c-chip">{{ badge }}</span>
      </div>

      <div v-if="hero.image" class="c-hero-media">
        <img :src="hero.image" alt="" loading="lazy" />
      </div>
    </div>
  </section>
</template>
