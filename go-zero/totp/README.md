# go-zero TOTP 示例

这个示例演示如何在 go-zero HTTP 服务中使用 `github.com/pquerna/otp/totp` 生成和校验动态验证码。

## 启动

```bash
go mod tidy
go run totp.go
```

默认监听 `http://127.0.0.1:8889`。

## 1. 生成 TOTP 密钥

```bash
curl -X POST "http://127.0.0.1:8889/api/totp/generate" \
  -H "Content-Type: application/json" \
  -d "{\"accountName\":\"demo@example.com\",\"issuer\":\"demo-go-zero\"}"
```

返回示例：

```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "issuer": "demo-go-zero",
  "accountName": "demo@example.com",
  "otpauthUrl": "otpauth://totp/demo-go-zero:demo@example.com?algorithm=SHA1&digits=6&issuer=demo-go-zero&period=30&secret=JBSWY3DPEHPK3PXP"
}
```

把返回的 `secret` 或 `otpauthUrl` 导入 Google Authenticator、Microsoft Authenticator 等应用后，即可得到 6 位动态码。

### 用二维码展示 otpauthUrl（推荐）

思路：把 **`otpauthUrl` 原文字符串** 交给二维码库编码成图片；认证器 App 扫码后，等价于手动输入密钥。

本项目提供两种方式（响应均为 **`Content-Type: image/png`**）：

1. **GET（适合简单页面 `<img src="...">`）**  
   把 URL 做查询参数传递，注意必须 **`encodeURIComponent`**，否则 `?、&` 会破坏 query：

   ```html
   <img
     alt="Scan to add TOTP"
     src="http://127.0.0.1:8889/api/totp/qrcode?url=这里是 encodeURIComponent(otpauthUrl)"
   />
   ```

   若 `otpauthUrl` 来自上一步 JSON，可在前端：

   ```javascript
   const src =
     "http://127.0.0.1:8889/api/totp/qrcode?url=" +
     encodeURIComponent(data.otpauthUrl);
   ```

2. **POST（更稳妥，otpauth 不会出现在访问日志的 query 里）**

   ```bash
   curl -s -X POST "http://127.0.0.1:8889/api/totp/qrcode" \
     -H "Content-Type: application/json" \
     -d "{\"otpauthUrl\":\"otpauth://totp/...\"}" \
     -o qr.png
   ```

   纯前端可把返回的 PNG 转成 Blob，再赋给 `<img src={blobUrl}>`。

**纯前端方案（不经服务端）**：用 `qrcode`（npm）等库在浏览器里对 `otpauthUrl` 生成 canvas/dataURL，不增加后端依赖；本项目实现的是 **服务端生成 PNG**，便于统一尺寸与缓存策略。

## 2. 校验验证码

```bash
curl -X POST "http://127.0.0.1:8889/api/totp/verify" \
  -H "Content-Type: application/json" \
  -d "{\"secret\":\"JBSWY3DPEHPK3PXP\",\"code\":\"123456\"}"
```

返回示例：

```json
{
  "valid": true
}
```
