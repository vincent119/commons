# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.2] - 2026-01-12

### Added

- **timex**: 新增 `FormatISO8601` 函式，支援帶毫秒與緊湊時區的 ISO 8601 格式 (e.g. `2026-01-12T18:09:11.000+0800`)。

## [v0.2.1] - 2026-01-08

### Added

- **graceful**: 新增 `WithClosers` API，批量註冊多個 `io.Closer` 資源。
- **ci**: 新增 GitHub Actions workflow，自動執行測試並更新覆蓋率徽章。

### Changed

- **graceful**: 將程式碼註解改為英文，提升國際化可讀性。

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
