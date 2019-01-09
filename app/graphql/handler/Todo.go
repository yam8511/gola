package handler

import (
	"encoding/json"
	"gola/app/common/data_struct"
	"time"

	"gola/app/business"
	"gola/app/common/error_code"

	"github.com/graphql-go/graphql"
)

// GetTodo 取Todo
func GetTodo(p graphql.ResolveParams) (interface{}, error) {
	argID, ok := p.Args["id"]
	if !ok {
		return nil, errorcode.GetGqlError(p, "param_required", nil)
	}

	var todoID int
	todoID, ok = argID.(int)
	if !ok {
		return nil, errorcode.GetGqlError(p, "param_invalid", nil)
	}
	todo := business.GetTodo(int64(todoID))

	return todo, nil
}

// GetTodos Todo列表
func GetTodos(p graphql.ResolveParams) (interface{}, error) {
	var done *bool
	argDone, ok := p.Args["done"]
	if ok {
		done = new(bool)
		*done, ok = argDone.(bool)
		if !ok {
			return nil, errorcode.GetGqlError(p, "param_invalid", nil)
		}
	}

	todos := business.GetTodos(done)

	return todos, nil
}

// AddTodo 新增Todo
func AddTodo(p graphql.ResolveParams) (interface{}, error) {
	argData, ok := p.Args["data"]
	if !ok {
		return nil, errorcode.GetGqlError(p, "param_required", nil)
	}

	argByte, err := json.Marshal(argData)
	if err != nil {
		return nil, errorcode.GetGqlError(p, "parse_err", err)
	}

	var input *datastruct.TodoInput
	err = json.Unmarshal(argByte, &input)
	if err != nil {
		return nil, errorcode.GetGqlError(p, "param_invalid", err)
	}

	if input == nil {
		return nil, errorcode.GetGqlError(p, "param_required", nil)
	}

	var expiredAt *time.Time
	if input.ExpiredAt != "" {
		*expiredAt, err = time.Parse("2006-01-02 15:04:05", input.ExpiredAt)
		if err != nil {
			return nil, errorcode.GetGqlError(p, "param_invalid", err)
		}
	}

	todo := business.AddTodo(input.Text, input.Done, expiredAt, input.Labels)
	return todo, nil
}

// RemoveTodo 移除Todo
func RemoveTodo(p graphql.ResolveParams) (interface{}, error) {
	argID, ok := p.Args["id"]
	if !ok {
		return nil, errorcode.GetGqlError(p, "param_required", nil)
	}

	var todoID int
	todoID, ok = argID.(int)
	if !ok {
		return nil, errorcode.GetGqlError(p, "param_invalid", nil)
	}

	todo := business.RemoveTodo(int64(todoID))
	return todo, nil
}
