// Package auth 定义账号认证能力的抽象。
// 低档实现为内建 JWT/session，高档实现为 OIDC（Zitadel/Fosite）。
package auth

import "context"

// User 是认证后返回的最小用户视图。
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// Credentials 是用户名密码凭证。
type Credentials struct {
	Username string
	Password string
}

// Auth 是账号认证接口。
type Auth interface {
	// Authenticate 校验凭证，返回访问令牌与用户。
	Authenticate(ctx context.Context, creds Credentials) (accessToken string, u *User, err error)
	// Validate 校验访问令牌，返回其所属用户。
	Validate(ctx context.Context, accessToken string) (*User, error)
	// Refresh 用刷新令牌换取新的访问令牌。
	Refresh(ctx context.Context, refreshToken string) (accessToken string, err error)
	// SignOut 使当前会话失效。
	SignOut(ctx context.Context, accessToken string) error
}
