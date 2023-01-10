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
	DeviceId     string    `json:"deviceId"`
	Name         string    `json:"name"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	Firmware     string    `json:"firmware"`
	Transport    string    `json:"transport"`
	Offline      uint8     `json:"offline"`
	Ip           string    `json:"ip"`
	Port         int8      `json:"port"`
	Expires      int8      `json:"expires"`
	RegisterTime time.Time `json:"registerTime"`
	Keepalive    time.Time `json:"keepalive"`
}
