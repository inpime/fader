package api

import (
	"interfaces"
	"net/url"

	"log"

	uuid "github.com/satori/go.uuid"
)

// filesByBucketID список файлов бакета
func filesByBucketID(id uuid.UUID) (res []*interfaces.File) {
	fileManager.(interfaces.FileImportManager).
		EachFile(func(item *interfaces.File) error {
			if uuid.Equal(id, item.BucketID) {
				res = append(res, item)
			}
			return nil
		})
	return
}

// listOfBuckets список бакетов
func listOfBuckets() (res []*interfaces.Bucket) {
	bucketManager.(interfaces.BucketImportManager).
		EachBucket(func(item *interfaces.Bucket) error {
			res = append(res, item)
			return nil
		})
	return
}

//////////////////////////////////////////////////////////
// Route pongo2
//////////////////////////////////////////////////////////

type RoutePongo2 struct {
	route interfaces.Route
}

func (r RoutePongo2) URLPath(pairs ...interface{}) *url.URL {
	var args = make([]string, len(pairs))

	for i, v := range pairs {
		switch v := v.(type) {
		case string:
			args[i] = v
		case uuid.UUID:
			args[i] = v.String()
		default:
			log.Printf("Pongo2 Route.URLPath: not expected type %T", v)
			return nil
		}
	}

	v, _ := r.route.URLPath(args...)

	return v
}

func (r RoutePongo2) Has(name string) bool {
	return r.GetName() == name
}

func (r RoutePongo2) GetName() string {
	return r.route.GetName()
}

func (r RoutePongo2) Options() interfaces.RequestHandler {
	return r.route.Options()
}
