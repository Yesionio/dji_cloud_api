package djicloudapi

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttProtocolVersion uint
type CloudApiErrHandler func(error)

const (
	MqttProtocolVersionV31  MqttProtocolVersion = 3
	MqttProtocolVersionV311 MqttProtocolVersion = 4
	MqttProtocolVersionV5   MqttProtocolVersion = 5
)

type CloudApiOption struct {
	MqttClientID        string
	MqttBroker          string
	MqttCleanSession    bool
	MqttProtocolVersion MqttProtocolVersion
	MqttEnableAuth      bool
	MqttUsername        string
	MqttPassword        string
	ErrHandler          CloudApiErrHandler
}

type DjiCloudApiCore struct {
	mqttOpt            *mqtt.ClientOptions
	mqttClient         mqtt.Client
	errHandler         CloudApiErrHandler
	OrgManageModule    *OrgManageModule
	DeviceManageModule *DeviceManageModule
	LiveStreamModule   *LiveStreamModule
	MediaManageModule  *MediaManageModule
	WaylineModule      *WaylineModule
	HMSModule          *HMSModule
	DebugModule        *DebugModule
}

func NewCloudApi(opt CloudApiOption) *DjiCloudApiCore {
	core := new(DjiCloudApiCore)
	mqttOpt := mqtt.NewClientOptions()
	mqttOpt.SetClientID(opt.MqttClientID)
	mqttOpt.SetCleanSession(opt.MqttCleanSession)
	mqttOpt.AddBroker(opt.MqttBroker)
	mqttOpt.SetProtocolVersion(uint(opt.MqttProtocolVersion))
	if opt.MqttEnableAuth {
		mqttOpt.SetUsername(opt.MqttUsername)
		mqttOpt.SetPassword(opt.MqttPassword)
	}

	mqttOpt.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		fmt.Println("mqtt连接丢失", err)
	})

	mqttOpt.SetOnConnectHandler(core.djiSubscribe)

	core.mqttOpt = mqttOpt

	core.mqttClient = mqtt.NewClient(mqttOpt)
	core.errHandler = opt.ErrHandler

	loadModule(core)

	return core
}

func loadModule(core *DjiCloudApiCore) {
	core.OrgManageModule = newOrgManageModule(core.mqttClient)
	core.DeviceManageModule = newDeviceManageModule(core.mqttClient)
	core.LiveStreamModule = newLiveStreamModule(core.mqttClient)
	core.MediaManageModule = newMediaManageModule(core.mqttClient)
	core.WaylineModule = newWaylineModule(core.mqttClient)
	core.HMSModule = newHMSModule(core.mqttClient)
	core.DebugModule = newDebugModule(core.mqttClient)
}

func (dji *DjiCloudApiCore) Run() error {
	if t := dji.mqttClient.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (dji *DjiCloudApiCore) djiSubscribe(c mqtt.Client) {
	fmt.Println("开始连接订阅相关通道")

	c.Subscribe("thing/product/+/osd", 0, dji.osdHandler)
	c.Subscribe("thing/product/+/state", 0, dji.stateHandler)
	c.Subscribe("thing/product/+/services_reply", 0, dji.serviceReplyHandler)
	c.Subscribe("thing/product/+/events", 0, dji.eventsHandler)
	c.Subscribe("thing/product/+/requests", 0, dji.requestHandler)
	c.Subscribe("sys/product/+/status", 0, dji.statusHandler)
	c.Subscribe("thing/product/+/property/set_reply", 0, func(c mqtt.Client, m mqtt.Message) {})
	c.Subscribe("thing/product/+/drc/up", 0, func(c mqtt.Client, m mqtt.Message) {})
}

func (dji *DjiCloudApiCore) requestHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "thing/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/requests", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	switch data.Method {
	case "airport_bind_status":
		if err = dji.OrgManageModule.airportBindStatus(msg); err != nil {
			dji.errHandler(err)
		}
	case "airport_organization_get":
		if err = dji.OrgManageModule.airportOrganizationGet(msg); err != nil {
			dji.errHandler(err)
		}
	case "airport_organization_bind":
		if err = dji.OrgManageModule.airportOrganizationBind(msg); err != nil {
			dji.errHandler(err)
		}
	case "storage_config_get":
		if err = dji.MediaManageModule.storageConfigGetHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "flighttask_resource_get":
		if err = dji.WaylineModule.resourceGetHandler(msg); err != nil {
			dji.errHandler(err)
		}
	}
}

func (dji *DjiCloudApiCore) osdHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "thing/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/osd", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	if err = dji.DeviceManageModule.osdHandler(msg); err != nil {
		dji.errHandler(err)
	}
}

func (dji *DjiCloudApiCore) stateHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "thing/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/state", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	if err = dji.DeviceManageModule.stateHandler(msg); err != nil {
		dji.errHandler(err)
	}
}

func (dji *DjiCloudApiCore) statusHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "sys/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/status", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	switch data.Method {
	case "update_topo":
		if err = dji.DeviceManageModule.topoHandler(msg); err != nil {
			dji.errHandler(err)
		}
	}
}

func (dji *DjiCloudApiCore) serviceReplyHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "sys/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/services_reply", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	switch data.Method {
	case "live_start_push", "live_stop_push", "live_set_quality", "live_lens_change":
		if err := dji.LiveStreamModule.replyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "upload_flighttask_media_prioritize":
		if err := dji.MediaManageModule.replyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "flighttask_prepare", "flighttask_execute", "flighttask_undo", "flighttask_pause", "flighttask_recovery", "return_home", "return_home_cancel":
		if err := dji.WaylineModule.replyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	}
}

func (dji *DjiCloudApiCore) eventsHandler(c mqtt.Client, m mqtt.Message) {
	topic := m.Topic()
	topic = strings.ReplaceAll(topic, "sys/product/", "")
	data := new(CommonRequestData)

	err := json.Unmarshal(m.Payload(), data)
	if err != nil {
		dji.errHandler(err)
	}
	sn := strings.ReplaceAll(topic, "/events", "")

	msg := &MessageData{
		SN:      sn,
		Topic:   topic,
		Payload: data.CommonData,
	}

	switch data.Method {
	case "file_upload_callback":
		if err := dji.MediaManageModule.fileUploadCallbackHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "highest_priority_upload_flighttask_media":
		if err := dji.MediaManageModule.highestPriorityUploadHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "device_exit_homing_notify":
		if err := dji.WaylineModule.deviceExitHomingNotifyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "flighttask_progress":
		if err := dji.WaylineModule.flighttaskProgressHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "flighttask_ready":
		if err := dji.WaylineModule.flightTaskReadyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "return_home_info":
		if err := dji.WaylineModule.returnHomeInfoHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "hms":
		if err := dji.HMSModule.hmsHandler(msg); err != nil {
			dji.errHandler(err)
		}
	}
}