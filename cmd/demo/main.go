package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	djicloudapi "github.com/Yesionio/dji_cloud_api"
	"github.com/google/uuid"
)

func main() {
	core := djicloudapi.NewCloudApi(djicloudapi.CloudApiOption{
		MqttClientID:        "jxlkjkljlkfjklsdf",
		MqttBroker:          "6136d917d16c5343.natapp.cc:12362",
		MqttCleanSession:    true,
		MqttProtocolVersion: djicloudapi.MqttProtocolVersionV5,
		MqttEnableAuth:      false,
		ErrHandler: func(err error) {
			log.Println("raise a err", err)
		},
	})

	core.OrgManageModule.Enable(func(snList []string) []*djicloudapi.OrgBindStatus {
		log.Println("记录请求", snList)
		rls := make([]*djicloudapi.OrgBindStatus, 0)
		for _, v := range snList {
			s := djicloudapi.OrgBindStatus{
				SN:                       v,
				IsDeviceBindOrganization: true,
				OrganizationID:           uuid.NewString(),
				OrganizationName:         v[:6] + "组织",
				DeviceCallsign:           v[:6] + "组织",
			}
			rls = append(rls, &s)
		}

		return rls
	}, func(boqp *djicloudapi.BindOrgQueryParam) string {
		log.Println("获取绑定信息", boqp)
		return "组织"
	}, func(btop []*djicloudapi.BindToOrgParam) []*djicloudapi.BindToOrgResult {
		log.Println("获取绑定结果", btop)
		return nil
	})

	log.Println("开始运行")
	if err := core.Run(); err != nil {
		panic(err)
	}
	log.Println("连接成功")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println("bye")
}
