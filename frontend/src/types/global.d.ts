import type { PublicSettings } from '@/types'

export type CodexConfigMode = 'standard' | 'websocket'

export interface CodexDesktopEnvironment {
  platform: string
  homeDir: string
  configDir: string
  configPath: string
  authPath: string
  hasConfig: boolean
  hasAuth: boolean
  existingBaseUrl: string
  existingMode: CodexConfigMode
  existingApiKey: string
  hasExistingApiKey: boolean
  hasChatGptTokens: boolean
}

export interface CodexConfigWriteInput {
  baseUrl: string
  apiKey: string
  mode?: CodexConfigMode
}

export interface CodexConfigWriteResult {
  platform: string
  homeDir: string
  configDir: string
  configPath: string
  authPath: string
  mode: CodexConfigMode
  configBackup: string | null
  authBackup: string | null
  configPreview: string
  authPreview: string
}

export interface GatewayUsageQueryInput {
  baseUrl: string
  apiKey: string
  startDate?: string
  endDate?: string
}

export interface Sub2APIDesktopBridge {
  isDesktop: true
  getEnvironment: () => Promise<CodexDesktopEnvironment>
  writeCodexConfig: (payload: CodexConfigWriteInput) => Promise<CodexConfigWriteResult>
  queryGatewayUsage: (payload: GatewayUsageQueryInput) => Promise<Record<string, any>>
  openCodexDownload: (platform?: 'mac' | 'windows' | 'linux') => Promise<string>
  openConfigDir: () => Promise<string>
  getServerUrl: () => Promise<string>
  setServerUrl: (url: string) => Promise<string>
}

declare global {
  interface Window {
    __APP_CONFIG__?: PublicSettings
    sub2apiDesktop?: Sub2APIDesktopBridge
  }
}

export {}
