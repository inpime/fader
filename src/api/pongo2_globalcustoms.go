package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/inpime/dbox"
	"net/url"
	"store"
	"strings"
	"utils"
)

func pongo2InitGlobalCustoms() {

	var emptyURL, _ = url.Parse("")

	// SearchFiles search files in bucket
	// (the request is formed in buildSearchQueryFilesByBycket)
	tpls.Globals["SearchFiles"] = func(
		bucketName,
		queryStr,
		page,
		perpage *pongo2.Value,
	) *pongo2.Value {

		filter := store.NewSearchFilter(strings.ToLower(bucketName.String()))
		filter.SetQueryString(queryStr.String())
		filter.SetPage(page.Integer())
		filter.SetPerPage(perpage.Integer())

		queryRaw := buildSearchQueryFilesByBycket(
			strings.ToLower(bucketName.String()),
			queryStr.String(),
			page.Integer(),
			perpage.Integer(),
		)
		filter.SetQueryRaw(queryRaw)

		return pongo2.AsValue(makeSearch(filter))
	}

	// LoadByID load file by ID
	tpls.Globals["LoadByID"] = func(
		bucketName,
		fileID *pongo2.Value,
	) *pongo2.Value {

		file, err := store.LoadOrNewFileID(
			strings.ToLower(bucketName.String()),
			fileID.String())

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(file)
	}

	// Load load file by name
	tpls.Globals["Load"] = func(
		bucketName,
		fileName *pongo2.Value,
	) *pongo2.Value {
		file, _ := store.LoadOrNewFile(
			strings.ToLower(bucketName.String()),
			fileName.String())

		// if err != nil {
		// 	return pongo2.AsValue(err)
		// }

		return pongo2.AsValue(file)
	}

	tpls.Globals["URLQuery"] = func(args ...*pongo2.Value) *pongo2.Value {
		if len(args) == 0 {
			return pongo2.AsValue(emptyURL)
		}

		_url, ok := args[0].Interface().(*url.URL)

		if !ok {
			logrus.Warningf("not expected type %T, want '*url.URL'", args[0].Interface())

			return pongo2.AsValue(emptyURL)
		}

		if (len(args)-1)%2 != 0 {
			logrus.Warningf("args expected in multiples of two, want %d", len(args)-1)
			return pongo2.AsValue(emptyURL)
		}

		queryParams := args[1:]
		urlQueryValues := _url.Query()

		for i := 0; i < len(queryParams); i += 2 {
			urlQueryValues.Add(queryParams[i].String(), queryParams[i+1].String())
		}

		_url.RawQuery = urlQueryValues.Encode()

		return pongo2.AsValue(_url)
	}

	// builds the path part of the URL
	tpls.Globals["URL"] = func(args ...*pongo2.Value) *pongo2.Value {
		if len(args) == 0 {
			return pongo2.AsValue(emptyURL)
		}

		route := router.Get(args[0].String())

		if route == nil {
			return pongo2.AsValue(emptyURL)
		}

		if (len(args)-1)%2 != 0 {
			logrus.Warningf("args expected in multiples of two, want %d", len(args)-1)
			return pongo2.AsValue(emptyURL)
		}

		stringArgs := []string{}

		for _, arg := range args[1:] {
			stringArgs = append(stringArgs, arg.String())
		}

		_url, err := route.URLPath(stringArgs...)

		if err != nil {
			logrus.WithError(err).Warning("build url")

			return pongo2.AsValue(emptyURL)
		}

		return pongo2.AsValue(_url)
	}

	// Load load file by name
	tpls.Globals["M"] = func() *pongo2.Value {

		return pongo2.AsValue(utils.Map())
	}

	tpls.Globals["A"] = func() *pongo2.Value {
		return pongo2.AsValue(utils.NewA([]string{}))
	}

	// CreateBucket special function (used only to create a bucket)
	tpls.Globals["CreateBucket"] = func(_opt *pongo2.Value) *pongo2.Value {
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

	/*
		// TODO: remove, existing stringformat


		tpls.Globals["findFiles"] = func(bucket, query, page, perpage *pongo2.Value) *pongo2.Value {

			return pongo2.AsValue(Searcher.SearchFiles(bucket.String(), query.String(), page.Integer(), perpage.Integer()))
		}

		tpls.Globals["loadFileByName"] = func(bucket, name *pongo2.Value) *pongo2.Value {
			file, err := Mng.GetFullDataFileByNames(bucket.String(), name.String())
			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(Mng.FileWrap(file))
		}

		tpls.Globals["loadFileByID"] = func(id *pongo2.Value) *pongo2.Value {
			file, err := Mng.GetFullDataFileByID(id.String())
			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(Mng.FileWrap(file))
		}

		tpls.Globals["updateStructDataFile"] = func(in *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("updateFileStructData: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			file, err := Mng.PUpdateStructDataFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(Mng.FileWrap(file))
		}

		tpls.Globals["updateRawDataFile"] = func(in *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("updateFileStructData: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			file, err := Mng.PUpdateFullDataFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(Mng.FileWrap(file))
		}

		// tpls.Globals["saveFileRawData"] = func(id, raw *pongo2.Value) *pongo2.Value {
		// 	dto := NewUpdateFileDTO()
		// 	dto.FileID = uuid.FromStringOrNil(id.String())

		// 	switch raw.Interface().(type) {
		// 	case string:
		// 		dto.TextData = raw.String()
		// 	case []byte:
		// 		dto.BinData = raw.Interface().([]byte)
		// 	default:
		// 		logrus.Warningf("not expected type data, %T", raw.Interface())
		// 		return pongo2.AsValue(fmt.Errorf("not expected type data"))
		// 	}

		// 	dto.BinData = []byte(raw.String())

		// 	file, err := Mng.UpdateRawDataFile(dto)

		// 	if err != nil {
		// 		return pongo2.AsValue(err)
		// 	}

		// 	return pongo2.AsValue(Mng.FileWrap(file))
		// }

		tpls.Globals["search"] = func(bucket, query, page, perpage *pongo2.Value) *pongo2.Value {

			return pongo2.AsValue(Searcher.SearchFiles(bucket.String(), query.String(), page.Integer(), perpage.Integer()))
		}

		// TODO: rename muted

		tpls.Globals["clear"] = func(args ...*pongo2.Value) *pongo2.Value {

			return pongo2.AsValue(nil)
		}

		tpls.Globals["struct"] = func() *pongo2.Value {

			return pongo2.AsValue(store.NewMap())
		}

		// System data transfer object

		// tpls.Globals["NewFileDTO"] = func() *pongo2.Value {
		// 	return pongo2.AsValue(NewCreateFileDTO())
		// }

		// tpls.Globals["UpdateFileDTO"] = func() *pongo2.Value {
		// 	return pongo2.AsValue(NewUpdateFileDTO())
		// }

		// ------------------------------
		// ------------------------------
		// ------------------------------

		tpls.Globals["UpdateFile"] = func(in *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool
			var err error

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("UpdateFile: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			_file.FileInterface, err = Mng.PUpdateFullDataFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(_file)
		}

		tpls.Globals["UpdateRawDataFile"] = func(in *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool
			var err error

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("updateFileStructData: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			_file.FileInterface, err = Mng.PUpdateFullDataFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(_file)
		}

		tpls.Globals["UpdateStructDataFile"] = func(in *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool
			var err error

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("updateFileStructData: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			_file.FileInterface, err = Mng.PUpdateStructDataFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(_file)
		}

		// Создание файла на основе *FIle
		// flagViaUploader - указывает на то что файл созданиется на основе загруженного через uplaoder. Имя файла будет взято из мени загружаемого файла
		tpls.Globals["CreateFile"] = func(in, flagViaUploader *pongo2.Value) *pongo2.Value {
			var _file *File
			var ok bool
			var err error

			if _file, ok = in.Interface().(*File); !ok {
				return pongo2.AsValue(fmt.Errorf("updateFileStructData: not expected type data, want %T, got %T", _file, in.Interface()))
			}

			if flagViaUploader.IsBool() && flagViaUploader.Bool() {
				_file.SetName(_file.FileOriginalName())
			}

			_file.FileInterface, err = Mng.PCreateFile(_file)

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(_file)
		}

		tpls.Globals["RenameFile"] = func(fileId, newName *pongo2.Value) *pongo2.Value {
			err := Mng.RenameFileByID(fileId.String(), newName.String())

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(nil)
		}

		tpls.Globals["RemoveFile"] = func(fileId *pongo2.Value) *pongo2.Value {
			file, err := Mng.DropFileByID(fileId.String())

			if err != nil {
				return pongo2.AsValue(err)
			}

			return pongo2.AsValue(Mng.FileWrap(file))
		}
	*/
}
