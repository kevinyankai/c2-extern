package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinRouter(t *testing.T) {
	//router := ginRouter()

	w := httptest.NewRecorder()
	data := make([]byte, 3)
	data[0] = '1'
	data[1] = '2'
	data[2] = '3'
	body := Commit2Request{
		MinerID: "f018888",
		SectorNumber: 1,
		C2In: data,
	}

	b, _ := json.Marshal(body)
	buf := new(bytes.Buffer)
	writer := gzip.NewWriter(buf)
	writer.Write(b)
	writer.Flush()
	writer.Close()

	//req, _ := http.NewRequest("POST", "/allocate_task", buf)
	//req, _ := http.NewRequest("POST", "/query_result", buf)
	//router.ServeHTTP(w, req)
	//req, err := http.NewRequest("POST", fmt.Sprint("http://192.168.0.167:8081/query_result"), buf)
	req, err := http.NewRequest("POST", fmt.Sprint("http://192.168.0.167:8081/allocate_task"), buf)
	if err != nil {
		t.Errorf("http create request failed. error: %+v", err)
		return
	}
	req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")

	res, err := http.DefaultClient.Do(req)
	retBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("read response body error: %+v", err)
		return
	}

	resp := &Commit2Response{}

	json.Unmarshal(retBody, resp)
	//result := resp.Data.([]byte)
	//t.Logf("Result: %+v", resp.Data)
	//r := []byte(resp.Data.(string))
	//r2 := make([]byte, 0)
	//json.Unmarshal(r, &r2)
	//t.Logf("Result %+v", r2)
	assert.Equal(t, http.StatusOK, w.Code)
}