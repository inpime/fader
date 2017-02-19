package api

import (
	"interfaces"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

// TODO: Тест роутинга
// TODO: Тест менеджера файлов: создать, найти, обновить, удалить
// TODO: Создать context для тестирования что бы передавать результат операции

func TestAddonsBasic_ContextCurrentFile(t *testing.T) {
	var L = lua.NewState()
	defer L.Close()
	SetupAddons()
	setupLuaModules(L)
	ctx := setupLuaContext("POST", "/", nil, L)
	// Регистрация типа файл
	mt := L.NewTypeMetatable(luaFileTypeName)
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

	ctx.CurrentFile = &interfaces.File{
		FileID:   uuid.NewV4(),
		BucketID: uuid.NewV4(),

		FileName: "filename",

		ContentType: "text/plain",
		Owners: []uuid.UUID{
			uuid.NewV4(),
		},
		IsPrivate:  true,
		IsReadOnly: false,

		MetaData: map[string]interface{}{},
		StructuralData: map[string]interface{}{
			"a":       "b",
			"int":     1,
			"float":   0.5,
			"boolean": true,
			"null":    nil,
			"string":  "foo bar",
			"array":   []string{"foo", "bar"},
			"object": map[string]interface{}{
				"foo": 1,
				"bar": 0.5,
			},
		},
		RawData: []byte("data data"),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := L.DoString(`
	cf = ctx():CurrentFile()
	d = cf:StructuralData()
	seq = 0
	if d.a == "b" then
		seq = seq + 1
	end 

	if d.int == 1 then
		seq = seq + 1
	end 

	if d.float == 0.5 then
		seq = seq + 1
	end 

	if d.boolean then
		seq = seq + 1 
	end 

	if d.null == nil then
		seq = seq + 1 
	end
	
	if d.string == "foo bar" then
		seq = seq + 1 
	end

	if d.array[1] == "foo" then
		seq = seq + 1 
	end

	if d.object.foo == 1 then
		seq = seq + 1 
	end

	d.newint = 2
	d.newfloat = 2.2
	d.newbool = true
	d.newnull = nil
	d.newstring = "new foo bar"
	d.newarray = {"a", 2}
	d.object = {
		a = "qw",
		i = 3.1,
	}

	cf:SetStructuralData(d)
	
	ctx():Set("res", seq)
`)

	assert.NoError(t, err, "execute lua")
	assert.EqualValues(t, 8, ctx.EchoCtx().Get("res"), "checksum")

	assert.EqualValues(t, 2, ctx.CurrentFile.StructuralData["newint"])
	assert.EqualValues(t, 2.2, ctx.CurrentFile.StructuralData["newfloat"])
	assert.EqualValues(t, true, ctx.CurrentFile.StructuralData["newbool"])
	assert.EqualValues(t, nil, ctx.CurrentFile.StructuralData["newnull"])
	assert.EqualValues(t, "new foo bar", ctx.CurrentFile.StructuralData["newstring"])
	assert.EqualValues(t, []interface{}{"a", 2.0}, ctx.CurrentFile.StructuralData["newarray"])
	assert.EqualValues(t, map[string]interface{}{"a": "qw", "i": 3.1}, ctx.CurrentFile.StructuralData["object"])
}

func TestAddonsBasic_UsedData_Lua(t *testing.T) {
	var L = lua.NewState()
	defer L.Close()
	L.PreloadModule(ADDONS_BASIC_NAME, NewBasicAddon().LuaLoader)
	ctx := L.NewUserData()
	ctx.Value = 1
	L.SetGlobal("ctx", ctx)
	err := L.DoString(`
local std = require("basic")
vvv = std.PrimaryIDsData
vvv:Add(std.PrimaryNamesData)
std.check(vvv)
`) // from loaded file
	assert.NoError(t, err)
	t.Logf("%v", ctx.Value)
}

// TestContext_basicTypes
// Простые типы int, float, string, bool, nil
// TODO: map, slice
func TestContext_basicTypes(t *testing.T) {
	var L = lua.NewState()
	defer L.Close()
	SetupAddons()
	setupLuaModules(L)
	ctx := setupLuaContext("POST", "/", nil, L)
	ctx.EchoCtx().Set("int", 1)
	ctx.EchoCtx().Set("float", 3.14)
	ctx.EchoCtx().Set("bool", true)
	ctx.EchoCtx().Set("string", "hello")
	ctx.EchoCtx().Set("nil", nil)
	ctx.EchoCtx().Set("arraystr", []string{"a", "b"})
	ctx.EchoCtx().Set("arrayfloat", []float64{1.1, 2.2})
	ctx.EchoCtx().Set("arrayint64", []float64{1, 2})

	err := L.DoString(`
if ctx():Get("int") == 1 then
	ctx():Set("intok", true)
end

if ctx():Get("float") == 3.14 then
	ctx():Set("floatok", true)
end

if ctx():Get("bool") == true then
	ctx():Set("boolok", true)
end

if ctx():Get("string") == "hello" then
	ctx():Set("stringok", true)
end

if ctx():Get("nil") == nil then
	ctx():Set("nilok", true)
end

-- check arr str

if #ctx():Get("arraystr") == 2 then 
	ctx():Set("arraystrlenok", true)

	arr = ctx():Get("arraystr")
	if arr[1] == "a" then
		if arr[2] == "b" then
			ctx():Set("arraystrcontainsok", true)
		end 
	end
end 

-- check arr float

if #ctx():Get("arrayfloat") == 2 then 
	ctx():Set("arrayfloatlenok", true)

	arr = ctx():Get("arrayfloat")
	if arr[1] == 1.1 then
		if arr[2] == 2.2 then
			ctx():Set("arrayfloatcontainsok", true)
		end 
	end
end 

ctx():Set("int", 2)
ctx():Set("float", 3.15)
ctx():Set("bool", true)
ctx():Set("string", "hello world")
ctx():Set("arraystr", {"aa", "bb", "cc"})
ass = {c = "a", b = "b", a = "c"}
ass[1] = 2
ass[2] = 3.3
ctx():Set("associative", ass)
`)
	assert.NoError(t, err, "execute lua")
	assert.EqualValues(t, true, ctx.EchoCtx().Get("intok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("floatok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("boolok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("stringok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("nilok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("arraystrlenok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("arraystrcontainsok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("arrayfloatlenok"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("arrayfloatcontainsok"))

	assert.EqualValues(t, 2, ctx.EchoCtx().Get("int"))
	assert.EqualValues(t, 3.15, ctx.EchoCtx().Get("float"))
	assert.EqualValues(t, true, ctx.EchoCtx().Get("bool"))
	assert.EqualValues(t, "hello world", ctx.EchoCtx().Get("string"))

	assert.Len(t, ctx.EchoCtx().Get("arraystr"), 3)
	assert.EqualValues(t, "aa", ctx.EchoCtx().Get("arraystr").([]interface{})[0])
	assert.EqualValues(t, "bb", ctx.EchoCtx().Get("arraystr").([]interface{})[1])
	assert.EqualValues(t, "cc", ctx.EchoCtx().Get("arraystr").([]interface{})[2])

	assert.Len(t, ctx.EchoCtx().Get("associative"), 5)
	assert.EqualValues(
		t,
		"c",
		ctx.EchoCtx().Get("associative").(map[string]interface{})["a"],
	)
	assert.EqualValues(
		t,
		"b",
		ctx.EchoCtx().Get("associative").(map[string]interface{})["b"],
	)
	assert.EqualValues(
		t,
		"a",
		ctx.EchoCtx().Get("associative").(map[string]interface{})["c"],
	)
	assert.EqualValues(
		t,
		2,
		ctx.EchoCtx().Get("associative").(map[string]interface{})["1"],
	)
	assert.EqualValues(
		t,
		3.3,
		ctx.EchoCtx().Get("associative").(map[string]interface{})["2"],
	)
}
