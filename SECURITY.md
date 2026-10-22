# Security Policy

本仓库（共享 Go 模块）的完整安全策略见 [Mythicove/docs](https://github.com/Mythicove/docs/blob/main/SECURITY.md)。

## 依赖基线

- 当前**零第三方依赖**（纯标准库，无 `go.sum`）。
- 新增依赖须通过 CI 的 `govulncheck`；若 `klauspost/compress` 被间接引入，须 ≥ v1.18.7。

## 上报

使用本仓库 Security 标签下的「Private vulnerability reporting」，或按 docs 仓库的 SECURITY.md 流程。
