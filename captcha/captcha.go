// Package captcha 定义人机验证（验证码）校验抽象。
// 实现：Cloudflare Turnstile / Google reCAPTCHA / hCaptcha 等，按配置选择。
package captcha

import "context"

// Captcha 校验一次人机验证 token。
type Captcha interface {
	// Verify 校验 token；remoteIP 为客户端 IP（部分服务商需要）。
	Verify(ctx context.Context, token, remoteIP string) (bool, error)
}
