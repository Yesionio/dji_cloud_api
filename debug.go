package djicloudapi

import mqtt "github.com/eclipse/paho.mqtt.golang"

type DebugModule struct {
	client mqtt.Client
}

func newDebugModule(client mqtt.Client) *DebugModule {
	r := new(DebugModule)
	r.client = client
	return r
}
