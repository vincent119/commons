.PHONY: all tidy lint test bench
all: tidy lint test

# 整理依賴
tidy:
	go mod tidy

# 執行 Lint（需安裝 golangci-lint）
lint:
	golangci-lint run --config .golangci.yml ./...

# 單元測試 + 資料競爭偵測 + 覆蓋率
test:
	go test ./... -race -coverprofile=coverage.out

# 基準測試
bench:
	go test ./... -run=^$ -bench=. -benchmem