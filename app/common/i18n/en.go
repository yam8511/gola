package i18n

import (
	"encoding/json"
	"gola/app/common/errorcode"

	"github.com/tidwall/gjson"
)

var enBook book

func loadDictEN() {
	dict := map[string]interface{}{
		"hello": "Hello~",
	}

	book, err := json.Marshal(dict)
	if err != nil {
		errorcode.GetAPIError("載入『英文』語系的字典檔失敗!!!", nil)
		return
	}
	enBook.content = book
}

func transEN(code string) string {
	enBook.Do(loadDictEN)
	result := gjson.GetBytes(enBook.content, code)
	return result.String()
}
