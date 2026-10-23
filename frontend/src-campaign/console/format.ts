/**
 * 控制台展示用的小工具。
 *
 * 金额与数字统一在这里格式化，避免三处面板各写一套导致口径不一致。
 */

export function formatMoney(value?: number | null, digits = 2): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return '$0.00'
  }
  return `$${value.toFixed(digits)}`
}

/** 小额费用保留更多位数：$0.0012 直接 toFixed(2) 会变成 $0.00。 */
export function formatCost(value?: number | null): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return '$0.00'
  }
  if (value === 0) {
    return '$0.00'
  }
  if (Math.abs(value) < 0.01) {
    return `$${value.toFixed(4)}`
  }
  return `$${value.toFixed(2)}`
}

export function formatNumber(value?: number | null): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return '0'
  }
  return value.toLocaleString('en-US')
}

/** 固定成 YYYY-MM-DD HH:mm，跨语言展示一致，不依赖运行环境 locale。 */
export function formatDateTime(value?: string | null): string {
  if (!value) {
    return '—'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  const pad = (input: number) => String(input).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(
    date.getHours(),
  )}:${pad(date.getMinutes())}`
}

/** sk-abcd…wxyz：只露头尾，够识别又不至于泄露。 */
export function maskKey(key?: string | null): string {
  if (!key) {
    return '—'
  }
  if (key.length <= 12) {
    return key
  }
  return `${key.slice(0, 7)}…${key.slice(-4)}`
}
