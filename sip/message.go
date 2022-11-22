package sip

import (
	"bytes"
	"fmt"
)

type requestMethod string

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
	INVITE    requestMethod = "INVITE"
	ACK       requestMethod = "ACK"
	CANCEL    requestMethod = "CANCEL"
	BYE       requestMethod = "BYE"
	REGISTER  requestMethod = "REGISTER"
	OPTIONS   requestMethod = "OPTIONS"
	SUBSCRIBE requestMethod = "SUBSCRIBE"
	NOTIFY    requestMethod = "NOTIFY"
	REFER     requestMethod = "REFER"
	INFO      requestMethod = "INFO"
	MESSAGE   requestMethod = "MESSAGE"
)

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
