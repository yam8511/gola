package handler

import (
	"fmt"
	"time"

	"gola/app/common/error_code"

	"github.com/graphql-go/graphql"
)

// Ping 測試連線
func Ping(p graphql.ResolveParams) (interface{}, error) {
	return fmt.Sprintf("pong at %s", time.Now().Format(time.RFC3339Nano)), nil
}

// TryError 測試錯誤
func TryError(p graphql.ResolveParams) (interface{}, error) {
	apiErr := errorcode.GetGqlError(p, "ping", nil)
	return fmt.Sprintf("pong at %s", time.Now().Format(time.RFC3339Nano)), apiErr
}
