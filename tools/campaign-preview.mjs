#!/usr/bin/env node
/**
 * 落地页本地预览服务（不依赖 PostgreSQL / Redis / Go）。
 *
 * 用途：改完 src-campaign 后想立刻看效果，又不想起整套后端时使用。
 *   node tools/campaign-preview.mjs            # 默认 4173
 *   PORT=5000 node tools/campaign-preview.mjs
 *
 * 行为：
 *   - 静态托管 backend/internal/web/dist（先跑一次 npx vite build）
 *   - 把 /api/v1/public/campaign 用 tools/campaign-preview-data.mjs 的示例数据应答
 *   - 模拟登录与一键建 Key，方便走通「点 CTA → 注册 → 拿到 Key」的完整流程
 *   - 给 campaign.html 注入 __APP_CONFIG__（真实环境由 Go 注入）
 *
 * 注意：这是纯前端联调工具，不要用于生产。
 */

import { createServer } from 'node:http'
import { readFile, stat } from 'node:fs/promises'
import { extname, join, normalize, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { appConfig, campaign } from './campaign-preview-data.mjs'

const ROOT = resolve(fileURLToPath(new URL('..', import.meta.url)))
const DIST = join(ROOT, 'backend/internal/web/dist')
const PORT = Number(process.env.PORT || 4173)
const ORIGIN = `http://localhost:${PORT}`

// ==================== 控制台假数据 ====================
// 有状态：新建/禁用/删除 Key 在预览里能看到结果，刷新页面后重置。

const consoleProfile = {
  id: 1,
  email: 'preview@example.com',
  username: 'preview',
  role: 'user',
  balance: 12.34,
}

const consoleStats = {
  total_requests: 1284,
  today_requests: 37,
  total_tokens: 3120000,
  today_tokens: 86400,
  total_actual_cost: 4.82,
  today_actual_cost: 0.21,
  total_api_keys: 2,
  active_api_keys: 1,
  rpm: 3.2,
  tpm: 840,
  average_duration_ms: 1240,
}

let keySeq = 3
const consoleKeys = [
  {
    id: 1,
    name: 'production',
    key: 'sk-live-9f2c4a71b8e35d6c0a4f7e21c8b9d3a5',
    status: 'active',
    created_at: '2026-10-12T08:31:00Z',
    last_used_at: '2026-10-23T01:12:00Z',
    quota: 0,
    quota_used: 3.11,
    group: { id: 31, name: 'All Chinese Models' },
  },
  {
    id: 2,
    name: 'staging',
    key: 'sk-test-1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d',
    status: 'inactive',
    created_at: '2026-10-05T03:20:00Z',
    last_used_at: null,
    quota: 0,
    quota_used: 0,
    group: { id: 31, name: 'All Chinese Models' },
  },
]

const consoleModels = ['deepseek-v4-pro', 'deepseek-v4.1-flash', 'kimi-k2.8-preview', 'glm-5.3']
let logSeq = 1
const consoleLogs = Array.from({ length: 43 }, (_, index) => {
  const input = 320 + ((index * 97) % 2400)
  const output = 80 + ((index * 53) % 620)
  const model = consoleModels[index % consoleModels.length]
  const price = model === 'deepseek-v4-pro' ? 0.00000099 : 0.0000003
  return {
    id: logSeq++,
    created_at: new Date(Date.now() - index * 7 * 60 * 1000).toISOString(),
    model,
    api_key_id: index % 5 === 0 ? 2 : 1,
    input_tokens: input,
    output_tokens: output,
    cache_read_tokens: Math.round(input * 0.4),
    total_tokens: input + output + Math.round(input * 0.4),
    actual_cost: Number((output * price + input * price * 0.5).toFixed(6)),
    duration_ms: 620 + ((index * 131) % 2400),
  }
})

function paginate(items, page, pageSize) {
  const total = items.length
  const pages = Math.max(1, Math.ceil(total / pageSize))
  const start = (page - 1) * pageSize
  return {
    items: items.slice(start, start + pageSize),
    total,
    page,
    page_size: pageSize,
    pages,
  }
}


const MIME = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.woff2': 'font/woff2',
  '.map': 'application/json; charset=utf-8',
}

function json(res, data, status = 200) {
  const body = JSON.stringify({ code: 0, message: 'success', data })
  res.writeHead(status, {
    'Content-Type': 'application/json; charset=utf-8',
    'Cache-Control': 'no-store',
  })
  res.end(body)
}

async function serveStatic(res, filePath) {
  try {
    const info = await stat(filePath)
    if (!info.isFile()) {
      throw new Error('not a file')
    }
  } catch {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' })
    res.end('not found')
    return
  }
  const content = await readFile(filePath)
  res.writeHead(200, {
    'Content-Type': MIME[extname(filePath)] || 'application/octet-stream',
    'Cache-Control': 'no-store',
  })
  res.end(content)
}

const server = createServer(async (req, res) => {
  const url = new URL(req.url || '/', ORIGIN)
  const path = url.pathname

  if (path === '/api/v1/public/campaign') {
    json(res, { ...campaign, api_base_url: campaign.api_base_url || ORIGIN })
    return
  }

  // ---------- 控制台 ----------
  if (path === '/api/v1/user/profile') {
    json(res, consoleProfile)
    return
  }
  if (path === '/api/v1/usage/dashboard/stats') {
    json(res, consoleStats)
    return
  }
  if (path === '/api/v1/usage') {
    const page = Number(url.searchParams.get('page') || 1)
    const size = Number(url.searchParams.get('page_size') || 20)
    json(res, paginate(consoleLogs, page, size))
    return
  }
  if (path === '/api/v1/keys' && req.method === 'GET') {
    const page = Number(url.searchParams.get('page') || 1)
    const size = Number(url.searchParams.get('page_size') || 20)
    json(res, paginate(consoleKeys, page, size))
    return
  }
  if (path === '/api/v1/keys' && req.method === 'POST') {
    let body = ''
    for await (const chunk of req) {
      body += chunk
    }
    const name = body ? JSON.parse(body).name : 'console'
    const created = {
      id: keySeq,
      name: name || `key-${keySeq}`,
      key: `sk-new-${String(keySeq).padStart(2, '0')}a1b2c3d4e5f6a7b8c9d0e1f2a3b4`,
      status: 'active',
      created_at: new Date().toISOString(),
      last_used_at: null,
      quota: 0,
      quota_used: 0,
      group: { id: 31, name: 'All Chinese Models' },
    }
    keySeq += 1
    consoleKeys.unshift(created)
    consoleStats.total_api_keys += 1
    consoleStats.active_api_keys += 1
    json(res, created)
    return
  }
  const keyMatch = path.match(/^\/api\/v1\/keys\/(\d+)$/)
  if (keyMatch) {
    const id = Number(keyMatch[1])
    const target = consoleKeys.find((item) => item.id === id)
    if (!target) {
      res.writeHead(404, { 'Content-Type': 'application/json; charset=utf-8' })
      res.end(JSON.stringify({ code: 404, message: 'key not found' }))
      return
    }
    if (req.method === 'PUT') {
      let body = ''
      for await (const chunk of req) {
        body += chunk
      }
      const status = body ? JSON.parse(body).status : target.status
      target.status = status
      consoleStats.active_api_keys = consoleKeys.filter((item) => item.status === 'active').length
      json(res, target)
      return
    }
    if (req.method === 'DELETE') {
      const index = consoleKeys.indexOf(target)
      consoleKeys.splice(index, 1)
      consoleStats.total_api_keys = consoleKeys.length
      consoleStats.active_api_keys = consoleKeys.filter((item) => item.status === 'active').length
      json(res, null)
      return
    }
  }

  if (path === '/api/v1/auth/login' || path === '/api/v1/auth/register') {
    json(res, {
      access_token: 'preview-token',
      refresh_token: 'preview-refresh',
      expires_in: 86400,
      token_type: 'Bearer',
      user: { id: 1, email: 'preview@example.com' },
    })
    return
  }
  if (path === '/api/v1/auth/send-verify-code') {
    json(res, { message: 'sent', countdown: 60 })
    return
  }
  if (/^\/api\/v1\/campaigns\/[^/]+\/key$/.test(path)) {
    json(res, {
      api_key_id: 1,
      api_key: 'sk-preview-00000000000000000000000000000000',
      name: `promo-${campaign.slug}`,
      base_url: campaign.api_base_url || ORIGIN,
      group_id: campaign.group_id,
      group_name: campaign.group?.name,
      created_at: new Date().toISOString(),
    })
    return
  }
  if (path.startsWith('/api/')) {
    res.writeHead(404, { 'Content-Type': 'application/json; charset=utf-8' })
    res.end(JSON.stringify({ code: 404, message: 'preview: api not mocked' }))
    return
  }

  // 模拟后端按 Host 返回 campaign.html 的行为：任何前端路由都落到落地页入口
  if (path === '/' || path === '/campaign.html' || !extname(path)) {
    try {
      const html = await readFile(join(DIST, 'campaign.html'), 'utf8')
      const injected = html.replace(
        '</head>',
        `<script>window.__APP_CONFIG__=${JSON.stringify(appConfig)};</script>\n</head>`,
      )
      res.writeHead(200, { 'Content-Type': MIME['.html'], 'Cache-Control': 'no-store' })
      res.end(injected)
    } catch {
      res.writeHead(500, { 'Content-Type': 'text/plain; charset=utf-8' })
      res.end('未找到构建产物，请先运行：cd frontend && npx vite build')
    }
    return
  }

  await serveStatic(res, join(DIST, normalize(path).replace(/^([/\\])+/, '')))
})

server.listen(PORT, '127.0.0.1', () => {
  console.log(`落地页预览：${ORIGIN}`)
  console.log(`静态目录：${DIST}`)
})
