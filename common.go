package djicloudapi

import "encoding/json"

type CommonData[T json.RawMessage | CommonOutPutData] struct {
	Tid       string `json:"tid"`
	Bid       string `json:"bid"`
	Timestamp int64  `json:"timestamp"`
	Gateway   string `json:"gateway"`
	Data      T      `json:"data"`
	Method    string `json:"method,omitempty"`
}

type CommonRequestData struct {
	CommonData[json.RawMessage]
}

type CommonOutPutData struct {
	Output any `json:"output,omitempty"`
	Result int `json:"result"`
}

type MessageData struct {
	SN      string
	Topic   string
	Payload CommonData[json.RawMessage]
}

type ReplyResultData struct {
	Result int `json:"result"`
}

type ReplyStatusData struct {
	Status string `json:"status"`
}

type CommonProgressData struct {
	Status   string               `json:"status"`
	Progress CommonProgressDetail `json:"progress"`
}

type CommonProgressDetail struct {
	Percent int    `json:"percent"`
	StepKey string `json:"step_key"`
}

type FnServiceCallbackResult func(result int)
type FnServiceCallbackStatus func(status string)
