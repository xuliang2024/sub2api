import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import PartnerView from '../PartnerView.vue'

const { getOverview, listUsers, updateUserRates, showError, showSuccess } = vi.hoisted(() => ({
  getOverview: vi.fn(),
  listUsers: vi.fn(),
  updateUserRates: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn(),
}))

vi.mock('@/api/partner', () => ({
  default: {
    getOverview,
    listUsers,
    listCommission: vi.fn(),
    updateUserRates,
  },
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError, showSuccess }),
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({ copyToClipboard: vi.fn() }),
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, unknown>) =>
        params ? `${key}|${JSON.stringify(params)}` : key,
    }),
  }
})

const overview = {
  user_id: 6322,
  email: 'partner@test.local',
  is_partner: true,
  min_rate_multiplier: 0.5001,
  total_commission: 0,
  month_commission: 0,
  pending_commission: 0,
  downline_count: 1,
  aff_code: 'BTK2YWEK9EL7',
}

function downline() {
  return {
    items: [
      {
        user_id: 6321,
        email: 'downline@test.local',
        username: 'downline',
        bound_source: 'admin',
        total_cost_30d: 0,
        commission_30d: 0,
        rates: [
          { group_id: 25, group_name: '聚合', default_rate: 1, rate_multiplier: null },
        ],
      },
    ],
  }
}

function mountView() {
  return mount(PartnerView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        StatCard: true,
        EmptyState: true,
        Pagination: true,
        Select: true,
        Icon: true,
      },
    },
  })
}

async function expandDownline(wrapper: ReturnType<typeof mountView>) {
  const trigger = wrapper
    .findAll('button')
    .find((btn) => btn.text().includes('downline@test.local'))
  expect(trigger, '下线行按钮应存在').toBeTruthy()
  await trigger!.trigger('click')
  await flushPromises()
}

describe('PartnerView 倍率滑块', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    getOverview.mockResolvedValue(overview)
    listUsers.mockResolvedValue(downline())
    updateUserRates.mockResolvedValue({ code: 0 })
  })

  it('用 range 滑块渲染，不再是 number 输入框', async () => {
    const wrapper = mountView()
    await flushPromises()
    await expandDownline(wrapper)

    expect(wrapper.find('input[type="range"]').exists()).toBe(true)
    expect(wrapper.find('input[type="number"]').exists()).toBe(false)
  })

  it('拖动滑块不会崩溃，并能显示数值', async () => {
    const wrapper = mountView()
    await flushPromises()
    await expandDownline(wrapper)

    const slider = wrapper.find('input[type="range"]')
    await slider.setValue('0.8')

    // 回归点：旧实现里 v-model 把值变成 number 后 .trim() 抛错，整页白屏。
    expect(wrapper.html()).not.toBe('')
    expect(wrapper.text()).toContain('×0.8')
    expect(showError).not.toHaveBeenCalled()
  })

  it('保存时只提交拖动过的分组', async () => {
    const wrapper = mountView()
    await flushPromises()
    await expandDownline(wrapper)

    await wrapper.find('input[type="range"]').setValue('0.8')
    const saveBtn = wrapper.findAll('button').find((btn) => btn.text() === 'common.save')
    expect(saveBtn, '保存按钮应存在').toBeTruthy()
    await saveBtn!.trigger('click')
    await flushPromises()

    expect(updateUserRates).toHaveBeenCalledWith(6321, [
      { group_id: 25, rate_multiplier: 0.8 },
    ])
  })

  it('全部跟随默认时提示先拖动，不发请求', async () => {
    const wrapper = mountView()
    await flushPromises()
    await expandDownline(wrapper)

    const saveBtn = wrapper.findAll('button').find((btn) => btn.text() === 'common.save')
    await saveBtn!.trigger('click')
    await flushPromises()

    expect(updateUserRates).not.toHaveBeenCalled()
    expect(showError).toHaveBeenCalled()
  })
})
