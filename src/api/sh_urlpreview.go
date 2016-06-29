package api

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"net/url"

	"fmt"
	"strings"

	"store"

	"bytes"
	"github.com/dyatlov/go-oembed/oembed"
	"github.com/dyatlov/go-url2oembed/url2oembed"
	"github.com/jeffail/tunny"
	"time"
)

var workerPool *tunny.WorkPool
var workers = make([]tunny.TunnyWorker, 100)

type workerData struct {
	Status int
	Data   string
}

type apiWorker struct {
	Parser *url2oembed.Parser
}

// Use this call to block further jobs if necessary
func (worker *apiWorker) TunnyReady() bool {
	return true
}

// This is where the work actually happens
func (worker *apiWorker) TunnyJob(data interface{}) interface{} {
	if u, ok := data.(string); ok {
		u = strings.Trim(u, "\r\n")

		logrus.Infof("Got url: %s", u)

		info := worker.Parser.Parse(u)

		if info == nil {
			logrus.Infof("No info for url: %s", u)

			return &workerData{Status: 404, Data: "{\"status\": \"error\", \"message\":\"Unable to retrieve information from provided url\"}"}
		}
		if info.Status < 300 {
			logrus.Infof("Url parsed: %s", u)

			return &workerData{Status: 200, Data: info.String()}
		}

		logrus.Infof("Something weird: %s", u)

		return &workerData{Status: 411, Data: fmt.Sprintf("{\"status\": \"error\", \"message\":\"Unable to obtain data. Status code: %d\"}", info.Status)}
	}

	return &workerData{Status: 500, Data: "{\"status\": \"error\", \"message\":\"Something weird happened\"}"}
}

func UrlPreviewWorkersInit() {
	file, err := store.LoadOrNewFile(SettingsBucketName, "urlpreview.providers.json")

	if err != nil {
		return
	}

	oe := oembed.NewOembed()
	oe.ParseProviders(bytes.NewReader(file.RawData().Bytes()))

	for i := range workers {
		p := url2oembed.NewParser(oe)
		p.MaxHTMLBodySize = 50000
		p.MaxBinaryBodySize = 4096
		p.WaitTimeout = time.Duration(7) * time.Second
		// p.BlacklistedIPNetworks = blackListNetworks
		// p.WhitelistedIPNetworks = whiteListNetworks
		workers[i] = &(apiWorker{Parser: p})
	}

	pool, err := tunny.CreateCustomPool(workers).Open()

	if err != nil {
		logrus.Fatal(err)
	}

	// defer pool.Close()

	workerPool = pool
}

func UrlPreview_SpecialHandler(ctx *ContextWrap) error {

	u := ctx.QueryParam("url")

	_, err := url.Parse(u)

	if err != nil {

		logrus.Warningf("Invalid URL provided: %s", u)

		return ctx.NoContent(http.StatusInternalServerError)
	}

	result, err := workerPool.SendWork(u)

	if data, ok := result.(*workerData); ok {

		return ctx.JSONBlob(data.Status, []byte(data.Data))
	}

	return nil
}
