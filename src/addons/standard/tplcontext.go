package standard

import (
	"api/config"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/inpime/dbox"
	"net/url"
	"store"
	"strings"
	"utils"
)

func (Extension) initTplContext() {
	pongo2.DefaultSet.Globals["NewUUID"] = func() *pongo2.Value {
		return pongo2.AsValue(dbox.NewUUID())
	}

	pongo2.DefaultSet.Globals["SectionAppConfig"] = func(sectionName *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(config.AppSettings().M(sectionName.String()))
	}

	// DeleteFile
	pongo2.DefaultSet.Globals["DeleteFile"] = func(bucketId, fileId *pongo2.Value) *pongo2.Value {
		if !bucketId.IsString() || !fileId.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		return pongo2.AsValue(store.DeleteFileID(bucketId.String(), fileId.String()))
	}

	pongo2.DefaultSet.Globals["NewFile"] = func(bucketName *pongo2.Value) *pongo2.Value {
		if !bucketName.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, _ := store.LoadOrNewFileID(bucketName.String(), "")

		return pongo2.AsValue(file)
	}

	// LoadByID load file by ID
	pongo2.DefaultSet.Globals["LoadByID"] = func(
		bucketName,
		fileID *pongo2.Value,
	) *pongo2.Value {

		if !bucketName.IsString() || !fileID.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, err := store.LoadOrNewFileID(
			strings.ToLower(bucketName.String()),
			fileID.String())

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(file)
	}

	// Load load file by name
	pongo2.DefaultSet.Globals["Load"] = func(
		bucketName,
		fileName *pongo2.Value,
	) *pongo2.Value {

		if !bucketName.IsString() || !fileName.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, _ := store.LoadOrNewFile(
			strings.ToLower(bucketName.String()),
			fileName.String())

		return pongo2.AsValue(file)
	}

	//

	pongo2.DefaultSet.Globals["URLQuery"] = func(args ...*pongo2.Value) *pongo2.Value {
		emptyUrl, _ := url.Parse("")

		if len(args) == 0 {
			return pongo2.AsValue(emptyUrl)
		}

		_url, ok := args[0].Interface().(*url.URL)

		if !ok {
			logrus.WithFields(logrus.Fields{
				"_service": addonName,
			}).Warningf("not expected type %T, want '*url.URL'", args[0].Interface())

			return pongo2.AsValue(emptyUrl)
		}

		if (len(args)-1)%2 != 0 {
			logrus.WithFields(logrus.Fields{
				"_service": addonName,
			}).Warningf("args expected in multiples of two, want %d", len(args)-1)
			return pongo2.AsValue(emptyUrl)
		}

		queryParams := args[1:]
		urlQueryValues := _url.Query()

		for i := 0; i < len(queryParams); i += 2 {
			urlQueryValues.Add(queryParams[i].String(), queryParams[i+1].String())
		}

		_url.RawQuery = urlQueryValues.Encode()

		return pongo2.AsValue(_url)
	}

	// // builds the path part of the URL
	// pongo2.DefaultSet.Globals["URL"] = func(args ...*pongo2.Value) *pongo2.Value {
	// 	emptyUrl, _ := url.Parse("")

	// 	if len(args) == 0 {
	// 		return pongo2.AsValue(emptyUrl)
	// 	}

	// 	routeName := args[0].String()
	// 	route := config.Router.Get(routeName)

	// 	if route == nil {
	// 		return pongo2.AsValue(emptyUrl)
	// 	}

	// 	if (len(args)-1)%2 != 0 {
	// 		logrus.WithFields(logrus.Fields{
	// 			"_service": addonName,
	// 		}).Warningf("args expected in multiples of two, want %d", len(args)-1)
	// 		return pongo2.AsValue(emptyUrl)
	// 	}

	// 	stringArgs := []string{}

	// 	for _, arg := range args[1:] {
	// 		stringArgs = append(stringArgs, arg.String())
	// 	}

	// 	_url, err := route.URLPath(stringArgs...)

	// 	if err != nil {
	// 		logrus.WithError(err).WithFields(logrus.Fields{
	// 			"_service": addonName,
	// 			"args":     stringArgs,
	// 		}).Warning("build url")

	// 		return pongo2.AsValue(emptyUrl)
	// 	}

	// 	return pongo2.AsValue(_url)
	// }

	// Load load file by name
	pongo2.DefaultSet.Globals["M"] = func() *pongo2.Value {

		return pongo2.AsValue(utils.Map())
	}

	pongo2.DefaultSet.Globals["A"] = func() *pongo2.Value {
		return pongo2.AsValue(utils.NewA([]string{}))
	}

	pongo2.DefaultSet.Globals["AIface"] = func() *pongo2.Value {
		return pongo2.AsValue([]interface{}{})
	}

	pongo2.DefaultSet.Globals["Validator"] = func() *pongo2.Value {

		return pongo2.AsValue(NewValidatorData())
	}

	// CreateBucket special function (used only to create a bucket)
	pongo2.DefaultSet.Globals["CreateBucket"] = func(_opt *pongo2.Value) *pongo2.Value {
		opt := utils.Map().LoadFrom(_opt.Interface())

		bucketName := opt.String("Name")
		bucket, err := store.BucketByName(bucketName)

		if opt.Bool("SameAsMetaStoreType") {
			bucket.InitInOneStore(dbox.StoreType(opt.String("MetaDataStoreType")))
		} else {

			bucket.InitMetaDataStore(
				dbox.StoreType(opt.String("MetaDataStoreType")),
				opt.Bool("MetaHaveSuffix")) // store key - <type>.<name>.meta
			bucket.InitMapDataStore(
				dbox.StoreType(opt.String("MapDataStoreType")),
				opt.Bool("MapDataHaveSuffix")) // store key - <type>.<name>.mapdata
			bucket.InitRawDataStore(
				dbox.StoreType(opt.String("RawDataStoreType")),
				opt.Bool("RawDataHaveSuffix")) // store key - <type>.<name>.rawdata
		}

		// Only create new bucket
		if err != dbox.ErrNotFound {
			return pongo2.AsValue(err)
		}

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().LoadFrom(opt.String("MappingMapDataFiles"))

		if err := bucket.UpdateMapping(); err != nil {
			logrus.WithError(err).Errorf("create new bucket %q: update mapping", bucketName)
			return pongo2.AsValue(err)
		}

		if err := bucket.Sync(); err != nil {
			logrus.WithError(err).Errorf("create new bucket %q: save bucket", bucketName)
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(bucket)
	}
}
