package router

import (
	"interfaces"
	"net/url"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var testhandler = interfaces.RequestHandler{
	Name:           "handlername",
	Path:           "/fc/{id:[a-zA-Z0-9._-]+}",
	AllowedMethods: []string{echo.GET, echo.POST},

	Bucket: "a",
	File:   "b",

	SpecialHandler: "specialhandler",
	HandlerArgs:    "arg1, arg2",
}

func TestVRoute_buildUrl_simple(t *testing.T) {
	router := NewRouter()

	router.Handle(testhandler.Path, testhandler).
		Methods(testhandler.AllowedMethods...).
		Name(testhandler.Name)

	url, err := router.Get(testhandler.Name).URLPath("id", "123")

	assert.NoError(t, err)
	assert.Equal(t, "/fc/123", url.String())

	_, err = router.Get(testhandler.Name).URLPath("id", "123!")

	assert.NotEmpty(t, err, "Build invalid url (regexp not matched)")
}

var _ interfaces.RouteMatcher = (*Router)(nil)

func TestVRoute_matcher_simple(t *testing.T) {
	router := NewRouter()

	router.Handle(testhandler.Path+"/1", testhandler).
		Methods(testhandler.AllowedMethods...).
		Name(testhandler.Name)

	var match = &interfaces.RouteMatch{}

	_url, _ := url.Parse("/fc/1234/1")

	matched := router.Match(
		interfaces.RequestParams{
			URL:    _url,
			Method: echo.GET,
		},
		match,
	)

	assert.True(t, matched)
	assert.Equal(t, match.Route.GetName(), testhandler.Name)
	assert.Equal(t, match.Vars["id"], "1234")
	assert.Equal(t, len(match.Vars), 1)
	assert.Equal(t, match.Handler.Bucket, testhandler.Bucket)
	assert.Equal(t, match.Handler.File, testhandler.File)
}

// func TestVRouter_asyncupdate_simple(t *testing.T) {
// 	router := NewRouter()

// 	router.Handle(testhandler.Path, testhandler).
// 		Methods(testhandler.AllowedMethods...).
// 		Name(testhandler.Name)

// 	wg := sync.WaitGroup{}

// 	wg.Add(1)
// 	go func() {
// 		for i := 0; i < 10000; i++ {
// 			router.Handle(testhandler.Path+"/"+strconv.Itoa(i), testhandler).
// 				Methods(testhandler.AllowedMethods...).
// 				Name(testhandler.Name + strconv.Itoa(i))
// 		}

// 		wg.Done()
// 	}()

// 	time.Sleep(time.Microsecond * 10)

// 	for thread := 0; thread < 5; thread++ {
// 		wg.Add(1)

// 		go func() {
// 			for i := 0; i < 1000; i++ {
// 				_url, _ := url.Parse("/fc/1234/" + strconv.Itoa(i))
// 				match := &interfaces.RouteMatch{}

// 				matched := router.Match(
// 					interfaces.RequestParams{
// 						URL:    _url,
// 						Method: echo.GET,
// 					},
// 					match,
// 				)

// 				if matched == false {
// 					wg.Done()
// 					t.Fatal("expected matched true, got", matched, thread, i)
// 				}
// 			}

// 			wg.Done()
// 		}()
// 	}

// 	wg.Wait()
// }
