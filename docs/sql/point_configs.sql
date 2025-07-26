-- 积分配置表
CREATE TABLE `point_configs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `rules_description` text NOT NULL COMMENT '积分规则说明',
  `invite_reward_points` int NOT NULL DEFAULT '10' COMMENT '邀请用户注册积分奖励',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分配置表';

-- 插入默认配置数据
INSERT INTO `point_configs` (`rules_description`, `invite_reward_points`) VALUES (
'# 积分规则说明

## 如何获得积分

- **邀请好友注册**：获得 **10** 积分
- **参与社区活动**：根据活动规则获得积分

## 积分使用说明

- 积分可用于兑换商城内的各种商品
- 兑换后积分立即扣除，不可退还
- 积分不可转让给其他用户

## 其他注意事项

- 积分有效期为获得后**1年**
- 商品数量有限，**先兑先得**
- 如有疑问请联系客服

---

> 💡 **提示**：积分是社区的通用货币，合理使用可以获得更多优质服务！', 
10
);