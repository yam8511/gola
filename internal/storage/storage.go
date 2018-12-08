package storage

import (
	"gola/internal/bootstrap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// CreateFile 建立檔案
func CreateFile(filename string, data []byte) error {
	filename = bootstrap.GetAppRoot() + "/storage/app/" + filename
	//檢查檔案是否存在
	if !CheckFileExists(filename) {
		//建立資料夾
		err := os.MkdirAll(filepath.Dir(filename), 0777)
		if err != nil {
			go bootstrap.WriteLog("ERROR", "CreateFile: 建立相關資料夾錯誤, "+err.Error())
			return err
		}
	}

	out, err := os.Create(filename)
	if err != nil {
		go bootstrap.WriteLog("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		return err
	}
	defer out.Close()

	return nil
}

// CopyFile 複製檔案
func CopyFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		go bootstrap.WriteLog("ERROR", "CopyFile: 開啟檔案錯誤, "+err.Error())
		return err
	}
	defer src.Close()

	filename := bootstrap.GetAppRoot() + "/storage/app/" + dst
	//檢查檔案是否存在
	if !CheckFileExists(filename) {
		//建立資料夾
		err = os.MkdirAll(filepath.Dir(filename), 0777)
		if err != nil {
			go bootstrap.WriteLog("ERROR", "CopyFile: 建立相關資料夾錯誤, "+err.Error())
			return err
		}
	}

	out, err := os.Create(filename)
	if err != nil {
		go bootstrap.WriteLog("ERROR", "CopyFile: 建立檔案錯誤, "+err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		go bootstrap.WriteLog("ERROR", "CopyFile: 複製檔案錯誤, "+err.Error())
		return err
	}
	return nil
}

// DeleteFile 刪除檔案
func DeleteFile(src string) {
	filename := bootstrap.GetAppRoot() + "/storage/app/" + src
	if CheckFileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
			go bootstrap.WriteLog("WARNING", "DeleteFile: 刪除檔案錯誤, "+err.Error())
		}
	}
}

// GetStorageFilePath 取檔案路徑
func GetStorageFilePath(filename string) string {
	filename = strings.TrimSpace(filename)
	if strings.HasPrefix(filename, "/") {
		filename = bootstrap.GetAppRoot() + "/storage" + filename
	} else {
		filename = bootstrap.GetAppRoot() + "/storage/" + filename
	}
	return filename
}

// GetStorageAppFilePath 取app底下的檔案路徑
func GetStorageAppFilePath(filename string) string {
	filename = strings.TrimSpace(filename)
	if strings.HasPrefix(filename, "/") {
		filename = bootstrap.GetAppRoot() + "/storage/app" + filename
	} else {
		filename = bootstrap.GetAppRoot() + "/storage/app/" + filename
	}
	return filename
}

// CheckFileExists 確認檔案是否存在
func CheckFileExists(filename string) (ok bool) {
	_, err := os.Stat(filename)
	if err == nil {
		ok = true
		return
	} else if !os.IsNotExist(err) {
		go bootstrap.WriteLog("WARNING", "CheckFileExists: 檢查檔案失敗, "+err.Error())
	}
	return
}
