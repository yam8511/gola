package router

import (
	"gola/app/graphql"

	"github.com/gin-gonic/gin"
)

// RouteProvider 路由提供者
func RouteProvider(r *gin.Engine) {
	LoadWebRouter(r)
	LoadAPIRouter(r)

	schema, err := graphql.SetupGraphQLSchema()
	if err != nil {
		panic(err)
	}
	LoadGraphQLRouter(r, schema)
}
