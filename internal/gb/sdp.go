package gb

import (
	sdp "github.com/panjjo/gosdp"
	"net"
	"time"
)

func createSdpInfo() string {
	origin := sdp.Origin{
		// TODO 发起者的国标id，后续修改
		Username:       "44010200491318000001",
		SessionID:      0,
		SessionVersion: 0,
		// Internet
		NetworkType: "IN",
		// ipv4
		AddressType: "IP4",
		// TODO 流媒体服务ip
		Address: "192.168.1.224",
	}

	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "video",
			Port:     30002,
			Protocol: "RTP/RTCP",
			Formats:  []string{"96", "98", "97"},
		},
		Connection: sdp.ConnectionData{
			NetworkType: "IN",
			AddressType: "IP4",
			// TODO 流媒体服务IP
			IP:  net.ParseIP("192.168.1.224"),
			TTL: 0,
		},
	}
	video.AddAttribute("recvonly")
	video.AddAttribute("rtpmap", "96", "PS/90000")
	video.AddAttribute("rtpmap", "98", "H264/90000")
	video.AddAttribute("rtpmap", "97", "MPEG4/90000")

	msg := sdp.Message{
		Version: 0,
		Origin:  origin,
		Name:    "Play",
		Medias:  sdp.Medias{video},
		Timing:  []sdp.Timing{sdp.Timing{Start: time.Time{}, End: time.Time{}}},
		SSRC:    "0102003583",
	}
	session := msg.Append(sdp.Session{})
	bytes := session.AppendTo([]byte{})
	return string(bytes)
}
