// Package storage 定义对象存储抽象。
// 低档实现为本地磁盘，高档实现为 MinIO / S3（R2、OSS 等兼容端点）。
package storage

import (
	"context"
	"io"
	"net/url"
)

// Object 是一次上传的输入。
type Object struct {
	Key         string
	Size        int64
	ContentType string
	Body        io.Reader
}

// Storage 是文件读写与直出 URL 签名的抽象。
type Storage interface {
	// Put 上传对象，返回存储键。
	Put(ctx context.Context, o Object) (key string, err error)
	// Get 读取对象内容。
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete 删除对象。
	Delete(ctx context.Context, key string) error
	// PresignGet 生成带签名的直出 URL（供 CDN/直链下载）。
	PresignGet(ctx context.Context, key string) (*url.URL, error)
}
