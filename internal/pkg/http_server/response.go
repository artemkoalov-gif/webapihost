package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var StatusNotOK = "unavailable"
var StatusOK = "available"

type Url struct {
	Url        string `json:"url"`
	Aviability int    `json:"avialability"`
}

type Errors struct {
	Error []string `json:"error"`
}
type Response struct {
	Result string `json:"result"`
	Errors bool   `json:"has_errors"`
	Links  []Url  `json:"links"`
}
type ResponseLinks struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_summ"`
}

type Links struct {
	mu   sync.RWMutex
	Data map[int]ResponseLinks `json:"data"`
}

func httpOk(w http.ResponseWriter) http.ResponseWriter {
	resp := Response{Result: "OK", Errors: false}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return w
}

func httpError(w http.ResponseWriter, err string) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
	return w
}

func (r *ResponseLinks) responseLinks(w http.ResponseWriter) error {
	marshal, _ := json.Marshal(r)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(marshal)
	if err != nil {
		return fmt.Errorf("error marshal %w", err)
	}
	return nil
}
