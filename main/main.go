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
	MinerID      uint64 `json:"minerID"`
}

type Commit2Response struct {
	RetCode 	int64					`json:"retCode"`
	ErrMsg		string					`json:"errMsg"`
	Addition    interface{}				`json:"addition"`
	Proof    	[]byte					`json:"proof"`
}

func main() {
	r := ginRouter()
	r.Run(":8008")
}

func ginRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/seal", sealCommit2)
	r.POST("/result", getResult)

	return r
}

func sealCommit2(c *gin.Context) {
	fmt.Println("Enter sealCommit2 function...............")
	req := &Commit2Request{}
	preHandleRequest(c, req)

	// TODO: handle business

	// return
	fmt.Println("Exiting sealCommit2 function...............")
	c.Data(http.StatusOK, "application/json", createResponse(RetCodeOK, "Successful", nil, nil))
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
	c.Data(http.StatusOK, "application/json", createResponse(RetCodeOK, "Successful", nil, proof))
}

func preHandleRequest(c *gin.Context, req *Commit2Request) {
	b, e := c.GetRawData()
	if e != nil {
		c.Data(http.StatusBadRequest, "application/json", createResponse(RetCodeBadRequest, fmt.Sprintf("bad request. get row data failed. error: %+v", e.Error()), nil, nil))
	}

	reader, e := gzip.NewReader(bytes.NewReader(b))
	defer reader.Close()
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. unzip data failed. error: %+v", e.Error()), nil, nil))
	}

	data, e := ioutil.ReadAll(reader)
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. unzip data failed. error: %+v", e.Error()), nil, nil))
	}

	e = json.Unmarshal(data, req)
	if e != nil {
		c.Data(http.StatusInternalServerError, "application/json", createResponse(RetCodeInternalServerError, fmt.Sprintf("interval server error. parse data failed. error: %+v", e.Error()), nil, nil))
	}
}

func createResponse(retCode int64, errMsg string, addition interface{}, proof []byte) []byte {
	resp := Commit2Response{
		RetCode: 	retCode,
		ErrMsg:  	errMsg,
		Addition: 	addition,
		Proof:		proof,
	}

	bytes, _ := json.Marshal(resp)
	return bytes
}
