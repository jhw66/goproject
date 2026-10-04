# go-zero MySQL 2FA 示例

这个服务现在是一个更接近正式项目结构的用户名密码 + TOTP 2FA 示例。

## 数据库变更

先执行 `mysql.sql`，当前 `user_info` 表至少包含这些字段：

- `username`
- `password_hash`
- `totp_secret`
- `totp_pending_secret`
- `totp_enabled`

`totp_secret` 与 `totp_pending_secret` 保存的是经过应用层加密后的内容。

## 启动

```bash
go run mysql.go
```

默认地址：`http://127.0.0.1:8888`

## 接口流程

### 1. 注册

```bash
curl -X POST "http://127.0.0.1:8888/api/users/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"demo\",\"password\":\"Passw0rd!\"}"
```

### 2. 首次登录（尚未启用 2FA）

```bash
curl -X POST "http://127.0.0.1:8888/api/users/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"demo\",\"password\":\"Passw0rd!\"}"
```

未启用 2FA 时会直接返回：

- `accessToken`
- `accessExpire`
- `need2fa=false`

### 3. 发起 TOTP 绑定

```bash
curl -X POST "http://127.0.0.1:8888/api/users/2fa/totp/enroll/start" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{}"
```

返回：

- `otpauthUrl`
- `qrCodeDataUrl`
- `issuer`
- `accountName`

把 `otpauthUrl` 或 `qrCodeDataUrl` 展示给认证器扫码，然后输入认证器里当前的 6 位动态码。

### 4. 确认启用 2FA

```bash
curl -X POST "http://127.0.0.1:8888/api/users/2fa/totp/enroll/confirm" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"code\":\"123456\"}"
```

成功后返回：

```json
{
  "enabled": true
}
```

### 5. 启用 2FA 后再次登录

```bash
curl -X POST "http://127.0.0.1:8888/api/users/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"demo\",\"password\":\"Passw0rd!\"}"
```

此时不会直接返回正式登录 token，而是返回：

- `need2fa=true`
- `challengeToken`

### 6. 提交第二步验证码

```bash
curl -X POST "http://127.0.0.1:8888/api/users/2fa/verify" \
  -H "Content-Type: application/json" \
  -d "{\"challengeToken\":\"CHALLENGE_TOKEN\",\"code\":\"123456\"}"
```

成功后返回正式：

- `accessToken`
- `accessExpire`

### 7. 关闭 2FA

```bash
curl -X POST "http://127.0.0.1:8888/api/users/2fa/totp/disable" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"Passw0rd!\",\"code\":\"123456\"}"
```

## 配置说明

`etc/mysql-api.yaml` 新增了两组配置：

- `Auth`
  - `AccessSecret`
  - `AccessExpire`
  - `ChallengeSecret`
  - `ChallengeExpire`
- `Totp`
  - `Issuer`
  - `SecretCipherKey`
  - `QRCodeSize`

正式环境里请替换示例 secret，不要直接使用仓库中的默认值。
