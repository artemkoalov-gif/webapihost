package http_server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (list *Links) HandleAdd(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload DataPayload
		var respLinks ResponseLinks

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log(logrus.WarnLevel, http.StatusBadRequest, "Failed to read body", err)
			return
		}

		if err := json.Unmarshal(body, &payload); err != nil {
			logger.Log(logrus.WarnLevel, "Invalid JSON", err)
			return
		}

		if len(payload.Url) == 0 {
			logger.Log(logrus.WarnLevel, "Missing required list of urls")
			httpError(w, "Missing required list of urls")
			return
		}

		m := make(map[string]string)
		for _, url := range payload.Url {
			logger.Log(logrus.InfoLevel, "Start check url: "+url)
			s, _ := httpRequest(url)
			m[url] = fmt.Sprintf("%s", s)
		}

		respLinks.Links = m
		respLinks.LinksNum = len(list.Data) + 1
		err = respLinks.responseLinks(w)
		if err != nil {
			logger.Log(logrus.WarnLevel, err)
		}

		list.mu.Lock()
		list.Data[respLinks.LinksNum] = respLinks
		list.mu.Unlock()
	}
}

func (list *Links) HandleGet(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ReqLinks

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log(logrus.WarnLevel, http.StatusBadRequest, "Failed to read body", err)
			return
		}

		if err := json.Unmarshal(body, &payload); err != nil {
			logger.Log(logrus.WarnLevel, "Invalid JSON", err)
			return
		}

		list.mu.Lock()
		dataCopy := make(map[int]ResponseLinks, len(list.Data))
		for _, idList := range payload.LinksList {
			l, ok := list.Data[idList]
			if ok {
				dataCopy[idList] = l
			}
		}
		list.mu.Unlock()

		lines := buildLinesForPDF(dataCopy, payload.LinksList)
		if len(lines) == 0 {
			httpError(w, "Links for provided ids not found")
			return
		}

		pdfBytes, err := buildPDF(lines)
		if err != nil {
			logger.Log(logrus.ErrorLevel, "Failed to build PDF", err)
			httpError(w, "Failed to build PDF")
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"links.pdf\"")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(pdfBytes); err != nil {
			logger.Log(logrus.WarnLevel, "Failed to write PDF response", err)
		}
	}
}
