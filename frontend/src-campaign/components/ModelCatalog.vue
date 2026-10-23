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
            <tr v-for="model in filtered" :key="`${model.platform}-${model.model}`">
              <td>
                <div class="c-table-name">{{ model.display_name }}</div>
                <div class="c-table-id">{{ model.model }}</div>
              </td>
              <td class="c-table-price">
                <s v-if="priceCell(model).official" class="c-table-official">
                  {{ priceCell(model).official }}
                </s>
                <span class="c-table-ours">{{ priceCell(model).ours }}</span>
              </td>
              <td class="c-table-price">
                <s v-if="outputCell(model).official" class="c-table-official">
                  {{ outputCell(model).official }}
                </s>
                <span class="c-table-ours">{{ outputCell(model).ours }}</span>
              </td>
              <td class="c-table-price">
                <span class="c-table-ours">
                  {{
                    model.cache_read_price != null
                      ? formatPricePerMillion(model.cache_read_price)
                      : '—'
                  }}
                </span>
              </td>
              <td>
                <span v-if="saving(model)" class="c-discount">
                  {{ t('models.save', { percent: saving(model) }) }}
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
