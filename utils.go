package djicloudapi

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func sendOutput(client mqtt.Client, msg *MessageData, data any, result int) error {
	pubData := CommonData[CommonOutPutData]{
		Tid:       msg.Payload.Tid,
		Bid:       msg.Payload.Bid,
		Timestamp: time.Now().Unix(),
		Gateway:   msg.Payload.Gateway,
		Method:    msg.Payload.Method,
		Data: CommonOutPutData{
			Output: data,
			Result: result,
		},
	}

	marData, err := json.Marshal(pubData)
	if err != nil {
		return err
	}

	t := client.Publish(msg.Topic+REPLY_SUFFIX, 0, false, marData)
	if ok := t.WaitTimeout(time.Second * 10); !ok || t.Error() != nil {
		return t.Error()
	}

	return nil
}

func sendEventReply(client mqtt.Client, msg *MessageData, result int) error {
	pubData := CommonData[CommonOutPutData]{
		Tid:       msg.Payload.Tid,
		Bid:       msg.Payload.Bid,
		Timestamp: time.Now().Unix(),
		Gateway:   msg.Payload.Gateway,
		Method:    msg.Payload.Method,
		Data: CommonOutPutData{
			Result: result,
		},
	}

	marData, err := json.Marshal(pubData)
	if err != nil {
		return err
	}

	t := client.Publish(msg.Topic+REPLY_SUFFIX, 0, false, marData)
	if ok := t.WaitTimeout(time.Second * 10); !ok || t.Error() != nil {
		return t.Error()
	}
	return nil
}

func sendService(client mqtt.Client, sn, tid, method string, data any) error {
	r, err := json.Marshal(data)
	if err != nil {
		return err
	}
	pubData := CommonData[json.RawMessage]{
		Tid:       tid,
		Bid:       uuid.NewString(),
		Timestamp: time.Now().Unix(),
		Method:    method,
		Data:      r,
	}

	marData, err := json.Marshal(pubData)
	if err != nil {
		return err
	}

	t := client.Publish("thing/product/"+sn+"services", 0, false, marData)
	if ok := t.WaitTimeout(time.Second * 10); !ok || t.Error() != nil {
		return t.Error()
	}

	return nil
}
