package model

type DeviceControl struct {
	// 设备id
	DeviceId string `json:"deviceId,omitempty"`
	// 通道id
	ChannelId string `json:"channelId,omitempty"`
	// 控制的命令，取值为：left、right、down、up、downright、downleft、upright、upleft、zoomin、zoomout
	Command string `json:"command,omitempty"`
	// 水平方向移动速度，取值：0-255
	HorizonSpeed int `json:"horizonSpeed,omitempty"`
	// 垂直方向移动速度，取值：0-255
	VerticalSpeed int `json:"verticalSpeed,omitempty"`
	// 变倍控制速度，取值：0-255
	ZoomSpeed int `json:"zoomSpeed,omitempty"`
}
