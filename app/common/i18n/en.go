package i18n

import (
	"encoding/json"
)

var enBook = NewBook("英文", func() ([]byte, error) {
	dict := map[string]interface{}{
		"hello": "Hello~",
	}
	return json.Marshal(dict)
})
