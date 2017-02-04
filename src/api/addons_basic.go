package api

import (
	"addons"
	"interfaces"
	"log"

	"github.com/flosch/pongo2"
	uuid "github.com/satori/go.uuid"
	"github.com/yuin/gopher-lua"
)

const (
	ADDONS_BASIC_VERSION     = "0.1"
	ADDONS_BASIC_NAME        = "basic"
	ADDONS_BASIC_AUTHOR      = "Fader"
	ADDONS_BASIC_DESCRIPTION = `Example of an addon for learning`
)

var (
	_ addons.Addon = (*AddonBasic)(nil)
)

func NewBasicAddon() *AddonBasic {
	return &AddonBasic{}
}

type AddonBasic struct {
}

func (a AddonBasic) Version() string {
	return ADDONS_BASIC_VERSION
}

func (a AddonBasic) Name() string {
	return ADDONS_BASIC_NAME
}

func (a AddonBasic) Author() string {
	return ADDONS_BASIC_AUTHOR
}

func (a AddonBasic) Description() string {
	return ADDONS_BASIC_DESCRIPTION
}

func (a *AddonBasic) LuaLoader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)

	// register other stuff
	L.SetField(mod, "version", lua.LString(ADDONS_BASIC_VERSION))
	L.SetField(mod, "name", lua.LString(ADDONS_BASIC_NAME))
	L.SetField(mod, "author", lua.LString(ADDONS_BASIC_AUTHOR))
	L.SetField(mod, "description", lua.LString(ADDONS_BASIC_DESCRIPTION))

	// returns the module
	L.Push(mod)

	////////////////////////////////////
	// Custom Types
	////////////////////////////////////

	// FormFile
	mt := L.NewTypeMetatable(luaFormFileTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), formFileMethods))

	// Route
	mt = L.NewTypeMetatable(luaRouteTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), routeMethods))

	// File
	mt = L.NewTypeMetatable(luaFileTypeName)
	L.SetGlobal("File", mt)
	L.SetField(
		mt,
		"new",
		L.NewFunction(func(L *lua.LState) int {
			L.Push(FileAsLuaFile(L, nil))
			return 1
		}),
	)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), fileMethods))

	// Session
	mt = L.NewTypeMetatable(luaSessionTypeName)
	L.SetGlobal("Session", mt)
	L.SetField(
		mt,
		"empty",
		L.NewFunction(func(L *lua.LState) int {
			L.Push(NewSession(L))
			return 1
		}),
	)
	L.SetField(
		mt,
		"bySessionID",
		L.NewFunction(func(L *lua.LState) int {
			sid := L.CheckString(1)
			L.Push(NewSessionFromSID(L, sid))
			return 1
		}),
	)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), sessionMethods))

	// Bucket
	mt = L.NewTypeMetatable(luaBucketTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), bucketMethods))

	// DataUsed
	// type
	mt = L.NewTypeMetatable(luaDataUsed)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), dataUsedMethods))
	// consts
	newDataUsed := func(
		L *lua.LState,
		v interfaces.DataUsed,
	) *lua.LUserData {
		ud := L.NewUserData()
		ud.Value = v
		L.SetMetatable(ud, L.GetTypeMetatable(luaDataUsed))
		return ud
	}
	L.SetField(
		mod,
		"PrimaryIDsData",
		newDataUsed(L, interfaces.PrimaryIDsData),
	)
	L.SetField(
		mod,
		"PrimaryNamesData",
		newDataUsed(L, interfaces.PrimaryNamesData),
	)
	L.SetField(
		mod,
		"ContentTypeData",
		newDataUsed(L, interfaces.ContentTypeData),
	)
	L.SetField(
		mod,
		"OwnersData",
		newDataUsed(L, interfaces.OwnersData),
	)
	L.SetField(
		mod,
		"AccessStatusData",
		newDataUsed(L, interfaces.AccessStatusData),
	)
	L.SetField(
		mod,
		"LuaScript",
		newDataUsed(L, interfaces.LuaScript),
	)
	L.SetField(
		mod,
		"MetaData",
		newDataUsed(L, interfaces.MetaData),
	)
	L.SetField(
		mod,
		"StructuralData",
		newDataUsed(L, interfaces.StructuralData),
	)
	L.SetField(
		mod,
		"RawData",
		newDataUsed(L, interfaces.RawData),
	)
	L.SetField(
		mod,
		"BucketStoreNames",
		newDataUsed(L, interfaces.BucketStoreNames),
	)
	L.SetField(
		mod,
		"FullFile",
		newDataUsed(L, interfaces.FullFile),
	)

	L.SetField(mod, "check", L.NewFunction(func(L *lua.LState) int {
		v := checkDataUsed(L)
		log.Println("PrimaryIDsData", v&interfaces.PrimaryIDsData != 0)
		log.Println("PrimaryNamesData", v&interfaces.PrimaryNamesData != 0)
		log.Println("ContentTypeData", v&interfaces.ContentTypeData != 0)
		log.Println("OwnersData", v&interfaces.OwnersData != 0)
		log.Println("AccessStatusData", v&interfaces.AccessStatusData != 0)
		log.Println("LuaScript", v&interfaces.LuaScript != 0)
		log.Println("MetaData", v&interfaces.MetaData != 0)
		log.Println("StructuralData", v&interfaces.StructuralData != 0)
		log.Println("RawData", v&interfaces.RawData != 0)
		log.Println("BucketStoreNames", v&interfaces.BucketStoreNames != 0)

		return 0
	}))

	return 1
}

func (a *AddonBasic) ExtContextPongo2(_ctx pongo2.Context) error {
	ctx := make(pongo2.Context)
	// ctx["ContextFunction"] = func() *pongo2.Value {
	// 	return pongo2.AsValue("context function")
	// }

	////////////////////////////////////////////////////////////////////////////
	// file & bucket manager
	////////////////////////////////////////////////////////////////////////////

	ctx["TPrimaryIDsData"] = interfaces.PrimaryIDsData
	ctx["TPrimaryNamesData"] = interfaces.PrimaryNamesData
	ctx["TContentTypeData"] = interfaces.ContentTypeData
	ctx["TOwnersData"] = interfaces.OwnersData
	ctx["TAccessStatusData"] = interfaces.AccessStatusData
	ctx["TLuaScript"] = interfaces.LuaScript
	ctx["TMetaData"] = interfaces.MetaData
	ctx["TStructuralData"] = interfaces.StructuralData
	ctx["TRawData"] = interfaces.RawData
	ctx["TBucketStoreNames"] = interfaces.BucketStoreNames

	ctx["TFullFile"] = interfaces.FullFile
	ctx["TFullFileWithoutRawData"] = interfaces.FileWithoutRawData
	ctx["TFullBucket"] = interfaces.FullBucket

	ctx["FindFileBuName"] = func(
		bname,
		fname,
		dtype *pongo2.Value,
	) *pongo2.Value {
		return pongo2.AsValue(nil)
	}
	ctx["FindFile"] = func(fid, dtype *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(nil)
	}
	ctx["FindBucketByName"] = func(
		bname,
		dtype *pongo2.Value,
	) *pongo2.Value {
		return pongo2.AsValue(nil)
	}
	ctx["FindBucket"] = func(bid, dtype *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(nil)
	}
	// манипуляторы в lua
	// collections
	ctx["ListBuckets"] = func() *pongo2.Value {
		return pongo2.AsValue(listOfBuckets())
	}
	ctx["ListFilesByBucketID"] = func(bucketID *pongo2.Value) *pongo2.Value {
		var bid uuid.UUID
		switch v := bucketID.Interface().(type) {
		case uuid.UUID:
			bid = v
		case string:
			bid = uuid.FromStringOrNil(v)
		}
		return pongo2.AsValue(filesByBucketID(bid))
	}
	// utils
	ctx["FileIsImage"] = func(f *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(false)
	}
	ctx["FileIsText"] = func(f *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(false)
	}
	ctx["FileIsRaw"] = func(f *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(false)
	}

	// Router

	ctx["Route"] = func(name *pongo2.Value) *pongo2.Value {
		foundRoute := vrouter.Get(name.String())
		if foundRoute == nil {
			return pongo2.AsValue(nil)
		}

		return pongo2.AsValue(
			&RoutePongo2{
				vrouter.Get(name.String()),
			},
		)
	}

	_ctx.Update(ctx)
	return nil
}

func (a *AddonBasic) ExtTagsFiltersPongo2(
	addf func(name string, fn pongo2.FilterFunction),
	repf func(name string, fn pongo2.FilterFunction),
	addt func(name string, fn pongo2.TagParser),
	rapt func(name string, fn pongo2.TagParser),
) error {

	pongo2.RegisterFilter(
		"btos",
		func(
			in *pongo2.Value,
			param *pongo2.Value,
		) (
			*pongo2.Value,
			*pongo2.Error,
		) {
			v, ok := in.Interface().([]byte)
			if !ok {
				return pongo2.AsValue(""), nil
			}
			return pongo2.AsValue(string(v)), nil
		},
	)

	return nil
}
