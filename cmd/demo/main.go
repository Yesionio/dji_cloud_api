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
	// startServer()
	_ = runDjiCloud()
	log.Println("连接成功")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println("bye")
}

// func startServer() error {
// 	// app := fiber.New()
// 	// app.Get("/index.html", func(c *fiber.Ctx) error {
// 	// 	return c.SendFile("index.html")
// 	// })
// 	djiSN := "7CTDM2100BH29T"

// 	core := runDjiCloud()

// 	cbChan := make(chan string, 1024)
// 	osdChan := make(chan string, 1024)

// 	// app.Get("/api/debug/:cmd", func(c *fiber.Ctx) error {
// 	// 	cmd := c.Params("cmd")

// 	// 	cs := uuid.NewString()
// 	// 	switch cmd {
// 	// 	case "debug_open":

// 	// 		core.DebugModule.DebugModeOpenAsync(djiSN, func(status string) {
// 	// 			cbChan <- status
// 	// 		})
// 	// 	}

// 	// 	return c.SendString(cs)
// 	// })

// 	// app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
// 	// 	fmt.Println("ws连接进入")
// 	// 	for {
// 	// 		select {
// 	// 		case val := <-cbChan:
// 	// 			s := strings.Split(val, ",")
// 	// 			c.WriteJSON(&map[string]any{
// 	// 				"status": s[1],
// 	// 				"uid":    s[0],
// 	// 			})
// 	// 		case val := <-osdChan:
// 	// 			c.WriteMessage(websocket.TextMessage, []byte(val))
// 	// 		}
// 	// 	}
// 	// }))

// 	return nil
// }

func runDjiCloud() *djicloudapi.DjiCloudApiCore {
	// drone sn = 1581F6Q8D242U00C613B
	core := djicloudapi.NewCloudApi(djicloudapi.CloudApiOption{
		MqttClientID:        "jxlkjkljlkfjklsdfx",
		MqttBroker:          "172.10.40.63:31883",
		MqttCleanSession:    true,
		MqttProtocolVersion: djicloudapi.MqttProtocolVersionV5,
		MqttEnableAuth:      false,
		ErrHandler: func(err error) {
			log.Println("raise a err", err)
		},
		NtpServerHost: "ntp.aliyun.com",
		AppID:         "143886",
		AppKey:        "357d609f3faa01d612ab64806383363",
		AppLicense:    "Hkks6Dvmyb/te+dXnhyIXf8X9yu/FZRYCFqKK6dy4SP+JsVG+WJziutQH7wDiKSyQgicDRwOUVIgoeHkHOMwz9uSA3wfQrs/KGizN3J2dMPsscPX7khn7cngcl1LHdzaZ6dLdHYxMjp0Kq0vxjnJcp8w+mdpl/emqgBuC09q7ks=",
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

	core.DeviceManageModule.SubDockOsdData("7CTDM2100BH29T", func(sn string, data *djicloudapi.DockOsdData) {
		fmt.Println("osd data is")
		fmt.Println(data)
	})

	core.DeviceManageModule.SubDroneOsdData("1581F6Q8D242U00C613B", func(sn string, data *djicloudapi.DroneOsdData) {
		fmt.Println("无人机的数据有", data)
	})

	log.Println("开始运行")
	if err := core.Run(); err != nil {
		panic(err)
	}

	return core
}
