package djicloudapi

type WaylineExitHomingNotify struct {
	SN     string `json:"sn"`
	Action int    `json:"action"`
	Reason int    `json:"reason"`
}

type WaylineFlightProgress struct {
	Ext      WaylineFlightProgressExt  `json:"ext"`
	Status   string                    `json:"status"`
	Progress WaylineFlightProgressInfo `json:"progress"`
}

type WaylineFlightProgressExt struct {
	CurrentWaypointIndex int                             `json:"current_waypoint_index"`
	WaylineMissionState  int                             `json:"wayline_mission_state"`
	MediaCount           int                             `json:"media_count"`
	TrackID              string                          `json:"track_id"`
	FlightID             string                          `json:"flight_id"`
	BreakPoint           WaylineFlightProgressBreakPoint `json:"break_point"`
	WaylineID            int                             `json:"wayline_id"`
}

type WaylineFlightProgressBreakPoint struct {
	Index        int     `json:"index"`
	State        int     `json:"state"`
	Progress     float64 `json:"progress"`
	WaylineID    int     `json:"wayline_id"`
	BreakReason  int     `json:"break_reason"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Height       float64 `json:"height"`
	AttitudeHead float64 `json:"attitude_head"`
}

type WaylineFlightProgressInfo struct {
	CurrentStep int `json:"current_step"`
	Percent     int `json:"percent"`
}

type WaylineFlightTaskReady struct {
	FlightIDs []string `json:"flight_ids"`
}

type WaylineReturnHomeInfo struct {
	PlannedPathPoints WaylinePlannedPathPoint `json:"planned_path_points"`
	LastPointType     int                     `json:"last_point_type"`
	FlightID          string                  `json:"flight_id"`
}

type WaylinePlannedPathPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"`
}

type ParamFlightTaskPrepare struct {
	FlightID              string                         `json:"flight_id"`
	ExecuteTime           int64                          `json:"execute_time"`
	TaskType              int                            `json:"task_type"`
	File                  ParamFlightTaskPrepareFile     `json:"file"`
	ReadyConditions       ParamFlightTaskReadyCond       `json:"ready_conditions"`
	ExecutableConditions  ParamFlightTaskExecCond        `json:"executable_conditions"`
	BreakPoint            ParamFlightTaskBreakPoint      `json:"break_point"`
	RthAltitude           int                            `json:"rth_altitude"`
	RthMode               int                            `json:"rth_mode"`
	OutOfControlAction    int                            `json:"out_of_control_action"`
	ExitWaylineWhenRcLost int                            `json:"exit_wayline_when_rc_lost"`
	WaylinePrecisionType  int                            `json:"wayline_precision_type"`
	SimulateMission       ParamFlightTaskSimulateMission `json:"simulate_mission"`
}

type ParamFlightTaskPrepareFile struct {
	URL         string `json:"url"`
	Fingerprint string `json:"fingerprint"`
}

type ParamFlightTaskReadyCond struct {
	BatteryCapacity int `json:"battery_capacity"`
	BeginTime       int `json:"begin_time"`
	EndTime         int `json:"end_time"`
}

type ParamFlightTaskExecCond struct {
	StorageCapacity int `json:"storage_capacity"`
}

type ParamFlightTaskBreakPoint struct {
	Index     int     `json:"index"`
	State     int     `json:"state"`
	Progress  float64 `json:"progress"`
	WaylineID int     `json:"wayline_id"`
}

type ParamFlightTaskSimulateMission struct {
	IsEnable  int     `json:"is_enable"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WaylineResourceGetData struct {
	File WaylineResourceGetFile `json:"file"`
}

type WaylineResourceGetFile struct {
	URL         string `json:"url"`
	Fingerprint string `json:"fingerprint"`
}
