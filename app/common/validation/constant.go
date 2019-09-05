package validation

// 驗證用的正則表達
const (
	RegexpEmail    = "^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	RegexpPassword = "^[\\w\\d]{6,14}$"
	RegexpPhone    = "^[0-9]{10}$"
)
