import { afterEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn().mockResolvedValue(true)
  })
}))

import UseKeyModal from '../UseKeyModal.vue'

describe('UseKeyModal', () => {
  afterEach(() => {
    delete (window as any).sub2apiDesktop
    Object.defineProperty(window.navigator, 'userAgent', {
      value: 'Mozilla/5.0 jsdom',
      configurable: true,
    })
  })

  it('renders GPT-5.4 mini entry in OpenCode config', async () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'openai'
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const opencodeTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.opencode')
    )

    expect(opencodeTab).toBeDefined()
    await opencodeTab!.trigger('click')
    await nextTick()

    const codeBlock = wrapper.find('pre code')
    expect(codeBlock.exists()).toBe(true)
    expect(codeBlock.text()).toContain('"name": "GPT-5.4 Mini"')
    expect(codeBlock.text()).not.toContain('"name": "GPT-5.4 Nano"')
  })

  it('hides desktop actions when the desktop bridge is unavailable', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com',
        platform: 'openai'
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    expect(wrapper.text()).not.toContain('desktopCodex.desktopActions')
  })

  it('shows disabled desktop actions in Electron when the bridge is unavailable', () => {
    Object.defineProperty(window.navigator, 'userAgent', {
      value: 'Mozilla/5.0 Electron/42.0.1',
      configurable: true,
    })

    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com',
        platform: 'openai'
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('desktopCodex.desktopActions')
    expect(wrapper.text()).toContain('desktopCodex.bridgeUnavailable')
    expect(wrapper.findAll('button').some((button) =>
      button.text().includes('desktopCodex.oneClickWrite') && button.attributes('disabled') !== undefined
    )).toBe(true)
  })

  it('writes Codex config through the desktop bridge when available', async () => {
    const writeCodexConfig = vi.fn().mockResolvedValue({})
    const queryGatewayUsage = vi.fn().mockResolvedValue({ status: 'active' })
    ;(window as any).sub2apiDesktop = {
      isDesktop: true,
      writeCodexConfig,
      queryGatewayUsage,
      openCodexDownload: vi.fn(),
      openConfigDir: vi.fn(),
      getEnvironment: vi.fn(),
      getServerUrl: vi.fn(),
      setServerUrl: vi.fn(),
    }

    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com',
        platform: 'openai'
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const writeButton = wrapper.findAll('button').find((button) =>
      button.text().includes('desktopCodex.oneClickWrite')
    )

    expect(writeButton).toBeDefined()
    await writeButton!.trigger('click')
    await Promise.resolve()
    await nextTick()

    expect(writeCodexConfig).toHaveBeenCalledWith({
      baseUrl: 'https://example.com',
      apiKey: 'sk-test',
      mode: 'standard'
    })
    expect(queryGatewayUsage).toHaveBeenCalledWith({
      baseUrl: 'https://example.com',
      apiKey: 'sk-test'
    })
  })
})
