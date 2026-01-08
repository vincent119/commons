# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.0] - 2026-01-08

### Added
- **graceful**: 新增 `graceful` 套件，提供應用程式優雅關機與生命週期管理功能。
    - 支援 `Run` 與 `HTTPTask` 統一管理啟動與關機。
    - 支援 `WithTimeout` 設定關機超時。
    - 支援 `WithCleanup` (LIFO) 與 `WithCloser` 資源釋放。
    - 支援 `WithLogger` (支援 `log/slog`)。
    - 完整繁體中文註解與文件。

## [v0.1.5] - Previous Releases
- 包含 `awsx`, `httpx`, `slicex`, `sqlx`, `stringx` 等通用工具庫。
