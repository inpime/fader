package api

import (
	"interfaces"
	"net/url"

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

func (r RoutePongo2) URLPath(pairs ...string) *url.URL {
	v, _ := r.route.URLPath(pairs...)

	return v
}

func (r RoutePongo2) GetName() string {
	return r.route.GetName()
}

func (r RoutePongo2) Options() interfaces.RequestHandler {
	return r.route.Options()
}
