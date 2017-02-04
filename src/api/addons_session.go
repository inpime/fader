package api

import (
	"interfaces"
	"sync"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	uuid "github.com/satori/go.uuid"
	lua "github.com/yuin/gopher-lua"
)

/*
Сессия относительно пользователя
Аноним тоже пользователь
*/

var luaSessionTypeName = "session"

func NewSession(
	L *lua.LState,
) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &luaSession{
		Props: make(map[string]interface{}),
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaSessionTypeName))
	return ud
}

func NewSessionFromUserID(
	L *lua.LState,
	userID string,
) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &luaSession{
		SessionID: "generate",
		UserID:    uuid.FromStringOrNil(userID),
		Props:     make(map[string]interface{}),
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaSessionTypeName))
	return ud
}

func NewSessionFromSID(
	L *lua.LState,
	sid string,
) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &luaSession{
		SessionID: sid,
		Props:     make(map[string]interface{}),
	}
	// TODO: load from SID
	return ud
}

func checkSession(L *lua.LState) *luaSession {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaSession); ok {
		return v
	}
	L.ArgError(1, "session expected")
	return nil
}

type luaSession struct {
	SessionID string

	UserID       uuid.UUID
	User         *interfaces.File
	userLoadOnce sync.Once
	userProvider interface{}

	Props map[string]interface{}
}

var sessionMethods = map[string]lua.LGFunction{
	"SessionID": func(L *lua.LState) int {
		s := checkSession(L)
		L.Push(lua.LString(s.SessionID))
		return 1
	},
	"IsEmpty": func(L *lua.LState) int {
		s := checkSession(L)
		isEmpty := s.SessionID == "" || uuid.Equal(uuid.Nil, s.UserID)
		L.Push(lua.LBool(isEmpty))
		return 1
	},
	"CreateForUser": func(L *lua.LState) int {
		s := checkSession(L)
		s.UserID = uuid.NewV4()
		s.SessionID = "session:" + s.UserID.String()
		return 0
	},
	"User": func(L *lua.LState) int { return 0 },
	"UserID": func(L *lua.LState) int {
		s := checkSession(L)
		L.Push(lua.LString(s.UserID.String()))
		return 1
	},
	"Set": func(L *lua.LState) int {
		s := checkSession(L)
		key := L.CheckString(2)
		lv := L.CheckAny(3)
		s.Props[key] = ToValueFromLValue(lv)
		return 0
	},
	"Get": func(L *lua.LState) int {
		s := checkSession(L)
		key := L.CheckString(2)
		L.Push(ToLValueOrNil(s.Props[key], L))
		return 1
	},
	"Del": func(L *lua.LState) int {
		s := checkSession(L)
		key := L.CheckString(2)
		delete(s.Props, key)
		return 0
	},
}

// serializer

func (s luaSession) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(
		s.SessionID,
		s.UserID,
		s.Props,
	)
}

func (s *luaSession) UnmarshalMsgpack(_b []byte) error {
	return msgpack.Unmarshal(_b,
		&s.SessionID,
		&s.UserID,
		&s.Props,
	)
}
