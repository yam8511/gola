package router

import (
	"gola/app/common/error_code"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	graphqlH "github.com/graphql-go/handler"
)

// LoadGraphQLRouter 載入GraphQL的路由
func LoadGraphQLRouter(r *gin.Engine, schema graphql.Schema) {
	// 設置 graphiql
	h := graphqlH.New(&graphqlH.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
		// 可以取得連線進來的請求
		// RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
		// 	fmt.Println("=======================")
		// 	fmt.Println("Request", r)
		// 	fmt.Println("=======================")
		// 	return nil
		// },
		// 可以取請求參數與回傳結果，像是Logger
		// ResultCallbackFn: func(ctx context.Context, params *graphql.Params, result *graphql.Result, responseBody []byte) {
		// },
		// 客製化錯誤回傳
		FormatErrorFn: func(err error) (f gqlerrors.FormattedError) {
			if err != nil {
				f = gqlerrors.FormatError(err)
				apiErr := errorcode.GetAPIError(err.Error(), nil)
				f.Message = apiErr.Error()
				f.Extensions = map[string]interface{}{
					"error_code": apiErr.ErrorCode(),
					"error_text": apiErr.ErrorText(),
				}
				if len(f.Path) > 0 {
					f.Extensions["path"] = f.Path[0]
				}
			}
			return
		},
	})

	r.Any("/graphql", gin.WrapH(h))
}
