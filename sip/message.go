package sip

import (
	"bytes"
	"fmt"
	"net"
)

type RequestMethod string

// MessageID MessageID
type MessageID string

type Port uint16

func (p *Port) equals(other interface{}) bool {
	if obj, ok := other.(*Port); ok {
		if p == obj {
			return true
		}

		return *p == *obj
	}

	return false
}

func (p *Port) Clone() *Port {
	if p == nil {
		return nil
	}
	newPort := *p
	return &newPort
}

const (
	INVITE    RequestMethod = "INVITE"
	ACK       RequestMethod = "ACK"
	CANCEL    RequestMethod = "CANCEL"
	BYE       RequestMethod = "BYE"
	REGISTER  RequestMethod = "REGISTER"
	OPTIONS   RequestMethod = "OPTIONS"
	SUBSCRIBE RequestMethod = "SUBSCRIBE"
	NOTIFY    RequestMethod = "NOTIFY"
	REFER     RequestMethod = "REFER"
	INFO      RequestMethod = "INFO"
	MESSAGE   RequestMethod = "MESSAGE"
)

const (
	DefaultProtocol   = "udp"
	DefaultSipVersion = "SIP/2.0"
)

// Message 将SIP协议中的一个消息抽象成接口
type Message interface {
	MessageID() MessageID
	StartLine() string
	String() string
	SipVersion() string
	SetSipVersion(version string)
	Body() []byte
	SetBody(body []byte, setContentLength bool)
	Transport() string
	Source() net.Addr
	SetSource(src net.Addr)
	Destination() net.Addr
	SetDestination(dest net.Addr)

	Via() (ViaHeader, bool)
	ViaParams() (*viaParams, bool)
	From() (*From, bool)
	To() (*To, bool)
	CallID() (*CallId, bool)
	CSeq() (*CSeq, bool)
	ContentLength() (*ContentLength, bool)
	ContentType() (*ContentType, bool)
	Contact() (*Contact, bool)
	Headers() []Header
	GetHeader(name string) []Header
	AppendHeader(header Header)
	RemoveHeader(name string)

	Clone() Message
	IsCancel() bool
	IsAck() bool
}

type message struct {
	*headers

	messageID  MessageID
	sipVersion string
	body       []byte
	source     net.Addr
	dest       net.Addr
	startLine  func() string
}

func (m *message) MessageID() MessageID {
	return m.messageID
}

func (m *message) StartLine() string {
	return m.startLine()
}

func (m message) String() string {
	var buf bytes.Buffer
	buf.WriteString(m.startLine() + "\r\n")
	buf.WriteString(m.headers.String() + "\r\n")
	buf.Write(m.body)
	return buf.String()
}

func (m *message) SipVersion() string {
	return m.sipVersion
}

func (m *message) SetSipVersion(version string) {
	m.sipVersion = version
}

func (m *message) Body() []byte {
	return m.body
}

func (m *message) SetBody(body []byte, setContentLength bool) {
	m.body = body
	if setContentLength {
		hdrs := m.GetHeader("Content-Length")
		if len(hdrs) == 0 {
			length := ContentLength(len(body))
			m.AppendHeader(&length)
		} else {
			length := ContentLength(len(body))
			hdrs[0] = &length
		}
	}
}

// Transport  Transport
func (m *message) Transport() string {
	if viaHop, ok := m.ViaParams(); ok {
		return viaHop.transport
	}
	return DefaultProtocol
}

// Source Source
func (m *message) Source() net.Addr {
	return m.source
}

// SetSource SetSource
func (m *message) SetSource(src net.Addr) {
	m.source = src
}

// Destination Destination
func (m *message) Destination() net.Addr {
	return m.dest
}

// SetDestination SetDestination
func (m *message) SetDestination(dest net.Addr) {
	m.dest = dest
}

// URI object of sip message
type URI struct {
	// True if and only if the URI is a SIPS URI.
	isEncrypted bool

	// FUser part in sip message request line
	FUser MaybeString

	// FHost part in sip message request line
	FHost string

	FPassword  MaybeString
	FPort      *Port
	FUriParams Params
	FHeaders   Params
}

func (uri *URI) String() string {
	var buffer bytes.Buffer

	// Compulsory protocol identifier.
	if uri.isEncrypted {
		buffer.WriteString("sips")
		buffer.WriteString(":")
	} else {
		buffer.WriteString("sip")
		buffer.WriteString(":")
	}

	// 可选的用户部分
	if user, ok := uri.FUser.(String); ok && user.String() != "" {
		buffer.WriteString(uri.FUser.String())
		if pass, ok := uri.FPassword.(String); ok && pass.String() != "" {
			buffer.WriteString(":")
			buffer.WriteString(pass.String())
		}
		buffer.WriteString("@")
	}

	buffer.WriteString(uri.FHost)
	if uri.FPort != nil {
		buffer.WriteString(fmt.Sprintf(":%d", *uri.FPort))
	}

	if (uri.FUriParams != nil) && uri.FUriParams.Length() > 0 {
		buffer.WriteString(";")
		buffer.WriteString(uri.FUriParams.ToString(';'))
	}

	if (uri.FHeaders != nil) && uri.FHeaders.Length() > 0 {
		buffer.WriteString("?")
		buffer.WriteString(uri.FHeaders.ToString('&'))
	}

	return buffer.String()
}

func (uri *URI) Equals(other interface{}) bool {
	return true
}

func (uri *URI) Clone() *URI {
	if uri == nil {
		var newURI *URI
		return newURI
	}

	cloneWithNil := func(params Params) Params {
		if params == nil {
			return NewParams()
		}
		return params.Clone()
	}

	newURI := &URI{
		isEncrypted: uri.isEncrypted,
		FUser:       uri.FUser,
		FHost:       uri.FHost,
		FPassword:   uri.FPassword,
		FUriParams:  cloneWithNil(uri.FUriParams),
		FHeaders:    cloneWithNil(uri.FHeaders),
	}

	if uri.FPort != nil {
		newURI.FPort = uri.FPort.Clone()
	}
	return newURI
}
