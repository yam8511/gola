package graphql

import (
	"github.com/graphql-go/graphql"
)

// SetupGraphQLSchema 配置 GraphQL 結構語法
func SetupGraphQLSchema() (graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:        "RootQuery",
			Fields:      setupQueryFields(),
			Description: "讀取資料的API接口",
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:        "RootMutation",
			Fields:      setupMutationFields(),
			Description: "寫入資料的API接口",
		}),
	}
	return graphql.NewSchema(schemaConfig)
}
