/** 落地页展示用的小工具：价格换算与示例代码占位替换。 */

/** 价格单位是「美元 / token」，页面统一展示成「美元 / 百万 token」。 */
export function formatPricePerMillion(value?: number): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return '—'
  }
  const perMillion = value * 1_000_000
  if (perMillion === 0) {
    return '$0'
  }
  // 极小值（低于百万分之一美元）才用科学计数。
  const abs = Math.abs(perMillion)
  if (abs < 0.000001) {
    return `$${perMillion.toExponential(1)}`
  }
  // 按量级给精度：$2.06 比 $2.063649 好读，$0.0015 也不会被抹成 $0。
  let digits = 6
  if (abs >= 1) {
    digits = 2
  } else if (abs >= 0.1) {
    digits = 3
  } else if (abs >= 0.01) {
    digits = 4
  } else if (abs >= 0.001) {
    digits = 5
  }
  return `$${trimZeros(perMillion.toFixed(digits))}`
}

/** 按次计费的价格（美元 / 次）。 */
export function formatPerRequest(value?: number): string {
  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return '—'
  }
  return `$${trimZeros(value.toFixed(4))}`
}

/** 倍率展示：1 → 1x，1.5 → 1.5x。 */
export function formatMultiplier(value?: number): string {
  if (typeof value !== 'number' || !Number.isFinite(value) || value <= 0) {
    return ''
  }
  return `${trimZeros(value.toFixed(2))}x`
}

/**
 * 相对官方价的折扣百分比（整数）。
 *
 * 只在「有官方价且确实更便宜」时返回正数，否则返回 0——宁可不展示，
 * 也不要凭空造一个折扣数字。
 */
export function discountPercent(ours?: number, official?: number): number {
  if (!ours || !official || official <= 0) {
    return 0
  }
  const saved = (1 - ours / official) * 100
  if (!Number.isFinite(saved) || saved < 1) {
    return 0
  }
  return Math.round(saved)
}

function trimZeros(input: string): string {
  return input.replace(/\.?0+$/, '')
}

export interface SampleContext {
  baseUrl?: string
  apiKey?: string
  model?: string
}

const PLACEHOLDER_API_KEY = 'YOUR_API_KEY'

/**
 * 替换示例代码里的占位符。
 *
 * 未拿到 Key 时用 YOUR_API_KEY 占位，让用户先看到完整可读的调用方式；
 * 拿到 Key 后再替换成真实值，复制即用。
 */
export function fillSample(code: string, context: SampleContext): string {
  return code
    .split('{{BASE_URL}}')
    .join(context.baseUrl || 'https://your-domain.example')
    .split('{{API_KEY}}')
    .join(context.apiKey || PLACEHOLDER_API_KEY)
    .split('{{MODEL}}')
    .join(context.model || '')
}

/** 首页主推模型（featured 优先，其次按后端返回顺序），用于示例代码默认值。 */
export function primaryModelName(models: { model: string; featured?: boolean }[]): string {
  if (!models.length) {
    return ''
  }
  const featured = models.find((item) => item.featured)
  return (featured || models[0]).model
}

export async function copyText(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
      return true
    }
  } catch {
    // 非安全上下文（http）下 clipboard 不可用，走下面的兜底。
  }
  try {
    const area = document.createElement('textarea')
    area.value = text
    area.style.position = 'fixed'
    area.style.opacity = '0'
    document.body.appendChild(area)
    area.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(area)
    return ok
  } catch {
    return false
  }
}
