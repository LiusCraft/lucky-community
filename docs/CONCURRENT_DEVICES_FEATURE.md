# 同时在线设备数量限制功能

## 功能概述

本功能实现了对用户同时在线设备数量的限制，默认每个用户只能有1台设备同时在线。管理员可以针对特定用户调整这个限制，管理员本身没有设备数量限制。

## 主要特性

1. **设备数量限制**：默认用户只能1台设备同时在线
2. **管理员无限制**：管理员用户没有设备数量限制
3. **自动踢出机制**：当用户登录设备超过限制时，自动踢出最早登录的设备
4. **设备管理**：用户可以查看当前在线设备并主动踢出指定设备
5. **管理员配置**：管理员可以为用户设置不同的设备数量限制

## 数据库变更

### 1. 用户表新增字段
```sql
ALTER TABLE users 
ADD COLUMN max_concurrent_devices INT NOT NULL DEFAULT 1 COMMENT '最大同时在线设备数量，默认为1';
```

### 2. 新增在线会话表
```sql
CREATE TABLE user_online_sessions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    device_info VARCHAR(500),
    ip_address VARCHAR(45),
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_heartbeat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_last_heartbeat (last_heartbeat),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## API 接口

### 用户接口

#### 1. 查看在线设备
```
GET /community/user/devices
```

响应示例：
```json
{
  "code": 200,
  "data": {
    "devices": [
      {
        "sessionId": "abc123",
        "deviceInfo": "Mozilla/5.0...",
        "ipAddress": "192.168.1.100",
        "loginTime": "2025-07-24T10:00:00Z",
        "lastHeartbeat": "2025-07-24T10:30:00Z"
      }
    ],
    "maxDevices": 1,
    "current": 1
  }
}
```

#### 2. 踢出设备
```
DELETE /community/user/devices/{sessionId}
```

### 管理员接口

#### 1. 查看用户在线设备
```
GET /community/admin/user/{id}/devices
```

#### 2. 设置用户最大设备数
```
PUT /community/admin/user/{id}/max-devices
Content-Type: application/json

{
  "maxDevices": 3
}
```

#### 3. 踢出用户设备
```
DELETE /community/admin/user/{userId}/devices/{sessionId}
```

## 工作原理

### 1. 登录流程
1. 用户登录成功后，系统生成会话ID（基于JWT token的MD5哈希）
2. 检查当前用户在线设备数量
3. 如果超过限制且非管理员，踢出最早登录的设备
4. 注册新设备会话

### 2. 心跳机制
1. 前端定期发送心跳请求
2. 更新设备的最后心跳时间和IP地址
3. 检查会话有效性
4. 保持原有的IP一致性检查作为安全备用

### 3. 设备清理
1. 系统自动清理超过5分钟没有心跳的会话
2. 过期的token会被加入黑名单
3. 用户主动踢出设备时，对应token也会被加入黑名单

## 部署步骤

1. **执行数据库迁移**
   ```bash
   mysql -u username -p database_name < docs/sql/concurrent_devices_migration.sql
   ```

2. **重新编译应用**
   ```bash
   go build -o build/community cmd/community/main.go
   ```

3. **重启服务**
   ```bash
   ./build/community
   ```

## 配置说明

### 默认配置
- 普通用户最大设备数：1
- 管理员最大设备数：999999（实际无限制）
- 会话过期时间：5分钟无心跳

### 自定义配置
管理员可以通过管理后台为特定用户设置不同的设备数量限制。

## 注意事项

1. **向后兼容**：现有用户的max_concurrent_devices字段会自动设置为1
2. **管理员识别**：通过invite_codes和member_infos表关联判断管理员身份
3. **安全性**：保持了原有的IP一致性检查作为额外安全措施
4. **性能**：会话清理是自动进行的，不需要额外的定时任务

## 故障排除

### 常见问题

1. **用户无法登录**
   - 检查数据库连接
   - 确认user_online_sessions表是否正确创建
   - 查看应用日志中的错误信息

2. **设备没有被踢出**
   - 检查心跳机制是否正常工作
   - 确认会话清理逻辑是否执行
   - 查看token黑名单是否生效

3. **管理员功能异常**
   - 确认管理员权限判断逻辑
   - 检查member_infos表中是否有name='admin'的记录

### 日志监控
关注以下日志信息：
- 设备注册失败
- 心跳更新失败
- 设备踢出操作
- 管理员配置变更

## 扩展功能

未来可以考虑添加：
1. 设备类型识别（手机、电脑、平板等）
2. 地理位置信息
3. 设备别名设置
4. 登录通知功能
5. 异常登录检测
