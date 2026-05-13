import { describe, expect, it } from 'vitest'
import { createRequire } from 'node:module'

const require = createRequire(import.meta.url)
const {
  DEFAULT_SERVER_URL,
  DEV_SERVER_URL,
  isDesktopBridgePathAllowed,
  isTrustedSenderUrl,
  normalizeServerUrl,
  resolveInitialServerUrl,
} = require('../security.cjs')

describe('desktop security helpers', () => {
  it('normalizes server URLs and rejects unsupported protocols', () => {
    expect(normalizeServerUrl(' https://codex.apiz.ai/// ')).toBe('https://codex.apiz.ai')
    expect(() => normalizeServerUrl('file:///tmp/app.html')).toThrow(/http/)
  })

  it('resolves packaged and development defaults', () => {
    expect(resolveInitialServerUrl({ isPackaged: true })).toBe(DEFAULT_SERVER_URL)
    expect(resolveInitialServerUrl({ isPackaged: false })).toBe(DEV_SERVER_URL)
  })

  it('allows only the configured origin to call desktop IPC', () => {
    expect(isTrustedSenderUrl('https://codex.apiz.ai/keys', 'https://codex.apiz.ai')).toBe(true)
    expect(isTrustedSenderUrl('https://evil.example/keys', 'https://codex.apiz.ai')).toBe(false)
  })

  it('blocks auth and payment pages from the desktop bridge', () => {
    expect(isDesktopBridgePathAllowed('https://codex.apiz.ai/keys')).toBe(true)
    expect(isDesktopBridgePathAllowed('https://codex.apiz.ai/auth/linuxdo/callback')).toBe(false)
    expect(isDesktopBridgePathAllowed('https://codex.apiz.ai/payment/result')).toBe(false)
    expect(isTrustedSenderUrl('https://codex.apiz.ai/payment/result', 'https://codex.apiz.ai')).toBe(false)
  })
})
