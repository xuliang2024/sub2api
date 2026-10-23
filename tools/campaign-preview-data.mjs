/**
 * 预览用示例数据 —— 模型与价格取自生产库 channel 8（分组 31「国产全模型」）。
 *
 * 这份数据同时是「后台该填什么」的参考：
 *   - allModels 与渠道定价逐条对应（USD/token，与计费同口径，前端乘 1e6 展示）
 *   - aliases 是历史别名（同价承接旧请求），对外不展示，用 hidden_models 过滤
 *   - featured 是主推卡片，其余模型进「模型广场」全量价目表
 */

/** 历史别名：deepseek-v4.1-flash 的同价别名，后台渠道页也合并显示。 */
const aliases = ['deepseek-flash', 'deepseek-v4-flash', 'deepseek-v4-1-flash-260910']

const allModels = [
  {
    model: 'deepseek-v4-pro',
    display_name: 'DeepSeek V4 Pro',
    tagline: 'Flagship reasoning for code and long documents',
    platform: 'deepseek',
    available: true,
    input_price: 0.00000033,
    output_price: 0.00000099,
    cache_read_price: 0.000000011,
    official_input_price: 0.00000132,
    official_output_price: 0.00000396,
    time_pricing: [
      { start_time: '09:00', end_time: '12:00', multiplier: 2 },
      { start_time: '14:00', end_time: '18:00', multiplier: 2 },
    ],
    time_pricing_timezone: 'Asia/Shanghai',
  },
  {
    model: 'deepseek-v4.1-flash',
    display_name: 'DeepSeek V4.1 Flash',
    tagline: 'High throughput, low latency for batch jobs',
    platform: 'deepseek',
    available: true,
    input_price: 0.000000075,
    output_price: 0.0000003,
    cache_read_price: 0.0000000015,
    official_input_price: 0.0000003,
    official_output_price: 0.0000012,
    time_pricing: [
      { start_time: '09:00', end_time: '12:00', multiplier: 2 },
      { start_time: '14:00', end_time: '18:00', multiplier: 2 },
    ],
    time_pricing_timezone: 'Asia/Shanghai',
  },
  {
    model: 'glm-5.2',
    display_name: 'GLM 5.2',
    tagline: 'General-purpose model, balanced cost',
    platform: 'zhipu',
    available: true,
    input_price: 0.000000589614,
    output_price: 0.000002063649,
    cache_read_price: 0.000000147403,
    official_input_price: 0.000001179228,
    official_output_price: 0.000004127298,
  },
  {
    model: 'glm-5.3',
    display_name: 'GLM 5.3',
    tagline: 'Balanced general-purpose model',
    platform: 'zhipu',
    available: true,
    input_price: 0.0000007,
    output_price: 0.0000022,
    cache_read_price: 0.00000013,
    official_input_price: 0.0000014,
    official_output_price: 0.0000044,
  },
  {
    model: 'glm-5.3-flash',
    display_name: 'GLM 5.3 Flash',
    tagline: 'Lowest cost in the GLM family',
    platform: 'zhipu',
    available: true,
    input_price: 0.000000058997,
    output_price: 0.00000020649,
    cache_read_price: 0.000000016962,
  },
  {
    model: 'glm-5v-turbo',
    display_name: 'GLM 5V Turbo',
    tagline: 'Vision input with turbo latency',
    platform: 'zhipu',
    available: true,
    input_price: 0.000000368732,
    output_price: 0.000001622419,
    cache_read_price: 0.000000088496,
    official_input_price: 0.000000737464,
    official_output_price: 0.000003244838,
  },
  {
    model: 'hy3',
    display_name: 'HY3',
    tagline: 'Compact model for routine tasks',
    platform: 'hunyuan',
    available: true,
    input_price: 0.000000073746,
    output_price: 0.000000294985,
    cache_read_price: 0.000000018437,
  },
  {
    model: 'hy4-preview',
    display_name: 'HY4 Preview',
    tagline: 'Next-generation preview with stronger reasoning',
    platform: 'hunyuan',
    available: true,
    input_price: 0.000000442478,
    output_price: 0.000001327434,
    cache_read_price: 0.000000022124,
  },
  {
    model: 'kimi-k2.6',
    display_name: 'Kimi K2.6',
    tagline: 'Long context, stable general performance',
    platform: 'kimi',
    available: true,
    input_price: 0.000000477,
    output_price: 0.0000019805,
    cache_read_price: 0.0000000805,
    official_input_price: 0.000000954,
    official_output_price: 0.000003961,
  },
  {
    model: 'kimi-k2.7',
    display_name: 'Kimi K2.7',
    tagline: 'Improved instruction following',
    platform: 'kimi',
    available: true,
    input_price: 0.000000479351,
    output_price: 0.00000199115,
    cache_read_price: 0.00000009587,
    official_input_price: 0.000000958702,
    official_output_price: 0.0000039823,
  },
  {
    model: 'kimi-k2.8-preview',
    display_name: 'Kimi K2.8 Preview',
    tagline: 'Very long context — read an entire repository',
    platform: 'kimi',
    available: true,
    input_price: 0.000000479351,
    output_price: 0.00000199115,
    cache_read_price: 0.00000009587,
    official_input_price: 0.000000958702,
    official_output_price: 0.0000039823,
  },
  {
    model: 'kimi-k3-1',
    display_name: 'Kimi K3.1',
    tagline: 'Top-tier reasoning, highest price band',
    platform: 'kimi',
    available: true,
    input_price: 0.000001474926,
    output_price: 0.000007374631,
    cache_read_price: 0.000000147493,
    official_input_price: 0.000002949852,
    official_output_price: 0.000014749262,
  },
  {
    model: 'mimo-v2.6-pro',
    display_name: 'MiMo V2.6 Pro',
    tagline: 'Cost-efficient for high-volume generation',
    platform: 'mimo',
    available: true,
    input_price: 0.000000221239,
    output_price: 0.000000442478,
    cache_read_price: 0.000000001844,
  },
  {
    model: 'minimax-m3',
    display_name: 'MiniMax M3',
    tagline: 'Long-form writing and summarisation',
    platform: 'minimax',
    available: true,
    input_price: 0.000000154867,
    output_price: 0.000000619469,
    cache_read_price: 0.000000030973,
  },
  {
    model: 'minimax-m3.1-flash',
    display_name: 'MiniMax M3.1 Flash',
    tagline: 'Cheapest option for classification workloads',
    platform: 'minimax',
    available: true,
    input_price: 0.000000073746,
    output_price: 0.000000294985,
    cache_read_price: 0.000000007375,
  },
  {
    model: 'space-bunny',
    display_name: 'Space Bunny',
    tagline: 'Lightweight model for chat and extraction',
    platform: 'space',
    available: true,
    input_price: 0.000000073746,
    output_price: 0.000000294985,
    cache_read_price: 0.000000007375,
  },
  {
    model: 'step-5-preview',
    display_name: 'Step 5 Preview',
    tagline: 'Multimodal preview, image input supported',
    platform: 'step',
    available: true,
    input_price: 0.000000516224,
    output_price: 0.000001474926,
    cache_read_price: 0.000000051622,
  },
]

/** 主推卡片：只挑 6 个，其余进模型广场。 */
const featured = [
  'deepseek-v4-pro',
  'deepseek-v4.1-flash',
  'kimi-k2.8-preview',
  'glm-5.3',
  'kimi-k3-1',
  'minimax-m3',
]

const visibleModels = allModels.filter((model) => !aliases.includes(model.model))

const featuredModels = featured
  .map((name) => visibleModels.find((model) => model.model === name))
  .filter((model) => Boolean(model))
  .map((model, index) => ({
    ...model,
    featured: index === 0,
    badge: index === 0 ? 'Most popular' : model.badge,
  }))

export const campaign = {
  slug: 'cn-frontier',
  name: 'CN Frontier Models — Global',
  locale: 'en',
  brand: {
    site_name: 'Frontier CN',
    primary: '#4f46e5',
    accent: '#06b6d4',
    footer: 'Local preview · sample data · not a live service',
    console_url: '',
  },
  hero: {
    eyebrow: '17 Chinese frontier models · one API',
    title: 'Chinese frontier models, at half the list price',
    subtitle:
      'DeepSeek, Kimi, GLM, Hunyuan and MiniMax behind a single OpenAI-compatible endpoint. Automatic failover across upstream accounts, so a dead provider account never becomes your outage.',
    badges: [
      'OpenAI · Anthropic · Gemini compatible',
      'Automatic upstream failover',
      '50% of official list price',
    ],
    highlights: [
      'One integration for every model — swap the model name, not the SDK',
      'Unhealthy upstream accounts are pulled out of rotation automatically',
      'Priced at 50% of each vendor’s list rate, billed per token',
    ],
  },
  content_md: `## Why teams switch

Most teams want Chinese models for cost, but not the four separate contracts, SDKs and billing accounts that come with them. This endpoint aggregates them:

- **One key, one base URL.** Same request shape for DeepSeek, Kimi, GLM, Hunyuan and MiniMax.
- **Failover built in.** When an upstream account degrades it is removed from rotation; your request goes to the next healthy one instead of failing.
- **List-price billing.** You are charged per token at 50% of the vendor's published rate — no seat fee, no minimum.

## Get started in three steps

1. Create an account and copy the API key we generate for you.
2. Set \`base_url\` to the endpoint shown on this page.
3. Call any model with your existing SDK.

\`\`\`bash
curl $BASE_URL/v1/chat/completions \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"deepseek-v4-pro","messages":[{"role":"user","content":"Hello"}]}'
\`\`\`

## Billing notes

| Item | Detail |
| --- | --- |
| Unit | USD per 1M tokens, input and output priced separately |
| Cached input | Billed at roughly 1/30 of the input rate |
| Peak hours | DeepSeek models are ×2 during 09:00–12:00 and 14:00–18:00 Asia/Shanghai |
| Minimum | None — prepaid credits, no monthly fee |

> Cached input applies automatically when your prompt prefix repeats. No extra parameter needed.`,
  models: featuredModels,
  catalog: visibleModels,
  hidden_models: aliases,
  code_samples: [
    {
      id: 'openai-python',
      label: 'OpenAI SDK (Python)',
      lang: 'python',
      code: `from openai import OpenAI

client = OpenAI(
    base_url="{{BASE_URL}}/v1",
    api_key="{{API_KEY}}",
)

resp = client.chat.completions.create(
    model="{{MODEL}}",
    messages=[{"role": "user", "content": "Explain KV cache in one sentence"}],
)
print(resp.choices[0].message.content)`,
    },
    {
      id: 'anthropic-python',
      label: 'Anthropic SDK (Python)',
      lang: 'python',
      code: `from anthropic import Anthropic

client = Anthropic(
    base_url="{{BASE_URL}}",
    api_key="{{API_KEY}}",
)

msg = client.messages.create(
    model="{{MODEL}}",
    max_tokens=512,
    messages=[{"role": "user", "content": "Hello"}],
)
print(msg.content[0].text)`,
    },
    {
      id: 'curl',
      label: 'curl',
      lang: 'bash',
      code: `curl {{BASE_URL}}/v1/chat/completions \\
  -H "Authorization: Bearer {{API_KEY}}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "{{MODEL}}",
    "messages": [{"role": "user", "content": "Hello"}]
  }'`,
    },
    {
      id: 'node',
      label: 'Node.js',
      lang: 'javascript',
      code: `import OpenAI from 'openai'

const client = new OpenAI({
  baseURL: '{{BASE_URL}}/v1',
  apiKey: '{{API_KEY}}',
})

const resp = await client.chat.completions.create({
  model: '{{MODEL}}',
  messages: [{ role: 'user', content: 'Hello' }],
})
console.log(resp.choices[0].message.content)`,
    },
  ],
  cta: {
    mode: 'key',
    label: 'Get your API key',
    sublabel: 'One key covers every model. Free credits on sign-up, no card required.',
    topup_url: '#',
  },
  trial: { balance: 1, note: 'New accounts start with $1 in credits — enough for ~2M input tokens.' },
  seo: {
    title: 'Chinese frontier models · one OpenAI-compatible API',
    description:
      'DeepSeek, Kimi, GLM, Hunyuan and MiniMax through a single endpoint at 50% of list price. OpenAI, Anthropic and Gemini compatible.',
    keywords: ['DeepSeek API', 'Kimi API', 'GLM API', 'MiniMax API', 'OpenAI compatible'],
  },
  translations: {
    zh: {
      hero_title: '中国前沿模型，官方价五折起',
      hero_subtitle:
        'DeepSeek、Kimi、GLM、混元、MiniMax 收在一个 OpenAI 兼容地址后面。上游账号自动故障转移，单个账号挂掉不会变成你的故障。',
      cta_label: '立即领取 API Key',
      content_md: `## 为什么换过来

想要国产模型的成本，但不想同时签四份合同、接四套 SDK、对四个账单——这里把它们聚合在一起：

- **一把 Key、一个地址。** DeepSeek、Kimi、GLM、混元、MiniMax 用同一个请求格式。
- **自带故障转移。** 上游账号异常会被自动移出轮转，请求继续打到下一个健康账号。
- **按刊例价五折计费。** 按 token 计费，没有席位费，没有最低消费。

## 三步接入

1. 注册并复制系统为你生成的 API Key
2. 把 base_url 指向本页显示的地址
3. 用现有 SDK 调用任意模型`,
    },
    ja: {
      hero_title: '中国のフロンティアモデルを、定価の半額で',
      hero_subtitle:
        'DeepSeek、Kimi、GLM、Hunyuan、MiniMax を 1 つの OpenAI 互換エンドポイントに集約。上流アカウントの障害は自動で切り離され、単一アカウントの停止がサービス停止になりません。',
      cta_label: 'APIキーを取得',
      content_md: `## 乗り換える理由

中国製モデルのコストは欲しいが、4 つの契約・4 つの SDK・4 つの請求は要らない。ここでは 1 つにまとめています。

- **キーは 1 つ、URL も 1 つ。** モデル名を変えるだけで、SDK はそのまま。
- **フェイルオーバー内蔵。** 異常な上流アカウントはローテーションから外れます。
- **定価の 50% で課金。** トークン単位、席料なし、最低利用額なし。`,
    },
    ko: {
      hero_title: '중국 프런티어 모델을 정가의 절반으로',
      hero_subtitle:
        'DeepSeek, Kimi, GLM, Hunyuan, MiniMax를 하나의 OpenAI 호환 엔드포인트로. 업스트림 계정 장애는 자동으로 제외되어 단일 계정 문제가 서비스 중단으로 이어지지 않습니다.',
      cta_label: 'API 키 받기',
      content_md: `## 옮겨오는 이유

중국 모델의 비용은 원하지만 네 개의 계약과 SDK, 청구서는 원하지 않습니다. 여기서 하나로 묶습니다.

- **키 하나, URL 하나.** 모델 이름만 바꾸면 기존 SDK 그대로 사용.
- **페일오버 내장.** 비정상 업스트림 계정은 자동으로 로테이션에서 빠집니다.
- **정가의 50% 과금.** 토큰 단위, 좌석료 없음, 최소 사용액 없음.`,
    },
    th: {
      hero_title: 'โมเดลจีนระดับแนวหน้า ในราคาครึ่งเดียวของราคาปกติ',
      hero_subtitle:
        'DeepSeek, Kimi, GLM, Hunyuan และ MiniMax รวมอยู่ในปลายทาง OpenAI-compatible เดียว ระบบสลับไปยังบัญชีต้นทางที่พร้อมใช้งานอัตโนมัติ บัญชีเดียวล่มจึงไม่กลายเป็นเหตุขัดข้องของคุณ',
      cta_label: 'รับ API key',
      content_md: `## ทำไมทีมถึงย้ายมา

อยากได้ต้นทุนของโมเดลจีน แต่ไม่อยากได้สัญญาสี่ฉบับ SDK สี่ชุด และใบเรียกเก็บสี่ใบ — ที่นี่รวมไว้ในที่เดียว

- **คีย์เดียว URL เดียว** เปลี่ยนแค่ชื่อโมเดล ใช้ SDK เดิมได้ทันที
- **มีระบบสลับอัตโนมัติ** บัญชีต้นทางที่ผิดปกติจะถูกนำออกจากการหมุนเวียน
- **คิดราคา 50% ของราคาปกติ** คิดตามโทเคน ไม่มีค่าที่นั่ง ไม่มีขั้นต่ำ`,
    },
  },
  group: { id: 31, name: 'All Chinese Models', platform: 'openai', rate_multiplier: 1 },
  group_id: 31,
  api_base_url: '',
}

export const appConfig = {
  registration_enabled: true,
  // 对齐线上：注册需要邮箱验证码（预览里发码接口是 mock，任何 6 位码都能通过）
  email_verify_enabled: true,
  promo_code_enabled: true,
  turnstile_enabled: false,
  site_name: 'Frontier CN',
}
