package model

type Code struct {
	Code int `json:"code"`
}

type Message struct {
	Msg string `json:"msg"`
}

type CodeMessage struct {
	Code
	Message
}

type GetRtpInfoResp struct {
	Code
	Exist     bool   `json:"exist,omitempty"`
	PeerIp    string `json:"peer_ip,omitempty"`
	PeerPort  int    `json:"peer_port,omitempty"`
	LocalIp   string `json:"local_ip,omitempty"`
	LocalPort int    `json:"local_port,omitempty"`
}
