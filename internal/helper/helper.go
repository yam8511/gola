package helper

import (
	"gola/internal/bootstrap"
	"strconv"
	"time"
)

// TaipeiZone 台北時區
func TaipeiZone() *time.Location {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		bootstrap.WriteLog("ERROR", "TaipeiZone: 載入時區錯誤, "+err.Error())
	}
	return loc
}

// StrToInt64 字串轉數字
func StrToInt64(str string) (int64, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	return num, err
}

// ParseTime 解析時間字串
func ParseTime(st string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, st)
	if err != nil {
		var loc *time.Location
		loc, err = time.LoadLocation("Asia/Taipei")
		if err != nil {
			return
		}
		t, err = time.ParseInLocation("2006-01-02 15:04:05", st, loc)
	}
	return
}

// UniqueSliceString Unique字串陣列
func UniqueSliceString(slice []string) (unique []string) {
	tmp := map[string]byte{}
	for i := range slice {
		data := slice[i]
		_, ok := tmp[data]
		if !ok {
			tmp[data] = 0
			unique = append(unique, data)
		}
	}
	return
}

// UniqueSliceInt64 Unique int64陣列
func UniqueSliceInt64(slice []int64) (unique []int64) {
	tmp := map[int64]byte{}
	for i := range slice {
		data := slice[i]
		_, ok := tmp[data]
		if !ok {
			tmp[data] = 0
			unique = append(unique, data)
		}
	}
	return
}

// InArrayInt64 指定值是否在陣列中
func InArrayInt64(array []int64, target int64) bool {
	for i := range array {
		if target == array[i] {
			return true
		}
	}
	return false
}

// InArrayInt 指定值是否在陣列中
func InArrayInt(array []int, target int) bool {
	for i := range array {
		if target == array[i] {
			return true
		}
	}
	return false
}
