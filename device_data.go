package djicloudapi

// DockOsdData 机场定频数据
type DockOsdData struct {
	Longitude                   float64                                `json:"longitude"`
	Latitude                    float64                                `json:"latitude"`
	ModeCode                    int                                    `json:"mode_code"`
	FlighttaskStepCode          int                                    `json:"flighttask_step_code"`
	SubDevice                   DockOsdSubDeviceData                   `json:"sub_device"`
	CoverState                  int                                    `json:"cover_state"`
	PutterState                 int                                    `json:"putter_state"`
	SupplementLightState        int                                    `json:"supplement_light_state"`
	NetworkState                DockOsdNetworkStateData                `json:"network_state"`
	DroneInDock                 int                                    `json:"drone_in_dock"`
	JobNumber                   int                                    `json:"job_number"`
	MediaFileDetail             DockOsdMediaFileDetailData             `json:"media_file_detail"`
	WirelessLink                DockOsdWirelessLinkData                `json:"wireless_link"`
	Rainfall                    int                                    `json:"rainfall"`
	WindSpeed                   float64                                `json:"wind_speed"`
	EnvironmentTemperature      float64                                `json:"environment_temperature"`
	Temperature                 float64                                `json:"temperature"`
	Humidity                    float64                                `json:"humidity"`
	ElectricSupplyVoltage       int                                    `json:"electric_supply_voltage"`
	WorkingVoltage              int                                    `json:"working_voltage"`
	WorkingCurrent              float64                                `json:"working_current"`
	Storage                     DockOsdStorageData                     `json:"storage"`
	FirstPowerOn                int64                                  `json:"first_power_on"`
	AlternateLandPoint          DockOsdAlternateLandPointData          `json:"alternate_land_point"`
	Height                      float64                                `json:"height"`
	ActivationTime              int                                    `json:"activation_time"`
	AirConditioner              DockOsdAirConditionerData              `json:"air_conditioner"`
	BatteryStoreMode            int                                    `json:"battery_store_mode"`
	AlarmState                  int                                    `json:"alarm_state"`
	DroneBatteryMaintenanceInfo DockOsdDroneBatteryMaintenanceInfoData `json:"drone_battery_maintenance_info"`
	BackupBattery               DockOsdBackupBatteryData               `json:"backup_battery"`
	DroneChargeState            DockOsdDroneChargeStateData            `json:"drone_charge_state"`
	EmergencyStopState          int                                    `json:"emergency_stop_state"`
	PositionState               DockOsdPositionStateData               `json:"position_state"`
	MaintainStatus              DockOsdMaintainStatusData              `json:"maintain_status"`
	DrcState                    int                                    `json:"drc_state"`
}

type DockOsdSubDeviceData struct {
	DeviceSN           string `json:"device_sn"`
	ProductType        string `json:"product_type"`
	DeviceOnlineStatus int    `json:"device_online_status"`
	DevicePaired       int    `json:"device_paired"`
}

type DockOsdNetworkStateData struct {
	Type    int     `json:"type"`
	Quality int     `json:"quality"`
	Rate    float64 `json:"rate"`
}

type DockOsdMediaFileDetailData struct {
	RemainUpload int `json:"remain_upload"`
}

type DockOsdWirelessLinkData struct {
	DongleNumber  int     `json:"dongle_number"`
	X4gLinkState  int     `json:"4g_link_state"`
	SdrLinkState  int     `json:"sdr_link_state"`
	LinkWorkmode  int     `json:"link_workmode"`
	SdrQuality    int     `json:"sdr_quality"`
	X4gQuality    int     `json:"4g_quality"`
	X4gUavQuality int     `json:"4g_uav_quality"`
	X4gGndQuality int     `json:"4g_gnd_quality"`
	SdrFreqBand   float64 `json:"sdr_freq_band"`
	X4gFreqBand   float64 `json:"4g_freq_band"`
}
type DockOsdStorageData struct {
	Total int `json:"total"`
	Used  int `json:"used"`
}

type DockOsdAlternateLandPointData struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	SafeLandHeight float64 `json:"safe_land_height"`
	IsConfigured   int     `json:"is_configured"`
}

type DockOsdAirConditionerData struct {
	AirConditionerState int `json:"air_conditioner_state"`
	SwitchTime          int `json:"switch_time"`
}

type DockOsdDroneBatteryMaintenanceInfoData struct {
	MaintenanceState    int                                               `json:"maintenance_state"`
	MaintenanceTimeLeft int                                               `json:"maintenance_time_left"`
	HeatState           int                                               `json:"heat_state"`
	Batteries           []DockOsdDroneBatteryMaintenanceInfoBatteryesData `json:"batteries"`
}

type DockOsdDroneBatteryMaintenanceInfoBatteryesData struct {
	Index           int     `json:"index"`
	CapacityPercent int     `json:"capacity_percent"`
	Voltage         int     `json:"voltage"`
	Temperature     float64 `json:"temperature"`
}

type DockOsdBackupBatteryData struct {
	Voltage     int     `json:"voltage"`
	Temperature float64 `json:"temperature"`
	Switch      int     `json:"switch"`
}

type DockOsdDroneChargeStateData struct {
	State           int `json:"state"`
	CapacityPercent int `json:"capacity_percent"`
}

type DockOsdPositionStateData struct {
	IsCalibration int `json:"is_calibration"`
	IsFixed       int `json:"is_fixed"`
	Quality       int `json:"quality"`
	GpsNumber     int `json:"gps_number"`
	RtkNumber     int `json:"rtk_number"`
}

type DockOsdMaintainStatusData struct {
	MaintainStatusArray []DockOsdMaintainStatusArrayData `json:"maintain_status_array"`
}

type DockOsdMaintainStatusArrayData struct {
	State                   int   `json:"state"`
	LastMaintainType        int   `json:"last_maintain_type"`
	LastMaintainTime        int64 `json:"last_maintain_time"`
	LastMaintainWorkSorties int   `json:"last_maintain_work_sorties"`
}

// DockStateData 机场状态变化数据
type DockStateData struct {
	FirmwareVersion           string                      `json:"firmware_version"`
	FirmwareUpgradeStatus     int                         `json:"firmware_upgrade_status"`
	LiveStatus                []DockStateLiveStatusData   `json:"live_status"`
	LiveCapacity              []DockStateLiveCapacityData `json:"live_capacity"`
	AccTime                   int64                       `json:"acc_time"`
	CompatibleStatus          int                         `json:"compatible_status"`
	WpmzVersion               string                      `json:"wpmz_version"`
	UserExperienceImprovement int                         `json:"user_experience_improvement"`
}

type DockStateLiveStatusData struct {
	VideoID      string `json:"video_id"`
	VideoType    string `json:"video_type"`
	VideoQuality int    `json:"video_quality"`
	Status       int    `json:"status"`
	ErrorStatus  int    `json:"error_status"`
}

type DockStateLiveCapacityData struct {
	AvailableVideoNumber  int                               `json:"available_video_number"`
	CoexistVideoNumberMax int                               `json:"coexist_video_number_max"`
	DeviceList            []DockStateLiveCapacityDeviceData `json:"device_list"`
}

type DockStateLiveCapacityDeviceData struct {
	SN                    string                            `json:"sn"`
	AvailableVideoNumber  int                               `json:"available_video_number"`
	CoexistVideoNumberMax int                               `json:"coexist_video_number_max"`
	CameraList            []DockStateLiveCapacityCameraData `json:"camera_list"`
}

type DockStateLiveCapacityCameraData struct {
	AvailableVideoNumber  int                              `json:"availableVideoNumber"`
	CoexistVideoNumberMax int                              `json:"coexistVideoNumberMax"`
	CameraIndex           string                           `json:"cameraIndex"`
	VideoList             []DockStateLiveCapacityVideoData `json:"videoList"`
}

type DockStateLiveCapacityVideoData struct {
	VideoIndex string `json:"videoIndex"`
	VideoType  string `json:"videoType"`
}

// DroneOsdData 无人机定频数据
type DroneOsdData struct {
	ModeCode            int                           `json:"mode_code"`
	Cameras             []DroneOsdCameraData          `json:"cameras"`
	Gear                int                           `json:"gear"`
	HorizontalSpeed     float64                       `json:"horizontal_speed"`
	VerticalSpeed       float64                       `json:"vertical_speed"`
	Longitude           float64                       `json:"longitude"`
	Latitude            float64                       `json:"latitude"`
	Height              float64                       `json:"height"`
	Elevation           float64                       `json:"elevation"`
	AttitudePitch       float64                       `json:"attitude_pitch"`
	AttitudeRoll        float64                       `json:"attitude_roll"`
	AttitudeHead        int                           `json:"attitude_head"`
	HomeDistance        float64                       `json:"home_distance"`
	WindSpeed           float64                       `json:"wind_speed"`
	WindDirection       int                           `json:"wind_direction"`
	TotalFlightTime     int                           `json:"total_flight_time"`
	TotalFlightDistance int                           `json:"total_flight_distance"`
	Battery             DroneOsdBatteryData           `json:"battery"`
	Storage             DroneOsdStorageData           `json:"storage"`
	PositionState       DroneOsdPositionData          `json:"position_state"`
	TrackId             string                        `json:"track_id"`
	TotalFlightSorties  int                           `json:"total_flight_sorties"`
	MaintainStatus      DroneOsdMaintainData          `json:"maintain_status"`
	ActivationTime      int                           `json:"activation_time"`
	NightLightsState    int                           `json:"night_lights_state"`
	HeightLimit         int                           `json:"height_limit"`
	IsNearHeightLimit   int                           `json:"is_near_height_limit"`
	IsNearAreaLimit     int                           `json:"is_near_area_limit"`
	ObstacleAvoidance   DroneOsdObstacleAvoidanceData `json:"obstacle_avoidance"`
}

type DroneOsdCameraData struct {
	RemainPhotoNum       int                        `json:"remain_photo_num"`
	RemainRecordDuration int                        `json:"remain_record_duration"`
	RecordTime           int                        `json:"record_time"`
	PayloadIndex         string                     `json:"payload_index"`
	CameraMode           int                        `json:"camera_mode"`
	PhotoState           int                        `json:"photo_state"`
	ScreenSplitEnable    bool                       `json:"screen_split_enable"`
	RecordingState       int                        `json:"recording_state"`
	ZoomFactor           float64                    `json:"zoom_factor"`
	IrZoomFactor         float64                    `json:"ir_zoom_factor"`
	LiveviewWorldRegion  DroneOsdCameraLiveViewData `json:"liveview_world_region"`
	PhotoStorageSettings []string                   `json:"photo_storage_settings"`
	VideoStorageSettings []string                   `json:"video_storage_settings"`
}
type DroneOsdCameraLiveViewData struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}

type DroneOsdBatteryData struct {
	CapacityPercent  int                       `json:"capacity_percent"`
	LandingPower     int                       `json:"landing_power"`
	RemainFlightTime int                       `json:"remain_flight_time"`
	ReturnHomePower  int                       `json:"return_home_power"`
	Batteries        []DroneOsdBatteryInfoData `mapstructure:"batteries"`
}

type DroneOsdBatteryInfoData struct {
	CapacityPercent        int     `json:"capacity_percent"`
	Index                  int     `json:"index"`
	Sn                     string  `json:"sn"`
	Type                   int     `json:"type"`
	SubType                int     `json:"sub_type"`
	FirmwareVersion        string  `json:"firmware_version"`
	LoopTimes              int     `json:"loop_times"`
	Voltage                int     `json:"voltage"`
	Temperature            float64 `json:"temperature"`
	HighVoltageStorageDays int     `json:"high_voltage_storage_days"`
}

type DroneOsdMaintainData struct {
	MaintainStatusArray []DroneOsdMaintainArray `json:"maintain_status_array"`
}

type DroneOsdMaintainArray struct {
	LastMaintainFlightTime    int `json:"last_maintain_flight_time"`
	LastMaintainTime          int `json:"last_maintain_time"`
	LastMaintainType          int `json:"last_maintain_type"`
	State                     int `json:"state"`
	LastMaintainFlightSorties int `json:"last_maintain_flight_sorties"`
}

type DroneOsdObstacleAvoidanceData struct {
	Downside int `json:"downside"`
	Horizon  int `json:"horizon"`
	Upside   int `json:"upside"`
}

type DroneOsdPositionData struct {
	IsFixed   int `json:"is_fixed"`
	Quality   int `json:"quality"`
	RtkNumber int `json:"rtk_number"`
	GpsNumber int `json:"gps_number"`
}

type DroneOsdStorageData struct {
	Total int `json:"total"`
	Used  int `json:"used"`
}

// DroneStateData 无人机状态变化数据
type DroneStateData struct {
	CameraWatermarkSettings           DroneStateWatermarkData `json:"camera_watermark_settings"`
	CommanderModeLostAction           int                     `json:"commander_mode_lost_action"`
	CommanderFlightMode               int                     `json:"commander_flight_mode"`
	CurrentCommanderFlightMode        int                     `json:"current_commander_flight_mode"`
	CommanderFlightHeight             float64                 `json:"commander_flight_height"`
	ModeCodeReason                    int                     `json:"mode_code_reason"`
	FirmwareVersion                   string                  `json:"firmware_version"`
	CompatibleStatus                  int                     `json:"compatible_status"`
	FirmwareUpgradeStatus             int                     `json:"firmware_upgrade_status"`
	HomeLongitude                     float64                 `json:"home_longitude"`
	HomeLatitude                      float64                 `json:"home_latitude"`
	ControlSource                     string                  `json:"control_source"`
	LowBatteryWarningThreshold        int                     `json:"low_battery_warning_threshold"`
	SeriousLowBatteryWarningThreshold int                     `json:"serious_low_battery_warning_threshold"`
	RthMode                           int                     `json:"rth_mode"`
	CurrentRthMode                    int                     `json:"current_rth_mode"`
	DongleInfos                       []DroneStateDongleData  `json:"dongle_infos"`
	OfflineMapEnable                  bool                    `json:"offline_map_enable"`
}

type DroneStateWatermarkData struct {
	GlobalEnable           int    `json:"global_enable"`
	DroneTypeEnable        int    `json:"drone_type_enable"`
	DroneSnEnable          int    `json:"drone_sn_enable"`
	DatetimeEnable         int    `json:"datetime_enable"`
	GpsEnable              int    `json:"gps_enable"`
	UserCustomStringEnable int    `json:"user_custom_string_enable"`
	UserCustomString       string `json:"user_custom_string"`
	Layout                 int    `json:"layout"`
}

type DroneStateDongleData struct {
	Imei              string                `json:"imei"`
	DongleType        int                   `json:"dongle_type"`
	Eid               string                `json:"eid"`
	EsimActivateState int                   `json:"esim_activate_state"`
	SimCardState      int                   `json:"sim_card_state"`
	SimSlot           int                   `json:"sim_slot"`
	EsimInfos         []DroneStateEsimData  `json:"esim_infos"`
	SimInfo           DroneStateSimInfoData `json:"sim_info"`
}

type DroneStateEsimData struct {
	TelecomOperator int    `json:"telecom_operator"`
	Enabled         bool   `json:"enabled"`
	Iccid           string `json:"iccid"`
}

type DroneStateSimInfoData struct {
	TelecomOperator int    `json:"telecom_operator"`
	SimType         int    `json:"sim_type"`
	Iccid           string `json:"iccid"`
}

type DeviceTopoData struct {
	Domain       string                    `json:"domain"`
	Type         int                       `json:"type"`
	SubType      int                       `json:"sub_type"`
	DeviceSecret string                    `json:"device_secret"`
	Nonce        string                    `json:"nonce"`
	ThingVersion string                    `json:"thing_version"`
	SubDevices   []DeviceTopoSubDeviceData `json:"sub_devices"`
}

type DeviceTopoSubDeviceData struct {
	SN           string `json:"sn"`
	Domain       string `json:"domain"`
	Type         int    `json:"type"`
	SubType      int    `json:"sub_type"`
	Index        string `json:"index"`
	DeviceSecret string `json:"device_secret"`
	Nonce        string `json:"nonce"`
	ThingVersion string `json:"thing_version"`
}
