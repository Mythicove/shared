# Mythicove Shared

**English** | [简体中文](./README.zh-CN.md)

Shared Go module for the Mythicove ecosystem. Hosts cross-product capability **interfaces** and domain types; product repositories (`mythic-mirror`, `mythic-hoard`, etc.) depend on it as a Go module.

## Design principles

- **Interfaces, pure domain types, and zero-dependency in-process implementations only** (implementations with external dependencies — DB/S3/Redis/Meilisearch — stay in the product repositories).
- Capabilities are abstracted behind interfaces, so deployments can switch between low-tier and high-tier implementations (10 → 10k scale) without touching business code.
- Package name = capability; see each package's doc comments for semantics.

## Packages

| Package | Capability | Low-tier impl | High-tier impl |
|---|---|---|---|
| `storage` | Object storage | Local disk | MinIO → S3 (R2/OSS) |
| `search` | Search | PostgreSQL `pg_trgm` | Meilisearch |
| `cache` | Cache | In-process memory | Redis |
| `ratelimit` | Rate limiting | In-process fixed window (`ratelimit.Memory`, built in) | Redis (planned) |
| `clientip` | Client IP resolution | Trusted-proxy aware (`clientip.Resolver`, configurable header order / proxy CIDRs, built in) | — |
| `queue` | Task queue | Inline execution | asynq (Redis) |
| `mailer` | Email | SMTP | Provider API |
| `i18n` | Internationalization | go-i18n (Go) / i18next (TS) | — |
| `mc` | Minecraft domain | Pure functions: offline UUIDs, player-name validation, etc. | — |

## Local development

When working across repositories, link local checkouts with a Go workspace (no tagging needed):

```
go work init
go work use . ../mythic-mirror ../mythic-hoard
```

## License

[Elastic License 2.0 (ELv2)](./LICENSE) — core product code, consistent with the product repositories.

This repository was previously released under AGPL-3.0 (v0.1.1–v0.1.2) and relicensed to ELv2 on 2026-10-26; see [LICENSE.history](./LICENSE.history) for details (including which license applies to old versions on the Go module proxy).
