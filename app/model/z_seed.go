package model

import (
	"fmt"
	"gola/app/common/helper"
	"strings"
	"time"
)

// UserSeed 使用者的種子
func UserSeed() (err error) {
	users := []User{
		{Username: "Admin", Name: "Admin", Password: "123456", Enable: true},
	}

	// 收集要新增資料的ID
	userID := []int64{}
	for i := range users {
		user := &users[i]
		userID = append(userID, user.ID)
	}

	db, err := User{}.Database(true).New()
	if err != nil {
		return err
	}

	// 先撈取目前資料表中已有的使用者
	existUser := []User{}
	if err = db.Where("id IN (?)", userID).Find(&existUser).Error; err != nil {
		return
	}

	// 如果資料都已經存在，則直接回傳
	if len(existUser) == len(users) {
		return
	}

	// 將已經存在的使用者ID，重新存進變數
	userID = []int64{}
	for i := range existUser {
		user := &existUser[i]
		userID = append(userID, user.ID)
	}

	// 如果需要新增的資料，不在ID陣列內，則存進變數
	needCreateUsers := []User{}
	for i := range users {
		user := &users[i]
		if !helper.InSliceInt64(userID, user.ID) {
			needCreateUsers = append(needCreateUsers, *user)
		}
	}

	sql := new(strings.Builder)
	sql.WriteString("INSERT INTO `users` (`username`,`name`,`password`,`enable`,`created_at`,`updated_at`) VALUES ")
	userLen := len(needCreateUsers)
	now := time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < userLen; i++ {
		user := &needCreateUsers[i]
		sql.WriteString(fmt.Sprintf(
			"('%s', '%s', '%s', %v, '%s', '%s')",
			user.Username,
			user.Name,
			user.Password,
			user.Enable,
			now,
			now,
		))
		if i != userLen-1 {
			sql.WriteString(",")
		}
	}
	tx := db.Begin()
	err = db.Exec(sql.String()).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
