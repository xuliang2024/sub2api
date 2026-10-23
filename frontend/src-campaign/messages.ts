/**
 * 落地页界面文案（UI chrome）。
 *
 * 默认 en：这套页面的主要场景是国产模型卖给海外用户。
 * 支持 en / zh / ja / ko / th 五种界面语言，非英文条目缺失时自动回落英文，
 * 因此新增语言不必一次补齐全部 key。
 *
 * 注意：标题、说明、模型卖点这类「内容」不在这里，它们来自 campaign 配置
 * （含按语言的覆盖，见 types 里的 translations）。
 */

import { inject, type InjectionKey } from 'vue'

export type CampaignLocale = 'en' | 'zh' | 'ja' | 'ko' | 'th'

type Dict = Record<string, string>

export interface LocaleOption {
  code: CampaignLocale
  /** 语言自称，下拉框里不做翻译。 */
  label: string
  /** <html lang> 用值 */
  htmlLang: string
}

export const SUPPORTED_LOCALES: LocaleOption[] = [
  { code: 'en', label: 'English', htmlLang: 'en' },
  { code: 'zh', label: '中文', htmlLang: 'zh-CN' },
  { code: 'ja', label: '日本語', htmlLang: 'ja' },
  { code: 'ko', label: '한국어', htmlLang: 'ko' },
  { code: 'th', label: 'ไทย', htmlLang: 'th' },
]

const en: Dict = {
  'nav.login': 'Log in',
  'nav.logout': 'Log out',
  'nav.console': 'Console',
  'nav.models': 'Models',
  'nav.pricing': 'Pricing',
  'nav.docs': 'Docs',
  'nav.language': 'Language',
  'tabs.overview': 'Overview',
  'tabs.catalog': 'Model catalog',
  'tabs.console': 'Console',
  'console.overview': 'Overview',
  'console.keys': 'API keys',
  'console.logs': 'Usage logs',
  'console.balance': 'Balance',
  'console.todayRequests': 'Requests today',
  'console.todayCost': 'Spend today',
  'console.totalRequests': 'Total requests',
  'console.todayTokens': 'Tokens today',
  'console.activeKeys': 'Active keys',
  'console.newKey': 'New key',
  'console.keyName': 'Name',
  'console.create': 'Create',
  'console.creating': 'Creating…',
  'console.copy': 'Copy',
  'console.copied': 'Copied',
  'console.disable': 'Disable',
  'console.enable': 'Enable',
  'console.remove': 'Delete',
  'console.confirmRemove': 'Delete this key? Calls using it will stop working.',
  'console.statusActive': 'Active',
  'console.statusInactive': 'Disabled',
  'console.statusOther': 'Limited',
  'console.emptyKeys': 'No keys yet. Create one to start calling.',
  'console.createdAt': 'Created',
  'console.lastUsed': 'Last used',
  'console.never': 'Never',
  'console.time': 'Time',
  'console.model': 'Model',
  'console.key': 'Key',
  'console.tokens': 'Tokens',
  'console.cost': 'Cost',
  'console.emptyLogs': 'No calls recorded yet.',
  'console.page': 'Page {page} / {pages}',
  'console.prev': 'Previous',
  'console.next': 'Next',
  'console.loadFailed': 'Could not load data',
  'console.today': 'Today',
  'models.viewAll': 'View all {count} models',

  'hero.secondary': 'View models & pricing',
  'hero.trialFallback': 'New accounts get free credits to start.',

  'models.title': 'Models & pricing',
  'models.subtitle':
    "Per 1M tokens, with your plan's rate already applied. You are billed on actual input and output tokens.",
  'models.group': 'Plan: {name}',
  'models.unavailable': 'Coming soon',
  'models.input': 'Input',
  'models.output': 'Output',
  'models.cacheRead': 'Cached input',
  'models.perRequest': 'Per call',
  'models.perCallUnit': '{price} / call',
  'models.peak': 'Peak {periods}',
  'models.officialPrice': 'List price',
  'models.ourPrice': 'You pay',
  'models.save': 'Save {percent}%',
  'models.offpeak': 'Off-peak',
  'models.peakShort': 'Peak',
  'models.empty': 'Model list is being configured. Check back shortly.',

  'catalog.title': 'All models in this plan',
  'catalog.subtitle': 'The full catalogue and price list, including models not featured above.',
  'catalog.search': 'Search model',
  'catalog.model': 'Model',
  'catalog.discount': 'Discount',
  'catalog.count': '{count} models',
  'catalog.empty': 'No models match your search.',

  'usage.title': 'How to use',

  'compat.title': 'One endpoint, three protocols',
  'compat.openai': 'OpenAI Chat Completions',
  'compat.anthropic': 'Anthropic Messages',
  'compat.responses': 'OpenAI Responses',
  'compat.gemini': 'Gemini',
  'compat.note': 'Same key, same base URL. Switch protocols without switching vendors.',

  'integration.title': 'Integration',
  'integration.subtitle':
    'Point your existing SDK at our base URL. No client rewrite, no new provider abstraction.',
  'integration.copy': 'Copy',
  'integration.copied': 'Copied',

  'agents.title': 'Coding agents',
  'agents.subtitle':
    'Codex, Claude Code, Aider and Cline run on this key too. Copy a snippet, start the agent.',
  'agents.note.codex':
    'Reads providers from ~/.codex/config.toml, using the Chat Completions wire format — tools work as usual.',
  'agents.note.claude':
    'Speaks the Anthropic Messages protocol. Three env vars, no config file to edit.',
  'agents.note.aider': 'OpenAI protocol in two env vars; name the model on the command line.',
  'agents.note.cline': 'In the VS Code extension pick OpenAI Compatible and fill these four fields.',

  'cta.defaultLabel': 'Get your API key',
  'cta.defaultSublabel': 'Sign up, get a key instantly, start calling in minutes.',

  'auth.loginTitle': 'Log in to continue',
  'auth.registerTitle': 'Create an account',
  'auth.haveAccount': 'I have an account',
  'auth.newUser': 'New here',
  'auth.email': 'Email',
  'auth.password': 'Password',
  'auth.passwordHint': 'At least 6 characters',
  'auth.verifyCode': 'Email code',
  'auth.verifyPlaceholder': '6-digit code',
  'auth.send': 'Send',
  'auth.sendIn': '{seconds}s',
  'auth.login': 'Log in',
  'auth.register': 'Create account',
  'auth.processing': 'Working…',
  'auth.close': 'Close',
  'auth.registrationClosed':
    'Sign-ups are closed on this site. Please log in with an existing account.',
  'auth.captchaNotice':
    'Human verification is enabled. If sign-up fails here, create the account on the main site and log in again.',
  'auth.emailRequired': 'Enter your email and password',
  'auth.verifyRequired': 'Enter the email code',
  'auth.loginFailed': 'Login failed',
  'auth.registerFailed': 'Sign-up failed',
  'auth.codeSentFailed': 'Could not send the code',

  'key.title': 'Your API key',
  'key.desc': 'A dedicated key was created for this plan. Clicking again returns the same key.',
  'key.group': ' (plan: {name})',
  'key.apiKey': 'API key',
  'key.baseUrl': 'Base URL',
  'key.copy': 'Copy',
  'key.copied': 'Copied',
  'key.hint':
    'Paste it into any OpenAI-compatible client. Top up in the console when you need more.',
  'key.topup': 'Add credits',
  'key.continue': 'Keep browsing',
  'key.retry': 'Retry',

  'error.title': 'This page is unavailable',
  'error.reload': 'Reload',
  'error.loadFailed': 'Could not load this page',
  'error.claimFailed': 'Could not create your key',
  'error.network': 'Network error. Check your connection and try again.',
  'error.timeout': 'Request timed out. Please try again.',
  'error.generic': 'Something went wrong. Please try again.',

  'footer.default': 'All product names are trademarks of their respective owners.',
}

const zh: Dict = {
  'nav.login': '登录',
  'nav.logout': '退出登录',
  'nav.console': '控制台',
  'nav.models': '模型',
  'nav.pricing': '价格',
  'nav.docs': '文档',
  'nav.language': '语言',
  'tabs.overview': '概览',
  'tabs.catalog': '模型广场',
  'tabs.console': '控制台',
  'console.overview': '概览',
  'console.keys': 'API Key',
  'console.logs': '调用日志',
  'console.balance': '余额',
  'console.todayRequests': '今日请求',
  'console.todayCost': '今日消耗',
  'console.totalRequests': '累计请求',
  'console.todayTokens': '今日 Tokens',
  'console.activeKeys': '可用 Key',
  'console.newKey': '新建 Key',
  'console.keyName': '名称',
  'console.create': '创建',
  'console.creating': '创建中…',
  'console.copy': '复制',
  'console.copied': '已复制',
  'console.disable': '禁用',
  'console.enable': '启用',
  'console.remove': '删除',
  'console.confirmRemove': '确定删除该 Key？使用它的调用会立即失效。',
  'console.statusActive': '可用',
  'console.statusInactive': '已禁用',
  'console.statusOther': '受限',
  'console.emptyKeys': '还没有 Key，新建一把即可开始调用。',
  'console.createdAt': '创建于',
  'console.lastUsed': '最近调用',
  'console.never': '从未',
  'console.time': '时间',
  'console.model': '模型',
  'console.key': 'Key',
  'console.tokens': 'Tokens',
  'console.cost': '费用',
  'console.emptyLogs': '还没有调用记录。',
  'console.page': '第 {page} / {pages} 页',
  'console.prev': '上一页',
  'console.next': '下一页',
  'console.loadFailed': '数据加载失败',
  'console.today': '今日',
  'models.viewAll': '查看全部 {count} 个模型',

  'hero.secondary': '查看模型与价格',
  'hero.trialFallback': '新账号注册即送试用额度。',

  'models.title': '模型与价格',
  'models.subtitle': '按每百万 token 计价，已含你的方案倍率；实际扣费以调用时的上下行 token 为准。',
  'models.group': '当前分组：{name}',
  'models.unavailable': '待开通',
  'models.input': '输入',
  'models.output': '输出',
  'models.cacheRead': '缓存命中',
  'models.perRequest': '按次调用',
  'models.perCallUnit': '{price} / 次',
  'models.peak': '忙时 {periods}',
  'models.officialPrice': '官方价',
  'models.ourPrice': '你的价格',
  'models.save': '省 {percent}%',
  'models.offpeak': '谷时',
  'models.peakShort': '忙时',
  'models.empty': '模型清单正在配置中，请稍后再来。',

  'catalog.title': '本分组全部模型',
  'catalog.subtitle': '完整模型与价目表，包含上方未展示的模型。',
  'catalog.search': '搜索模型',
  'catalog.model': '模型',
  'catalog.discount': '优惠',
  'catalog.count': '{count} 个模型',
  'catalog.empty': '没有匹配的模型。',

  'usage.title': '使用说明',

  'compat.title': '一个地址，三套协议',
  'compat.openai': 'OpenAI Chat Completions',
  'compat.anthropic': 'Anthropic Messages',
  'compat.responses': 'OpenAI Responses',
  'compat.gemini': 'Gemini',
  'compat.note': '同一把 Key、同一个 Base URL，换协议不用换供应商。',

  'integration.title': '接入示例',
  'integration.subtitle': '把现有 SDK 的 Base URL 指向这里即可，不需要改客户端、不需要新的抽象层。',
  'integration.copy': '复制',
  'integration.copied': '已复制',

  'agents.title': '编码 Agent 接入',
  'agents.subtitle': 'Codex、Claude Code、Aider、Cline 用同一把 Key。复制配置，直接启动。',
  'agents.note.codex':
    '从 ~/.codex/config.toml 读取 provider，走 Chat Completions 协议，工具调用照常。',
  'agents.note.claude': '走 Anthropic Messages 协议，三个环境变量即可，不用改配置文件。',
  'agents.note.aider': '两个环境变量接入 OpenAI 协议，模型在命令行里指定。',
  'agents.note.cline': 'VS Code 插件里选 OpenAI Compatible，填这四项即可。',

  'cta.defaultLabel': '立即领取 API Key',
  'cta.defaultSublabel': '注册即得 Key，几分钟内跑通第一个请求。',

  'auth.loginTitle': '登录后继续',
  'auth.registerTitle': '创建账号',
  'auth.haveAccount': '已有账号',
  'auth.newUser': '新用户注册',
  'auth.email': '邮箱',
  'auth.password': '密码',
  'auth.passwordHint': '至少 6 位',
  'auth.verifyCode': '邮箱验证码',
  'auth.verifyPlaceholder': '6 位验证码',
  'auth.send': '发送',
  'auth.sendIn': '{seconds}s',
  'auth.login': '登录',
  'auth.register': '注册并继续',
  'auth.processing': '处理中…',
  'auth.close': '关闭',
  'auth.registrationClosed': '当前站点暂未开放注册，请使用已有账号登录。',
  'auth.captchaNotice': '本站已开启人机校验，如注册失败请前往主站注册后返回登录。',
  'auth.emailRequired': '请填写邮箱与密码',
  'auth.verifyRequired': '请填写邮箱验证码',
  'auth.loginFailed': '登录失败',
  'auth.registerFailed': '注册失败',
  'auth.codeSentFailed': '验证码发送失败',

  'key.title': '你的 API Key',
  'key.desc': '已为本方案创建专用 Key，重复点击不会重复发放。',
  'key.group': '（分组：{name}）',
  'key.apiKey': 'API Key',
  'key.baseUrl': 'Base URL',
  'key.copy': '复制',
  'key.copied': '已复制',
  'key.hint': '粘贴到任意 OpenAI 兼容客户端即可使用；额度不足时到控制台充值。',
  'key.topup': '去充值',
  'key.continue': '继续浏览',
  'key.retry': '重试',

  'error.title': '页面暂时不可用',
  'error.reload': '重新加载',
  'error.loadFailed': '页面加载失败',
  'error.claimFailed': '创建 Key 失败',
  'error.network': '网络异常，请检查网络后重试。',
  'error.timeout': '请求超时，请稍后重试。',
  'error.generic': '操作失败，请稍后重试。',

  'footer.default': '各产品名称归其各自所有者所有。',
}

const ja: Dict = {
  'nav.login': 'ログイン',
  'nav.logout': 'ログアウト',
  'nav.console': 'コンソール',
  'nav.models': 'モデル',
  'nav.pricing': '料金',
  'nav.docs': 'ドキュメント',
  'nav.language': '言語',
  'tabs.overview': '概要',
  'tabs.catalog': 'モデル一覧',
  'tabs.console': 'コンソール',
  'console.overview': '概要',
  'console.keys': 'APIキー',
  'console.logs': '利用ログ',
  'console.balance': '残高',
  'console.todayRequests': '本日のリクエスト',
  'console.todayCost': '本日の利用額',
  'console.totalRequests': '累計リクエスト',
  'console.todayTokens': '本日トークン',
  'console.activeKeys': '有効キー',
  'console.newKey': 'キーを作成',
  'console.keyName': '名前',
  'console.create': '作成',
  'console.creating': '作成中…',
  'console.copy': 'コピー',
  'console.copied': 'コピーしました',
  'console.disable': '無効化',
  'console.enable': '有効化',
  'console.remove': '削除',
  'console.confirmRemove': 'このキーを削除しますか？利用中の呼び出しは停止します。',
  'console.statusActive': '有効',
  'console.statusInactive': '無効',
  'console.statusOther': '制限',
  'console.emptyKeys': 'キーがありません。作成して開始してください。',
  'console.createdAt': '作成日',
  'console.lastUsed': '最終利用',
  'console.never': 'なし',
  'console.time': '日時',
  'console.model': 'モデル',
  'console.key': 'キー',
  'console.tokens': 'トークン',
  'console.cost': '料金',
  'console.emptyLogs': '呼び出し履歴がありません。',
  'console.page': '{page} / {pages} ページ',
  'console.prev': '前へ',
  'console.next': '次へ',
  'console.loadFailed': 'データを読み込めませんでした',
  'console.today': '本日',
  'models.viewAll': '全 {count} モデルを見る',

  'hero.secondary': 'モデルと料金を見る',
  'hero.trialFallback': '新規アカウントには無料クレジットを付与します。',

  'models.title': 'モデルと料金',
  'models.subtitle':
    '100万トークンあたりの価格です。プランのレートは反映済みで、実際の入出力トークンに応じて課金されます。',
  'models.group': 'プラン：{name}',
  'models.unavailable': '準備中',
  'models.input': '入力',
  'models.output': '出力',
  'models.cacheRead': 'キャッシュ入力',
  'models.perRequest': '1回あたり',
  'models.perCallUnit': '{price} / 回',
  'models.peak': 'ピーク {periods}',
  'models.officialPrice': '定価',
  'models.ourPrice': 'ご利用価格',
  'models.save': '{percent}% お得',
  'models.offpeak': 'オフピーク',
  'models.peakShort': 'ピーク',
  'models.empty': 'モデル一覧を準備中です。しばらくしてからご確認ください。',

  'catalog.title': 'このプランの全モデル',
  'catalog.subtitle': '上に表示していないモデルを含む、全モデルと料金の一覧です。',
  'catalog.search': 'モデルを検索',
  'catalog.model': 'モデル',
  'catalog.discount': '割引',
  'catalog.count': '{count} モデル',
  'catalog.empty': '一致するモデルがありません。',

  'usage.title': '使い方',

  'compat.title': '1つのエンドポイント、3つのプロトコル',
  'compat.openai': 'OpenAI Chat Completions',
  'compat.anthropic': 'Anthropic Messages',
  'compat.responses': 'OpenAI Responses',
  'compat.gemini': 'Gemini',
  'compat.note': '同じキー、同じベースURL。プロトコルを変えてもベンダーは変わりません。',

  'integration.title': '導入手順',
  'integration.subtitle':
    '既存のSDKのベースURLを差し替えるだけ。クライアントの書き換えは不要です。',
  'integration.copy': 'コピー',
  'integration.copied': 'コピーしました',

  'agents.title': 'コーディングエージェント',
  'agents.subtitle':
    'Codex、Claude Code、Aider、Cline も同じキーで動きます。設定をコピーして起動するだけです。',
  'agents.note.codex':
    '~/.codex/config.toml からプロバイダーを読み込み、Chat Completions で接続します。',
  'agents.note.claude':
    'Anthropic Messages プロトコル。環境変数3つだけで、設定ファイルの編集は不要です。',
  'agents.note.aider': '環境変数2つで OpenAI プロトコルに接続し、モデルはコマンドラインで指定します。',
  'agents.note.cline': 'VS Code 拡張機能で OpenAI Compatible を選び、この4項目を入力します。',

  'cta.defaultLabel': 'APIキーを取得',
  'cta.defaultSublabel': '登録後すぐにキーを発行。数分で最初の呼び出しまで。',

  'auth.loginTitle': '続行するにはログイン',
  'auth.registerTitle': 'アカウントを作成',
  'auth.haveAccount': 'アカウントをお持ちの方',
  'auth.newUser': '新規登録',
  'auth.email': 'メールアドレス',
  'auth.password': 'パスワード',
  'auth.passwordHint': '6文字以上',
  'auth.verifyCode': '確認コード',
  'auth.verifyPlaceholder': '6桁のコード',
  'auth.send': '送信',
  'auth.sendIn': '{seconds}秒',
  'auth.login': 'ログイン',
  'auth.register': '登録して続行',
  'auth.processing': '処理中…',
  'auth.close': '閉じる',
  'auth.registrationClosed':
    '現在このサイトでは新規登録を受け付けていません。既存のアカウントでログインしてください。',
  'auth.captchaNotice':
    '人間確認が有効です。登録に失敗する場合はメインサイトで登録したうえでログインしてください。',
  'auth.emailRequired': 'メールアドレスとパスワードを入力してください',
  'auth.verifyRequired': '確認コードを入力してください',
  'auth.loginFailed': 'ログインに失敗しました',
  'auth.registerFailed': '登録に失敗しました',
  'auth.codeSentFailed': 'コードを送信できませんでした',

  'key.title': 'APIキー',
  'key.desc': 'このプラン専用のキーを発行しました。再度クリックしても同じキーが返ります。',
  'key.group': '（プラン：{name}）',
  'key.apiKey': 'APIキー',
  'key.baseUrl': 'ベースURL',
  'key.copy': 'コピー',
  'key.copied': 'コピーしました',
  'key.hint':
    'OpenAI互換のクライアントに貼り付けるだけで使えます。残高不足時はコンソールでチャージしてください。',
  'key.topup': 'チャージ',
  'key.continue': '閲覧を続ける',
  'key.retry': '再試行',

  'error.title': 'このページは現在ご利用いただけません',
  'error.reload': '再読み込み',
  'error.loadFailed': 'ページを読み込めませんでした',
  'error.claimFailed': 'キーの発行に失敗しました',
  'error.network': 'ネットワークエラーです。接続を確認して再試行してください。',
  'error.timeout': 'タイムアウトしました。もう一度お試しください。',
  'error.generic': '問題が発生しました。もう一度お試しください。',

  'footer.default': '各製品名は各社の商標または登録商標です。',
}

const ko: Dict = {
  'nav.login': '로그인',
  'nav.logout': '로그아웃',
  'nav.console': '콘솔',
  'nav.models': '모델',
  'nav.pricing': '가격',
  'nav.docs': '문서',
  'nav.language': '언어',
  'tabs.overview': '개요',
  'tabs.catalog': '모델 목록',
  'tabs.console': '콘솔',
  'console.overview': '개요',
  'console.keys': 'API 키',
  'console.logs': '호출 로그',
  'console.balance': '잔액',
  'console.todayRequests': '오늘 요청',
  'console.todayCost': '오늘 사용액',
  'console.totalRequests': '누적 요청',
  'console.todayTokens': '오늘 토큰',
  'console.activeKeys': '활성 키',
  'console.newKey': '새 키',
  'console.keyName': '이름',
  'console.create': '생성',
  'console.creating': '생성 중…',
  'console.copy': '복사',
  'console.copied': '복사됨',
  'console.disable': '비활성화',
  'console.enable': '활성화',
  'console.remove': '삭제',
  'console.confirmRemove': '이 키를 삭제할까요? 사용 중인 호출이 중단됩니다.',
  'console.statusActive': '활성',
  'console.statusInactive': '비활성',
  'console.statusOther': '제한',
  'console.emptyKeys': '키가 없습니다. 만들어 시작하세요.',
  'console.createdAt': '생성일',
  'console.lastUsed': '최근 사용',
  'console.never': '없음',
  'console.time': '시간',
  'console.model': '모델',
  'console.key': '키',
  'console.tokens': '토큰',
  'console.cost': '비용',
  'console.emptyLogs': '호출 기록이 없습니다.',
  'console.page': '{page} / {pages} 페이지',
  'console.prev': '이전',
  'console.next': '다음',
  'console.loadFailed': '데이터를 불러오지 못했습니다',
  'console.today': '오늘',
  'models.viewAll': '전체 {count}개 모델 보기',

  'hero.secondary': '모델 및 가격 보기',
  'hero.trialFallback': '신규 계정에는 무료 크레딧이 지급됩니다.',

  'models.title': '모델 및 가격',
  'models.subtitle':
    '100만 토큰당 가격이며 플랜 배율이 이미 적용되어 있습니다. 실제 입출력 토큰 기준으로 과금됩니다.',
  'models.group': '플랜: {name}',
  'models.unavailable': '준비 중',
  'models.input': '입력',
  'models.output': '출력',
  'models.cacheRead': '캐시 입력',
  'models.perRequest': '호출당',
  'models.perCallUnit': '{price} / 회',
  'models.peak': '피크 {periods}',
  'models.officialPrice': '정가',
  'models.ourPrice': '이용 요금',
  'models.save': '{percent}% 절약',
  'models.offpeak': '비수기',
  'models.peakShort': '성수기',
  'models.empty': '모델 목록을 준비 중입니다. 잠시 후 다시 확인해 주세요.',

  'catalog.title': '이 플랜의 전체 모델',
  'catalog.subtitle': '위에 소개되지 않은 모델을 포함한 전체 모델과 가격 목록입니다.',
  'catalog.search': '모델 검색',
  'catalog.model': '모델',
  'catalog.discount': '할인',
  'catalog.count': '{count}개 모델',
  'catalog.empty': '일치하는 모델이 없습니다.',

  'usage.title': '사용 방법',

  'compat.title': '하나의 엔드포인트, 세 가지 프로토콜',
  'compat.openai': 'OpenAI Chat Completions',
  'compat.anthropic': 'Anthropic Messages',
  'compat.responses': 'OpenAI Responses',
  'compat.gemini': 'Gemini',
  'compat.note': '동일한 키, 동일한 Base URL. 프로토콜만 바꾸면 됩니다.',

  'integration.title': '연동',
  'integration.subtitle':
    '기존 SDK의 Base URL만 변경하면 됩니다. 클라이언트 수정이나 새로운 추상화가 필요 없습니다.',
  'integration.copy': '복사',
  'integration.copied': '복사됨',

  'agents.title': '코딩 에이전트',
  'agents.subtitle':
    'Codex, Claude Code, Aider, Cline도 같은 키로 동작합니다. 설정을 복사해 바로 실행하세요.',
  'agents.note.codex':
    '~/.codex/config.toml에서 프로바이더를 읽고 Chat Completions로 연결합니다.',
  'agents.note.claude':
    'Anthropic Messages 프로토콜. 환경 변수 3개면 되고 설정 파일은 건드리지 않습니다.',
  'agents.note.aider': '환경 변수 2개로 OpenAI 프로토콜에 연결하고 모델은 명령줄에서 지정합니다.',
  'agents.note.cline': 'VS Code 확장에서 OpenAI Compatible을 선택하고 네 항목만 채우면 됩니다.',

  'cta.defaultLabel': 'API 키 받기',
  'cta.defaultSublabel': '가입 즉시 키가 발급되며 몇 분 안에 첫 호출이 가능합니다.',

  'auth.loginTitle': '계속하려면 로그인',
  'auth.registerTitle': '계정 만들기',
  'auth.haveAccount': '기존 계정',
  'auth.newUser': '신규 가입',
  'auth.email': '이메일',
  'auth.password': '비밀번호',
  'auth.passwordHint': '6자 이상',
  'auth.verifyCode': '이메일 인증 코드',
  'auth.verifyPlaceholder': '6자리 코드',
  'auth.send': '전송',
  'auth.sendIn': '{seconds}초',
  'auth.login': '로그인',
  'auth.register': '가입하고 계속',
  'auth.processing': '처리 중…',
  'auth.close': '닫기',
  'auth.registrationClosed': '현재 이 사이트는 회원가입이 닫혀 있습니다. 기존 계정으로 로그인하세요.',
  'auth.captchaNotice':
    '사람 확인이 켜져 있습니다. 가입이 실패하면 메인 사이트에서 가입한 뒤 다시 로그인하세요.',
  'auth.emailRequired': '이메일과 비밀번호를 입력하세요',
  'auth.verifyRequired': '인증 코드를 입력하세요',
  'auth.loginFailed': '로그인 실패',
  'auth.registerFailed': '가입 실패',
  'auth.codeSentFailed': '코드를 전송하지 못했습니다',

  'key.title': 'API 키',
  'key.desc': '이 플랜 전용 키가 발급되었습니다. 다시 클릭해도 같은 키가 반환됩니다.',
  'key.group': ' (플랜: {name})',
  'key.apiKey': 'API 키',
  'key.baseUrl': 'Base URL',
  'key.copy': '복사',
  'key.copied': '복사됨',
  'key.hint':
    'OpenAI 호환 클라이언트에 붙여넣으면 바로 사용할 수 있습니다. 잔액이 부족하면 콘솔에서 충전하세요.',
  'key.topup': '충전',
  'key.continue': '계속 둘러보기',
  'key.retry': '다시 시도',

  'error.title': '이 페이지를 사용할 수 없습니다',
  'error.reload': '새로고침',
  'error.loadFailed': '페이지를 불러오지 못했습니다',
  'error.claimFailed': '키 발급에 실패했습니다',
  'error.network': '네트워크 오류입니다. 연결을 확인하고 다시 시도하세요.',
  'error.timeout': '요청 시간이 초과되었습니다. 다시 시도하세요.',
  'error.generic': '문제가 발생했습니다. 다시 시도하세요.',

  'footer.default': '각 제품명은 해당 소유자의 상표입니다.',
}

const th: Dict = {
  'nav.login': 'เข้าสู่ระบบ',
  'nav.logout': 'ออกจากระบบ',
  'nav.console': 'คอนโซล',
  'nav.models': 'โมเดล',
  'nav.pricing': 'ราคา',
  'nav.docs': 'เอกสาร',
  'nav.language': 'ภาษา',
  'tabs.overview': 'ภาพรวม',
  'tabs.catalog': 'รายการโมเดล',
  'tabs.console': 'คอนโซล',
  'console.overview': 'ภาพรวม',
  'console.keys': 'API key',
  'console.logs': 'บันทึกการเรียกใช้',
  'console.balance': 'เครดิตคงเหลือ',
  'console.todayRequests': 'คำขอวันนี้',
  'console.todayCost': 'ค่าใช้จ่ายวันนี้',
  'console.totalRequests': 'คำขอทั้งหมด',
  'console.todayTokens': 'โทเคนวันนี้',
  'console.activeKeys': 'คีย์ที่ใช้งาน',
  'console.newKey': 'สร้างคีย์',
  'console.keyName': 'ชื่อ',
  'console.create': 'สร้าง',
  'console.creating': 'กำลังสร้าง…',
  'console.copy': 'คัดลอก',
  'console.copied': 'คัดลอกแล้ว',
  'console.disable': 'ปิดใช้งาน',
  'console.enable': 'เปิดใช้งาน',
  'console.remove': 'ลบ',
  'console.confirmRemove': 'ลบคีย์นี้หรือไม่ การเรียกที่ใช้คีย์นี้จะหยุดทำงาน',
  'console.statusActive': 'ใช้งานได้',
  'console.statusInactive': 'ปิดใช้งาน',
  'console.statusOther': 'จำกัด',
  'console.emptyKeys': 'ยังไม่มีคีย์ สร้างหนึ่งรายการเพื่อเริ่มต้น',
  'console.createdAt': 'สร้างเมื่อ',
  'console.lastUsed': 'ใช้ล่าสุด',
  'console.never': 'ไม่มี',
  'console.time': 'เวลา',
  'console.model': 'โมเดล',
  'console.key': 'คีย์',
  'console.tokens': 'โทเคน',
  'console.cost': 'ค่าใช้จ่าย',
  'console.emptyLogs': 'ยังไม่มีบันทึกการเรียกใช้',
  'console.page': 'หน้า {page} / {pages}',
  'console.prev': 'ก่อนหน้า',
  'console.next': 'ถัดไป',
  'console.loadFailed': 'โหลดข้อมูลไม่สำเร็จ',
  'console.today': 'วันนี้',
  'models.viewAll': 'ดูทั้ง {count} โมเดล',

  'hero.secondary': 'ดูโมเดลและราคา',
  'hero.trialFallback': 'บัญชีใหม่รับเครดิตฟรีเพื่อเริ่มต้น',

  'models.title': 'โมเดลและราคา',
  'models.subtitle':
    'ราคาต่อ 1 ล้านโทเคน รวมอัตราของแพ็กเกจแล้ว และคิดตามจำนวนโทเคนเข้า-ออกจริง',
  'models.group': 'แพ็กเกจ: {name}',
  'models.unavailable': 'เร็ว ๆ นี้',
  'models.input': 'อินพุต',
  'models.output': 'เอาต์พุต',
  'models.cacheRead': 'อินพุตแคช',
  'models.perRequest': 'ต่อครั้ง',
  'models.perCallUnit': '{price} / ครั้ง',
  'models.peak': 'ช่วงพีค {periods}',
  'models.officialPrice': 'ราคาปกติ',
  'models.ourPrice': 'ราคาที่คุณจ่าย',
  'models.save': 'ประหยัด {percent}%',
  'models.offpeak': 'นอกช่วงพีค',
  'models.peakShort': 'ช่วงพีค',
  'models.empty': 'กำลังเตรียมรายการโมเดล กรุณากลับมาดูใหม่ภายหลัง',

  'catalog.title': 'โมเดลทั้งหมดในแพ็กเกจนี้',
  'catalog.subtitle': 'รายการโมเดลและราคาทั้งหมด รวมโมเดลที่ไม่ได้แสดงด้านบน',
  'catalog.search': 'ค้นหาโมเดล',
  'catalog.model': 'โมเดล',
  'catalog.discount': 'ส่วนลด',
  'catalog.count': '{count} โมเดล',
  'catalog.empty': 'ไม่พบโมเดลที่ตรงกัน',

  'usage.title': 'วิธีใช้งาน',

  'compat.title': 'ปลายทางเดียว สามโปรโตคอล',
  'compat.openai': 'OpenAI Chat Completions',
  'compat.anthropic': 'Anthropic Messages',
  'compat.responses': 'OpenAI Responses',
  'compat.gemini': 'Gemini',
  'compat.note': 'คีย์เดียวกัน Base URL เดียวกัน เปลี่ยนโปรโตคอลได้โดยไม่ต้องเปลี่ยนผู้ให้บริการ',

  'integration.title': 'การเชื่อมต่อ',
  'integration.subtitle':
    'เพียงเปลี่ยน Base URL ของ SDK เดิม ไม่ต้องแก้โค้ดไคลเอนต์หรือเพิ่มชั้น abstraction ใหม่',
  'integration.copy': 'คัดลอก',
  'integration.copied': 'คัดลอกแล้ว',

  'cta.defaultLabel': 'รับ API key',
  'cta.defaultSublabel': 'สมัครแล้วรับคีย์ทันที เริ่มเรียกใช้ได้ในไม่กี่นาที',

  'auth.loginTitle': 'เข้าสู่ระบบเพื่อดำเนินการต่อ',
  'auth.registerTitle': 'สร้างบัญชี',
  'auth.haveAccount': 'มีบัญชีอยู่แล้ว',
  'auth.newUser': 'สมัครใหม่',
  'auth.email': 'อีเมล',
  'auth.password': 'รหัสผ่าน',
  'auth.passwordHint': 'อย่างน้อย 6 ตัวอักษร',
  'auth.verifyCode': 'รหัสยืนยันทางอีเมล',
  'auth.verifyPlaceholder': 'รหัส 6 หลัก',
  'auth.send': 'ส่ง',
  'auth.sendIn': '{seconds} วินาที',
  'auth.login': 'เข้าสู่ระบบ',
  'auth.register': 'สมัครและดำเนินการต่อ',
  'auth.processing': 'กำลังดำเนินการ…',
  'auth.close': 'ปิด',
  'auth.registrationClosed': 'ขณะนี้ปิดรับสมัครสมาชิก กรุณาเข้าสู่ระบบด้วยบัญชีเดิม',
  'auth.captchaNotice':
    'เปิดการยืนยันตัวตนแล้ว หากสมัครไม่สำเร็จ กรุณาสมัครที่เว็บหลักแล้วกลับมาเข้าสู่ระบบ',
  'auth.emailRequired': 'กรุณากรอกอีเมลและรหัสผ่าน',
  'auth.verifyRequired': 'กรุณากรอกรหัสยืนยัน',
  'auth.loginFailed': 'เข้าสู่ระบบไม่สำเร็จ',
  'auth.registerFailed': 'สมัครไม่สำเร็จ',
  'auth.codeSentFailed': 'ส่งรหัสไม่สำเร็จ',

  'key.title': 'API key ของคุณ',
  'key.desc': 'สร้างคีย์เฉพาะสำหรับแพ็กเกจนี้แล้ว กดซ้ำจะได้คีย์เดิม',
  'key.group': ' (แพ็กเกจ: {name})',
  'key.apiKey': 'API key',
  'key.baseUrl': 'Base URL',
  'key.copy': 'คัดลอก',
  'key.copied': 'คัดลอกแล้ว',
  'key.hint':
    'วางในไคลเอนต์ที่รองรับ OpenAI ได้ทันที หากเครดิตไม่พอให้เติมในคอนโซล',
  'key.topup': 'เติมเครดิต',
  'key.continue': 'ดูต่อ',
  'key.retry': 'ลองใหม่',

  'error.title': 'หน้านี้ใช้งานไม่ได้ในขณะนี้',
  'error.reload': 'โหลดใหม่',
  'error.loadFailed': 'โหลดหน้าไม่สำเร็จ',
  'error.claimFailed': 'สร้างคีย์ไม่สำเร็จ',
  'error.network': 'เครือข่ายขัดข้อง กรุณาตรวจสอบการเชื่อมต่อแล้วลองใหม่',
  'error.timeout': 'หมดเวลา กรุณาลองใหม่อีกครั้ง',
  'error.generic': 'เกิดข้อผิดพลาด กรุณาลองใหม่',

  'footer.default': 'ชื่อผลิตภัณฑ์เป็นเครื่องหมายการค้าของเจ้าของแต่ละราย',
}

const DICTS: Record<CampaignLocale, Dict> = { en, zh, ja, ko, th }

export type Translate = (key: string, vars?: Record<string, string | number>) => string

export function normalizeCampaignLocale(value?: string): CampaignLocale {
  const raw = String(value || '')
    .trim()
    .toLowerCase()
  if (raw.startsWith('zh')) return 'zh'
  if (raw.startsWith('ja')) return 'ja'
  if (raw.startsWith('ko')) return 'ko'
  if (raw.startsWith('th')) return 'th'
  return 'en'
}

export function isSupportedLocale(value: string): value is CampaignLocale {
  return SUPPORTED_LOCALES.some((item) => item.code === value)
}

export function createMessages(locale: CampaignLocale): Translate {
  const dict = DICTS[locale] ?? en
  return (key, vars) => {
    let text = dict[key] ?? en[key] ?? key
    if (vars) {
      for (const [name, value] of Object.entries(vars)) {
        text = text.split(`{${name}}`).join(String(value))
      }
    }
    return text
  }
}

/**
 * 通过 provide/inject 下发当前语言。
 *
 * 落地页只有几个组件，用 inject 比逐层传 locale + 再建字典更省事；
 * 提供的闭包在渲染时读取 computed，语言切换能正常触发重渲染。
 */
export const CampaignMessagesKey: InjectionKey<Translate> = Symbol('campaign-messages')

export function useCampaignMessages(): Translate {
  const translate = inject(CampaignMessagesKey)
  if (!translate) {
    return createMessages('en')
  }
  return translate
}
