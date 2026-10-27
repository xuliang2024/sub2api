/**
 * 模型定价页「折扣」体系的前端开关。
 *
 * 关闭（false）时页面不展示任何折扣信息：列表折扣列、展开子表的「对比官方」列、
 * 高于官价的标黄、顶部「平均折扣」、仅看有折扣筛选、折扣排序。
 * 官方价与本站价照常展示，接口也照旧返回 official_pricing —— 改回 true 即恢复。
 */
export const SHOW_OFFICIAL_DISCOUNT = false
