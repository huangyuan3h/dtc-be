package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
)

// 定义 Product 类型
type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

// 模拟产品数据
var products = []Product{
    {ID: "1", Name: "T-Shirt", Description: "Comfortable cotton T-shirt", Price: 19.99},
    {ID: "2", Name: "Jeans", Description: "Classic denim jeans", Price: 49.99},
    {ID: "3", Name: "Sneakers", Description: "Stylish and comfortable sneakers", Price: 79.99},
}

// 创建 Product 对象类型
var productType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Product",
    Fields: graphql.Fields{
        "id":          &graphql.Field{Type: graphql.String},
        "name":        &graphql.Field{Type: graphql.String},
        "description": &graphql.Field{Type: graphql.String},
        "price":       &graphql.Field{Type: graphql.Float},
    },
})

// 获取 Product 查询字段
func getProductQueryFields() graphql.Fields {
    return graphql.Fields{
        // 查询所有产品
        "products": &graphql.Field{
            Type: graphql.NewList(productType),
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                return products, nil
            },
        },
        // 根据 ID 查询单个产品
        "product": &graphql.Field{
            Type: productType,
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id, ok := p.Args["id"].(string)
                if !ok {
                    return nil, errors.New("invalid product ID")
                }
                for _, product := range products {
                    if product.ID == id {
                        return product, nil
                    }
                }
                return nil, nil // 未找到产品
            },
        },
    }
}
