package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type HMSData struct {
	List []HMSListData `json:"list"`
}

type HMSListData struct {
	Level      int    `json:"level"`
	Module     int    `json:"module"`
	InTheSky   int    `json:"in_the_sky"`
	Code       string `json:"code"`
	DeviceType string `json:"device_type"`
	Imminent   int    `json:"imminent"`
}

type HMSListArgsData struct {
	ComponentIndex int `json:"component_index"`
	SensorIndex    int `json:"sensor_index"`
}

type FnEvtHMS func(sn string, data *HMSData)

type HMSModule struct {
	client mqtt.Client
	evtHMS FnEvtHMS
}

func (hm *HMSModule) SetOnEvtHMS(fn FnEvtHMS) {
	hm.evtHMS = fn
}

func (hm *HMSModule) hmsHandler(msg *MessageData) error {
	if hm.evtHMS != nil {
		data := &HMSData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		hm.evtHMS(msg.SN, data)

		err = sendEventReply(hm.client, msg, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func newHMSModule(client mqtt.Client) *HMSModule {
	r := new(HMSModule)
	r.client = client
	return r
}
