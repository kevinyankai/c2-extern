package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const RetCodeOK = 0
const RetCodeBadRequest = 400
const RetCodeInternalServerError = 500
const RetCodeUnknownError = -1

type Commit2Request struct {
	C2In         []byte `json:"c2In"`
	SectorNumber uint64 `json:"sectorNumber"`
	MinerID      string `json:"minerID"`
}

type Commit2Response struct {
	RetCode 	int					`json:"code"`
	Message		string				`json:"msg"`
	Data    interface{}				`json:"data"`
}

func main() {
	r := ginRouter()
	r.Run("192.168.0.167:8081")
}

func ginRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/allocate_task", sealCommit2)
	r.POST("/query_result", getResult)

	return r
}

func sealCommit2(c *gin.Context) {
	fmt.Println("Enter sealCommit2 function...............")
	req := &Commit2Request{}
	preHandleRequest(c, req)

	// TODO: handle business

	// return
	fmt.Println("Exiting sealCommit2 function...............")
	c.Data(http.StatusOK, "application/json", createResponse(RetCodeOK, "Successful", nil))
}

func getResult(c *gin.Context) {
	fmt.Println("Enter getResult function...............")
	req := &Commit2Request{}
	preHandleRequest(c, req)

	// TODO: handle business
	proof := make([]byte, 3)
	proof[0] = 'a'
	proof[1] = 'b'
	proof[2] = 'c'

	fmt.Println("Exiting getResult function...............")
	c.Data(http.StatusOK, "application/json", createResponse(RetCodeOK, "Successful", proof))
}

func preHandleRequest(c *gin.Context, req *Commit2Request) {
	b, e := c.GetRawData()
	if e != nil {
		c.Data(http.StatusBadRequest, "application/json", createResponse(RetCodeBadRequest, fmt.Sprintf("bad request. get row data failed. error: %+v", e.Error()), nil))
	}

	reader, e := gzip.NewReader(bytes.NewReader(b))
	defer reader.Close()
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. unzip data failed. error: %+v", e.Error()), nil))
	}

	data, e := ioutil.ReadAll(reader)
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. unzip data failed. error: %+v", e.Error()), nil))
	}

	e = json.Unmarshal(data, req)
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. parse data failed. error: %+v", e.Error()), nil))
	}
}

func createResponse(retCode int, errMsg string, data interface{}) []byte {
	resp := Commit2Response{
		RetCode: 	retCode,
		Message:  	errMsg,
		Data: 		data,
	}

	bytes, _ := json.Marshal(resp)
	return bytes
}
