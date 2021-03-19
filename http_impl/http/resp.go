package http

import (
	"encoding/json"
	"net/http"
)

type DataResp struct {
	ErrorCode int         `json:"error_code"` //0成功 1失败
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
}

func Response(w http.ResponseWriter, data interface{}, errorCode int, msg string) {
	Respdata := &DataResp{
		ErrorCode: errorCode,
		Msg:       msg,
		Data:      data,
	}
	dataByte, _ := json.Marshal(Respdata)
	w.Write(dataByte)
}
