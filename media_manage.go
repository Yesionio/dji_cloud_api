package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type FnEvtFileUploadCallback func(sn string, data *FileUploadData) error
type FnEvtHighestPriorityUpload func(sn string, flightID string) error
type FnStorageConfigGet func(sn string, module int) (*StorageConfigGetData, error)

type MediaManageModule struct {
	client                   mqtt.Client
	callbackList             map[string]FnServiceCallbackResult
	evtFileUploadCallback    FnEvtFileUploadCallback
	evtHighestPriorityUpload FnEvtHighestPriorityUpload
	fnStorageConfigGet       FnStorageConfigGet
}

// SetOnEvtFileUpload 订阅文件上传回调
func (mmm *MediaManageModule) SetOnEvtFileUpload(fn FnEvtFileUploadCallback) {
	mmm.evtFileUploadCallback = fn
}

// SetOnEvtHighestPriorityUpload 当文件优先级上报时调用
func (mmm *MediaManageModule) SetOnEvtHighestPriorityUpload(fn FnEvtHighestPriorityUpload) {
	mmm.evtHighestPriorityUpload = fn
}

// UploadFlighttaskMediaPrioritizeAsync 设置文件上传优先级
func (mmm *MediaManageModule) UploadFlighttaskMediaPrioritizeAsync(sn string, flightID string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	mmm.callbackList[tid] = callback

	param := map[string]any{
		"flight_id": flightID,
	}

	err := sendService(mmm.client, sn, tid, "upload_flighttask_media_prioritize", &param)
	if err != nil {
		return err
	}

	return nil
}

func (mmm *MediaManageModule) SetOnStorageConfigGet(fn FnStorageConfigGet) {
	if fn != nil {
		mmm.fnStorageConfigGet = fn
	}
}

func (mmm *MediaManageModule) fileUploadCallbackHandler(msg *MessageData) error {
	if mmm.evtFileUploadCallback != nil {
		data := &FileUploadData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}
		err = mmm.evtFileUploadCallback(msg.SN, data)
		if err != nil {
			sendEventReply(mmm.client, msg, 1)
		} else {
			sendEventReply(mmm.client, msg, 0)
		}
	}
	return nil
}

type paramHighPriorityUpload struct {
	FlightID string `json:"flight_id"`
}

func (mmm *MediaManageModule) highestPriorityUploadHandler(msg *MessageData) error {
	if mmm.evtHighestPriorityUpload != nil {
		data := &paramHighPriorityUpload{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}
		err = mmm.evtHighestPriorityUpload(msg.SN, data.FlightID)
		if err != nil {
			sendEventReply(mmm.client, msg, 1)
		} else {
			sendEventReply(mmm.client, msg, 0)
		}
	}
	return nil
}

type paramStorageConfigGet struct {
	Module int `json:"module"`
}

func (mmm *MediaManageModule) storageConfigGetHandler(msg *MessageData) error {
	if mmm.fnStorageConfigGet != nil {
		data := &paramStorageConfigGet{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		rls, err := mmm.fnStorageConfigGet(msg.SN, data.Module)
		if err != nil {
			sendOutput(mmm.client, msg, rls, 1)
			return err
		}

		err = sendOutput(mmm.client, msg, rls, 0)
		if err != nil {
			return err
		}

	}
	return nil
}

func (mmm *MediaManageModule) replyHandler(msg *MessageData) error {
	result := &ReplyResultData{}
	err := json.Unmarshal(msg.Payload.Data, result)
	if err != nil {
		return err
	}

	if fn, has := mmm.callbackList[msg.Payload.Tid]; has {
		fn(result.Result)
		delete(mmm.callbackList, msg.Payload.Tid)
	}
	return nil
}

func newMediaManageModule(client mqtt.Client) *MediaManageModule {
	r := new(MediaManageModule)
	r.client = client
	r.callbackList = make(map[string]FnServiceCallbackResult)
	return r
}
