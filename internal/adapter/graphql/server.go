
package graphqlsrv

import (
	"encoding/json"
	"net/http"

	"github.com/example/clean-arch-orders/internal/domain"
	"github.com/example/clean-arch-orders/internal/usecase"
	"github.com/graphql-go/graphql"
)

type Server struct {
	ListUC *usecase.ListOrders
}

func NewServer(list *usecase.ListOrders) *Server { return &Server{ListUC: list} }

func (s *Server) Handler() http.Handler {
	orderType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id":            &graphql.Field{Type: graphql.String},
			"customer_name": &graphql.Field{Type: graphql.String},
			"total_amount":  &graphql.Field{Type: graphql.Float},
			"created_at":    &graphql.Field{Type: graphql.String},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					res, err := s.ListUC.Execute()
					if err != nil { return nil, err }
					orders := res.([]domain.Order)
					return orders, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params struct{ Query string `json:"query"` }
		if r.Method == http.MethodGet {
			params.Query = r.URL.Query().Get("query")
		} else {
			_ = json.NewDecoder(r.Body).Decode(&params)
		}
		result := graphql.Do(graphql.Params{Schema: schema, RequestString: params.Query})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}
