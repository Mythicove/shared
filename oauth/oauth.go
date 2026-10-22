// Package oauth 定义第三方 OAuth 登录抽象。
// 实现：GitHub / Google / Microsoft / Discord / QQ / 微信 / 微博 等，按配置启用多个。
package oauth

import "context"

// Profile 是第三方登录返回的统一用户资料。
type Profile struct {
	Provider string // github / google / qq ...
	UID      string // 平台内唯一标识
	Email    string
	Nickname string
}

// Provider 是单个第三方登录源。
type Provider interface {
	// AuthURL 生成授权跳转 URL。
	AuthURL(state string) string
	// Exchange 用回调 code 换取用户资料。
	Exchange(ctx context.Context, code string) (*Profile, error)
}
