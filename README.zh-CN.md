# Mythicove Shared

[English](./README.md) | **简体中文**

Mythicove 生态的共享 Go 模块，承载跨产品的能力接口与领域类型。各产品仓库（mythic-mirror、mythic-hoard 等）通过 Go module 引用本仓库。

## 设计原则

- 只放接口、纯领域类型和零依赖的进程内实现。带外部依赖的实现（DB、S3、Redis、Meilisearch 等）留在各产品仓库。
- 能力通过接口抽象。配合部署分级（10 人到 10k 人），可以在低档和高档实现之间切换，不改业务代码。
- 包名即能力，语义见各包注释。

## 包

| 包 | 能力 | 低档实现 | 高档实现 |
|---|---|---|---|
| storage | 对象存储 | 本地磁盘 | MinIO，再到 S3（R2、OSS） |
| search | 搜索 | PostgreSQL pg_trgm | Meilisearch |
| cache | 缓存 | 进程内内存 | Redis |
| ratelimit | 限流 | 进程内固定窗口（ratelimit.Memory，已内置） | Redis（各产品仓 internal/redislimit 实现接口） |
| clientip | 客户端 IP 解析 | 可信代理感知（clientip.Resolver，头序和代理网段可配，已内置） | 无 |
| queue | 任务队列 | 内联执行 | asynq（Redis） |
| mailer | 邮件 | SMTP | 服务商 API |
| captcha | 验证码 | 多驱动校验（Turnstile、reCAPTCHA、hCaptcha、Cap） | 同左 |
| cdn | CDN 刷新 | 钩子接口 | 按服务商实现 |
| oauth | 第三方 OAuth | GitHub/Microsoft profile 归一 | 同左 |
| i18n | 国际化 | go-i18n（Go）、i18next（TS） | 无 |
| mc | Minecraft 领域 | 离线 UUID、玩家名校验等纯函数 | 无 |

## 本地联调

多仓库同时开发时，用 Go workspace 关联本地 checkout，免打 tag：

```
go work init
go work use . ../mythic-mirror ../mythic-hoard
```

## 许可证

[Elastic License 2.0 (ELv2)](./LICENSE)，核心产品代码，与产品仓库一致：源码可得，非 OSI 开源。许可理念与预期用法 FAQ 见文档站：[中文](https://docs.gov.skin/community/licensing/) · [English](https://docs.gov.skin/en/community/licensing/)。

引用、修改、再分发本模块全部免费，包括商业自部署；ELv2 只限制三件事：把软件本身作为托管/管理服务提供给第三方、绕过 license key、移除版权声明。宽松许可的接入点见 [api](https://github.com/Mythicove/api) 仓库（MIT）。

本仓库曾以 AGPL-3.0 发布（v0.1.1、v0.1.2），2026-10-26 起换证为 ELv2。历史与影响（含 Go module proxy 旧版本归属）见 [LICENSE.history](./LICENSE.history)。
