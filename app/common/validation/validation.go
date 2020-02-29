package validation

import (
	"fmt"
	"gola/internal/logger"
	"regexp"
	"strings"
)

// IsEmailPassed 驗證Email
func IsEmailPassed(email string) (isPassed bool) {
	email = strings.TrimSpace(email)
	if email == "" {
		isPassed = true
		return
	}
	reg, err := regexp.Compile(RegexpEmail)
	if err != nil {
		go logger.Warn(fmt.Sprintf("編譯[信箱]的正則表達式錯誤 /%s/ ---> %s", RegexpEmail, err.Error()))
		return
	}
	isPassed = reg.MatchString(email)
	return
}

// IsPhonePassed 驗證Phone
func IsPhonePassed(phone string) (isPassed bool) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		isPassed = true
		return
	}
	reg, err := regexp.Compile(RegexpPhone)
	if err != nil {
		go logger.Warn(fmt.Sprintf("編譯[手機]的正則表達式錯誤 /%s/ ---> %s", RegexpPhone, err.Error()))
		return
	}
	isPassed = reg.MatchString(phone)
	return
}

// IsPasswordPassed 驗證Password
func IsPasswordPassed(password string) (isPassed bool) {
	password = strings.TrimSpace(password)
	if password == "" {
		isPassed = true
		return
	}
	reg, err := regexp.Compile(RegexpPassword)
	if err != nil {
		go logger.Warn(fmt.Sprintf("編譯[密碼]的正則表達式錯誤 /%s/ ---> %s", RegexpPassword, err.Error()))
		return
	}
	isPassed = reg.MatchString(password)
	return
}
