// Package i18n 定义国际化抽象。
// 实现可用 go-i18n（Go 后端）/ i18next（TS 前端）。
package i18n

import "context"

// Localizer 按 locale 解析消息键，处理模板参数与复数。
type Localizer interface {
	// T 翻译键，args 为模板参数。
	T(ctx context.Context, locale, key string, args map[string]any) string
	// Plural 处理复数形式。
	Plural(ctx context.Context, locale, key string, n int, args map[string]any) string
}

// DefaultLocale 是缺省回退 locale。
const DefaultLocale = "en-US"

// Locales 是起步支持的 locale 列表（BCP 47）。
var Locales = []string{"zh-CN", "en-US"}
