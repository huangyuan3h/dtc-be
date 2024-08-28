package schema

import (
	"github.com/graphql-go/graphql"
)

// 创建 Schema
func CreateSchema() (graphql.Schema, error) {
    fields := graphql.Fields{
        "hello": &graphql.Field{
            Type: graphql.String,
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                return "world", nil
            },
        },
    }

    // 添加 Product 查询字段
    for key, value := range getProductQueryFields() {
        fields[key] = value
    }

    // 创建 schema
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
    schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
    return graphql.NewSchema(schemaConfig)
}