import { describe, expect, it } from 'vitest'
import { createRequire } from 'node:module'
import fs from 'node:fs/promises'
import os from 'node:os'
import path from 'node:path'

const require = createRequire(import.meta.url)
const {
  buildConfigBlock,
  buildGatewayUsageUrl,
  getEnvironment,
  normalizeAndValidateInput,
  queryGatewayUsage,
  stripManagedConfig,
  writeCodexConfig,
} = require('../configWriter.cjs')

async function tempHome() {
  return fs.mkdtemp(path.join(os.tmpdir(), 'sub2api-codex-desktop-'))
}

describe('desktop configWriter', () => {
  it('normalizes config input and defaults to standard mode', () => {
    expect(normalizeAndValidateInput({
      baseUrl: ' https://codex.apiz.ai/// ',
      apiKey: ' sk-test ',
    })).toEqual({
      baseUrl: 'https://codex.apiz.ai',
      apiKey: 'sk-test',
      mode: 'standard',
    })

    expect(() => normalizeAndValidateInput({ baseUrl: 'ftp://example.test', apiKey: 'sk-test' }))
      .toThrow(/http/)
  })

  it('builds standard and websocket Codex config blocks', () => {
    const standard = buildConfigBlock('https://codex.apiz.ai', 'standard')
    expect(standard).toContain('model = "gpt-5.4"')
    expect(standard).not.toContain('responses_websockets_v2')

    const websocket = buildConfigBlock('https://codex.apiz.ai', 'websocket')
    expect(websocket).toContain('supports_websockets = true')
    expect(websocket).toContain('[features]')
    expect(websocket).toContain('responses_websockets_v2 = true')
  })

  it('strips managed config while preserving unrelated sections', () => {
    const existing = [
      'model = "old-model"',
      'network_access = "disabled"',
      '',
      '[model_providers.OpenAI]',
      'base_url = "https://old.example"',
      '',
      '[features]',
      'responses_websockets_v2 = true',
      'experimental_client = true',
      '',
      '[profiles.dev]',
      'approval_policy = "never"',
    ].join('\n')

    const stripped = stripManagedConfig(existing)
    expect(stripped).not.toContain('old-model')
    expect(stripped).not.toContain('model_providers.OpenAI')
    expect(stripped).not.toContain('responses_websockets_v2')
    expect(stripped).toContain('experimental_client = true')
    expect(stripped).toContain('[profiles.dev]')
  })

  it('writes config and auth files with backups and preserved auth fields', async () => {
    const homeDir = await tempHome()
    const configDir = path.join(homeDir, '.codex')
    const configPath = path.join(configDir, 'config.toml')
    const authPath = path.join(configDir, 'auth.json')
    await fs.mkdir(configDir, { recursive: true })
    await fs.writeFile(configPath, 'model = "old"\n[profiles.dev]\napproval_policy = "never"\n')
    await fs.writeFile(authPath, JSON.stringify({ EXISTING_TOKEN: 'keep-me' }, null, 2))

    const result = await writeCodexConfig(
      { baseUrl: 'https://codex.apiz.ai/', apiKey: 'sk-1234567890', mode: 'websocket' },
      { homeDir, platform: 'linux' },
    )

    const config = await fs.readFile(configPath, 'utf8')
    const auth = JSON.parse(await fs.readFile(authPath, 'utf8'))
    const environment = await getEnvironment({ homeDir, platform: 'linux' })

    expect(result.configBackup).toBeTruthy()
    expect(result.authBackup).toBeTruthy()
    expect(config.startsWith(buildConfigBlock('https://codex.apiz.ai', 'websocket'))).toBe(true)
    expect(config).toContain('[profiles.dev]')
    expect(auth.EXISTING_TOKEN).toBe('keep-me')
    expect(auth.OPENAI_API_KEY).toBe('sk-1234567890')
    expect(environment.existingMode).toBe('websocket')
  })

  it('queries gateway usage with bearer auth', async () => {
    let request: any
    const result = await queryGatewayUsage(
      {
        baseUrl: 'https://codex.apiz.ai',
        apiKey: 'sk-test',
      },
      {
        fetchImpl: async (url: string, options: any) => {
          request = { url, options }
          return {
            ok: true,
            status: 200,
            text: async () => JSON.stringify({ status: 'active', available: 10 }),
          }
        },
      },
    )

    expect(request.url).toBe('https://codex.apiz.ai/v1/usage')
    expect(request.options.headers.authorization).toBe('Bearer sk-test')
    expect(result.available).toBe(10)
    expect(buildGatewayUsageUrl('https://codex.apiz.ai/v1')).toBe('https://codex.apiz.ai/v1/usage')
  })
})
