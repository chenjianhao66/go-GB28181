package gbsip

import "encoding/xml"

// notify消息
type (
	// Keepalive xml解析心跳包结构
	Keepalive struct {
		CmdType  string `xml:"CmdType"`
		SN       int    `xml:"SN"`
		DeviceID string `xml:"DeviceID"`
		Status   string `xml:"Status"`
		Info     string `xml:"Info"`
	}
)

// message消息
type (
	// CmdType 命令类型
	CmdType struct {
		CmdType string `xml:"CmdType"`
	}

	// DeviceID 设备id
	DeviceID struct {
		DeviceID string `xml:"DeviceID"`
	}

	SN struct {
		SN string `xml:"SN"`
	}

	R struct {
		Result string `xml:"Result"`
	}

	Mata struct {
		CmdType
		DeviceID
		SN
	}

	// DeviceInfo 设备信息
	DeviceInfo struct {
		CmdType      string `xml:"CmdType"`
		SN           string `xml:"SN"`
		DeviceID     string `xml:"DeviceID"`
		Result       string `xml:"Result"`
		DeviceName   string `xml:"DeviceName"`
		Manufacturer string `xml:"Manufacturer"`
		Model        string `xml:"Model"`
		Firmware     string `xml:"Firmware"`
	}

	// DeviceCatalogResponse 设备目录响应
	DeviceCatalogResponse struct {
		Name xml.Name `xml:"Response"`
		CmdType
		SN
		DeviceID
		SunNum        string     `xml:"SunNum"`
		DeviceListNum string     `xml:"Num,attr"`
		DeviceList    DeviceList `xml:"DeviceList"`
	}

	// DeviceList 目录列表
	DeviceList struct {
		Name  xml.Name      `xml:"DeviceList"`
		Items []CatalogItem `xml:"Item"`
	}

	// CatalogItem 设备目录中的每一项结构体
	CatalogItem struct {
		XmlName xml.Name `xml:"Item"`
		DeviceID
		Name         string `xml:"Name"`
		Manufacturer string `xml:"Manufacturer"`
		Model        string `xml:"Model"`
		Owner        string `xml:"Owner"`
		CivilCode    string `xml:"CivilCode"`
		Address      string `xml:"Address"`
		Parental     string `xml:"Parental"`
		ParentID     string `xml:"ParentID"`
		SafetyWay    string `xml:"SafetyWay"`
		RegisterWay  string `xml:"RegisterWay"`
		Secrecy      string `xml:"Secrecy"`
		Status       string `xml:"Status"`
	}

	// DeviceBasicConfigResp 设备基本配置查询返回结构体
	DeviceBasicConfigResp struct {
		Mata
		R
		BasicParam BasicParam `xml:"BasicParam"`
	}

	// BasicParam 设备基本配置Basic配置项结构体
	BasicParam struct {
		Name string `xml:"Name"`
		DeviceID
		SIPServerID       string `xml:"SIPServerID"`
		SIPServerIP       string `xml:"SIPServerIP"`
		SIPServerPort     string `xml:"SIPServerPort"`
		DomainName        string `xml:"DomainName"`
		Expiration        string `xml:"Expiration"`
		Password          string `xml:"Password"`
		HeartBeatInterval int    `xml:"HeartBeatInterval"`
		HeartBeatCount    int    `xml:"HeartBeatCount"`
	}
)
