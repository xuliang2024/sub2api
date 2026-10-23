import { createApp } from 'vue'
import '../src/style.css'
import './campaign.css'
import CampaignApp from './CampaignApp.vue'
import LegalPage from './LegalPage.vue'

// 落地页不需要路由：整页是单屏营销页 + 弹层，用 hash 锚点即可。
// 这样能省掉 vue-router 的运行时代价，也让首屏更快（投放页的 LCP 直接换钱）。
// 例外是法务页：Google OAuth 品牌审核要能直接打开 /privacy 与 /terms，
// 所以按 pathname 单独挂载一个静态页。
const path = window.location.pathname.replace(/\/+$/, '')
const isLegalPage = path === '/privacy' || path === '/terms'

createApp(isLegalPage ? LegalPage : CampaignApp).mount('#app')
