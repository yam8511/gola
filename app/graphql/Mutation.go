package graphql

import (
	"gola/app/graphql/handler"
	"gola/app/graphql/scalar"

	"github.com/graphql-go/graphql"
)

// 配置Mutation用的欄位
func setupMutationFields() graphql.Fields {
	return graphql.Fields{
		"add_todo": &graphql.Field{
			Description: "新增Todo",
			Type:        scalar.ScalarTodo,
			Resolve:     handler.AddTodo,
			Args: graphql.FieldConfigArgument{
				"data": &graphql.ArgumentConfig{
					Type:        (scalar.ScalarTodoInput),
					Description: "輸入Todo資料",
				},
			},
		},
		"remove_todo": &graphql.Field{
			Description: "移除Todo",
			Type:        scalar.ScalarTodo,
			Resolve:     handler.RemoveTodo,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Todo的ID",
				},
			},
		},
		// "customer_edit": &graphql.Field{
		// 	Description: "編輯客戶資料",
		// 	Type:        graphql.NewList(scalar.ScalarCustomer),
		// 	Resolve:     handler.EditCustomerResolveFn,
		// 	Args: graphql.FieldConfigArgument{
		// 		"data": &graphql.ArgumentConfig{
		// 			Type:        graphql.NewNonNull(graphql.NewList(scalar.ScalarCustomerEditInput)),
		// 			Description: "輸入客戶資料",
		// 		},
		// 	},
		// },
		// "customer_delete": &graphql.Field{
		// 	Description: "刪除客戶資料",
		// 	Type:        graphql.NewList(scalar.ScalarCustomer),
		// 	Resolve:     handler.DeleteCustomerResolveFn,
		// 	Args: graphql.FieldConfigArgument{
		// 		"id": &graphql.ArgumentConfig{
		// 			Type:        graphql.NewNonNull(graphql.NewList(graphql.Int)),
		// 			Description: "輸入客戶ID",
		// 		},
		// 	},
		// },
	}
}
