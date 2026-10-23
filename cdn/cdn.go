// Package cdn 定义 CDN 缓存刷新抽象。
// 实现：Cloudflare / 阿里云 CDN / 腾讯云 CDN / BunnyCDN 等，按配置选择。
package cdn

import "context"

// CDN 刷新指定路径的缓存。
type CDN interface {
	// Purge 刷新 paths（URL 或路径前缀，语义由实现约定）。
	Purge(ctx context.Context, paths []string) error
}
