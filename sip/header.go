package sip

import (
	"bytes"
	"fmt"
	"strings"
)

// Header sip Header
type Header interface {
	Name() string
	Clone() Header
	String() string
	Equals(other interface{}) bool
}

// =================== ContentLengthHead ===================

// ContentLength ContentLength part in SIP Header
type ContentLength uint32

func (c *ContentLength) Name() string {
	return "Content-Length"
}

func (c *ContentLength) Clone() Header {
	return c
}

func (c ContentLength) String() string {
	return fmt.Sprintf("Content-Length: %d", uint32(c))
}

func (c *ContentLength) Equals(other interface{}) bool {
	// 本身等于nil的话直接返回false
	if c == nil {
		return false
	}
	// 进行类型断言，进行值的判断
	if o, ok := other.(ContentLength); ok {
		return *c == o
	}

	// 进行指针的类型断言，对指针地址进行判断，对指针所对应的值进行判断
	if o, ok := other.(*ContentLength); ok {
		if c == o {
			return true
		}
		return *c == *o
	}

	// 如果类型断言都失败了，则返回false
	return false
}

// =================== Via Header ===================

// ViaHeader 请求中的Via请求头
type ViaHeader []*viaParams

func (v ViaHeader) Name() string {
	return "Via"
}

func (v ViaHeader) Clone() Header {
	if v == nil {
		var newVia ViaHeader
		return newVia
	}

	params := make([]*viaParams, len(v))
	for _, v := range v {
		params = append(params, v.Clone())
	}
	return ViaHeader(params)
}

func (v ViaHeader) String() string {
	var buf bytes.Buffer
	buf.WriteString("Via: ")
	for i, params := range v {
		buf.WriteString(params.String())
		if i != len(v) {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func (v ViaHeader) Equals(other interface{}) bool {
	if obj, ok := other.(ViaHeader); ok {
		if len(obj) != len(v) {
			return false
		}
		for i, params := range v {
			if !params.Equals(obj[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// sip请求的via字段抽象结构体
type viaParams struct {
	protocolName    string
	protocolVersion string
	transport       string
	host            string
	// 可选字段，所以是指针类型，因为指针类型可以判断是否为nil
	port *Port
	// 不规则参数
	Params Params
}

func (v *viaParams) String() string {
	var buf bytes.Buffer
	buf.WriteString(
		fmt.Sprintf("%s/%s/%s %s",
			v.protocolName,
			v.protocolVersion,
			v.transport,
			v.host,
		),
	)
	if v.port != nil {
		buf.WriteString(fmt.Sprintf(":%d;", v.port))
	}

	return buf.String()
}

func (v *viaParams) Equals(other interface{}) bool {
	if v == nil {
		return false
	}

	if p, ok := other.(*viaParams); ok {
		if p == v {
			return true
		}

		result := v.protocolName == p.protocolName && v.transport == p.transport && v.host == p.host && v.port.equals(p.port)
		return result
	}

	return false
}

func (v *viaParams) Clone() *viaParams {
	var new *viaParams
	if v == nil {
		return new
	}

	new = &viaParams{
		protocolName:    v.protocolName,
		protocolVersion: v.protocolVersion,
		transport:       v.transport,
		host:            v.host,
	}

	if v.port != nil {
		new.port = v.port.Clone()
	}

	if v.Params != nil {
		new.Params = v.Params.Clone()
	}
	return new
}

// =================== from Header ===================

type From struct {
	Address *URI

	Params Params
}

func (f *From) Name() string {
	return "From"
}

func (f *From) Clone() Header {
	return f
}

func (f *From) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("From: ")

	//if displayName, ok := f.DisplayName.(String); ok && displayName.String() != "" {
	//	buffer.WriteString(fmt.Sprintf("\"%s\" ", displayName))
	//}

	buffer.WriteString(fmt.Sprintf("<%s>", f.Address))

	if f.Params.Length() > 0 {
		buffer.WriteString(";")
		buffer.WriteString(f.Params.ToString(';'))
	}

	return buffer.String()
}

func (f *From) Equals(other interface{}) bool {
	// TODO: fix this
	return true
}

// =================== Params interface in Header ===================

// Params Generic list of parameters on a Header.
type Params interface {
	Get(key string) (MaybeString, bool)
	Add(key string, val MaybeString) Params
	Clone() Params
	Equals(params interface{}) bool
	ToString(sep uint8) string
	String() string
	Length() int
	Items() map[string]MaybeString
	Keys() []string
	Has(key string) bool
}

// MaybeString  wrapper
type MaybeString interface {
	String() string
	Equals(other interface{}) bool
}

// String 包装请求头中的参数结构体
type String struct {
	Str string
}

func (str String) String() string {
	return str.Str
}

// Equals Equals
func (str String) Equals(other interface{}) bool {
	if v, ok := other.(String); ok {
		return str.Str == v.Str
	}

	return false
}

// Params implementation.
type headerParams struct {
	params     map[string]MaybeString
	paramOrder []string
}

// NewParams Create an empty set of parameters.
func NewParams() Params {
	return &headerParams{
		params:     make(map[string]MaybeString),
		paramOrder: []string{},
	}
}

func (h *headerParams) Get(key string) (MaybeString, bool) {
	v, ok := h.params[key]
	return v, ok
}

func (h *headerParams) Add(key string, val MaybeString) Params {
	if _, ok := h.params[key]; !ok {
		h.paramOrder = append(h.paramOrder, key)
	}
	h.params[key] = val
	return h
}

func (h *headerParams) Clone() Params {
	if h == nil {
		var dump *headerParams
		return dump
	}
	newParams := NewParams()
	for _, key := range h.Keys() {
		if val, ok := h.Get(key); ok {
			newParams.Add(key, val)
		}
	}

	return newParams
}

func (h *headerParams) Equals(params interface{}) bool {
	v, ok := params.(*headerParams)
	if !ok {
		return false
	}

	if v == h {
		return true
	}

	if v == nil && h != nil || v != nil && h == nil {
		return false
	}

	if v.Length() == 0 && h.Length() == 0 {
		return true
	}

	if v.Length() != h.Length() {
		return false
	}

	for key, pVal := range h.Items() {
		qVal, ok := v.Get(key)
		if !ok {
			return false
		}
		if pVal != qVal {
			return false
		}
	}

	return true
}

func (h *headerParams) ToString(sep uint8) string {
	if h == nil {
		return ""
	}

	var buffer bytes.Buffer
	// 用于判断是否是第一次
	first := true

	for _, key := range h.Keys() {
		val, ok := h.Get(key)
		if !ok {
			continue
		}

		// 第一次不用执行，其余都执行
		if !first {
			buffer.WriteString(fmt.Sprintf("%c", sep))
		}
		first = false

		buffer.WriteString(key)

		if val, ok := val.(String); ok {
			if strings.ContainsAny(val.String(), tab) {
				// 如果该字段包含制表符，则添加\
				buffer.WriteString(fmt.Sprintf("=\"%s\"", val.String()))
			} else {
				// 反之不添加
				buffer.WriteString(fmt.Sprintf("=%s", val.String()))
			}
		}
	}

	return buffer.String()
}

func (h *headerParams) String() string {
	if h == nil {
		return ""
	}
	return h.ToString('&')
}

func (h *headerParams) Length() int {
	return len(h.paramOrder)
}

func (h *headerParams) Items() map[string]MaybeString {
	return h.params
}

func (h *headerParams) Keys() []string {
	return h.paramOrder
}

func (h *headerParams) Has(key string) bool {
	_, ok := h.params[key]
	return ok
}
