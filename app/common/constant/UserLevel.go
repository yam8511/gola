package constant

import "strconv"

// UserLevel 使用者等級
type UserLevel int64

// 使用者等級清單，在下面 IsValid 記得補上
const (
	MemberUserLevel = UserLevel(0) // 會員
	AdminUserLevel  = UserLevel(1) // 管理者
	SuperUserLevel  = UserLevel(2) // 系統管理者
)

// IsValid 是否有效的等級
func (level UserLevel) IsValid() bool {
	switch level {
	case MemberUserLevel, // 會員
		AdminUserLevel, // 管理者
		SuperUserLevel: // 系統管理者
		return true
	default:
		return false
	}
}

func (level UserLevel) String() string {
	return strconv.FormatInt(int64(level), 10)
}
