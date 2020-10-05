package i18n

import (
	"fmt"
	"gola/internal/logger"
	"sync"

	"github.com/tidwall/gjson"
)

// 翻譯書
type Book struct {
	sync.Once
	name    string
	content gjson.Result
	load    func() ([]byte, error)
}

// 新建一個翻譯書
func NewBook(name string, load func() ([]byte, error)) *Book {
	return &Book{
		name: name,
		load: load,
	}
}

// 翻譯
func (b *Book) Trans(code string) string {
	b.Do(func() {
		if b.load != nil {
			content, err := b.load()
			if err != nil {
				logger.Error(fmt.Errorf("載入[%s]字典檔失敗: %w", b.name, err))
			}
			b.content = gjson.GetBytes(content, "..0")
		}
	})
	return b.content.Get(code).String()
}
