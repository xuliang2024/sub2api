<script setup lang="ts">
import { computed, ref } from 'vue'
import type { CampaignModel } from '../types'
import { discountPercent, formatPerRequest, formatPricePerMillion } from '../format'
import { useCampaignMessages } from '../messages'

const props = defineProps<{ models: CampaignModel[] }>()

const t = useCampaignMessages()
const keyword = ref('')

const filtered = computed(() => {
  const query = keyword.value.trim().toLowerCase()
  if (!query) {
    return props.models
  }
  return props.models.filter(
    (model) =>
      model.model.toLowerCase().includes(query) ||
      (model.display_name || '').toLowerCase().includes(query) ||
      (model.platform || '').toLowerCase().includes(query),
  )
})

/** 官方价只在确实比我们贵时展示（与卡片区同一口径）。 */
function official(ours?: number, list?: number): string {
  if (ours == null || list == null || !(ours < list)) {
    return ''
  }
  return formatPricePerMillion(list)
}

function priceCell(model: CampaignModel): { official: string; ours: string } {
  if (model.per_request_price != null) {
    return { official: '', ours: formatPerRequest(model.per_request_price) }
  }
  return {
    official: official(model.input_price, model.official_input_price),
    ours: formatPricePerMillion(model.input_price),
  }
}

function outputCell(model: CampaignModel): { official: string; ours: string } {
  if (model.per_request_price != null) {
    return { official: '', ours: '—' }
  }
  return {
    official: official(model.output_price, model.official_output_price),
    ours: formatPricePerMillion(model.output_price),
  }
}

function saving(model: CampaignModel): number {
  return discountPercent(model.input_price, model.official_input_price)
}

/** 分时模型的高峰倍率；无分时或倍率 <=1 时返回 1。 */
function tierMultiplier(model: CampaignModel): number {
  const periods = model.time_pricing ?? []
  let max = 1
  for (const period of periods) {
    if (typeof period.multiplier === 'number' && period.multiplier > max) {
      max = period.multiplier
    }
  }
  return max
}

interface CatalogRow {
  key: string
  model: CampaignModel
  tier?: 'offpeak' | 'peak'
  inputOfficial: string
  inputOurs: string
  outputOfficial: string
  outputOurs: string
  cacheOurs: string
  saving: number
}

/**
 * 分时模型拆成两行（谷时 / 忙时）：官方价是不随档位变化的标价，
 * 折扣按「相对标价」算，谷时更深、忙时等于常规折扣。
 */
const rows = computed<CatalogRow[]>(() => {
  const out: CatalogRow[] = []
  for (const model of filtered.value) {
    const mult = tierMultiplier(model)
    if (mult <= 1) {
      out.push({
        key: `${model.platform}-${model.model}`,
        model,
        inputOfficial: priceCell(model).official,
        inputOurs: priceCell(model).ours,
        outputOfficial: outputCell(model).official,
        outputOurs: outputCell(model).ours,
        cacheOurs:
          model.cache_read_price != null ? formatPricePerMillion(model.cache_read_price) : '—',
        saving: saving(model),
      })
      continue
    }
    for (const tier of ['offpeak', 'peak'] as const) {
      const factor = tier === 'peak' ? mult : 1
      const scale = (value?: number | null) => (typeof value === 'number' ? value * factor : undefined)
      out.push({
        key: `${model.platform}-${model.model}-${tier}`,
        model,
        tier,
        inputOfficial: official(model.input_price, model.official_input_price),
        inputOurs: formatPricePerMillion(scale(model.input_price)),
        outputOfficial: official(model.output_price, model.official_output_price),
        outputOurs: formatPricePerMillion(scale(model.output_price)),
        cacheOurs:
          model.cache_read_price != null
            ? formatPricePerMillion(scale(model.cache_read_price))
            : '—',
        saving: discountPercent(scale(model.input_price), model.official_input_price),
      })
    }
  }
  return out
})
</script>

<template>
  <section v-if="models.length" class="c-section" id="catalog">
    <div class="campaign-container">
      <div class="c-section-head">
        <h2 class="c-section-title">{{ t('catalog.title') }}</h2>
        <p class="c-section-desc">
          {{ t('catalog.subtitle') }}
          <span class="c-catalog-count">{{ t('catalog.count', { count: models.length }) }}</span>
        </p>
      </div>

      <div class="c-catalog-toolbar">
        <input
          v-model="keyword"
          class="c-input c-catalog-search"
          type="search"
          :placeholder="t('catalog.search')"
        />
      </div>

      <div class="c-table-wrap">
        <table class="c-table">
          <thead>
            <tr>
              <th scope="col">{{ t('catalog.model') }}</th>
              <th scope="col">{{ t('models.input') }}</th>
              <th scope="col">{{ t('models.output') }}</th>
              <th scope="col">{{ t('models.cacheRead') }}</th>
              <th scope="col">{{ t('catalog.discount') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.key">
              <td>
                <div class="c-table-name">{{ row.model.display_name }}</div>
                <div class="c-table-id">
                  {{ row.model.model }}
                  <span v-if="row.tier" class="c-tier-tag">
                    {{ row.tier === 'peak' ? t('models.peakShort') : t('models.offpeak') }}
                  </span>
                </div>
              </td>
              <td class="c-table-price">
                <s v-if="row.inputOfficial" class="c-table-official">{{ row.inputOfficial }}</s>
                <span class="c-table-ours">{{ row.inputOurs }}</span>
              </td>
              <td class="c-table-price">
                <s v-if="row.outputOfficial" class="c-table-official">{{ row.outputOfficial }}</s>
                <span class="c-table-ours">{{ row.outputOurs }}</span>
              </td>
              <td class="c-table-price">
                <span class="c-table-ours">{{ row.cacheOurs }}</span>
              </td>
              <td>
                <span v-if="row.saving" class="c-discount">
                  {{ t('models.save', { percent: row.saving }) }}
                </span>
                <span v-else class="c-table-muted">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <p v-if="!filtered.length" class="c-note mt-4">{{ t('catalog.empty') }}</p>
    </div>
  </section>
</template>
