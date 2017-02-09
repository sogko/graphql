package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

// to avoid possible compiler optimization
var benchmark_result *graphql.Result

func benchmarkQuery(b *testing.B, schema graphql.Schema, query string) {
	params := graphql.Params{
		Schema:        testutil.StarWarsSchema,
		RequestString: query,
	}

	var result *graphql.Result
	for n := 0; n < b.N; n++ {
		result = graphql.Do(params)
	}
	benchmark_result = result
}

func BenchmarkQuery_BasicQuery(b *testing.B) {
	query := `
		{
			hero {
				name
			}
		}
	`
	benchmarkQuery(b, testutil.StarWarsSchema, query)
}

func BenchmarkQuery_BasicQueryWithList(b *testing.B) {
	query := `
		{
			hero {
				id
				name
				friends {
					id
					name
				}
				appearsIn
			}
		}
	`
	benchmarkQuery(b, testutil.StarWarsSchema, query)
}

func BenchmarkQuery_DeeperQueryWithList(b *testing.B) {
	query := `
		{
			hero {
				id
				name
				friends {
					id
					name
					friends {
						id
						name
						appearsIn
					}
					appearsIn
				}
				appearsIn
			}
		}
	`
	benchmarkQuery(b, testutil.StarWarsSchema, query)
}

func BenchmarkQuery_QueryWithManyRootFieldsAndList(b *testing.B) {
	query := `
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
			leia: human(id: "1003") {
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
			artoo: droid(id: "2001") {
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
	`
	benchmarkQuery(b, testutil.StarWarsSchema, query)
}

func BenchmarkQuery_DeeperQueryWithManyRootFieldsAndList(b *testing.B) {
	query := `
		{
			hero {
				id
				name
				friends {
					name
					friends {
						id
						name
						appearsIn
					}
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
					friends {
						id
						name
						appearsIn
					}
				}
				appearsIn
				homePlanet
			}
			leia: human(id: "1003") {
				id
				name
				friends {
					id
					name
					appearsIn
					friends {
						id
						name
						appearsIn
					}
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
					friends {
						id
						name
						appearsIn
					}
				}
				appearsIn
				primaryFunction
			}
			artoo: droid(id: "2001") {
				id
				name
				friends {
					id
					name
					appearsIn
					friends {
						id
						name
						appearsIn
					}
				}
				appearsIn
				primaryFunction
			}
		}
	`
	benchmarkQuery(b, testutil.StarWarsSchema, query)
}
