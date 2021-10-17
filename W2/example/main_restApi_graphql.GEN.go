package main

import (
	"github.com/graphql-go/graphql"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/rqAuth"
)

// can be hit using with /graphql

var graphqlQueries = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		`UserList`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserListOut`,
				Fields: graphql.Fields{
					`Limit`: &graphql.Field{
						Type: graphql.Int, // uint32
					},
					`Offset`: &graphql.Field{
						Type: graphql.Int, // uint32
					},
					`Total`: &graphql.Field{
						Type: graphql.Int, // uint32
					},
					`Users`: &graphql.Field{
						Type: graphql.NewList(rqAuth.GraphqlTypeUsers), //  []rqAuth.Users
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`Limit`: &graphql.ArgumentConfig{
					Type: graphql.Int, // uint32
				},
				`Offset`: &graphql.ArgumentConfig{
					Type: graphql.Int, // uint32
				},
			}, // UserListIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserProfile`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserProfileOut`,
				Fields: graphql.Fields{
					`User`: &graphql.Field{
						Type: rqAuth.GraphqlTypeUsers, // rqAuth.Users
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`sessionToken`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserProfileIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
	},
})

var graphqlMutations = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		`UserChangeEmail`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   `UserChangeEmailOut`,
				Fields: graphql.Fields{},
			}),
			Args: graphql.FieldConfigArgument{}, // UserChangeEmailIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserChangePassword`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserChangePasswordOut`,
				Fields: graphql.Fields{
					`UpdatedAt`: &graphql.Field{
						Type: graphql.Int, // int64
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`Password`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`NewPassword`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`sessionToken`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserChangePasswordIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserConfirmEmail`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   `UserConfirmEmailOut`,
				Fields: graphql.Fields{},
			}),
			Args: graphql.FieldConfigArgument{}, // UserConfirmEmailIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserForgotPassword`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserForgotPasswordOut`,
				Fields: graphql.Fields{
					`Ok`: &graphql.Field{
						Type: graphql.Boolean, // bool
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`Email`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`ChangePassCallback`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserForgotPasswordIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserLogin`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserLoginOut`,
				Fields: graphql.Fields{
					`WalletId`: &graphql.Field{
						Type: graphql.String, // string
					},
					`sessionToken`: &graphql.Field{
						Type: graphql.String, // string
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`Email`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`Password`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserLoginIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserLogout`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserLogoutOut`,
				Fields: graphql.Fields{
					`LoggedOut`: &graphql.Field{
						Type: graphql.Boolean, // bool
					},
					`sessionToken`: &graphql.Field{
						Type: graphql.String, // string
					},
				},
			}),
			Args: graphql.FieldConfigArgument{}, // UserLogoutIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserRegister`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserRegisterOut`,
				Fields: graphql.Fields{
					`CreatedAt`: &graphql.Field{
						Type: graphql.Int, // int64
					},
					`UserId`: &graphql.Field{
						Type: graphql.Int, // uint64
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`UserName`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`Email`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`Password`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserRegisterIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		`UserResetPassword`: &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: `UserResetPasswordOut`,
				Fields: graphql.Fields{
					`Ok`: &graphql.Field{
						Type: graphql.Boolean, // bool
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				`Password`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`SecretCode`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
				`Hash`: &graphql.ArgumentConfig{
					Type: graphql.String, // string
				},
			}, // UserResetPasswordIn
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
	},
})

var graphqlSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    graphqlQueries,
	Mutation: graphqlMutations,
})

func handleGraphql(in *GraphqlRequest, out *GraphqlResponse) {
	params := graphql.Params{Schema: graphqlSchema, RequestString: in.Query}
	out.Result = graphql.Do(params)
}
