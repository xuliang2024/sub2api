/**
 * 法务文本（隐私政策 / 服务条款）。
 *
 * 主站（codex.apiz.ai）与落地页（cn.apiz.ai）共用同一份内容，避免两处表述不一致。
 * 以中英双语维护；落地页只做 zh/en 两档，其它语言回落英文。
 *
 * 注意：这是产品可用的基础模板，不构成法律意见。涉及具体司法辖区、退款政策、
 * 数据处理者角色等条款，上线前应让律师过一遍。
 */

export interface LegalSection {
  heading: string
  paras?: string[]
  bullets?: string[]
}

export interface LegalDocument {
  title: string
  updated: string
  intro: string
  sections: LegalSection[]
}

export type LegalLocale = 'zh' | 'en'

const CONTACT = 'admin@apiz.ai'
const UPDATED = '2026-10-23'

const privacyZh: LegalDocument = {
  title: '隐私政策',
  updated: UPDATED,
  intro:
    '本政策说明 apiz.ai（以下称"我们"）在提供服务时收集哪些信息、如何使用与保存，以及你可以如何行使相关权利。使用我们的服务即表示你已阅读本政策。',
  sections: [
    {
      heading: '1. 我们收集的信息',
      bullets: [
        '账号信息：注册邮箱、密码哈希（第三方登录时为对方返回的账号标识）、创建时间。',
        '登录方式：使用 Google 等第三方登录时，我们会收到该平台返回的邮箱、昵称与头像地址。我们不会获得你在该平台的密码。',
        '调用元数据：模型名称、输入与输出 token 数、请求耗时、状态码、使用的 API Key、调用时间。',
        '计费与支付：充值金额、支付渠道返回的订单号与状态。我们不保存银行卡号、CVV 等支付凭据。',
        '设备与访问日志：IP 地址、User-Agent、请求路径与时间戳，用于安全审计与故障排查。',
      ],
    },
    {
      heading: '2. 关于你的模型输入与输出',
      paras: [
        '你通过 API 发送的内容会经由我们的网关转发到你选择的上游模型供应商，由该供应商按其自身条款处理。',
        '我们默认只记录第 1 条所述的元数据，不持久化保存请求与响应的正文。为排查故障或响应监管要求，我们可能在有限时间内调取请求摘要，且仅限必要人员访问。',
        '我们不会使用你的请求或响应内容训练模型，也不会将其出售给第三方。',
      ],
    },
    {
      heading: '3. 我们如何使用这些信息',
      bullets: [
        '提供与计费服务：路由请求、统计用量、生成账单。',
        '保障安全：识别滥用、盗刷、异常调用与账号风险。',
        '改进服务：以聚合、去标识的方式分析可用性与性能。',
        '通知：发送验证码、余额提醒、服务变更等必要通知。',
      ],
    },
    {
      heading: '4. 共享与第三方',
      paras: [
        '我们不会出售你的个人信息。仅在以下情形共享必要信息：',
      ],
      bullets: [
        '上游模型供应商：转发你的请求以完成推理。',
        '身份与邮件服务商：Google 登录、验证码邮件投递。',
        '支付渠道：完成充值收款与对账。',
        '云与基础设施服务商：承载服务运行与数据存储。',
        '法律要求：依有效法律程序或监管要求披露。',
      ],
    },
    {
      heading: '5. Cookie 与本地存储',
      paras: [
        '我们使用浏览器本地存储保存登录凭据与语言偏好，用于维持会话。OAuth 登录过程中会写入短期有效的校验 Cookie，登录完成后自动清除。你可以清除本地存储，但会被登出。',
      ],
    },
    {
      heading: '6. 数据保留',
      paras: [
        '账号信息在账号存续期间保留，注销后按备份周期清理。调用元数据与账单记录按财务与合规要求保留。访问日志保留期较短，仅用于安全分析。',
      ],
    },
    {
      heading: '7. 安全',
      paras: [
        '全站强制 HTTPS，密码以不可逆哈希存储，API Key 支持随时吊销。但互联网传输与存储不存在绝对安全，请妥善保管你的凭据；若怀疑泄露，请立即吊销并联系我们。',
      ],
    },
    {
      heading: '8. 你的权利',
      paras: [
        '你可以访问、更正账号信息，导出调用记录，或申请注销账号并删除个人信息（法律法规要求保留的除外）。',
        `行使上述权利请发邮件至 ${CONTACT}，我们会在合理期限内答复。`,
      ],
    },
    {
      heading: '9. 未成年人',
      paras: ['本服务不面向未满 18 周岁的用户。若发现未成年人注册，我们会停用该账号并删除相关信息。'],
    },
    {
      heading: '10. 政策更新',
      paras: ['我们可能更新本政策。重大变更会在站内公告或通过邮件通知；继续使用服务视为接受更新后的政策。'],
    },
    {
      heading: '11. 联系我们',
      paras: [`隐私相关问题请联系 ${CONTACT}。`],
    },
  ],
}

const privacyEn: LegalDocument = {
  title: 'Privacy Policy',
  updated: UPDATED,
  intro:
    'This policy explains what apiz.ai ("we", "us") collects when you use the service, how that information is used and stored, and what rights you have. By using the service you acknowledge this policy.',
  sections: [
    {
      heading: '1. Information we collect',
      bullets: [
        'Account data: email address, password hash (or the account identifier returned by a third-party identity provider), and account creation time.',
        'Sign-in data: when you sign in with Google, we receive the email address, display name and avatar URL returned by Google. We never receive your Google password.',
        'Request metadata: model name, input and output token counts, latency, status codes, the API key used, and timestamps.',
        'Billing data: top-up amounts, and the order ID and status returned by the payment provider. We do not store card numbers or CVVs.',
        'Access logs: IP address, user agent, request path and timestamp, kept for security auditing and troubleshooting.',
      ],
    },
    {
      heading: '2. Your model inputs and outputs',
      paras: [
        'Content you send through the API is forwarded by our gateway to the upstream model vendor you selected, and is handled under that vendor’s own terms.',
        'By default we record only the metadata listed above and do not persist request or response bodies. For troubleshooting or to comply with a lawful request, we may retrieve a request summary for a limited period, accessible only to staff who need it.',
        'We do not use your prompts or completions to train models, and we do not sell them.',
      ],
    },
    {
      heading: '3. How we use information',
      bullets: [
        'To provide and bill the service: routing requests, metering usage, issuing invoices.',
        'To keep the service safe: detecting abuse, credential theft, anomalous traffic and account risk.',
        'To improve the service: aggregated, de-identified analysis of availability and performance.',
        'To notify you: verification codes, low-balance alerts and material service changes.',
      ],
    },
    {
      heading: '4. Sharing and third parties',
      paras: ['We do not sell personal information. We share only what is necessary with:'],
      bullets: [
        'Upstream model vendors — to perform inference on your behalf.',
        'Identity and email providers — Google sign-in and delivery of verification emails.',
        'Payment providers — to collect and reconcile top-ups.',
        'Cloud and infrastructure providers — to host the service and store data.',
        'Authorities — where required by valid legal process.',
      ],
    },
    {
      heading: '5. Cookies and local storage',
      paras: [
        'We use browser local storage to keep your session token and language preference. The OAuth flow sets a short-lived state cookie that is cleared once sign-in completes. Clearing local storage signs you out.',
      ],
    },
    {
      heading: '6. Retention',
      paras: [
        'Account data is kept while the account is active and removed on deletion subject to backup cycles. Usage metadata and billing records are retained as required for accounting and compliance. Access logs are kept for a shorter period and used only for security analysis.',
      ],
    },
    {
      heading: '7. Security',
      paras: [
        'The service is HTTPS-only, passwords are stored as irreversible hashes, and API keys can be revoked at any time. No transmission or storage is perfectly secure: keep your credentials safe, and revoke and contact us immediately if you suspect exposure.',
      ],
    },
    {
      heading: '8. Your rights',
      paras: [
        'You may access and correct your account data, export your usage records, or request deletion of your account and personal data (except where retention is legally required).',
        `To exercise these rights, email ${CONTACT}. We respond within a reasonable period.`,
      ],
    },
    {
      heading: '9. Minors',
      paras: ['The service is not intended for anyone under 18. If we learn a minor has registered, we will disable the account and delete the associated data.'],
    },
    {
      heading: '10. Changes',
      paras: ['We may update this policy. Material changes are announced on the site or by email; continued use constitutes acceptance.'],
    },
    {
      heading: '11. Contact',
      paras: [`For privacy questions, contact ${CONTACT}.`],
    },
  ],
}

const termsZh: LegalDocument = {
  title: '服务条款',
  updated: UPDATED,
  intro:
    '本条款是你与 apiz.ai 之间关于使用本服务的协议。注册账号或调用 API 即表示你同意本条款；如不同意，请停止使用。',
  sections: [
    {
      heading: '1. 服务说明',
      paras: [
        '我们提供统一的模型调用网关：将你的请求转发至第三方模型供应商，并统一鉴权、计费与用量统计。模型本身由各供应商提供，我们不控制其输出。',
      ],
    },
    {
      heading: '2. 账号与凭据',
      bullets: [
        '你应提供真实有效的邮箱，并对账号下的一切活动负责。',
        'API Key 等同于密码，请勿提交到公开代码仓库或前端代码中；泄露后应立即吊销。',
        '严禁将账号转借他人用于规避限额或规避本条款。',
      ],
    },
    {
      heading: '3. 计费与付款',
      bullets: [
        '服务采用预付余额、按量计费。单价以站内价格页为准，可能随上游调价而变更，变更不影响已发生的用量。',
        '余额不足时调用会被拒绝，请在余额耗尽前充值。',
        '除法律要求或我们另行书面同意外，已消费金额不予退还。',
        '若因计价错误导致多扣，经核实后我们会退回差额。',
      ],
    },
    {
      heading: '4. 合理使用',
      paras: ['你不得利用本服务从事以下行为：'],
      bullets: [
        '违反适用法律法规，或生成、传播违法内容。',
        '攻击、扫描、压测我们的服务或他人系统，或绕过限流与鉴权。',
        '未经授权转售、二次分发服务，或以本服务搭建面向公众的同质产品。',
        '侵犯他人知识产权、隐私或其他合法权益。',
      ],
    },
    {
      heading: '5. 第三方与商标',
      paras: [
        'DeepSeek、Kimi、GLM、Hunyuan、MiniMax 等名称与标识归各自权利人所有。本服务与这些厂商无隶属或代言关系，仅作为技术接入方提供转发。模型的价格与能力以其官方说明为准。',
      ],
    },
    {
      heading: '6. 服务可用性',
      paras: [
        '我们尽合理努力保障可用性，但不承诺具体的 SLA。上游供应商维护、限流、政策调整或网络故障都可能导致中断或延迟。我们会尽量提前公告计划内维护。',
      ],
    },
    {
      heading: '7. 责任限制',
      paras: [
        '在法律允许的最大范围内，我们不对间接损失、利润损失、数据损失或第三方索赔承担责任。就任一事由，我们的累计赔偿责任不超过你在该事由发生前三个月实际支付的费用。',
      ],
    },
    {
      heading: '8. 终止',
      paras: [
        '你可以随时停止使用并申请注销。若你严重违反本条款、存在欺诈或危害系统安全的行为，我们可以暂停或终止服务；未消费余额的处理按适用法律执行。',
      ],
    },
    {
      heading: '9. 条款变更',
      paras: ['我们可能更新本条款，重大变更会在站内公告。变更生效后继续使用服务视为接受。'],
    },
    {
      heading: '10. 适用法律',
      paras: ['本条款适用运营方注册地的法律。因本条款产生的争议，双方应先友好协商；协商不成的，提交运营方所在地有管辖权的法院。'],
    },
    {
      heading: '11. 联系我们',
      paras: [`条款相关问题请联系 ${CONTACT}。`],
    },
  ],
}

const termsEn: LegalDocument = {
  title: 'Terms of Service',
  updated: UPDATED,
  intro:
    'These terms govern your use of apiz.ai. By creating an account or calling the API you agree to them; if you do not agree, please stop using the service.',
  sections: [
    {
      heading: '1. The service',
      paras: [
        'We operate a unified model gateway: your requests are forwarded to third-party model vendors, with authentication, billing and usage metering handled by us. The models are provided by those vendors and we do not control their output.',
      ],
    },
    {
      heading: '2. Accounts and credentials',
      bullets: [
        'Provide a valid email address and remain responsible for all activity under your account.',
        'API keys are equivalent to passwords. Never commit them to public repositories or ship them in client-side code; revoke immediately if exposed.',
        'Do not share accounts to circumvent rate limits or these terms.',
      ],
    },
    {
      heading: '3. Billing',
      bullets: [
        'The service is prepaid and metered. Unit prices are published on the pricing page and may change with upstream pricing; changes never apply retroactively to usage already incurred.',
        'Calls are rejected when your balance is insufficient — top up before it runs out.',
        'Except where required by law or agreed by us in writing, consumed amounts are non-refundable.',
        'If a pricing error overcharges you, we will refund the difference after verification.',
      ],
    },
    {
      heading: '4. Acceptable use',
      paras: ['You must not use the service to:'],
      bullets: [
        'Violate applicable law, or generate or distribute unlawful content.',
        'Attack, scan or stress our systems or anyone else’s, or bypass rate limits or authentication.',
        'Resell or redistribute the service without authorization, or build a public product that merely reskins it.',
        'Infringe intellectual property, privacy or other rights of third parties.',
      ],
    },
    {
      heading: '5. Third-party marks',
      paras: [
        'DeepSeek, Kimi, GLM, Hunyuan, MiniMax and other names and marks belong to their respective owners. We are an independent technical access point and are not affiliated with, or endorsed by, those vendors. Model capabilities and list prices are as published by them.',
      ],
    },
    {
      heading: '6. Availability',
      paras: [
        'We make reasonable efforts to keep the service available but do not warrant a specific SLA. Upstream maintenance, rate limiting, policy changes or network faults may cause interruption or latency. We announce planned maintenance where practical.',
      ],
    },
    {
      heading: '7. Limitation of liability',
      paras: [
        'To the maximum extent permitted by law, we are not liable for indirect or consequential loss, lost profits, lost data or third-party claims. Our aggregate liability for any matter will not exceed the amount you paid in the three months preceding that matter.',
      ],
    },
    {
      heading: '8. Termination',
      paras: [
        'You may stop using the service and request deletion at any time. We may suspend or terminate your access for serious breach, fraud or conduct that endangers the system; treatment of unused balance follows applicable law.',
      ],
    },
    {
      heading: '9. Changes to these terms',
      paras: ['We may update these terms and will announce material changes on the site. Continued use after the effective date constitutes acceptance.'],
    },
    {
      heading: '10. Governing law',
      paras: [
        'These terms are governed by the laws of the jurisdiction in which the operator is registered. Disputes should first be raised with us informally; unresolved disputes go to the competent courts at the operator’s location.',
      ],
    },
    {
      heading: '11. Contact',
      paras: [`Questions about these terms: ${CONTACT}.`],
    },
  ],
}

export const legalDocuments: Record<'privacy' | 'terms', Record<LegalLocale, LegalDocument>> = {
  privacy: { zh: privacyZh, en: privacyEn },
  terms: { zh: termsZh, en: termsEn },
}

/** 把任意语言代码收敛到我们维护的两档，其它语言回落英文。 */
export function normalizeLegalLocale(locale: string | undefined): LegalLocale {
  return (locale || '').toLowerCase().startsWith('zh') ? 'zh' : 'en'
}
