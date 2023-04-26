package model

import (
	"time"
)

// Meta 结构元数据，
type Meta struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Device 设备表entity
type Device struct {
	Meta
	// 设备的sip唯一id
	DeviceId string `json:"deviceId" gorm:"column:deviceId"`

	// 设备的sip域名
	Domain string `json:"domain" gorm:"column:domain"`

	// 设备名
	Name string `json:"name" gorm:"column:name"`

	// 制造厂商
	Manufacturer string `json:"manufacturer" gorm:"column:manufacturer"`

	// 设备型号
	Model string `json:"model" gorm:"column:model"`

	// 固件版本
	Firmware string `json:"firmware" gorm:"column:firmware"`

	// 传输模式
	Transport string `json:"transport" gorm:"column:transport"`

	// 是否在线
	Offline uint8 `json:"offline" gorm:"column:offline"`

	// ip地址
	Ip string `json:"ip" gorm:"column:ip"`

	// 传输端口
	Port string `json:"port" gorm:"column:port"`

	// 心跳过期时间
	Expires string `json:"expires" gorm:"column:expires"`

	// 注册时间
	RegisterTime time.Time `json:"registerTime" gorm:"column:register_time"`

	// 上次心跳时间
	Keepalive time.Time `json:"keepalive" gorm:"column:keepalive"`
}
