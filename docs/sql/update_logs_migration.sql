-- 更新日志系统数据库迁移脚本
-- 创建时间: 2025-01-20
-- 功能: 为社区平台添加更新日志功能

-- 创建 update_logs 表
CREATE TABLE `update_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(255) NOT NULL COMMENT '更新标题',
  `description` text NOT NULL COMMENT '更新描述',
  `content` longtext COMMENT '更新详细内容（支持Markdown）',
  `version` varchar(50) DEFAULT NULL COMMENT '版本号',
  `type` varchar(50) DEFAULT 'feature' COMMENT '更新类型：新功能、修复、改进、安全、其他',
  `status` varchar(20) DEFAULT 'active' COMMENT '状态：活跃、停用',
  `publish_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '发布日期',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_type` (`type`),
  KEY `idx_publish_date` (`publish_date`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='更新日志表';

-- 插入示例数据
INSERT INTO `update_logs` (`title`, `description`, `content`, `version`, `type`, `status`, `publish_date`) VALUES
('新增AI问答功能', '集成ChatGPT API，支持技术问答，用户可以直接在平台内进行AI对话交流...', '# 新增AI问答功能\n\n我们很高兴地宣布，平台现已集成**ChatGPT API**，为用户提供智能技术问答服务！\n\n## ✨ 主要功能\n\n- **智能问答**：支持各种技术问题的AI回答\n- **代码高亮**：自动识别和高亮显示代码片段  \n- **多轮对话**：保持上下文连贯的对话体验\n- **快速响应**：秒级响应，提升用户体验\n\n## 🔧 使用方式\n\n1. 点击页面右下角的AI助手图标\n2. 输入您的技术问题\n3. 获得专业的AI回答和建议\n\n期待这个功能能帮助大家更高效地解决技术问题！', 'v2.1.0', 'feature', 'active', '2025-01-15 10:00:00'),

('修复评论表情显示问题', '优化移动端表情选择器交互体验，修复了表情在某些设备上无法正常显示的问题...', '# 修复评论表情显示问题\n\n针对用户反馈的表情显示问题，我们进行了全面的优化和修复。\n\n## 🐛 修复内容\n\n- 修复移动端表情选择器的溢出问题\n- 解决部分设备上表情无法正常显示的兼容性问题\n- 优化表情加载性能，减少加载时间\n- 改进表情选择器的触摸交互体验\n\n## 📱 移动端优化\n\n- 调整表情选择器尺寸，适配不同屏幕\n- 优化触摸反馈，提升操作流畅度\n- 修复键盘弹出时的布局问题', 'v2.0.5', 'bugfix', 'active', '2025-01-10 14:30:00'),

('提升页面加载速度', '重构了图片加载策略，使用懒加载技术，整体页面加载速度提升50%...', '# 提升页面加载速度\n\n通过多项性能优化措施，显著提升了平台的整体加载性能。\n\n## ⚡ 优化内容\n\n- **图片懒加载**：实现图片按需加载，减少初始加载时间\n- **资源压缩**：优化图片和静态资源，减小文件体积\n- **缓存策略**：改进浏览器缓存机制，提升二次访问速度\n- **代码分割**：优化JavaScript打包策略，减少首屏加载时间\n\n## 📊 性能提升\n\n- 首页加载时间减少 **50%**\n- 图片加载效率提升 **60%**  \n- 整体用户体验评分提升 **30%**', 'v2.0.4', 'improvement', 'active', '2025-01-08 09:15:00'),

('优化课程视频播放', '支持倍速播放和进度记忆功能，增强了视频学习体验...', '# 优化课程视频播放\n\n为了提供更好的在线学习体验，我们对视频播放器进行了全面升级。\n\n## 🎥 新增功能\n\n- **倍速播放**：支持0.5x、1.25x、1.5x、2x倍速\n- **进度记忆**：自动保存观看进度，下次继续观看\n- **快捷键支持**：空格暂停/播放，左右键快进/快退\n- **画质选择**：根据网络状况自动或手动调整画质\n\n## 📱 移动端优化\n\n- 优化全屏播放体验\n- 支持手势控制音量和亮度\n- 改进移动网络下的播放稳定性', 'v2.0.3', 'feature', 'active', '2025-01-05 16:45:00'),

('新增文章分享功能', '支持一键分享文章到微信、QQ等社交平台，方便用户传播优质内容...', '# 新增文章分享功能\n\n现在用户可以轻松分享优质文章到各大社交平台了！\n\n## 📤 分享功能\n\n- **多平台支持**：微信、QQ、微博、钉钉等\n- **自定义分享**：可编辑分享标题和描述\n- **分享统计**：查看文章分享数据\n- **分享奖励**：分享获得积分奖励\n\n让知识传播得更广更远！', 'v2.0.2', 'feature', 'active', '2025-01-03 11:20:00'),

('增强安全性配置', '升级了用户认证系统，添加了双因子认证和更严格的密码策略', '# 增强安全性配置\n\n为了保护用户账户安全，我们对平台的安全策略进行了全面升级。\n\n## 🔒 安全增强\n\n- **双因子认证**：支持短信和邮箱验证码\n- **密码策略**：更严格的密码复杂度要求\n- **登录保护**：异常登录行为检测和通知\n- **会话管理**：更安全的用户会话处理\n\n## 🛡️ 隐私保护\n\n- 数据传输加密升级\n- 用户信息脱敏处理\n- 操作日志完善', 'v2.0.1', 'security', 'active', '2025-01-01 10:00:00');

-- 验证数据插入
SELECT COUNT(*) as total_count FROM update_logs;
SELECT type, COUNT(*) as count FROM update_logs GROUP BY type;
SELECT status, COUNT(*) as count FROM update_logs GROUP BY status;