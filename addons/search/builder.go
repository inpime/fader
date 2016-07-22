package search

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/inpime/fader/utils/sdata"
)

func BuildAdvancedQuery(queryString, queryFieldsString string, page, perpage int, filter, sortFields string) map[string]interface{} {
	var query = sdata.NewStringMap()
	query.Set("fields", []string{})

	from := perpage * page
	if page <= 0 {
		from = 0
	}

	size := perpage

	// 1 size ,from
	query.Set("size", size)
	query.Set("from", from)

	// 2. query
	if len(queryString) > 0 {
		var queryFields = strings.Fields(queryFieldsString)
		var queryTerms = []interface{}{}

		queryTerms = append(queryTerms, map[string]interface{}{"multi_match": map[string]interface{}{
			"query":            queryString,
			"fields":           queryFields,
			"operator":         "and",
			"zero_terms_query": "all",
		}})

		query.Set("query", map[string]interface{}{
			"bool": map[string]interface{}{
				"must": queryTerms,
			},
		})
	}

	// 3. filters
	var filterQueryTerms = []interface{}{}
	var fopts = strings.Fields(filter)

	// examples:
	//  prefix category category_name
	//  range count gt,qwd
	if len(fopts)%3 == 0 {
		for i := 0; i < len(fopts); i += 3 {
			var opName,
				fieldName,
				fieldValue = fopts[i],
				fopts[i+1],
				fopts[i+2]

			switch opName {
			case "range", "numeric_range":
				/*
				   "range": {
				       "FIELD": {
				           "gt": 10,
				           "gte": 10,
				           "lte": 20,
				           "lt": 20,
				       }
				   }
				*/
				_fieldValue := strings.Split(fieldValue, ",")

				if len(_fieldValue) == 2 {
					_operation, _value := _fieldValue[0], _fieldValue[1]
					filterQueryTerms = append(filterQueryTerms, map[string]interface{}{opName: map[string]interface{}{fieldName: map[string]interface{}{_operation: _value}}})
				} else {
					logrus.WithField("_api", NAME).Warningf("invalid range options %q", fieldValue)
				}
			case "terms":
				/*
				   "terms": {
				       "FIELD": [
				           "VALUE1",
				           "VALUE2"
				       ]
				   }
				*/
				_fieldValues := strings.Split(fieldValue, ",")
				if len(_fieldValues) > 0 {
					filterQueryTerms = append(filterQueryTerms, map[string]interface{}{opName: map[string]interface{}{fieldName: _fieldValues}})
				} else {
					logrus.WithField("_api", NAME).Warningf("invalid terms options %q", fieldValue)
				}
			default:
				filterQueryTerms = append(filterQueryTerms, map[string]interface{}{opName: map[string]interface{}{fieldName: fieldValue}})
			}
		}
	}

	if len(filterQueryTerms) > 0 {
		query.Set("filter", map[string]interface{}{
			"bool": map[string]interface{}{
				"must": filterQueryTerms,
			},
		})
	}

	// 4. sort

	var sortopts = strings.Fields(sortFields)

	if len(sortopts) > 0 {
		var sortOptions = []interface{}{}
		for i := 0; i < len(sortopts); i++ {
			_fieldOptions := strings.Split(sortopts[i], ",")
			if len(_fieldOptions) == 2 {
				fieldName, mode := _fieldOptions[0], _fieldOptions[1]
				sortOptions = append(sortOptions, map[string]interface{}{fieldName: mode})
			} else {
				logrus.WithField("_api", NAME).Warningf("invalid sort options %q", sortopts[i])
			}

		}
		query.Set("sort", sortOptions)
	}

	// 5. aggr

	// request.Set("aggs", map[string]interface{}{
	// 	"genders": map[string]interface{}{
	// 		"filters": map[string]interface{}{
	// 			"filters": map[string]interface{}{
	// 				"filter_name": map[string]interface{}{
	// 					"and": map[string]interface{}{
	// 						"filters": terms,
	// 					},
	// 				},
	// 			},
	// 		},
	// 		"aggs": map[string]interface{}{
	// 			"categories": map[string]interface{}{
	// 				"terms": map[string]interface{}{
	// 					"field": "category",
	// 					"size":  100, // magic number
	// 				},
	// 			},
	// 		},
	// 	},
	// })

	return query.ToMap()
}
