package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

// TODO: Тест роутинга
// TODO: Тест менеджера файлов: создать, найти, обновить, удалить
// TODO: Создать context для тестирования что бы передавать результат операции

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
