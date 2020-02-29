.PHONY: docs

docs: # 產生swagger文檔
	rm -r docs || exit 0
	swag init -g main.go -d .
	rm docs/swagger.*
	sed "s/^func init/func Init/" docs/docs.go > docs/gola.go
	mv docs/gola.go docs/docs.go

setup: # 開啟DB, Cache資料庫
	docker-compose up -d db cache
