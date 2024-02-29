package djicloudapi

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type OrgBindStatus struct {
	SN                       string `json:"sn"`
	IsDeviceBindOrganization bool   `json:"is_device_bind_organization"`
	OrganizationID           string `json:"organization_id"`
	OrganizationName         string `json:"organization_name"`
	DeviceCallsign           string `json:"device_callsign"`
}

type BindOrgQueryParam struct {
	DeviceBindingCode string `json:"device_binding_code"`
	OrganizationID    string `json:"organization_id"`
}

type BindToOrgParam struct {
	DeviceBindingCode string `json:"device_binding_code"`
	DeviceCallsign    string `json:"device_callsign"`
	DeviceModelKey    string `json:"device_model_key"`
	OrganizationID    string `json:"organization_id"`
	SN                string `json:"sn"`
}

type BindToOrgResult struct {
	SN      string `json:"sn"`
	ErrCode int    `json:"err_code"`
}

type FnGetBindOrg func([]string) []*OrgBindStatus
type FnQueryBindOrgName func(*BindOrgQueryParam) string
type FnBindToOrg func([]*BindToOrgParam) []*BindToOrgResult

type OrgManageModule struct {
	client               mqtt.Client
	isEnable             bool
	funcGetBindOrg       FnGetBindOrg
	funcQueryBindOrgName FnQueryBindOrgName
	funcBindToOrg        FnBindToOrg
}

func (omm *OrgManageModule) Enable(getBindOrg FnGetBindOrg, queryBindOrgName FnQueryBindOrgName, bindToOrg FnBindToOrg) {
	omm.isEnable = true
	omm.funcGetBindOrg = getBindOrg
	omm.funcQueryBindOrgName = queryBindOrgName
	omm.funcBindToOrg = bindToOrg
}

type ParamBindStatus struct {
	Devices []*ParamBindStatusData `json:"devices"`
}

type ParamBindStatusData struct {
	SN string `json:"sn"`
}

type RespBindStatus struct {
	BindStatus []*OrgBindStatus `json:"bind_status"`
}

func (omm *OrgManageModule) airportBindStatus(msg *MessageData) error {
	if omm.isEnable {
		data := &ParamBindStatus{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		p := make([]string, 0)
		for _, pbsd := range data.Devices {
			p = append(p, pbsd.SN)
		}
		rls := omm.funcGetBindOrg(p)
		resp := &RespBindStatus{
			BindStatus: rls,
		}

		err = sendOutput(omm.client, msg, resp, 0)
		if err != nil {
			return err
		}

	}
	return nil
}

type RespOrganizationGet struct {
	OrganizationName string `json:"organization_name"`
}

func (omm *OrgManageModule) airportOrganizationGet(msg *MessageData) error {
	if omm.isEnable {
		data := &BindOrgQueryParam{}
		log.Println("")
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		rls := omm.funcQueryBindOrgName(data)
		resp := &RespOrganizationGet{
			OrganizationName: rls,
		}

		err = sendOutput(omm.client, msg, resp, 0)
		if err != nil {
			return err
		}

	}
	return nil
}

type paramOranizationBind struct {
	BindDevices []*BindToOrgParam `json:"bind_devices"`
}

type RespOranizationBind struct {
	ErrInfos []*BindToOrgResult `json:"err_infos"`
}

func (omm *OrgManageModule) airportOrganizationBind(msg *MessageData) error {
	if omm.isEnable {
		data := &paramOranizationBind{}
		err := json.Unmarshal(msg.Payload.Data, data)
		if err != nil {
			return err
		}

		rls := omm.funcBindToOrg(data.BindDevices)
		resp := &RespOranizationBind{
			ErrInfos: rls,
		}

		err = sendOutput(omm.client, msg, resp, 0)
		if err != nil {
			return err
		}

	}
	return nil
}

func newOrgManageModule(client mqtt.Client) *OrgManageModule {
	r := new(OrgManageModule)
	r.client = client
	return r
}
