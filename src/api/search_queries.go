package api

// import (
// 	"strings"
// 	"utils"
// )

// // buildSearchQueryFilesByBycket
// // поиск всех файлов в указанном бакете по текстовому запросу (по умолчанию по RawData)
// // Sorted by CreatedAt desc
// func buildSearchQueryFilesByBycket(bucketName, qstr string, page, perpage int) map[string]interface{} {
// 	query := utils.Map()

// 	// filter arguments

// 	qstr = strings.TrimSpace(qstr)
// 	bucketName = strings.TrimSpace(bucketName)

// 	from := perpage * page
// 	if page <= 0 {
// 		from = 0
// 	}

// 	size := perpage

// 	// prepare arguments of query

// 	queryFileter := utils.Map().
// 		Set("term", utils.Map().Set("Bucket", bucketName))

// 	querySort := utils.Map().
// 		Set("Meta.CreatedAt", utils.Map().Set("order", "desc"))

// 	_query := utils.Map().
// 		Set("query_string", utils.Map().
// 			Set("default_field", "TextData").
// 			Set("query", qstr))

// 	// build query

// 	// 1. sort
// 	// 2. from
// 	// 3. size
// 	// 4. filter
// 	// 5. fields
// 	// 6. query

// 	// 1. sort
// 	query.Set("sort", []interface{}{querySort})

// 	// 2. from
// 	query.Set("from", from)

// 	// 3. size
// 	query.Set("size", size)

// 	// 4. filter
// 	query.Set("filter", queryFileter)

// 	// 5. filter
// 	query.Set("fields", []string{})

// 	// 6. filter
// 	if len(qstr) > 0 {
// 		query.Set("query", _query)
// 	}

// 	return query
// }
