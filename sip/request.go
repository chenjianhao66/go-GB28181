package sip

import (
	"bytes"
	"fmt"
	"github.com/gofrs/uuid"
)

// Request Request
type Request struct {
	message
	method    RequestMethod
	recipient *URI
}

// NewRequest 根据消息id 请求方法 uri sip版本 请求头和消息体构造一个request对象
// 其中消息id
func NewRequest(messID MessageID, method RequestMethod, uri *URI, sipVersion string, headers []Header, body []byte) *Request {
	req := new(Request)
	if messID == "" {
		// 随即生成一个
		req.messageID = MessageID(uuid.Must(uuid.NewV4()).String())
	} else {
		req.messageID = messID
	}
	req.SetMethod(method)
	req.SetRecipient(uri)
	req.SetSipVersion(sipVersion)
	if len(body) != 0 {
		req.SetBody(body, true)
	}
	req.headers = newHeader(headers)
	return req
}
func (r *Request) StartLine() string {
	var buf bytes.Buffer
	buf.WriteString(
		fmt.Sprintf(
			"%s %s %s",
			string(r.method),
			r.Recipient(),
			r.SipVersion(),
		),
	)
	return buf.String()
}

func (r *Request) Method() RequestMethod {
	return r.method
}

func (r *Request) SetMethod(method RequestMethod) {
	r.method = method
}

func (r *Request) Recipient() *URI {
	return r.recipient
}

func (r *Request) SetRecipient(uri *URI) {
	r.recipient = uri
}

func (r *Request) Clone() Message {
	return &Request{
		message:   message{},
		method:    "",
		recipient: nil,
	}
}

func (r *Request) IsCancel() bool {
	return r.Method() == CANCEL
}

func (r *Request) IsAck() bool {
	return r.Method() == ACK
}
