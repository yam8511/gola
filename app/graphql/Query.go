package graphql

import (
	"gola/app/graphql/handler"
	"gola/app/graphql/scalar"

	"github.com/graphql-go/graphql"
)

// 配置Query用的欄位
func setupQueryFields() graphql.Fields {
	return graphql.Fields{
		"ping": &graphql.Field{
			Description: "測試連線",
			Type:        graphql.String,
			Resolve:     handler.Ping,
		},
		"try_error": &graphql.Field{
			Description: "測試GraphQL的錯誤",
			Type:        graphql.String,
			Resolve:     handler.TryError,
		},
		"todo": &graphql.Field{
			Description: "指定ID的Todo",
			Type:        scalar.ScalarTodo,
			Resolve:     handler.GetTodo,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Todo的ID",
				},
			},
		},
		"todos": &graphql.Field{
			Description: "Todo列表",
			Type:        graphql.NewList(scalar.ScalarTodo),
			Resolve:     handler.GetTodos,
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "過濾條件：是否完成",
				},
			},
		},
	}
}
