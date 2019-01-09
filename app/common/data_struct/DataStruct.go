package datastruct

import "time"

// Todo Todo物件
type Todo struct {
	ID        int64                  `json:"id"`
	Text      string                 `json:"text"`
	Done      bool                   `json:"done"`
	ExpiredAt *time.Time             `json:"expired_at"`
	Labels    map[string]interface{} `json:"labels"`
}

// TodoInput Todo物件
type TodoInput struct {
	ID        int64                  `json:"id"`
	Text      string                 `json:"text"`
	Done      bool                   `json:"done"`
	ExpiredAt string                 `json:"expired_at"`
	Labels    map[string]interface{} `json:"labels"`
}
