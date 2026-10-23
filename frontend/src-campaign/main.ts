import { createApp } from 'vue'
import '../src/style.css'
import './campaign.css'
import CampaignApp from './CampaignApp.vue'

// 落地页不需要路由：整页是单屏营销页 + 弹层，用 hash 锚点即可。
// 这样能省掉 vue-router 的运行时代价，也让首屏更快（投放页的 LCP 直接换钱）。
createApp(CampaignApp).mount('#app')
