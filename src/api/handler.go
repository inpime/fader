package api

import (
	"api/router"
	"context"
	"interfaces"
	"net/http"

	"addons"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/yuin/gopher-lua"
)

var (
	NotFoundHandler = func(c echo.Context) error {

		return c.String(http.StatusNotFound, "Not Found")
	}

	ForbiddenHandler = func(c echo.Context) error {

		return c.String(http.StatusForbidden, "Forbidden")
	}

	MaintenanceHandler = func(c echo.Context) error {
		c.Response().Header().Set("Retry-After", "3600") // retry after 1 hourse
		return c.String(http.StatusServiceUnavailable, "Service Unavailable")
	}

	InternalErrorHandler = func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
)

func FaderHandler(ctx echo.Context) error {
	// Routing ------------------------------------------------------

	route := router.MatchVRouteFromContext(ctx)
	logger.Printf("\t[DEBUG] route: name %q", route.Route.GetName())

	// TODO: if special handler

	// Load the handled file ----------------------------------------
	file, err := fileLoaderForRouting.File(
		ctx.StdContext(),
		route.Handler.Bucket,
		route.Handler.File,
	)

	if err != nil {
		logger.Printf("[ERR] load file %q, %q, %s",
			route.Handler.Bucket,
			route.Handler.File,
			err,
		)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// PreRender ----------------------------------------------------

	var L = lua.NewState()
	defer L.Close()

	// Lua addons ---------------------------------------------------
	for _, addon := range addons.Addons {
		L.PreloadModule(addon.Name(), addon.LuaLoader)
	}

	_ctx := ContextLuaExecutor(L, ctx)
	_luaScript := string(file.LuaScript)

	// TODO: if empty lua script - skip lua executor

	err = L.DoString(_luaScript) // from loaded file

	// render from lua script
	if _ctx.Rendered || err != nil || _ctx.Err != nil {
		if err != nil || _ctx.Err != nil {
			logger.Println("[ERR] lua script", err)
			logger.Println("[ERR] context executor", _ctx.Err)
		}

		return err
	}

	// TODO: flush session
	// TODO: redirect if specified
	// TODO: custom response type

	// View -------------------------------------------------------------------

	var tpl *pongo2.Template

	pongo2.DefaultSet.Debug = true // TODO: set from config

	// if Debug true then recompile tpl on any request
	tplName := route.Handler.Bucket + "/" + route.Handler.File
	tpl, err = pongo2.FromCache(tplName)

	if err != nil {
		logger.Printf("[ERR] get tempalte file %q, %s\n", tplName, err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	res, err := tpl.Execute(pongo2.Context{
		"ctx": ctx, // TODO: контект для pongo2
	})

	if err != nil {
		logger.Printf("[ERR] execute template %q, %s\n", tplName, err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// response status specified in the lua script
	responseStatus := http.StatusOK
	if _ctx.ResponseStatus > 0 {
		responseStatus = _ctx.ResponseStatus
	}

	return ctx.HTML(responseStatus, res)
}

// internal components

var (
	_ interfaces.FileLoader = (*fileProvider)(nil)
)

func NewFileProvider(
	manager interfaces.FileManager,
	flags interfaces.DataUsed,
) *fileProvider {
	return &fileProvider{
		filemanager: manager,
		flags:       flags,
	}
}

type fileProvider struct {
	flags       interfaces.DataUsed
	filemanager interfaces.FileManager
}

func (p *fileProvider) File(
	ctx context.Context,
	bucketName, fileName string,
) (
	*interfaces.File, error,
) {
	ctx, cancel := context.WithTimeout(ctx, settings.TimeoutFileProvider)
	done := make(chan error, 1)
	defer func() {
		cancel()
		close(done)
	}()

	var file *interfaces.File

	go func() {
		var err error
		file, err = p.filemanager.FindFileByName(
			bucketName, fileName,
			p.flags,
		)
		done <- err
	}()

	select {
	case <-ctx.Done():
		<-done
		return nil, ctx.Err()
	case err := <-done:
		return file, err
	}
}
