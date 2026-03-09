package model

type Channel struct {
	Meta
	// 设备唯一sipid
	DeviceId string `json:"deviceId,omitempty" gorm:"column:deviceId;unique;comment:通道id"`

	// 通道名称
	Name string `json:"name,omitempty" gorm:"column:name;comment:通道名称"`

	// 设备制造厂商
	Manufacturer string `json:"manufacturer,omitempty" gorm:"column:manufacturer;comment:当为设备时，设备厂商"`

	// 设备型号
	Model string `json:"model,omitempty" gorm:"column:model;comment:当为设备时，设备型号"`

	// 设备归属
	Owner string `json:"owner,omitempty" gorm:"column:owner;comment:当为设备时，设备归属"`

	// 行政区域
	CivilCode string `json:"civilCode,omitempty" gorm:"column:civilCode;comment:行政区域"`

	// 安装地址
	Address string `json:"address,omitempty" gorm:"column:address;comment:当为设备时，安装地址"`

	// 是否有子设备，1有、0没有
	Parental string `json:"parental,omitempty" gorm:"column:parental;comment:当为设备时，是否有子设备，1有，0没有"`

	// 父设备/区域/系统ID
	ParentID string `json:"parentID,omitempty" gorm:"column:parentId;comment:父设备/区域/系统ID"`

	// 信令安全模式，0不采用、2 S/MIME签名方式、3 S/MIME加密他签名同时采用方式、4 数字摘要方式
	SafetyWay string `json:"safetyWay,omitempty" gorm:"column:safetyWay;comment:信令安全模式，0不采用、2 S/MIME签名方式、3 S/MIME加密他签名同时采用方式、4 数字摘要方式"`

	// 注册方式，1 标准认证注册模式 、2 基于口令的双向认证模式、3 基于数字证书的双向认证注册模式
	RegisterWay string `json:"registerWay,omitempty" gorm:"column:registerWay;comment:注册方式，1 标准认证注册模式 、2 基于口令的双向认证模式、3 基于数字证书的双向认证注册模式"`

	// 保密属性，0不涉密、1涉密
	Secrecy string `json:"secrecy,omitempty" gorm:"column:secrecy;comment:保密属性，0不涉密、1涉密"`

	// 设备状态
	Status string `json:"status,omitempty" gorm:"column:status;comment:设备状态"`
}

type CameraExpand struct {
}

func NewChannelMust(deviceId string) Channel {
	return Channel{
		DeviceId: deviceId,
	}
}
