package errorcode
import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)
// GetGqlError 由錯誤碼取得Gql回傳
func GetGqlError(p graphql.ResolveParams, text string, err error) gqlerrors.FormattedError {

	// type FormattedError struct {
	// 	Message       string                    `json:"message"`
	// 	Locations     []location.SourceLocation `json:"locations"`
	// 	Path          []interface{}             `json:"path,omitempty"`
	// 	Extensions    map[string]interface{}    `json:"extensions,omitempty"`
	// 	originalError error
	// }

	apiErr := gqlerrors.NewFormattedError(text)
	if p.Info.Path != nil {
		apiErr.Path = []interface{}{p.Info.Path.Key}
	}

	return apiErr
}
