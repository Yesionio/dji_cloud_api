package djicloudapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	// 时间同步服务器(ntp)地址 -> 机巢
	NtpServerHost string
	// 在DJI开发者网站的AppID -> 机巢
	AppID string
	// 在DJI开发者网站的Key -> 机巢
	AppKey string
	// 在DJI开发者网站的License -> 机巢
	AppLicense string
}

type DjiCloudApiCore struct {
	mqttOpt            *mqtt.ClientOptions
	option             CloudApiOption
	mqttClient         mqtt.Client
	errHandler         CloudApiErrHandler
	OrgManageModule    *OrgManageModule
	DeviceManageModule *DeviceManageModule
	LiveStreamModule   *LiveStreamModule
	MediaManageModule  *MediaManageModule
	WaylineModule      *WaylineModule
	HMSModule          *HMSModule
	DebugModule        *DebugModule
	UpgradeModule      *UpgradeModule
	RemoteLogModule    *RemoteLogModule
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
	core.option = opt

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
	core.UpgradeModule = newUpgradeModule(core.mqttClient)
	core.RemoteLogModule = newRemoteLogModuleLog(core.mqttClient)
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

	dji.DebugModule.DebugModeOpenAsync("7CTDM2100BH29T", func(status string) {
		fmt.Println("open debug status is ", status)

		// dji.DebugModule.DroneCloseAsync("7CTDM2100BH29T", func(status string) {
		// 	fmt.Println("关闭飞行器", status)
		// dji.DebugModule.ChargeCloseAsync("7CTDM2100BH29T", func(status string) {
		// 	fmt.Println("关闭充电", status)
		// })
		dji.DebugModule.DroneOpenAsync("7CTDM2100BH29T", func(status string) {
			fmt.Println("打开飞机", status)

			dji.DebugModule.DebugModeCloseAsync("7CTDM2100BH29T", func(status string) {
				fmt.Println("close debug status is ", status)
			})
		})
		// })
	})
}

func (dji *DjiCloudApiCore) requestsConfigUpdate(msg *MessageData) error {
	fmt.Println("call config update", string(msg.Payload.Data))
	data := map[string]any{
		"ntp_server_host": dji.option.NtpServerHost,
		"app_id":          dji.option.AppID,
		"app_key":         dji.option.AppKey,
		"app_license":     dji.option.AppLicense,
	}
	r, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	pubData := CommonData[json.RawMessage]{
		Tid:       msg.Payload.Tid,
		Bid:       msg.Payload.Bid,
		Timestamp: time.Now().Unix(),
		Gateway:   msg.Payload.Gateway,
		Method:    msg.Payload.Method,
		Data:      r,
	}
	marData, err := json.Marshal(pubData)
	if err != nil {
		return err
	}
	t := dji.mqttClient.Publish(msg.Topic+REPLY_SUFFIX, 0, false, marData)
	if ok := t.WaitTimeout(time.Second * 10); !ok || t.Error() != nil {
		return t.Error()
	}
	return nil
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
		Topic:   m.Topic(),
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
	case "config":
		if err = dji.requestsConfigUpdate(msg); err != nil {
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
		Topic:   m.Topic(),
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
		Topic:   m.Topic(),
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
		Topic:   m.Topic(),
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
		Topic:   m.Topic(),
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
	case "debug_mode_open", "debug_mode_close", "supplement_light_open", "supplement_light_close",
		"battery_maintenance_switch", "air_conditioner_mode_switch", "alarm_state_switch",
		"battery_store_mode_switch", "device_reboot", "drone_open", "drone_close", "device_format",
		"drone_format", "cover_open", "cover_close", "charge_open", "charge_close", "sdr_workmode_switch",
		"esim_activate", "sim_slot_switch", "esim_operator_switch":
		if err := dji.DebugModule.replyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "ota_create":
		if err := dji.UpgradeModule.replyHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "fileupload_list", "fileupload_start", "fileupload_update":
		if err := dji.RemoteLogModule.replyHandler(msg); err != nil {
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
		Topic:   m.Topic(),
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
	case "drone_open", "drone_close", "device_reboot", "cover_close", "cover_open", "charge_open", "charge_close", "drone_format", "device_format", "esim_activate", "esim_operator_switch":
		if err := dji.DebugModule.debugEventHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "ota_progress":
		if err := dji.UpgradeModule.otaEventHandler(msg); err != nil {
			dji.errHandler(err)
		}
	case "fileupload_progress":
		if err := dji.RemoteLogModule.remoteLogEventHandler(msg); err != nil {
			dji.errHandler(err)
		}
	}
}
