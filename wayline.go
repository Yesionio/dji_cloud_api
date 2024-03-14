package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type FnEvtDeviceExitHomingNotify func(sn string, data *WaylineExitHomingNotify)
type FnEvtFlighttaskProgress func(sn string, data *WaylineFlightProgress)
type FnEvtFlightTaskReady func(sn string, data *WaylineFlightTaskReady)
type FnEvtReturnHomeInfo func(sn string, data *WaylineReturnHomeInfo)
type FnFlightTaskResourceGet func(flightID string) (*WaylineResourceGetData, error)

type WaylineModule struct {
	client                    mqtt.Client
	evtDeviceExitHomingNotify FnEvtDeviceExitHomingNotify
	evtFlighttaskProgress     FnEvtFlighttaskProgress
	evtFlightTaskReady        FnEvtFlightTaskReady
	evtReturnHomeInfo         FnEvtReturnHomeInfo
	callbackList              map[string]FnServiceCallbackResult
	callbackStatusList        map[string]FnServiceCallbackStatus
	fnFlightTaskResourceGet   FnFlightTaskResourceGet
}

// SetOnEvtExitHomingNotify 当发出退出返航时调用
func (wm *WaylineModule) SetOnEvtExitHomingNotify(fn FnEvtDeviceExitHomingNotify) {
	wm.evtDeviceExitHomingNotify = fn
}

// SetOnEvtExitHomingNotify 当发出退出返航时调用
func (wm *WaylineModule) SetOnEvtFlighttaskProgress(fn FnEvtFlighttaskProgress) {
	wm.evtFlighttaskProgress = fn
}

// SetOnEvtFlightTaskReady 当有任务就绪时调用
func (wm *WaylineModule) SetOnEvtFlightTaskReady(fn FnEvtFlightTaskReady) {
	wm.evtFlightTaskReady = fn
}

// SetOnEvtFlightTaskReady 当有任务就绪时调用
func (wm *WaylineModule) SetOnEvtReturnHomeInfo(fn FnEvtReturnHomeInfo) {
	wm.evtReturnHomeInfo = fn
}

func (wm *WaylineModule) SetOnResourceGet(fn FnFlightTaskResourceGet) {
	wm.fnFlightTaskResourceGet = fn
}

// FlighttaskPrepare 下发航线任务
func (wm *WaylineModule) FlighttaskPrepareAsync(sn string, param ParamFlightTaskPrepare, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	err := sendService(wm.client, sn, tid, "flighttask_prepare", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskExecute 执行航线任务
func (wm *WaylineModule) FlighttaskExecuteAsync(sn string, flightID string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	param := map[string]any{
		"flight_id": flightID,
	}

	err := sendService(wm.client, sn, tid, "flighttask_execute", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskUndo 取消航线任务
func (wm *WaylineModule) FlighttaskUndoAsync(sn string, flightIDs []string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	param := map[string]any{
		"flight_ids": flightIDs,
	}

	err := sendService(wm.client, sn, tid, "flighttask_undo", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskPause 暂停航线任务
func (wm *WaylineModule) FlighttaskPauseAsync(sn string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	param := map[string]any{}

	err := sendService(wm.client, sn, tid, "flighttask_pause", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskRecovery 恢复航线任务
func (wm *WaylineModule) FlighttaskRecoveryAsync(sn string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	param := map[string]any{}

	err := sendService(wm.client, sn, tid, "flighttask_recovery", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskReturnHome 一键返航
func (wm *WaylineModule) FlighttaskReturnHomeAsync(sn string, callback FnServiceCallbackStatus) error {
	tid := uuid.NewString()
	wm.callbackStatusList[tid] = callback

	param := map[string]any{}

	err := sendService(wm.client, sn, tid, "return_home", &param)
	if err != nil {
		return err
	}
	return nil
}

// FlighttaskReturnHomeCancel 取消返航
func (wm *WaylineModule) FlighttaskReturnHomeCancelAsync(sn string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	wm.callbackList[tid] = callback

	param := map[string]any{}

	err := sendService(wm.client, sn, tid, "return_home_cancel", &param)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WaylineModule) deviceExitHomingNotifyHandler(msg *MessageData) error {
	if wm.evtDeviceExitHomingNotify != nil {
		data := &WaylineExitHomingNotify{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		wm.evtDeviceExitHomingNotify(msg.SN, data)

		err = sendEventReply(wm.client, msg, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wm *WaylineModule) flighttaskProgressHandler(msg *MessageData) error {
	if wm.evtDeviceExitHomingNotify != nil {
		data := &WaylineFlightProgress{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		wm.evtFlighttaskProgress(msg.SN, data)

		err = sendEventReply(wm.client, msg, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wm *WaylineModule) flightTaskReadyHandler(msg *MessageData) error {
	if wm.evtDeviceExitHomingNotify != nil {
		data := &WaylineFlightProgress{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		wm.evtFlighttaskProgress(msg.SN, data)

		err = sendEventReply(wm.client, msg, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wm *WaylineModule) returnHomeInfoHandler(msg *MessageData) error {
	if wm.evtDeviceExitHomingNotify != nil {
		data := &WaylineReturnHomeInfo{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		wm.evtReturnHomeInfo(msg.SN, data)

		err = sendEventReply(wm.client, msg, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

type paramResourceGet struct {
	FlightID string `json:"flight_id"`
}

func (wm *WaylineModule) resourceGetHandler(msg *MessageData) error {
	if wm.fnFlightTaskResourceGet != nil {
		data := &paramResourceGet{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		rls, err := wm.fnFlightTaskResourceGet(data.FlightID)
		if err != nil {
			sendOutput(wm.client, msg, rls, 1)
		} else {
			sendOutput(wm.client, msg, rls, 0)
		}
	}
	return nil
}

func (wm *WaylineModule) replyHandler(msg *MessageData) error {

	if fn, has := wm.callbackList[msg.Payload.Tid]; has {
		result := &ReplyResultData{}
		err := json.Unmarshal(msg.Payload.Data, result)
		if err != nil {
			return err
		}
		fn(result.Result)
		delete(wm.callbackList, msg.Payload.Tid)
	}

	if fn, has := wm.callbackStatusList[msg.Payload.Tid]; has {
		status := &ReplyStatusData{}
		err := json.Unmarshal(msg.Payload.Data, status)
		if err != nil {
			return err
		}
		fn(status.Output.Status)
		delete(wm.callbackStatusList, msg.Payload.Tid)
	}
	return nil
}

func newWaylineModule(client mqtt.Client) *WaylineModule {
	r := new(WaylineModule)
	r.client = client
	r.callbackList = make(map[string]FnServiceCallbackResult)
	r.callbackStatusList = make(map[string]FnServiceCallbackStatus)
	return r
}
