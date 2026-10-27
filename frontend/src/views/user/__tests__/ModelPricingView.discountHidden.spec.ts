import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import ModelPricingView from '../ModelPricingView.vue'

const channelsAPI = vi.hoisted(() => ({ getAvailable: vi.fn() }))
const groupsAPI = vi.hoisted(() => ({ getUserGroupRates: vi.fn() }))

vi.mock('@/api/channels', () => ({ default: channelsAPI }))
vi.mock('@/api/groups', () => ({ default: groupsAPI }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showError: vi.fn() }) }))
vi.mock('vue-i18n', async (importOriginal) => ({
  ...(await importOriginal<typeof import('vue-i18n')>()),
  useI18n: () => ({ t }),
}))

enableAutoUnmount(afterEach)

/** 极简 t：未命中原样返回 key，便于断言「某个词条没有渲染」。 */
function t(key: string, _params?: Record<string, unknown>): string {
  return key
}

function pricing(input: number, output: number, cacheRead: number, cacheWrite: number) {
  return {
    billing_mode: 'token' as const,
    input_price: input,
    output_price: output,
    cache_read_price: cacheRead,
    cache_write_price: cacheWrite,
    image_input_price: null,
    image_output_price: null,
    per_request_price: null,
    intervals: [],
  }
}

function group(id: number, name: string, rate: number) {
  return {
    id,
    name,
    platform: 'openai',
    subscription_type: 'standard',
    rate_multiplier: rate,
    peak_rate_enabled: false,
    peak_start: '',
    peak_end: '',
    peak_rate_multiplier: 1,
    is_exclusive: false,
  }
}

const official = {
  input_price: 3e-7,
  output_price: 1.2e-6,
  cache_read_price: 6e-9,
  cache_write_price: 3e-7,
}

function oneChannel() {
  return [
    {
      name: '国产模型线路',
      description: '',
      platforms: [
        {
          platform: 'openai',
          groups: [group(31, '国产全模型', 1)],
          supported_models: [
            {
              name: 'deepseek-flash',
              platform: 'openai',
              pricing: pricing(2e-7, 6e-7, 6e-9, 2e-7),
              official_pricing: official,
            },
          ],
        },
      ],
    },
  ]
}

beforeEach(() => {
  vi.clearAllMocks()
  groupsAPI.getUserGroupRates.mockResolvedValue({})
})

describe('model pricing view without the discount column', () => {
  it('renders no discount column, no average discount and no discount filter', async () => {
    channelsAPI.getAvailable.mockResolvedValue(oneChannel())
    const wrapper = mount(ModelPricingView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          Icon: true,
          PlatformIcon: true,
          Select: true,
          Toggle: true,
        },
      },
    })
    await flushPromises()

    const header = wrapper.get('thead').text()
    expect(header).not.toContain('modelPricing.columns.discount')
    expect(header).not.toContain('modelPricing.stats.avgDiscount')
    expect(wrapper.text()).not.toContain('modelPricing.onlyDiscounted')

    // 模型 + 官方 3 + 本站 3 + 分组 + 展开 = 9 列
    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    expect(row.findAll('td')).toHaveLength(9)
  })

  it('keeps the expanded offer table free of the "vs official" column', async () => {
    channelsAPI.getAvailable.mockResolvedValue(oneChannel())
    const wrapper = mount(ModelPricingView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          Icon: true,
          PlatformIcon: true,
          Select: true,
          Toggle: true,
        },
      },
    })
    await flushPromises()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    await row.trigger('click')

    expect(wrapper.get('table table thead').text()).not.toContain('modelPricing.expand.vsOfficial')
    const subRow = wrapper.get('table table tbody tr')
    // 分组、渠道、倍率、输入、输出、缓存读、缓存写、末尾占位 = 8 列
    expect(subRow.findAll('td')).toHaveLength(8)
  })
})
