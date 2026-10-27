# Mythicove Shared

**English** | [简体中文](./README.zh-CN.md)

The shared Go module of the Mythicove ecosystem. It holds cross-product capability interfaces and domain types; product repositories (mythic-mirror, mythic-hoard and friends) depend on it as a Go module.

## Design principles

- Interfaces, pure domain types and zero-dependency in-process implementations only. Implementations with external dependencies (DB, S3, Redis, Meilisearch and so on) stay in the product repositories.
- Capabilities sit behind interfaces. Together with the deployment tiers (10 to 10k users), a deployment can switch between the low and high implementation without touching business code.
- The package name is the capability; see each package's doc comments for semantics.

## Packages

| Package | Capability | Low-tier implementation | High-tier implementation |
|---|---|---|---|
| storage | Object storage | Local disk | MinIO, then S3 (R2, OSS) |
| search | Search | PostgreSQL pg_trgm | Meilisearch |
| cache | Cache | In-process memory | Redis |
| ratelimit | Rate limiting | In-process fixed window (ratelimit.Memory, built in) | Redis (product repos implement the interface in internal/redislimit) |
| clientip | Client IP resolution | Trusted-proxy aware (clientip.Resolver, configurable header order and proxy CIDRs, built in) | none |
| queue | Task queue | Inline execution | asynq (Redis) |
| mailer | Email | SMTP | Provider API |
| captcha | CAPTCHA | Multi-driver verification (Turnstile, reCAPTCHA, hCaptcha, Cap) | same |
| cdn | CDN purge | Hook interface | per-provider implementation |
| oauth | Third-party OAuth | GitHub/Microsoft profile normalization | same |
| i18n | Internationalization | go-i18n (Go), i18next (TS) | none |
| mc | Minecraft domain | Pure functions: offline UUIDs, player name validation | none |

## Local development

When working across repositories, link local checkouts with a Go workspace so no tags are needed:

```
go work init
go work use . ../mythic-mirror ../mythic-hoard
```

## License

[Elastic License 2.0 (ELv2)](./LICENSE). Core product code, consistent with the product repositories: source-available, not OSI open source. The licensing philosophy and intended-use FAQ live on the docs site: [English](https://docs.gov.skin/en/community/licensing/) · [中文](https://docs.gov.skin/community/licensing/).

Importing, modifying and redistributing this module are all free, including in commercial self-hosted deployments; ELv2 only restricts offering the software itself to third parties as a hosted/managed service, circumventing license keys, and removing notices. For permissive-license integration points, see the [api](https://github.com/Mythicove/api) repository (MIT).

The repository was released under AGPL-3.0 (v0.1.1, v0.1.2) and moved to ELv2 on 2026-10-26. See [LICENSE.history](./LICENSE.history) for the history and what it means, including which license applies to old versions on the Go module proxy.
