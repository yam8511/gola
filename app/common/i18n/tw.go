package i18n

import (
	"encoding/json"
	errorcode "gola/app/common/errorcode"

	"github.com/tidwall/gjson"
)

var twBook book

func loadDictTW() {
	dict := map[string]interface{}{
		"hello": "你好",
	}

	book, err := json.Marshal(dict)
	if err != nil {
		errorcode.Code_Undefined.New("載入『繁體中文』語系的字典檔失敗!!! %w", err)
		return
	}
	twBook.content = book
}

func transTW(code string) string {
	twBook.Do(loadDictTW)
	result := gjson.GetBytes(twBook.content, code)
	return result.String()
}
