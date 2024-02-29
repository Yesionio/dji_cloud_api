package djicloudapi

type FileUploadData struct {
	File FileUploadInfoData `json:"file"`
}

type FileUploadInfoData struct {
	ObjectKey string                `json:"object_key"`
	Path      string                `json:"path"`
	Name      string                `json:"name"`
	Ext       FileUploadExtInfoData `json:"ext"`
	Metadata  FileUploadMetadata    `json:"metadata"`
}

type FileUploadExtInfoData struct {
	FlightID        string `json:"flight_id"`
	DroneModelKey   string `json:"drone_model_key"`
	PayloadModelKey string `json:"payload_model_key"`
	IsOriginal      bool   `json:"is_original"`
}

type FileUploadMetadata struct {
	GimbalYawDagree  float64                 `json:"gimbal_yaw_dagree"`
	AbsoluteAltitude float64                 `json:"absolute_altitude"`
	RelativeAltitude float64                 `json:"relative_altitude"`
	CreateTime       string                  `json:"create_time"`
	ShootPosition    FileUploadShootPosition `json:"shoot_position"`
}

type FileUploadShootPosition struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type StorageConfigGetData struct {
	Bucket          string                      `json:"bucket"`
	Credentials     StorageConfigGetCredentials `json:"credentials"`
	Endpoint        string                      `json:"endpoint"`
	Provider        string                      `json:"provider"`
	Region          string                      `json:"region"`
	ObjectKeyPrefix string                      `json:"object_key_prefix"`
}

type StorageConfigGetCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"security_token"`
	Expire          int    `json:"expire"`
}
