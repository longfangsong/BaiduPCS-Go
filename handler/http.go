package handler

import (
	"github.com/iikira/BaiduPCS-Go/internal/pcscommand"
	"github.com/iikira/BaiduPCS-Go/requester"
	"github.com/iikira/BaiduPCS-Go/requester/multipartreader"
	"net/http"
)

var ch = make(chan string, 1)

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return
	}
	pcs := pcscommand.GetBaiduPCS()
	info, err := pcs.LocateDownload("/lana/" + name)
	if err != nil {
		http.Error(w, "failed to get url", http.StatusInternalServerError)
	}
	url := info.URLs[0].URL
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

type HttpRequestReaderLen64 struct {
	Request *http.Request
}

func (reader *HttpRequestReaderLen64) Len() int64 {
	return reader.Request.ContentLength
}

func (reader *HttpRequestReaderLen64) Read(p []byte) (n int, err error) {
	return reader.Request.Body.Read(p)
}

func PostFileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return
	}
	pcs := pcscommand.GetBaiduPCS()
	err := pcs.Upload("/lana/"+name, func(uploadURL string, jar http.CookieJar) (resp *http.Response, err error) {
		mr := multipartreader.NewMultipartReader()
		mr.AddFormFile("file", "file", &HttpRequestReaderLen64{Request: r})
		mr.CloseMultipart()

		c := requester.NewHTTPClient()
		c.SetCookiejar(jar)
		return c.Req(http.MethodPost, uploadURL, mr, nil)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		PostFileHandler(w, r)
	case "GET":
		GetFileHandler(w, r)
	}
}
