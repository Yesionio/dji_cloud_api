package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type DebugProgressData struct {
	Status   string              `json:"status"`
	Progress DebugProgressDetail `json:"progress"`
}

type DebugProgressDetail struct {
	Percent int    `json:"percent"`
	StepKey string `json:"step_key"`
}

type FnEvtDebugProgress func(sn, event string, data *DebugProgressData)

type DebugModule struct {
	client             mqtt.Client
	callbackStatusList map[string]FnServiceCallbackStatus

	evtDebugProgress FnEvtDebugProgress
}

func (dgb *DebugModule) SetOnEvtDebugProgress(fn FnEvtDebugProgress) {
	dgb.evtDebugProgress = fn
}

// DebugModeOpenAsync 打开调试模式
func (dbg *DebugModule) DebugModeOpenAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "debug_mode_open", callback)
}

// DebugModeOpenAsync 关闭调试模式
func (dbg *DebugModule) DebugModeCloseAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "debug_mode_close", callback)
}

// SupplementLightOpenAsync 打开补光灯
func (dbg *DebugModule) SupplementLightOpenAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "supplement_light_open", callback)
}

// SupplementLightCloseAsync 打开补光灯
func (dbg *DebugModule) SupplementLightCloseAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "supplement_light_close", callback)
}

// BatteryMaintenanceSwitchAsync 电池保养模式切换
func (dbg *DebugModule) BatteryMaintenanceSwitchAsync(sn string, enable int, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "battery_maintenance_switch", callback, map[string]any{"action": enable})
}

// AirConditionerModeSwitchAsync 机场空调工作模式切换
func (dbg *DebugModule) AirConditionerModeSwitchAsync(sn string, action int, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "air_conditioner_mode_switch", callback, map[string]any{"action": action})
}

// AlarmStateSwitchAsync 机场声光报警开关
func (dbg *DebugModule) AlarmStateSwitchAsync(sn string, enable int, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "alarm_state_switch", callback, map[string]any{"action": enable})
}

// BatteryStoreModeSwitchAsync 电池运行模式切换
func (dbg *DebugModule) BatteryStoreModeSwitchAsync(sn string, action int, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "battery_store_mode_switch", callback, map[string]any{"action": action})
}

// SdrWorkmodeSwitchAsync 增强图传开关
func (dbg *DebugModule) SdrWorkmodeSwitchAsync(sn string, action int, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "sdr_workmode_switch", callback, map[string]any{"link_workmode": action})
}

type ParamEsimActivate struct {
	Imei       string `json:"imei"`
	DeviceType string `json:"device_type"`
}

// EsimActivateAsync eSIM激活
func (dbg *DebugModule) EsimActivateAsync(sn string, param ParamEsimActivate, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "esim_activate", callback, map[string]any{
		"imei":        param.Imei,
		"device_type": param.DeviceType,
	})
}

type ParamEsimSlotSwitch struct {
	Imei       string `json:"imei"`
	DeviceType string `json:"device_type"`
	SimSlot    int    `json:"sim_slot"`
}

// EsimActivateAsync eSIM 和 SIM 切换
func (dbg *DebugModule) EsimSlotSwitchAsync(sn string, param ParamEsimSlotSwitch, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "sim_slot_switch", callback, map[string]any{
		"imei":        param.Imei,
		"device_type": param.DeviceType,
		"sim_slot":    param.SimSlot,
	})
}

type ParamEsimOperatorSwitch struct {
	Imei            string `json:"imei"`
	DeviceType      string `json:"device_type"`
	TelecomOperator int    `json:"telecom_operator"`
}

// EsimActivateAsync eSIM激活
func (dbg *DebugModule) EsimOperatorSwitchAsync(sn string, param ParamEsimOperatorSwitch, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "sim_slot_switch", callback, map[string]any{
		"imei":             param.Imei,
		"device_type":      param.DeviceType,
		"telecom_operator": param.TelecomOperator,
	})
}

// DeviceRebootAsync 机场重启
func (dbg *DebugModule) DeviceRebootAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "device_reboot", callback)
}

// DroneOpenAsync 飞行器开机
func (dbg *DebugModule) DroneOpenAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "drone_open", callback)
}

// DroneCloseAsync 飞行器关机
func (dbg *DebugModule) DroneCloseAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "drone_close", callback)
}

// DeviceFormatAsync 机场数据格式化
func (dbg *DebugModule) DeviceFormatAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "device_format", callback)
}

// DroneFormatAsync 飞行器数据格式化
func (dbg *DebugModule) DroneFormatAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "drone_format", callback)
}

// CoverOpenAsync 打开舱盖
func (dbg *DebugModule) CoverOpenAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "cover_open", callback)
}

// CoverCloseAsync 关闭舱盖
func (dbg *DebugModule) CoverCloseAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "cover_close", callback)
}

// ChargeOpenAsync 打开充电
func (dbg *DebugModule) ChargeOpenAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "charge_open", callback)
}

// ChargeCloseAsync 关闭充电
func (dbg *DebugModule) ChargeCloseAsync(sn string, callback FnServiceCallbackStatus) error {
	return dbg.sendStructure(sn, "charge_close", callback)
}

func (dbg *DebugModule) sendStructure(sn, command string, callback FnServiceCallbackStatus, data ...map[string]any) error {
	tid := uuid.NewString()
	dbg.callbackStatusList[tid] = callback

	param := map[string]any{}
	if len(data) > 0 {
		param = data[0]
	}

	err := sendService(dbg.client, sn, tid, command, &param)
	if err != nil {
		return err
	}
	return nil
}

func (dbg *DebugModule) debugEventHandler(msg *MessageData) error {
	if dbg.evtDebugProgress != nil {
		data := &DebugProgressData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		dbg.evtDebugProgress(msg.SN, msg.Payload.Method, data)

		err = sendEventReply(dbg.client, msg, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dgb *DebugModule) replyHandler(msg *MessageData) error {

	if fn, has := dgb.callbackStatusList[msg.Payload.Tid]; has {
		status := &ReplyStatusData{}
		err := json.Unmarshal(msg.Payload.Data, status)
		if err != nil {
			return err
		}
		fn(status.Status)
		delete(dgb.callbackStatusList, msg.Payload.Tid)
	}
	return nil
}

func newDebugModule(client mqtt.Client) *DebugModule {
	r := new(DebugModule)
	r.client = client
	r.callbackStatusList = make(map[string]FnServiceCallbackStatus)
	return r
}
