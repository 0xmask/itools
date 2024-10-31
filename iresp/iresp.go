package iresp

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Debug bool //开启调试，错误信息会响应给前端

type Message struct {
	Status int         `json:"status"` //状态0错1对
	Code   int         `json:"code"`   //错误码
	Msg    interface{} `json:"msg"`    //提示信息
	Data   interface{} `json:"data"`   //响应数据
}

func GrpcError(code int, errs ...error) Message {
	var errMsg interface{}
	var errCode = code
	if len(errs) > 0 {
		sta := status.Convert(errs[0])
		errMsg = sta.Message()
		if sta.Code() != codes.OK {
			errCode = int(sta.Code())
		}
	}

	return Message{
		Status: 0,
		Code:   errCode,
		Msg:    errMsg,
		Data:   nil,
	}
}

// Error 格式化错误相应
// err 错误信息，支持一般错误和rpc格式化错误，需要配合resp.StatusError使用
func Error(code int, errs ...error) Message {
	var errMsg interface{}
	if len(errs) > 0 {
		errMsg = errs[0].Error()
	}

	return Message{
		Status: 0,
		Code:   code,
		Msg:    errMsg,
		Data:   nil,
	}
}

// Success 格式化正确信息
// data 相应数据
func Success(code int, data ...interface{}) Message {
	var respData interface{}
	if len(data) > 0 {
		respData = data[0]
	}
	return Message{
		Status: 1,
		Code:   code,
		Msg:    "success",
		Data:   respData,
	}
}
