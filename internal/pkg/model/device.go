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
	DeviceId string `json:"deviceId" gorm:"column:deviceId;comment:设备的sip唯一id"`

	// 设备的sip域名
	Domain string `json:"domain" gorm:"column:domain;comment:设备的sip域名"`

	// 设备名
	Name string `json:"name" gorm:"column:name;comment:设备名"`

	// 制造厂商
	Manufacturer string `json:"manufacturer" gorm:"column:manufacturer;comment:制造厂商"`

	// 设备型号
	Model string `json:"model" gorm:"column:model;comment:设备型号"`

	// 固件版本
	Firmware string `json:"firmware" gorm:"column:firmware;comment:固件版本"`

	// 传输模式
	Transport string `json:"transport" gorm:"column:transport;comment:传输模式"`

	// 是否在线
	Offline uint8 `json:"offline" gorm:"column:offline;comment:是否在线:0不在线1在线"`

	// ip地址
	Ip string `json:"ip" gorm:"column:ip;comment:ip地址"`

	// 传输端口
	Port string `json:"port" gorm:"column:port;comment:传输端口"`

	// 心跳过期时间
	Expires string `json:"expires" gorm:"column:expires;comment:心跳过期时间"`

	// 注册时间
	RegisterTime time.Time `json:"registerTime" gorm:"column:register_time;comment:注册时间"`

	// 上次心跳时间
	Keepalive time.Time `json:"keepalive" gorm:"column:keepalive;comment:上次心跳时间"`

	// 心跳间隔时间，范围值 5-255
	HeartBeatInterval int `json:"heartBeatInterval" gorm:"column:heartBeatInterval;comment:心跳间隔时间，5-255;default:5"`

	// 心跳超时次数，范围值：3-255
	HeartBeatCount int `json:"heartBeatCount" gorm:"column:heartBeatCount;comment:心跳超时次数，3-255;default:3"`
}
