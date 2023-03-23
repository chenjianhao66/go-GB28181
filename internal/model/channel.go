package model

type Channel struct {
	Meta
	DeviceId     string `json:"DeviceId,omitempty" gorm:"column:deviceId;comment:通道id"`
	Name         string `json:"Name,omitempty" gorm:"column:name;comment:通道名称"`
	Manufacturer string `json:"Manufacturer,omitempty" gorm:"column:manufacturer;comment:当为设备时，设备厂商"`
	Model        string `json:"Model,omitempty" gorm:"column:model;comment:当为设备时，设备型号"`
	Owner        string `json:"Owner,omitempty" gorm:"column:owner;comment:当为设备时，设备归属"`
	CivilCode    string `json:"CivilCode,omitempty" gorm:"column:civilCode;comment:行政区域"`
	Address      string `json:"Address,omitempty" gorm:"column:address;comment:当为设备时，安装地址"`
	Parental     string `json:"Parental,omitempty" gorm:"column:parental;comment:当为设备时，是否有子设备，1有，0没有"`
	ParentID     string `json:"ParentID,omitempty" gorm:"column:parentId;comment:父设备/区域/系统ID"`
	SafetyWay    string `json:"SafetyWay,omitempty" gorm:"column:safetyWay;comment:信令安全模式，0不采用、2 S/MIME签名方式、3 S/MIME加密他签名同时采用方式、4 数字摘要方式"`
	RegisterWay  string `json:"RegisterWay,omitempty" gorm:"column:registerWay;comment:注册方式，1 标准认证注册模式 、2 基于口令的双向认证模式、3 基于数字证书的双向认证注册模式"`
	Secrecy      string `json:"Secrecy,omitempty" gorm:"column:secrecy;comment:保密属性，0不涉密、1涉密"`
	Status       string `json:"Status,omitempty" gorm:"column:status;comment:设备状态"`
}

type CameraExpand struct {
}

func NewChannelMust(deviceId string) Channel {
	return Channel{
		DeviceId: deviceId,
	}
}
