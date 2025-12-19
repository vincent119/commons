.PHONY: all tidy lint test bench coverage vet fmt generate clean help

# 預設目標：整理、檢查、測試
all: tidy lint test

# =============================================================================
# 基本指令
# =============================================================================

## 整理依賴
tidy:
	go mod tidy

## 格式化程式碼
fmt:
	go fmt ./...
	@echo "✓ 程式碼格式化完成"

## 執行 go vet
vet:
	go vet ./...
	@echo "✓ go vet 檢查通過"

## 執行 go generate（若有使用）
generate:
	go generate ./...
	@echo "✓ go generate 完成"

# =============================================================================
# 品質檢查
# =============================================================================

## 執行 Lint（需安裝 golangci-lint）
lint:
	golangci-lint run --config .golangci.yml ./...
	@echo "✓ golangci-lint 檢查通過"

## 完整檢查（vet + lint）
check: vet lint
	@echo "✓ 所有檢查通過"

# =============================================================================
# 測試相關
# =============================================================================

## 單元測試 + 資料競爭偵測 + 覆蓋率
test:
	go test ./... -race -coverprofile=coverage.out
	@echo "✓ 測試完成，覆蓋率報告已產生: coverage.out"

## 顯示測試覆蓋率報告（終端機顯示）
coverage: test
	go tool cover -func=coverage.out

## 產生 HTML 覆蓋率報告並開啟瀏覽器
coverage-html: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ HTML 覆蓋率報告已產生: coverage.html"
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "請手動開啟 coverage.html"

## 基準測試
bench:
	go test ./... -run=^$$ -bench=. -benchmem

## 只執行特定套件的基準測試（使用 PKG=套件名稱）
bench-pkg:
	go test ./$(PKG)/... -run=^$$ -bench=. -benchmem

# =============================================================================
# 清理
# =============================================================================

## 清理產生的檔案
clean:
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "✓ 清理完成"

# =============================================================================
# 輔助指令
# =============================================================================

## 顯示說明
help:
	@echo "commons 共用模組 Makefile"
	@echo ""
	@echo "使用方式: make [目標]"
	@echo ""
	@echo "基本指令:"
	@echo "  tidy          整理 go.mod 依賴"
	@echo "  fmt           格式化程式碼"
	@echo "  vet           執行 go vet 檢查"
	@echo "  generate      執行 go generate"
	@echo ""
	@echo "品質檢查:"
	@echo "  lint          執行 golangci-lint"
	@echo "  check         執行 vet + lint"
	@echo ""
	@echo "測試相關:"
	@echo "  test          單元測試（含 race 偵測與覆蓋率）"
	@echo "  coverage      顯示覆蓋率報告（終端機）"
	@echo "  coverage-html 產生 HTML 覆蓋率報告"
	@echo "  bench         基準測試"
	@echo "  bench-pkg     特定套件基準測試（PKG=套件名）"
	@echo ""
	@echo "其他:"
	@echo "  clean         清理產生的檔案"
	@echo "  all           tidy + lint + test（預設）"
	@echo "  help          顯示此說明"
