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
	DeviceId     string    `json:"deviceId" gorm:"column:deviceId"`
	Domain       string    `json:"domain" gorm:"column:domain"`
	Name         string    `json:"name" gorm:"column:name"`
	Manufacturer string    `json:"manufacturer" gorm:"column:manufacturer"`
	Model        string    `json:"model" gorm:"column:model"`
	Firmware     string    `json:"firmware" gorm:"column:firmware"`
	Transport    string    `json:"transport" gorm:"column:transport"`
	Offline      uint8     `json:"offline" gorm:"column:offline"`
	Ip           string    `json:"ip" gorm:"column:ip"`
	Port         string    `json:"port" gorm:"column:port"`
	Expires      string    `json:"expires" gorm:"column:expires"`
	RegisterTime time.Time `json:"registerTime" gorm:"column:register_time"`
	Keepalive    time.Time `json:"keepalive" gorm:"column:keepalive"`
}
