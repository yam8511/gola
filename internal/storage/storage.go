package storage

import (
	"bytes"
	"errors"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// WriteFile 開啟檔案，若檔案不存在，系統會自動建立
func WriteFile(filename string, data []byte) error {
	filename = GetStorageAppFilePath(filename)
CREATE:
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		if os.IsNotExist(err) {
			//建立資料夾
			err := os.MkdirAll(filepath.Dir(filename), 0777)
			if err != nil {
				logger.Error(errors.New("storage.WriteFile: 建立資料夾錯誤: " + err.Error()))
				return err
			}
			goto CREATE
		}
		logger.Error(errors.New("storage.WriteFile: 開啟檔案錯誤: " + err.Error()))
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		if os.IsNotExist(err) {
			goto CREATE
		}
		logger.Error(errors.New("storage.WriteFile: 寫入檔案錯誤: " + err.Error()))
		goto CLOSE
	}
CLOSE:
	f.Close()

	return err
}

// CopyFile 複製檔案
func CopyFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		logger.Error(errors.New("storage.CopyFile: 開啟檔案錯誤: " + err.Error()))
		return err
	}
	defer src.Close()

	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, src)
	if err != nil {
		logger.Error(errors.New("storage.CopyFile: 複製檔案錯誤: " + err.Error()))
		return err
	}

	return WriteFile(dst, buf.Bytes())
}

// DeleteFile 刪除檔案
func DeleteFile(filename string) {
	filename = GetStorageAppFilePath(filename)
	if CheckFileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
			if !os.IsNotExist(err) {
				logger.Warn("storage.DeleteFile: 刪除檔案錯誤: " + err.Error())
			}
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
		logger.Warn("storage.CheckFileExists: 檢查檔案失敗: " + err.Error())
	}
	return
}
