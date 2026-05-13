import { afterEach, describe, expect, it, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const appStoreMock = vi.hoisted(() => ({
  apiBaseUrl: '',
  fetchPublicSettings: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn(),
}))

const apiMocks = vi.hoisted(() => ({
  getAvailable: vi.fn(),
  listKeys: vi.fn(),
  createKey: vi.fn(),
}))

vi.mock('vue-i18n', () => ({
  createI18n: () => ({
    global: {
      t: (key: string) => key,
      locale: { value: 'en' },
    },
    install: vi.fn(),
  }),
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => appStoreMock
}))

vi.mock('@/api', () => ({
  userGroupsAPI: {
    getAvailable: apiMocks.getAvailable,
  },
  keysAPI: {
    list: apiMocks.listKeys,
    create: apiMocks.createKey,
  },
}))

vi.mock('@/components/layout/AppLayout.vue', () => ({
  default: {
    template: '<div><slot /></div>',
  },
}))

vi.mock('@/components/icons/Icon.vue', () => ({
  default: {
    template: '<span />',
  },
}))

vi.mock('@/components/common/Select.vue', () => ({
  default: {
    props: ['modelValue', 'options', 'disabled', 'placeholder'],
    emits: ['update:modelValue'],
    template: '<div class="select-stub">{{ placeholder }}{{ options?.map((o) => o.label).join(",") }}</div>',
  },
}))

async function mountView() {
  const { default: CodexSetupView } = await import('../CodexSetupView.vue')
  return mount(CodexSetupView, {
    global: {},
  })
}

describe('CodexSetupView', () => {
  afterEach(() => {
    vi.clearAllMocks()
    delete (window as any).sub2apiDesktop
    appStoreMock.apiBaseUrl = ''
  })

  it('shows an empty OpenAI group state', async () => {
    appStoreMock.fetchPublicSettings.mockResolvedValue({ api_base_url: 'https://codex.apiz.ai' })
    apiMocks.getAvailable.mockResolvedValue([])
    apiMocks.listKeys.mockResolvedValue({ items: [] })

    const wrapper = await mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('desktopCodex.noOpenaiGroups')
  })

  it('creates a Codex Desktop key for the only OpenAI group', async () => {
    appStoreMock.fetchPublicSettings.mockResolvedValue({ api_base_url: 'https://codex.apiz.ai' })
    apiMocks.getAvailable.mockResolvedValue([
      {
        id: 7,
        name: 'OpenAI',
        platform: 'openai',
        status: 'active',
      },
    ])
    apiMocks.listKeys.mockResolvedValue({ items: [] })
    apiMocks.createKey.mockResolvedValue({
      id: 11,
      name: 'Codex Desktop',
      key: 'sk-test',
      group_id: 7,
      status: 'active',
    })

    const wrapper = await mountView()
    await flushPromises()

    const createButton = wrapper.findAll('button').find((button) =>
      button.text().includes('desktopCodex.createKey')
    )
    expect(createButton).toBeDefined()
    await createButton!.trigger('click')
    await flushPromises()

    expect(apiMocks.createKey).toHaveBeenCalledWith('Codex Desktop', 7)
  })

  it('writes and verifies Codex config with the selected key', async () => {
    const writeCodexConfig = vi.fn().mockResolvedValue({})
    const queryGatewayUsage = vi.fn().mockResolvedValue({ status: 'active', available: 10 })
    ;(window as any).sub2apiDesktop = {
      isDesktop: true,
      getServerUrl: vi.fn().mockResolvedValue('https://codex.apiz.ai'),
      getEnvironment: vi.fn().mockResolvedValue({
        platform: 'darwin',
        configPath: '/tmp/.codex/config.toml',
        authPath: '/tmp/.codex/auth.json',
        existingMode: 'standard',
        existingBaseUrl: '',
      }),
      writeCodexConfig,
      queryGatewayUsage,
      openCodexDownload: vi.fn(),
      openConfigDir: vi.fn(),
      setServerUrl: vi.fn(),
    }
    appStoreMock.fetchPublicSettings.mockResolvedValue({ api_base_url: 'https://codex.apiz.ai' })
    apiMocks.getAvailable.mockResolvedValue([
      {
        id: 7,
        name: 'OpenAI',
        platform: 'openai',
        status: 'active',
      },
    ])
    apiMocks.listKeys.mockResolvedValue({
      items: [
        {
          id: 11,
          name: 'Desktop Key',
          key: 'sk-test',
          group_id: 7,
          status: 'active',
        },
      ],
    })

    const wrapper = await mountView()
    await flushPromises()

    const writeButton = wrapper.findAll('button').find((button) =>
      button.text().includes('desktopCodex.writeConfig')
    )
    expect(writeButton).toBeDefined()
    await writeButton!.trigger('click')
    await flushPromises()

    expect(writeCodexConfig).toHaveBeenCalledWith({
      baseUrl: 'https://codex.apiz.ai',
      apiKey: 'sk-test',
      mode: 'standard',
    })
    expect(queryGatewayUsage).toHaveBeenCalledWith({
      baseUrl: 'https://codex.apiz.ai',
      apiKey: 'sk-test',
    })
  })

  it('auto-selects the codex group and creates a key before writing', async () => {
    const writeCodexConfig = vi.fn().mockResolvedValue({})
    const queryGatewayUsage = vi.fn().mockResolvedValue({ status: 'active', available: 10 })
    ;(window as any).sub2apiDesktop = {
      isDesktop: true,
      getServerUrl: vi.fn().mockResolvedValue('https://codex.apiz.ai'),
      getEnvironment: vi.fn().mockResolvedValue({
        platform: 'darwin',
        configPath: '/tmp/.codex/config.toml',
        authPath: '/tmp/.codex/auth.json',
        existingMode: 'standard',
        existingBaseUrl: '',
      }),
      writeCodexConfig,
      queryGatewayUsage,
      openCodexDownload: vi.fn(),
      openConfigDir: vi.fn(),
      setServerUrl: vi.fn(),
    }
    appStoreMock.fetchPublicSettings.mockResolvedValue({ api_base_url: 'https://codex.apiz.ai' })
    apiMocks.getAvailable.mockResolvedValue([
      {
        id: 3,
        name: 'OpenAI',
        platform: 'openai',
        status: 'active',
      },
      {
        id: 8,
        name: 'codex',
        platform: 'openai',
        status: 'active',
      },
    ])
    apiMocks.listKeys.mockResolvedValue({ items: [] })
    apiMocks.createKey.mockResolvedValue({
      id: 21,
      name: 'Codex Desktop',
      key: 'sk-created',
      group_id: 8,
      status: 'active',
    })

    const wrapper = await mountView()
    await flushPromises()

    const writeButton = wrapper.findAll('button').find((button) =>
      button.text().includes('desktopCodex.writeConfig')
    )
    expect(writeButton).toBeDefined()
    await writeButton!.trigger('click')
    await flushPromises()

    expect(apiMocks.createKey).toHaveBeenCalledWith('Codex Desktop', 8)
    expect(writeCodexConfig).toHaveBeenCalledWith({
      baseUrl: 'https://codex.apiz.ai',
      apiKey: 'sk-created',
      mode: 'standard',
    })
  })
})
