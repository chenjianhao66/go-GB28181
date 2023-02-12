package model

type C struct {
	Code int `json:"code"`
}

type Message struct {
	Msg string `json:"msg"`
}

type CodeMessage struct {
	C
	Message
}

type GetRtpInfoResp struct {
	C
	Message
	Exist     bool   `json:"exist,omitempty"`
	PeerIp    string `json:"peer_ip,omitempty"`
	PeerPort  int    `json:"peer_port,omitempty"`
	LocalIp   string `json:"local_ip,omitempty"`
	LocalPort int    `json:"local_port,omitempty"`
}

type CreateRtpServerResp struct {
	C
	Message
	Port int `json:"port"`
}
