-- 同时在线设备数量限制功能数据库迁移脚本
-- 创建时间：2025-07-24

-- 为用户表添加最大同时在线设备数量字段
ALTER TABLE users 
ADD COLUMN max_concurrent_devices INT NOT NULL DEFAULT 1 COMMENT '最大同时在线设备数量，默认为1';

-- 创建在线设备会话表
CREATE TABLE user_online_session (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id INT NOT NULL COMMENT '用户ID，关联 users 表',
    session_id VARCHAR(255) NOT NULL UNIQUE COMMENT '会话ID，通常使用JWT token的唯一标识',
    device_info VARCHAR(500) COMMENT '设备信息（User-Agent等）',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    last_heartbeat TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后心跳时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_last_heartbeat (last_heartbeat),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户在线会话表';

-- 为管理员用户设置无限制（设置为一个很大的数字）
UPDATE users 
SET max_concurrent_devices = 999999 
WHERE id IN (
    SELECT u.id 
    FROM users u
    JOIN invite_codes ic ON u.invite_code = ic.code
    JOIN member_infos mi ON mi.id = ic.member_id
    WHERE mi.name = 'admin'
);
