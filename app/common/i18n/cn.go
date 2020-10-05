package i18n

import (
	"encoding/json"
)

var cnBook = NewBook("簡體中文", func() ([]byte, error) {
	dict := map[string]interface{}{
		"hello": "你好",
	}
	return json.Marshal(dict)
})
