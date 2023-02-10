package model

// StreamInfo 流信息
type StreamInfo struct {
	MediaServerId string `json:"mediaServerId"`
	App           string `json:"app"`
	Ip            string `json:"ip"`
	DeviceID      string `json:"deviceID"`
	ChannelId     string `json:"channelId"`
	Stream        string `json:"stream"`
	Rtmp          string `json:"rtmp"`
	Rtsp          string `json:"rtsp"`
	Flv           string `json:"flv"`
	HttpsFlv      string `json:"httpsFlv"`
	WSFly         string `json:"WSFly"`
	WSSFly        string `json:"WSSFly"`
	Fmp4          string `json:"fmp4"`
	HttpsFmp4     string `json:"httpsFmp4"`
	WSFmpt4       string `json:"WSFmpt4"`
	Hls           string `json:"hls"`
	HttpsHls      string `json:"httpsHls"`
	WsHls         string `json:"wsHls"`
	Ts            string `json:"ts"`
	HttpsTs       string `json:"httpsTs"`
	WebsocketTs   string `json:"websocketTs"`
}
