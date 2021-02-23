package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinRouter(t *testing.T) {
	router := ginRouter()

	w := httptest.NewRecorder()
	data := make([]byte, 3)
	data[0] = '1'
	data[1] = '2'
	data[2] = '3'
	body := Commit2Request{
		MinerID: 1,
		SectorNumber: 2,
		C2In: data,
	}

	b, _ := json.Marshal(body)
	buf := new(bytes.Buffer)
	writer := gzip.NewWriter(buf)
	writer.Write(b)
	writer.Flush()
	writer.Close()

	req, _ := http.NewRequest("POST", "/result", buf)
	router.ServeHTTP(w, req)

	resp := &Commit2Response{}

	json.Unmarshal(w.Body.Bytes(), resp)
	assert.Equal(t, http.StatusOK, w.Code)
}