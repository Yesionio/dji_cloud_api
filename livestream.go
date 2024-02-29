package djicloudapi

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type LiveStreamModule struct {
	client       mqtt.Client
	callbackList map[string]FnServiceCallbackResult
}

type ParamLiveStartPush struct {
	UrlType      int    `json:"url_type"`
	Url          string `json:"url"`
	VideoID      string `json:"video_id"`
	VideoQuality int    `json:"video_quality"`
}

func (lsm *LiveStreamModule) LiveStartPushAsync(sn string, param ParamLiveStartPush, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	lsm.callbackList[tid] = callback

	err := sendService(lsm.client, sn, tid, "live_start_push", &param)
	if err != nil {
		return err
	}

	return nil
}

func (lsm *LiveStreamModule) LiveStopPushAsync(sn string, videoID string, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	lsm.callbackList[tid] = callback

	param := map[string]any{
		"video_id": videoID,
	}
	err := sendService(lsm.client, sn, tid, "live_stop_push", &param)
	if err != nil {
		return err
	}

	return nil
}

type ParamLiveSetQuality struct {
	VideoID      string
	VideoQuality int
}

func (lsm *LiveStreamModule) LiveSetQualityAsync(sn string, param ParamLiveSetQuality, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	lsm.callbackList[tid] = callback

	err := sendService(lsm.client, sn, tid, "live_set_quality", &param)
	if err != nil {
		return err
	}

	return nil
}

type ParamLiveLensChange struct {
	VideoID   string
	VideoType string
}

func (lsm *LiveStreamModule) LiveLensChangeAsync(sn string, param ParamLiveLensChange, callback FnServiceCallbackResult) error {
	tid := uuid.NewString()
	lsm.callbackList[tid] = callback

	err := sendService(lsm.client, sn, tid, "live_lens_change", &param)
	if err != nil {
		return err
	}

	return nil
}

func (lsm *LiveStreamModule) replyHandler(msg *MessageData) error {
	result := &ReplyResultData{}
	err := json.Unmarshal(msg.Payload.Data, result)
	if err != nil {
		return err
	}

	if fn, has := lsm.callbackList[msg.Payload.Tid]; has {
		fn(result.Result)
		delete(lsm.callbackList, msg.Payload.Tid)
	}
	return nil
}

func newLiveStreamModule(client mqtt.Client) *LiveStreamModule {
	r := new(LiveStreamModule)
	r.client = client
	r.callbackList = make(map[string]FnServiceCallbackResult)
	return r
}
