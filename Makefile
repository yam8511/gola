PWD := $(shell pwd)
COMMA := ,
FLAG_LIST := $(subst $(COMMA), ,$(FLAG))
APP := gola

.PHONY: docs

docs: # 產生swagger文檔
	rm -r docs || exit 0
	swag init -g main.go -d .
	rm docs/swagger.*
	sed "s/^func init/func Init/" docs/docs.go > docs/gola.go
	mv docs/gola.go docs/docs.go

setup: # 開啟DB, Cache資料庫
	docker-compose up -d db cache

# 單元測試
test:
	PROJECT_ROOT=$(PWD) go test $(FLAG_LIST) -coverprofile=./cover.out ./app/...
	go tool cover -html cover.out
	rm cover.out

# 編譯
build: vendor
	go build -mod vendor $(FLAG_LIST) -o $(APP)

# 依賴套件
vendor: go.sum
	rm go.sum
	GOPROXY= go mod tidy
	GOPROXY= go mod vendor

clean:
	go clean -cache
	go clean -modcache
	rm go.sum || exit 0
	rm -r vendor || exit 0

# Go套件管理
go.mod:
	go mod init $(APP)

go.sum: go.mod
	go mod tidy
