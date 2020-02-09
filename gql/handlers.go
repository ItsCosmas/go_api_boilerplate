package gql

import (
	"go_api_boilerplate/gql/gen"

	"go_api_boilerplate/services/authservice"
	"go_api_boilerplate/services/userservice"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(us userservice.UserService, as authservice.AuthService) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	conf := gen.Config{
		Resolvers: &Resolver{
			UserService: us,
			AuthService: as,
		},
	}
	exec := gen.NewExecutableSchema(conf)
	h := handler.GraphQL(exec)
	return func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
}

// PlaygroundHandler Defines the Playground handler to expose our playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := handler.Playground("GraphQL Playground", path)
	return func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
}
