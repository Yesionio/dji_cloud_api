package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type FnSubDockOsdData func(sn string, data *DockOsdData)
type FnSubDockStateData func(sn string, data *DockStateData)
type FnSubDroneOsdData func(sn string, data *DroneOsdData)
type FnSubDroneStateData func(sn string, data *DroneStateData)
type FnSubDeviceTopoUpdate func(sn string, data *DeviceTopoData)

type DeviceManageModule struct {
	client                 mqtt.Client
	subDockOsdMap          map[string][]FnSubDockOsdData
	subDockStateMap        map[string][]FnSubDockStateData
	subDroneOsdMap         map[string][]FnSubDroneOsdData
	subDroneStateMap       map[string][]FnSubDroneStateData
	subDeviceTopoUpdateMap map[string][]FnSubDeviceTopoUpdate
}

func (dm *DeviceManageModule) SubDockOsdData(sn string, fn FnSubDockOsdData) {
	if x, has := dm.subDockOsdMap[sn]; has {
		x = append(x, fn)
		dm.subDockOsdMap[sn] = x
	} else {
		dm.subDockOsdMap[sn] = []FnSubDockOsdData{fn}
	}
}

func (dm *DeviceManageModule) SubDockStateData(sn string, fn FnSubDockStateData) {
	if x, has := dm.subDockStateMap[sn]; has {
		x = append(x, fn)
		dm.subDockStateMap[sn] = x
	} else {
		dm.subDockStateMap[sn] = []FnSubDockStateData{fn}
	}
}

func (dm *DeviceManageModule) SubDroneOsdData(sn string, fn FnSubDroneOsdData) {
	if x, has := dm.subDroneOsdMap[sn]; has {
		x = append(x, fn)
		dm.subDroneOsdMap[sn] = x
	} else {
		dm.subDroneOsdMap[sn] = []FnSubDroneOsdData{fn}
	}
}

func (dm *DeviceManageModule) SubDroneStateData(sn string, fn FnSubDroneStateData) {
	if x, has := dm.subDroneStateMap[sn]; has {
		x = append(x, fn)
		dm.subDroneStateMap[sn] = x
	} else {
		dm.subDroneStateMap[sn] = []FnSubDroneStateData{fn}
	}
}

func (dm *DeviceManageModule) SubDeviceTopoUpdate(sn string, fn FnSubDeviceTopoUpdate) {
	if x, has := dm.subDeviceTopoUpdateMap[sn]; has {
		x = append(x, fn)
		dm.subDeviceTopoUpdateMap[sn] = x
	} else {
		dm.subDeviceTopoUpdateMap[sn] = []FnSubDeviceTopoUpdate{fn}
	}
}

func (dm *DeviceManageModule) osdHandler(msg *MessageData) error {
	if r, has := dm.subDockOsdMap[msg.SN]; has {
		data := &DockOsdData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		for _, fn := range r {
			fn(msg.SN, data)
		}
	}

	if r, has := dm.subDroneOsdMap[msg.SN]; has {
		data := &DroneOsdData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		for _, fn := range r {
			fn(msg.SN, data)
		}
	}
	return nil
}

func (dm *DeviceManageModule) stateHandler(msg *MessageData) error {
	if r, has := dm.subDockStateMap[msg.SN]; has {
		data := &DockStateData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		for _, fn := range r {
			fn(msg.SN, data)
		}
	}

	if r, has := dm.subDroneStateMap[msg.SN]; has {
		data := &DroneStateData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		for _, fn := range r {
			fn(msg.SN, data)
		}
	}
	return nil
}

func (dm *DeviceManageModule) topoHandler(msg *MessageData) error {
	if r, has := dm.subDeviceTopoUpdateMap[msg.SN]; has {
		data := &DeviceTopoData{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}
		for _, fn := range r {
			fn(msg.SN, data)
		}
	}
	return nil
}

func newDeviceManageModule(client mqtt.Client) *DeviceManageModule {
	r := new(DeviceManageModule)
	r.client = client
	r.subDockOsdMap = make(map[string][]FnSubDockOsdData)
	r.subDockStateMap = make(map[string][]FnSubDockStateData)
	r.subDroneOsdMap = make(map[string][]FnSubDroneOsdData)
	r.subDroneStateMap = make(map[string][]FnSubDroneStateData)
	r.subDeviceTopoUpdateMap = make(map[string][]FnSubDeviceTopoUpdate)
	return r
}
