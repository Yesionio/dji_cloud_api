package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type CbLogFileUploadData struct {
	Files  []CbLogFileUploadListData `json:"files"`
	Result int                       `json:"result"`
}

type CbLogFileUploadListData struct {
	DeviceSN string                        `json:"device_sn"`
	Result   int                           `json:"result"`
	Module   int                           `json:"module"`
	List     []CbLogFileUploadListInfoData `json:"list"`
}

type CbLogFileUploadListInfoData struct {
	BootIndex int `json:"boot_index"`
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
	Size      int `json:"size"`
}

type FnEvtLogFileUploadProgress func(sn string, data *CommonProgressData)
type FnServiceCallbackLogUpload func(data *CbLogFileUploadData)

type RemoteLogModule struct {
	client mqtt.Client

	callbackResultList       map[string]FnServiceCallbackResult
	callbackLogUploadList    map[string]FnServiceCallbackLogUpload
	evtLogFileUploadProgress FnEvtLogFileUploadProgress
}

// SetOnEvtLogFileUploadProgress 当日志上传进度更新时调用
func (rlm *RemoteLogModule) SetOnEvtLogFileUploadProgress(fn FnEvtLogFileUploadProgress) {
	rlm.evtLogFileUploadProgress = fn
}

// FileUploadListAsync 获取设备可上传的文件列表
func (rlm *RemoteLogModule) FileUploadListAsync(sn string, moduleList []string, callback FnServiceCallbackLogUpload) error {
	tid := uuid.NewString()
	rlm.callbackLogUploadList[tid] = callback

	param := map[string]any{
		"module_list": moduleList,
	}

	err := sendService(rlm.client, sn, tid, "fileupload_list", &param)
	if err != nil {
		return err
	}
	return nil
}

type ParamLogFileUploadStart struct {
	Bucket      string                      `json:"bucket"`
	Region      string                      `json:"region"`
	Credentials ParamLogFileUploadStartCred `json:"credentials"`
	Endpoint    string                      `json:"endpoint"`
	Provider    string                      `json:"provider"`
}

type ParamLogFileUploadStartCred struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Expire          int    `json:"expire"`
	SecurityToken   string `json:"security_token"`
}

type ParamLogFileUploadStartParam struct {
	Files []ParamLogFileUploadStartFiles `json:"files"`
}

type ParamLogFileUploadStartFiles struct {
	ObjectKey string `json:"object_key"`
	Module    string `json:"module"`
	List      []struct {
		BootIndex int `json:"boot_index"`
	} `json:"list"`
}

// FileUploadStartAsync 发起日志文件上传
func (rlm *RemoteLogModule) FileUploadStartAsync(sn string, param ParamLogFileUploadStart, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	rlm.callbackResultList[tid] = callback

	err := sendService(rlm.client, sn, tid, "fileupload_start", &param)
	if err != nil {
		return err
	}
	return nil
}

type ParamLogUploadUpdateData struct {
	Status     string   `json:"status"`
	ModuleList []string `json:"module_list"`
}

// FileUploadUpdateAsync 上传状态更新
func (rlm *RemoteLogModule) FileUploadUpdateAsync(sn string, param ParamLogUploadUpdateData, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	rlm.callbackResultList[tid] = callback

	err := sendService(rlm.client, sn, tid, "fileupload_update", &param)
	if err != nil {
		return err
	}
	return nil
}

func (rlm *RemoteLogModule) remoteLogEventHandler(msg *MessageData) error {
	if rlm.evtLogFileUploadProgress != nil {
		data := &CommonProgressData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		rlm.evtLogFileUploadProgress(msg.SN, data)

		err = sendEventReply(rlm.client, msg, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rlm *RemoteLogModule) replyHandler(msg *MessageData) error {

	if fn, has := rlm.callbackResultList[msg.Payload.Tid]; has {
		result := &ReplyResultData{}
		err := json.Unmarshal(msg.Payload.Data, result)
		if err != nil {
			return err
		}
		fn(result.Result)
		delete(rlm.callbackResultList, msg.Payload.Tid)
	}

	if fn, has := rlm.callbackLogUploadList[msg.Payload.Tid]; has {
		data := &CbLogFileUploadData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		fn(data)
		delete(rlm.callbackLogUploadList, msg.Payload.Tid)
	}
	return nil
}

func newRemoteLogModuleLog(client mqtt.Client) *RemoteLogModule {
	r := new(RemoteLogModule)
	r.client = client
	r.callbackResultList = make(map[string]FnServiceCallbackResult)

	return r
}
