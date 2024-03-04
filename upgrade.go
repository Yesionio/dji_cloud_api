package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type FnEvtOtaProgress func(sn string, data *CommonProgressData)

type UpgradeModule struct {
	client mqtt.Client

	callbackStatusList map[string]FnServiceCallbackStatus
	evtOtaProgress     FnEvtOtaProgress
}

func (um *UpgradeModule) otaEventHandler(msg *MessageData) error {
	if um.evtOtaProgress != nil {
		data := &CommonProgressData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		um.evtOtaProgress(msg.SN, data)

		err = sendEventReply(um.client, msg, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (um *UpgradeModule) SetOnEvtOtaProgress(fn FnEvtOtaProgress) {
	um.evtOtaProgress = fn
}

type ParamOtaCreate struct {
	SN                  string `json:"sn"`
	ProductVersion      string `json:"product_version"`
	FileURL             string `json:"file_url"`
	MD5                 string `json:"md5"`
	FileSize            int    `json:"file_size"`
	FileName            string `json:"file_name"`
	FirmwareUpgradeType int    `json:"firmware_upgrade_type"`
}

func (um *UpgradeModule) OtaCreateAsync(sn string, param []ParamOtaCreate, callback FnServiceCallbackStatus) error {
	tid := uuid.NewString()
	um.callbackStatusList[tid] = callback

	err := sendService(um.client, sn, tid, "ota_create", &param)
	if err != nil {
		return err
	}
	return nil
}

func (um *UpgradeModule) replyHandler(msg *MessageData) error {

	if fn, has := um.callbackStatusList[msg.Payload.Tid]; has {
		status := &ReplyStatusData{}
		err := json.Unmarshal(msg.Payload.Data, status)
		if err != nil {
			return err
		}
		fn(status.Status)
		delete(um.callbackStatusList, msg.Payload.Tid)
	}
	return nil
}

func newUpgradeModule(client mqtt.Client) *UpgradeModule {
	r := new(UpgradeModule)
	r.client = client
	r.callbackStatusList = make(map[string]FnServiceCallbackStatus)
	return r
}
