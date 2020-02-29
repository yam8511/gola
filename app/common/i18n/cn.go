package i18n

import (
	"encoding/json"
	"gola/app/common/errorcode"

	"github.com/tidwall/gjson"
)

var cnBook book

func loadDictCN() {
	dict := map[string]interface{}{
		"hello": "你好",
	}

	book, err := json.Marshal(dict)
	if err != nil {
		errorcode.GetAPIError("載入『簡體中文』語系的字典檔失敗!!!", nil)
		return
	}
	cnBook.content = book
}

func transCN(code string) string {
	cnBook.Do(loadDictCN)
	result := gjson.GetBytes(cnBook.content, code)
	return result.String()
}
