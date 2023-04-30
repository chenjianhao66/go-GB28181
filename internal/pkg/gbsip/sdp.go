package gbsip

import (
	sdp "github.com/panjjo/gosdp"
	"net"
	"time"
)

func createSdpInfo(mediaIp, channelId, ssrc string, rtpPort int) string {
	origin := sdp.Origin{
		Username:       channelId,
		SessionID:      0,
		SessionVersion: 0,
		// Internet
		NetworkType: "IN",
		// ipv4
		AddressType: "IP4",
		Address:     mediaIp,
	}

	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "video",
			Port:     rtpPort,
			Protocol: "RTP/RTCP",
			Formats:  []string{"96", "98", "97"},
		},
		Connection: sdp.ConnectionData{
			NetworkType: "IN",
			AddressType: "IP4",
			IP:          net.ParseIP(mediaIp),
			TTL:         0,
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
		SSRC:    ssrc,
	}
	session := msg.Append(sdp.Session{})
	bytes := session.AppendTo([]byte{})
	return string(bytes)
}
