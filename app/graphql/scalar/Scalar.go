package scalar

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// ScalarMap 定義map的基本型態
var ScalarMap = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Map",
	Description: "javascript 的 object, 資料為key:value的格式",
	Serialize: func(value interface{}) interface{} {
		// log.Printf("序列化資料格式 %T, 資料: %v\n", value, value)
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		// log.Printf("解析資料格式 %T, 資料: %v\n", value, value)
		return value
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// log.Printf("解析文字格式 %T, 資料: %v, kind: %v, loc: %v\n", valueAST.GetValue(), valueAST.GetValue(), valueAST.GetKind(), valueAST.GetLoc())
		switch valueAST := valueAST.(type) {
		case *ast.ObjectValue:
			data := map[string]interface{}{}
			for i := range valueAST.Fields {
				field := valueAST.Fields[i]
				key := field.Name.Value
				val := field.Value.GetValue()
				data[key] = val
			}
			return data
		}
		return nil
	},
})
