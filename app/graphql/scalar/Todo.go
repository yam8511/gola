package scalar

import (
	"github.com/graphql-go/graphql"
)

// ScalarTodo Todo物件
var ScalarTodo = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Todo",
	Description: "Todo物件",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Description: "ID",
			Type:        graphql.NewNonNull(graphql.Int),
		},
		"text": &graphql.Field{
			Description: "項目內容",
			Type:        graphql.String,
		},
		"done": &graphql.Field{
			Description: "是否完成",
			Type:        graphql.Boolean,
		},
		"expired_at": &graphql.Field{
			Description: "有效期限",
			Type:        graphql.DateTime,
		},
		"labels": &graphql.Field{
			Description: "標籤",
			Type:        ScalarMap,
		},
	},
})

// ScalarTodoInput TODO的輸入物件
var ScalarTodoInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "TodoInput",
	Description: "Todo的輸入物件",
	Fields: graphql.InputObjectConfigFieldMap{
		"text": &graphql.InputObjectFieldConfig{
			Description: "項目內容",
			Type:        graphql.NewNonNull(graphql.String),
		},
		"done": &graphql.InputObjectFieldConfig{
			Description: "是否完成",
			Type:        graphql.NewNonNull(graphql.Boolean),
		},
		"expired_at": &graphql.InputObjectFieldConfig{
			Description: "有效期限",
			Type:        graphql.String,
		},
		"labels": &graphql.InputObjectFieldConfig{
			Description: "標籤",
			Type:        ScalarMap,
		},
	},
})
