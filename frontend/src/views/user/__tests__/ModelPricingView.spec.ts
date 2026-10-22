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

const messages: Record<string, string> = {
  'modelPricing.discount.value': '{value} 折',
  'modelPricing.discount.save': '省 {percent}%',
  'modelPricing.discount.above': '高于官价 {percent}%',
  'modelPricing.discount.notListed': '官方未收录',
  'modelPricing.discount.range': '{min}–{max} 折',
  'modelPricing.discount.timePricing': '分时价 · 峰时 ×{multiplier}',
  'modelPricing.discount.outputDiscount': '输出 {value} 折',
  'modelPricing.blocks.unit': 'USD / 1M',
  'modelPricing.units.perRequestShort': '次',
  'modelPricing.units.perImage': '张',
  'modelPricing.cacheLabels.read': '读',
  'modelPricing.cacheLabels.write': '写',
}

/** 极简 t：只解析用到的词条，未命中原样返回 key，便于断言。 */
function t(key: string, params?: Record<string, unknown>): string {
  const template = messages[key] ?? key
  return template.replace(/\{(\w+)\}/g, (_m, name: string) => String(params?.[name] ?? ''))
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

/** 两条渠道都给同一个模型报价：合并后应只有一行。 */
function twoChannelsSameModel() {
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
              pricing: pricing(2.2e-7, 6.6e-7, 7e-9, 2.2e-7),
              official_pricing: official,
            },
          ],
        },
      ],
    },
    {
      name: '备用线路 A',
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

async function mountView() {
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
  return wrapper
}

beforeEach(() => {
  vi.clearAllMocks()
  groupsAPI.getUserGroupRates.mockResolvedValue({})
})

describe('model-centric price table', () => {
  it('merges the same model served by multiple channels into one row', async () => {
    channelsAPI.getAvailable.mockResolvedValue(twoChannelsSameModel())
    const wrapper = await mountView()

    const rows = wrapper.findAll('tbody tr').filter((tr) => tr.classes().includes('cursor-pointer'))
    expect(rows).toHaveLength(1)
    expect(rows[0].text()).toContain('deepseek-flash')
  })

  it('shows the cheapest grouped offer and its discount against the official price', async () => {
    channelsAPI.getAvailable.mockResolvedValue(twoChannelsSameModel())
    const wrapper = await mountView()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    const cells = row.findAll('td')
    // 本站输入价取最便宜的报价：2e-7 × 1e6 = 0.20（统一两位小数）
    expect(cells[4].text()).toBe('0.20')
    // 官方输入：3e-7 × 1e6 = 0.30
    expect(cells[1].text()).toBe('0.30')
    // 折扣 0.2 / 0.3 = 0.667 → 6.7 折
    expect(cells[7].text()).toContain('6.7 折')
  })

  it('flags offers priced above the official list price', async () => {
    channelsAPI.getAvailable.mockResolvedValue([
      {
        name: '贵线路',
        description: '',
        platforms: [
          {
            platform: 'openai',
            groups: [group(31, '国产全模型', 2)],
            supported_models: [
              {
                name: 'overpriced-model',
                platform: 'openai',
                pricing: pricing(2.2e-7, 6.6e-7, 7e-9, 2.2e-7),
                official_pricing: official,
              },
            ],
          },
        ],
      },
    ])
    const wrapper = await mountView()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    const cells = row.findAll('td')
    // 2.2e-7 × 倍率 2 = 4.4e-7 vs 官价 3e-7 → 14.7 折，且标注高于官价
    expect(cells[7].text()).toContain('14.7 折')
    expect(cells[7].text()).toContain('高于官价 47%')
  })

  it('expands a row into every group × channel offer', async () => {
    channelsAPI.getAvailable.mockResolvedValue(twoChannelsSameModel())
    const wrapper = await mountView()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    await row.trigger('click')

    const detail = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))
    expect(detail).toBeDefined()
    // 展开后出现子表，两条线路各一行
    const subRows = wrapper.findAll('table table tbody tr')
    expect(subRows).toHaveLength(2)
    // 按实付价升序：备用线路 A（0.2）在前，国产模型线路（0.22）在后
    expect(subRows[0].text()).toContain('备用线路 A')
    expect(subRows[1].text()).toContain('国产模型线路')
  })

  it('shows a discount range for models with time-based pricing', async () => {
    channelsAPI.getAvailable.mockResolvedValue([
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
                pricing: {
                  ...pricing(7.5e-8, 3e-7, 1.5e-9, 7.5e-8),
                  time_pricing: { timezone: 'Asia/Shanghai', periods: [{ start_time: '09:00', end_time: '12:00', multiplier: 2 }] },
                  max_time_multiplier: 2,
                },
                official_pricing: official,
              },
            ],
          },
        ],
      },
    ])
    const wrapper = await mountView()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    const discountCell = row.findAll('td')[7]
    // 谷时 0.075 / 0.30 = 2.5 折，峰时 0.15 / 0.30 = 5.0 折
    expect(discountCell.text()).toContain('2.5–5.0 折')
    expect(discountCell.text()).toContain('峰时 ×2.0')
  })

  it('marks models without an official price instead of inventing a discount', async () => {
    channelsAPI.getAvailable.mockResolvedValue([
      {
        name: 'no-official',
        description: '',
        platforms: [
          {
            platform: 'openai',
            groups: [group(31, '国产全模型', 1)],
            supported_models: [
              {
                name: 'spark-bunny',
                platform: 'openai',
                pricing: pricing(7.4e-8, 2.4e-7, 7e-9, 7.4e-8),
                official_pricing: null,
              },
            ],
          },
        ],
      },
    ])
    const wrapper = await mountView()

    const row = wrapper.findAll('tbody tr').find((tr) => tr.classes().includes('cursor-pointer'))!
    expect(row.findAll('td')[7].text()).toContain('官方未收录')
  })

  it('hides models without a discount only when the toggle is on', async () => {
    channelsAPI.getAvailable.mockResolvedValue([
      {
        name: 'at-list-price',
        description: '',
        platforms: [
          {
            platform: 'openai',
            groups: [group(31, '国产全模型', 1)],
            supported_models: [
              {
                name: 'list-price-model',
                platform: 'openai',
                pricing: pricing(3e-7, 1.2e-6, 6e-9, 3e-7),
                official_pricing: official,
              },
            ],
          },
        ],
      },
    ])
    const wrapper = await mountView()

    expect(wrapper.findAll('tbody tr').filter((tr) => tr.classes().includes('cursor-pointer'))).toHaveLength(1)
    // 点「仅看有折扣」文案同样会切换开关
    await wrapper.get('.cursor-pointer.select-none').trigger('click')
    await flushPromises()
    expect(wrapper.findAll('tbody tr').filter((tr) => tr.classes().includes('cursor-pointer'))).toHaveLength(0)
  })
})
