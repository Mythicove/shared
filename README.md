# Mythicove Shared

Mythicove 生态的共享 Go 模块。承载跨产品的能力**接口**与领域类型，各产品仓库（`mythic-mirror`、`mythic-hoard` 等）通过 Go module 引用本仓库。

## 设计原则

- **只放接口与纯领域类型，不放具体实现**（DB/S3/Redis/Meilisearch 的实现留在各产品或后续的实现包中）。
- 能力通过接口抽象，配合部署分级（10 → 10k）在低档/高档实现间切换，不改业务代码。
- 包名即能力，语义见各包注释。

## 包

| 包 | 能力 | 低档实现 | 高档实现 |
|---|---|---|---|
| `auth` | 账号认证 | 内建 JWT/session | OIDC（Zitadel/Fosite） |
| `storage` | 对象存储 | 本地磁盘 | MinIO → S3（R2/OSS） |
| `search` | 搜索 | PostgreSQL `pg_trgm` | Meilisearch |
| `cache` | 缓存 | 进程内内存 | Redis |
| `ratelimit` | 限流 | 进程内 | Redis |
| `queue` | 任务队列 | 内联执行 | asynq（Redis） |
| `mailer` | 邮件 | SMTP | 服务商 API |
| `mc` | Minecraft 领域 | 离线 UUID、玩家名校验等纯函数 | — |

## 本地联调

多仓库同时开发时用 Go workspace 关联本地 checkout（免打 tag）：

```
go work init
go work use . ../mythic-mirror ../mythic-hoard
```

## 许可证

[AGPL-3.0](./LICENSE)（核心产品代码，与产品仓库一致）。
