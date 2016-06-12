package graphql_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"golang.org/x/net/context"
)

type T struct {
	Query    string
	Schema   graphql.Schema
	Expected interface{}
}

var Tests = []T{}

func init() {
	Tests = []T{
		{
			Query: `
				query HeroNameQuery {
					hero {
						name
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"name": "R2-D2",
					},
				},
			},
		},
		{
			Query: `
				query HeroNameAndFriendsQuery {
					hero {
						id
						name
						friends {
							name
						}
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"hero": map[string]interface{}{
						"id":   "2001",
						"name": "R2-D2",
						"friends": []interface{}{
							map[string]interface{}{
								"name": "Luke Skywalker",
							},
							map[string]interface{}{
								"name": "Han Solo",
							},
							map[string]interface{}{
								"name": "Leia Organa",
							},
						},
					},
				},
			},
		},
		{
			Query: `
				{
					hero {
						id
						name
						friends {
							name
						}
						appearsIn
					}
					vader: human(id: "1001") {
						id
						name
						friends {
							id
							name
							appearsIn
						}
						appearsIn
						homePlanet
					}
					threepio: droid(id: "2000") {
						id
						name
						friends {
							id
							name
							appearsIn
						}
						appearsIn
						primaryFunction
					}
				}
			`,
			Schema: testutil.StarWarsSchema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"threepio": map[string]interface{}{
						"primaryFunction": "Protocol",
						"id":              "2000",
						"appearsIn": []interface{}{
							"NEWHOPE",
							"EMPIRE",
							"JEDI",
						},
						"name": "C-3PO",
						"friends": []interface{}{
							map[string]interface{}{
								"id":   "1000",
								"name": "Luke Skywalker",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
							},
							map[string]interface{}{
								"name": "Han Solo",
								"id":   "1002",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
							},
							map[string]interface{}{
								"name": "Leia Organa",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
								"id": "1003",
							},
							map[string]interface{}{
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
								"id":   "2001",
								"name": "R2-D2",
							},
						},
					},
					"vader": map[string]interface{}{
						"id": "1001",
						"appearsIn": []interface{}{
							"NEWHOPE",
							"EMPIRE",
							"JEDI",
						},
						"friends": []interface{}{
							map[string]interface{}{
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
								"name": "Han Solo",
								"id":   "1002",
							},
							map[string]interface{}{
								"id":   "1003",
								"name": "Leia Organa",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
							},
							map[string]interface{}{
								"id":   "2000",
								"name": "C-3PO",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
							},
							map[string]interface{}{
								"id": "2001",
								"appearsIn": []interface{}{
									"NEWHOPE",
									"EMPIRE",
									"JEDI",
								},
								"name": "R2-D2",
							},
							map[string]interface{}{
								"appearsIn": []interface{}{
									"NEWHOPE",
								},
								"id":   "1004",
								"name": "Wilhuff Tarkin",
							},
						},
						"homePlanet": "Tatooine",
						"name":       "Darth Vader",
					},
					"hero": map[string]interface{}{
						"name": "R2-D2",
						"id":   "2001",
						"appearsIn": []interface{}{
							"NEWHOPE",
							"EMPIRE",
							"JEDI",
						},
						"friends": []interface{}{
							map[string]interface{}{
								"name": "Luke Skywalker",
							},
							map[string]interface{}{
								"name": "Han Solo",
							},
							map[string]interface{}{
								"name": "Leia Organa",
							},
						},
					},
				},
			},
		},
	}
}

func TestQuery(t *testing.T) {
	for _, test := range Tests {
		params := graphql.Params{
			Schema:        test.Schema,
			RequestString: test.Query,
		}
		testGraphql(test, params, t)
	}
}

func testGraphql(test T, p graphql.Params, t *testing.T) {
	result := graphql.Do(p)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result, test.Expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", test.Query, testutil.Diff(test.Expected, result))
	}
}

func TestBasicGraphQLExample(t *testing.T) {
	// taken from `graphql-js` README

	helloFieldResolved := func(p graphql.ResolveParams) (interface{}, error) {
		return "world", nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQueryType",
			Fields: graphql.Fields{
				"hello": &graphql.Field{
					Description: "Returns `world`",
					Type:        graphql.String,
					Resolve:     helloFieldResolved,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := "{ hello }"
	var expected interface{}
	expected = map[string]interface{}{
		"hello": "world",
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}

func TestThreadsContextFromParamsThrough(t *testing.T) {
	extractFieldFromContextFn := func(p graphql.ResolveParams) (interface{}, error) {
		return p.Context.Value(p.Args["key"]), nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: extractFieldFromContextFn,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ value(key:"a") }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.TODO(), "a", "xyz"),
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]interface{}{"value": "xyz"}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}
