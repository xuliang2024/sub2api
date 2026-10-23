<script setup lang="ts">
import { computed } from 'vue'
import type { CampaignGroup, CampaignModel } from '../types'
import {
  discountPercent,
  formatPerRequest,
  formatPricePerMillion,
} from '../format'
import { useCampaignMessages } from '../messages'

const props = defineProps<{
  models: CampaignModel[]
  group?: CampaignGroup
  /** 分组展示名（来自 brand.group_label）；为空则不在标题里提分组。 */
  groupLabel?: string
  /** 分组全量模型数；大于卡片数时才提示「查看全部」。 */
  catalogCount?: number
}>()

const emit = defineEmits<{
  (event: 'open-catalog'): void
}>()

const t = useCampaignMessages()

const subtitle = computed(() => {
  const base = t('models.subtitle')
  const label = props.groupLabel?.trim()
  return label ? `${base} ${t('models.group', { name: label })}` : base
})

interface PriceRow {
  label: string
  /** 官方价；为空表示没有官方价、或我们并不更便宜——两种情况都不展示对比。 */
  official: string
  ours: string
}

/**
 * 官方价只在「确实比我们贵」时才展示。
 *
 * 反之（我们更贵或持平）留空：广告页放一个对我们不利的对比没有意义，
 * 也不能为了让对比好看去改数字。
 */
function officialCell(ours?: number, official?: number): string {
  if (ours == null || official == null || !(ours < official)) {
    return ''
  }
  return formatPricePerMillion(official)
}

function priceRows(model: CampaignModel): PriceRow[] {
  if (model.per_request_price != null) {
    return [
      {
        label: t('models.perRequest'),
        official: '',
        ours: formatPerRequest(model.per_request_price),
      },
    ]
  }

  const rows: PriceRow[] = [
    {
      label: t('models.input'),
      official: officialCell(model.input_price, model.official_input_price),
      ours: formatPricePerMillion(model.input_price),
    },
    {
      label: t('models.output'),
      official: officialCell(model.output_price, model.official_output_price),
      ours: formatPricePerMillion(model.output_price),
    },
  ]
  if (model.cache_read_price != null) {
    rows.push({
      label: t('models.cacheRead'),
      official: '',
      ours: formatPricePerMillion(model.cache_read_price),
    })
  }
  return rows
}

/** 没有任何价格（分组未配置该模型）时不渲染价目表，避免一张空表。 */
function hasPrice(model: CampaignModel): boolean {
  return (
    model.per_request_price != null ||
    model.input_price != null ||
    model.output_price != null ||
    model.cache_read_price != null
  )
}

function saving(model: CampaignModel): number {
  return discountPercent(model.input_price, model.official_input_price)
}

function peakNote(model: CampaignModel): string {
  if (!model.time_pricing?.length) {
    return ''
  }
  const periods = model.time_pricing
    .map((period) => `${period.start_time}-${period.end_time}×${period.multiplier}`)
    .join(' ')
  return t('models.peak', { periods })
}
</script>

<template>
  <section class="c-section c-section-soft" id="models">
    <div class="campaign-container">
      <div class="c-section-head">
        <h2 class="c-section-title">{{ t('models.title') }}</h2>
        <p class="c-section-desc">{{ subtitle }}</p>
        <button
          v-if="(catalogCount ?? 0) > models.length"
          class="c-link mt-2"
          type="button"
          @click="emit('open-catalog')"
        >
          {{ t('models.viewAll', { count: catalogCount ?? 0 }) }} →
        </button>
      </div>

      <div v-if="models.length" class="c-model-grid">
        <article v-for="model in models" :key="model.model" class="c-card c-model-card">
          <div class="c-model-top">
            <div class="min-w-0">
              <div class="c-model-name">{{ model.display_name }}</div>
              <div class="c-model-id">{{ model.model }}</div>
            </div>
            <span v-if="model.badge" class="c-badge">{{ model.badge }}</span>
            <span v-else-if="!model.available" class="c-badge c-badge-muted">
              {{ t('models.unavailable') }}
            </span>
          </div>

          <p v-if="model.tagline" class="c-model-tagline">{{ model.tagline }}</p>

          <p v-if="!hasPrice(model)" class="c-note">
            {{ t('models.unavailable') }}
          </p>

          <div v-if="hasPrice(model)" class="c-price-table">
            <div class="c-price-head">
              <span></span>
              <span>{{ t('models.officialPrice') }}</span>
              <span>{{ t('models.ourPrice') }}</span>
            </div>
            <div v-for="row in priceRows(model)" :key="row.label" class="c-price-row">
              <span class="c-price-label">{{ row.label }}</span>
              <span class="c-price-official" :class="{ 'is-blank': !row.official }">
                {{ row.official }}
              </span>
              <span class="c-price-value">{{ row.ours }}</span>
            </div>
          </div>

          <div class="c-model-foot">
            <span v-if="saving(model)" class="c-discount">
              {{ t('models.save', { percent: saving(model) }) }}
            </span>
            <span
              v-if="peakNote(model)"
              class="c-note"
              :title="model.time_pricing_timezone || ''"
            >
              {{ peakNote(model) }}
            </span>
          </div>
        </article>
      </div>

      <p v-else class="text-sm" style="color: var(--c-text-soft)">{{ t('models.empty') }}</p>
    </div>
  </section>
</template>
