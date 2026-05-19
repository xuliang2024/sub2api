import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import UserDashboardSpendingRanking from '../UserDashboardSpendingRanking.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => ({
        'dashboard.spendingRanking': 'User Spending Ranking',
        'dashboard.top10Users': 'Top 10 users',
        'dashboard.last7Days': 'Last 7 days',
        'dashboard.you': 'You',
        'dashboard.requests': 'Requests',
        'dashboard.tokens': 'Tokens',
        'dashboard.noRankingData': 'No ranking data',
        'dashboard.rankingDataHint': 'No records',
      }[key] ?? key)
    })
  }
})

vi.mock('@/utils/format', () => ({
  formatCostFixed: (amount: number) => amount.toFixed(4),
  formatNumberLocaleString: (num: number) => num.toLocaleString(),
  formatTokensK: (tokens: number) => {
    if (tokens >= 1_000_000) return `${(tokens / 1_000_000).toFixed(1)}M`
    if (tokens >= 1000) return `${(tokens / 1000).toFixed(1)}K`
    return tokens.toString()
  },
}))

describe('UserDashboardSpendingRanking', () => {
  it('renders ranking rows and marks the current user', () => {
    const wrapper = mount(UserDashboardSpendingRanking, {
      props: {
        loading: false,
        totalActualCost: 10,
        items: [
          {
            rank: 1,
            display_name: 'a***e@e***e.com',
            actual_cost: 7.5,
            requests: 12,
            tokens: 12000,
            is_current_user: true,
          },
          {
            rank: 2,
            display_name: 'b***b@d***n.test',
            actual_cost: 2.5,
            requests: 8,
            tokens: 3000,
            is_current_user: false,
          },
        ],
      },
      global: {
        stubs: {
          LoadingSpinner: true,
          EmptyState: true,
        },
      },
    })

    expect(wrapper.text()).toContain('User Spending Ranking')
    expect(wrapper.text()).toContain('a***e@e***e.com')
    expect(wrapper.text()).toContain('You')
    expect(wrapper.text()).toContain('$7.5000')
    expect(wrapper.text()).toContain('75.0%')
    expect(wrapper.text()).toContain('12.0K Tokens')
  })
})
