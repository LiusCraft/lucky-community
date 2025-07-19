-- 福利系统数据库表
-- 创建时间: 2025-01-09

-- 福利项目表
CREATE TABLE welfare_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL COMMENT '福利标题',
    description TEXT COMMENT '简短描述',
    detail_content TEXT COMMENT '详细说明(Markdown格式)',
    tag VARCHAR(50) NOT NULL COMMENT '标签: 云服务、开发工具、优惠券、学习资源',
    
    -- 价格相关
    price DECIMAL(10,2) DEFAULT 0.00 COMMENT '现价，0表示免费或优惠券',
    original_price DECIMAL(10,2) DEFAULT 0.00 COMMENT '原价',
    discount_text VARCHAR(50) COMMENT '折扣标签: 3折、免费、限时等',
    
    -- 操作
    action_text VARCHAR(50) DEFAULT '立即查看' COMMENT '按钮文本',
    
    -- 状态管理
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    sort_order INT DEFAULT 0 COMMENT '排序',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_status_sort (status, sort_order),
    INDEX idx_tag (tag)
) COMMENT='福利项目表';

-- 插入示例数据
INSERT INTO welfare_items (title, description, detail_content, tag, price, original_price, discount_text, action_text, sort_order) VALUES
(
    '腾讯云服务器', 
    '新用户专享，高性能云主机', 
    '# 腾讯云服务器优惠活动

## 活动详情
- **优惠力度**: 首年3折，续费6折
- **适用范围**: 新用户专享
- **服务器配置**: 
  - CPU: 2核
  - 内存: 4GB
  - 硬盘: 50GB SSD
  - 带宽: 5Mbps

## 使用说明
1. 注册腾讯云账户
2. 完成实名认证
3. 选择服务器配置
4. 使用优惠码享受折扣

> 💡 **提示**: 建议选择年付方案，性价比更高！

## 注意事项
- 优惠仅限新用户
- 每个账户限购一次
- 优惠不可与其他活动叠加使用', 
    '云服务', 
    99.00, 
    330.00, 
    '3折', 
    '立即购买', 
    100
),
(
    'JetBrains全家桶', 
    '学生认证免费使用', 
    '# JetBrains学生免费计划

## 包含软件
- **IntelliJ IDEA Ultimate** - Java开发神器
- **PyCharm Professional** - Python开发利器
- **WebStorm** - 前端开发首选
- **GoLand** - Go语言开发工具
- **PhpStorm** - PHP开发环境
- **Rider** - .NET开发IDE
- **CLion** - C/C++开发工具
- **DataGrip** - 数据库管理工具

## 申请条件
- 在校学生（本科/研究生/博士）
- 有效的学校邮箱(.edu邮箱)
- 学生证或在读证明

## 申请流程
1. 访问JetBrains学生页面
2. 使用学校邮箱注册
3. 上传学生证明材料
4. 等待审核通过（通常1-2天）
5. 下载并激活软件

**有效期**: 1年，可续期至毕业

## 续期说明
- 每年需要重新验证学生身份
- 可以在到期前30天申请续期
- 毕业后可享受40%优惠购买正版授权', 
    '开发工具', 
    0.00, 
    1999.00, 
    '免费', 
    '免费申请', 
    90
),
(
    '阿里云优惠券', 
    '新用户专享代金券', 
    '# 阿里云新用户优惠券

## 优惠券详情
- **代金券面额**: 
  - 满100减50元
  - 满500减200元
  - 满1000减400元
  - 满2000减800元
- **适用产品**: 云服务器ECS、RDS数据库、OSS存储、CDN等
- **有效期**: 领取后30天内使用

## 使用限制
- 仅限新用户使用
- 不可与其他优惠叠加
- 每个账户限领一次
- 不支持退款换现

## 推荐搭配
建议与阿里云学生机一起使用，性价比最高！

### 学生机配置
- 1核2GB内存
- 40GB系统盘
- 1Mbps带宽
- 月付9.5元

> 🎯 **温馨提示**: 优惠券数量有限，先到先得！

## 常见问题
**Q: 优惠券可以叠加使用吗？**
A: 不可以，每次订单只能使用一张优惠券。

**Q: 优惠券过期了怎么办？**
A: 过期的优惠券无法使用，请在有效期内及时使用。', 
    '优惠券', 
    0.00, 
    0.00, 
    '满减', 
    '立即领取', 
    80
),
(
    '编程学习资料', 
    '500G精选学习视频', 
    '# 编程学习资料包

## 资料内容

### 前端开发 (150G)
- **Vue.js全栈开发实战** - 从基础到项目实战
- **React企业级项目实战** - Hook + Redux + TypeScript
- **微信小程序开发教程** - 从入门到上线
- **TypeScript进阶指南** - 类型系统深度解析
- **Webpack构建优化** - 性能调优实战

### 后端开发 (200G)
- **Spring Boot微服务实战** - 分布式架构设计
- **Go语言从入门到精通** - 并发编程与微服务
- **Python爬虫与数据分析** - 实战项目案例
- **Node.js全栈开发** - Express + MongoDB
- **Docker容器化部署** - DevOps实践

### 数据库 (80G)
- **MySQL性能优化实战** - 索引优化与SQL调优
- **Redis缓存技术详解** - 集群搭建与实战
- **MongoDB实战教程** - 文档数据库应用
- **Elasticsearch搜索引擎** - 全文检索实战

### 其他 (70G)
- **算法与数据结构** - LeetCode刷题指南
- **系统设计面试指南** - 大厂面试真题
- **DevOps实践教程** - CI/CD流水线搭建
- **计算机网络原理** - 网络协议深度解析

## 获取方式
1. 完成支付
2. 添加客服微信：**codexy2025**
3. 发送支付截图
4. 获取网盘链接和提取码

## 资料特色
- ✅ 最新录制，技术栈版本新
- ✅ 项目驱动，实战性强
- ✅ 代码完整，可直接运行
- ✅ 永久更新，持续添加新内容

**🔥 限时优惠，原价199元，现价仅29元！**

> 💝 **额外福利**: 购买后可加入专属学习群，与千名程序员一起交流学习！', 
    '学习资源', 
    29.00, 
    199.00, 
    '限时', 
    '立即下载', 
    70
);