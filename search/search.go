// Package search 定义搜索抽象。
// 低档实现为 PostgreSQL pg_trgm，高档实现为 Meilisearch。
package search

import "context"

// Document 是被索引的文档。
type Document struct {
	ID     string
	Index  string
	Fields map[string]any
	// Facets 用于分类、loader、游戏版本等过滤维度。
	Facets map[string][]string
}

// Search 是索引与查询抽象。
type Search interface {
	// Index 批量写入文档。
	Index(ctx context.Context, index string, docs []Document) error
	// Delete 从索引中删除单条文档。
	Delete(ctx context.Context, index string, id string) error
	// Query 全文检索，返回命中的文档 ID 列表与总命中数（分页 total 用）。
	Query(ctx context.Context, index, q string, facets map[string][]string, limit, offset int) (ids []string, total int64, err error)
}
