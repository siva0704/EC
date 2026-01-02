package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// SearchQueryBuilder constructs the Elasticsearch DSL
func BuildGrocerySearchQuery(term string, filters map[string]string) map[string]interface{} {
	// 1. Base Query: Multi-match with Fuzziness
	shouldClause := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":     term,
				"fields":    []string{"name^3", "description", "category^2", "brand"}, // Boost name and category
				"fuzziness": "AUTO", // Handle typos like "banannas"
			},
		},
	}

	// 2. Ranking Boosts (Function Score)
	// Boost In-Stock and High-Margin items
	functionScore := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": shouldClause,
			},
		},
		"functions": []map[string]interface{}{
			{
				"filter": map[string]interface{}{"term": map[string]interface{}{"is_in_stock": true}},
				"weight": 10, // Massive boost for in-stock items
			},
			{
				"field_value_factor": map[string]interface{}{
					"field":    "margin_score", // 0.0 to 1.0 pre-calculated score
					"factor":   1.5,
					"modifier": "sqrt",
				},
			},
		},
		"boost_mode": "multiply",
	}

	// 3. Facets / Aggregations
	aggs := map[string]interface{}{
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{"field": "category.keyword"},
		},
		"brands": map[string]interface{}{
			"terms": map[string]interface{}{"field": "brand.keyword"},
		},
		"price_ranges": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "price",
				"ranges": []map[string]interface{}{
					{"to": 5.00, "key": "Under $5"},
					{"from": 5.00, "to": 20.00, "key": "$5 - $20"},
					{"from": 20.00, "key": "Over $20"},
				},
			},
		},
	}

    // 4. Filters (Post-Query)
    filterClauses := []map[string]interface{}{}
    for k, v := range filters {
        filterClauses = append(filterClauses, map[string]interface{}{
            "term": map[string]interface{}{k + ".keyword": v},
        })
    }
    
    // Combine into final Bool query if filters exist
    finalQuery := map[string]interface{}{}
    if len(filterClauses) > 0 {
        finalQuery = map[string]interface{}{
            "bool": map[string]interface{}{
                "must": functionScore,
                "filter": filterClauses,
            },
        }
    } else {
        finalQuery = functionScore
    }

	return map[string]interface{}{
		"query": finalQuery,
		"aggs":  aggs,
        "size": 20, // Page size
        "from": 0,
	}
}
