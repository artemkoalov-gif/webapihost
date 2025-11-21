package http_server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type DataPayload struct {
	Id     int      `json:"id"`
	Url    []string `json:"url"`
	Status int      `json:"status"`
}

type ReqLinks struct {
	LinksList []int `json:"links_list"`
}

func HttpMiddleware(entry *logrus.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status/v1/" && r.URL.Path != "/get/v1/" {
			entry.Log(logrus.WarnLevel, "Handling error request: "+r.URL.Path)
			httpError(w, "Handling error request: "+r.URL.Path)
			return
		}

		entry.Log(logrus.InfoLevel, "Handling request: "+r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}
