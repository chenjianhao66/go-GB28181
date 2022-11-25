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
	MaybeString
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
		if i != len(v)-1 {
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
		buf.WriteString(fmt.Sprintf(":%d;", *v.port))
	}

	if v.Params.Length() > 0 {
		//buf.WriteString(";")
		buf.WriteString(v.Params.ToString(';'))
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

// =================== From Header ===================

type From struct {
	DisplayName MaybeString

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
	if h, ok := other.(*From); ok {
		if f == h {
			return true
		}

		if f == nil && h != nil || f != nil && h == nil {
			return false
		}

		res := true

		// 判断displayName字段是否相等
		if f.DisplayName != h.DisplayName {
			// 如果f.displayName == nil，那就期望h.displayName == nil
			if f.DisplayName == nil {
				res = res && h.DisplayName == nil
			} else {
				res = res && f.DisplayName.Equals(h.DisplayName)
			}
		}

		// 判断address字段是否相等
		if f.Address != h.Address {
			// 如果f.address字段为空，那判断h.address字段是否为空，还需要判断res变量是否在之前被设置为false
			if f.Address == nil {
				res = res && h.Address == nil
			} else {
				// f.address字段不为空，那么就调用Equals方法去判断
				res = res && f.Address.Equals(h.Address)
			}
		}

		// 判断Params字段是否相等
		if f.Params != h.Params {
			if f.Params == nil {
				res = res && h.Params == nil
			} else {
				res = res && f.Params.Equals(h.Params)
			}
		}

		return res
	}

	return false
}

// =================== To Header ===================

type To struct {
	DisplayName MaybeString

	Address *URI

	Params Params
}

func (to *To) Name() string {
	return "To"
}

func (to *To) Clone() Header {
	if to == nil {
		var newTo *To
		return newTo
	}

	newTo := &To{
		DisplayName: to.DisplayName,
		Address:     to.Address.Clone(),
	}

	if to.Params != nil {
		newTo.Params = to.Params.Clone()
	}
	return newTo
}

func (to *To) String() string {
	var buf bytes.Buffer
	buf.WriteString("To: ")
	if displayName, ok := to.DisplayName.(String); ok && displayName.String() != "" {
		buf.WriteString(fmt.Sprintf("\"%s\" ", displayName))
	}

	buf.WriteString(fmt.Sprintf("<%s>", to.Address))

	if to.Params != nil && to.Params.Length() > 0 {
		buf.WriteString(";")
		buf.WriteString(to.Params.ToString(';'))
	}
	return buf.String()
}

func (to *To) Equals(other interface{}) bool {
	if h, ok := other.(*To); ok {
		if to == h {
			return true
		}
		if to == nil && h != nil || to != nil && h == nil {
			return false
		}

		res := true

		if to.DisplayName != h.DisplayName {
			if to.DisplayName == nil {
				res = res && h.DisplayName == nil
			} else {
				res = res && to.DisplayName.Equals(h.DisplayName)
			}
		}

		if to.Address != h.Address {
			if to.Address == nil {
				res = res && h.Address == nil
			} else {
				res = res && to.Address.Equals(h.Address)
			}
		}

		if to.Params != h.Params {
			if to.Params == nil {
				res = res && h.Params == nil
			} else {
				res = res && to.Params.Equals(h.Params)
			}
		}

		return res
	}

	return false
}

// =================== CallId Header ===================

type CallId string

func (callId *CallId) Name() string {
	return "Call-ID"
}

func (callId *CallId) Clone() Header {
	return callId
}

func (callId CallId) String() string {
	return fmt.Sprintf("Call-ID: %s", string(callId))
}

func (callId *CallId) Equals(other interface{}) bool {
	if h, ok := other.(CallId); ok {
		if callId == nil {
			return false
		}

		return *callId == h
	}
	if h, ok := other.(*CallId); ok {
		if callId == h {
			return true
		}
		if callId == nil && h != nil || callId != nil && h == nil {
			return false
		}

		return *callId == *h
	}

	return false
}

// =================== CallId Header ===================

type CSeq struct {
	SeqNumber  uint32
	MethodName RequestMethod
}

func (cSeq *CSeq) Name() string {
	return "CSeq"
}

func (cSeq *CSeq) Clone() Header {
	if cSeq == nil {
		var newCSeq *CSeq
		return newCSeq
	}

	return &CSeq{
		SeqNumber:  cSeq.SeqNumber,
		MethodName: cSeq.MethodName,
	}
}

func (cSeq *CSeq) String() string {
	return fmt.Sprintf("CSeq: %d %s", cSeq.SeqNumber, cSeq.MethodName)
}

func (cSeq *CSeq) Equals(other interface{}) bool {
	if seq, ok := other.(*CSeq); ok {
		if cSeq == seq {
			return true
		}

		if cSeq == nil && seq != nil || cSeq != nil && seq == nil {
			return false
		}

		return cSeq.SeqNumber == seq.SeqNumber && cSeq.MethodName == seq.MethodName
	}

	return false
}

// =================== CallId Header ===================

type Contact struct {
	DisplayName MaybeString
	Address     *URI
	Params      Params
}

func (contact *Contact) Name() string {
	return "Contact"
}

func (contact *Contact) Clone() Header {
	if contact == nil {
		var newContact *Contact
		return newContact
	}

	newContact := &Contact{
		DisplayName: contact.DisplayName,
		Address:     contact.Address.Clone(),
	}

	if contact.Params != nil {
		newContact.Params = contact.Params.Clone()
	}

	return newContact
}

func (contact *Contact) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Contact: ")

	if displayName, ok := contact.DisplayName.(String); ok && displayName.String() != "" {
		buffer.WriteString(fmt.Sprintf("\"%s\" ", displayName))
	}

	buffer.WriteString(fmt.Sprintf("<%s>", contact.Address.String()))

	if (contact.Params != nil) && (contact.Params.Length() > 0) {
		buffer.WriteString(";")
		buffer.WriteString(contact.Params.ToString(';'))
	}

	return buffer.String()
}

func (contact *Contact) Equals(other interface{}) bool {
	if h, ok := other.(*Contact); ok {
		if contact == h {
			return true
		}
		if contact == nil && h != nil || contact != nil && h == nil {
			return false
		}

		res := true

		if contact.DisplayName != h.DisplayName {
			if contact.DisplayName == nil {
				res = res && h.DisplayName == nil
			} else {
				res = res && contact.DisplayName.Equals(h.DisplayName)
			}
		}

		if contact.Address != h.Address {
			if contact.Address == nil {
				res = res && h.Address == nil
			} else {
				res = res && contact.Address.Equals(h.Address)
			}
		}

		if contact.Params != h.Params {
			if contact.Params == nil {
				res = res && h.Params == nil
			} else {
				res = res && contact.Params.Equals(h.Params)
			}
		}

		return res
	}

	return false
}

// =================== Max-Forwards Header ===================

type MaxForwards uint32

func (maxForwards *MaxForwards) Name() string {
	return "Max-Forwards"
}

func (maxForwards *MaxForwards) Clone() Header {
	return maxForwards
}

func (maxForwards MaxForwards) String() string {
	return fmt.Sprintf("Max-Forwards: %d", int(maxForwards))
}

func (maxForwards *MaxForwards) Equals(other interface{}) bool {
	if h, ok := other.(MaxForwards); ok {
		if maxForwards == nil {
			return false
		}

		return *maxForwards == h
	}
	if h, ok := other.(*MaxForwards); ok {
		if maxForwards == h {
			return true
		}
		if maxForwards == nil && h != nil || maxForwards != nil && h == nil {
			return false
		}

		return *maxForwards == *h
	}

	return false
}

// =================== User-Agent ===================

type UserAgent string

func (userAgent *UserAgent) Name() string {
	return "User-Agent"
}

func (userAgent *UserAgent) Clone() Header {
	return userAgent
}

func (userAgent UserAgent) String() string {
	return fmt.Sprintf("User-Agent: %s", string(userAgent))
}

func (userAgent *UserAgent) Equals(other interface{}) bool {
	if h, ok := other.(UserAgent); ok {
		if userAgent == nil {
			return false
		}

		return *userAgent == h
	}
	if h, ok := other.(*UserAgent); ok {
		if userAgent == h {
			return true
		}
		if userAgent == nil && h != nil || userAgent != nil && h == nil {
			return false
		}

		return *userAgent == *h
	}

	return false
}

// =================== Expires ===================

type Expires uint32

func (expires *Expires) Name() string {
	return "Expires"
}

func (expires *Expires) Clone() Header {
	return expires
}

func (expires Expires) String() string {
	return fmt.Sprintf("Expires: %d", int(expires))
}

func (expires *Expires) Equals(other interface{}) bool {
	if h, ok := other.(Expires); ok {
		if expires == nil {
			return false
		}

		return *expires == h
	}
	if h, ok := other.(*Expires); ok {
		if expires == h {
			return true
		}
		if expires == nil && h != nil || expires != nil && h == nil {
			return false
		}

		return *expires == *h
	}

	return false
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

// ==================   ContentTypeHeader   ================

// ContentType ContentType
type ContentType string

func (ct ContentType) String() string { return "Content-Type: " + string(ct) }

// Name Name
func (ct *ContentType) Name() string { return "Content-Type" }

// Clone Clone
func (ct *ContentType) Clone() Header { return ct }

// Equals Equals
func (ct *ContentType) Equals(other interface{}) bool {
	if h, ok := other.(ContentType); ok {
		if ct == nil {
			return false
		}

		return *ct == h
	}
	if h, ok := other.(*ContentType); ok {
		if ct == h {
			return true
		}
		if ct == nil && h != nil || ct != nil && h == nil {
			return false
		}

		return *ct == *h
	}

	return false
}

// =================== Params interface in Header ===================

// Params Generic list of parameters on a Header.
type Params interface {
	MaybeString
	Get(key string) (MaybeString, bool)
	Add(key string, val MaybeString) Params
	Clone() Params
	ToString(sep uint8) string
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

// 消息头内的特殊字段
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

		//buffer.WriteString(key)

		if val, ok := val.(String); ok {
			if val.String() == "" {
				buffer.WriteString(fmt.Sprintf("%s", key))
			} else if strings.ContainsAny(val.String(), tab) {
				// 如果该字段包含制表符，则添加\
				buffer.WriteString(fmt.Sprintf("%s=\"%s\"", key, val.String()))
			} else {
				// 反之不添加
				buffer.WriteString(fmt.Sprintf("%s=%s", key, val.String()))
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

// =================== headers ===================

// 一个SIP消息内存在的消息头以及消息头顺序
// 抽象类，实现了部分Message接口的方法
type headers struct {
	headers     map[string][]Header
	headerOrder []string
}

// 根据h去返回一个headers实例
func newHeader(h []Header) *headers {
	hs := new(headers)
	hs.headers = make(map[string][]Header)
	hs.headerOrder = make([]string, 0)
	for _, header := range h {
		hs.AppendHeader(header)
	}
	return hs
}

// CopyHeaders 从来源Message根据name去复制请求头到目标Message
func CopyHeaders(name string, source, target Message) {
	name = strings.ToLower(name)
	for _, h := range source.GetHeader(name) {
		target.AppendHeader(h.Clone())
	}
}

// GetHeader 根据请求头名获取请求头切片，如果指定请求头不存在则返回空切片
func (h *headers) GetHeader(name string) []Header {
	if h.headers == nil {
		h.headers = make(map[string][]Header)
		h.headerOrder = []string{}
	}
	if headers, ok := h.headers[name]; ok {
		return headers
	}
	return []Header{}
}

func (h *headers) RemoveHeader(name string) {
	name = strings.ToLower(name)
	delete(h.headers, name)
	for i, v := range h.headerOrder {
		if v == name {
			// 跳过当前i所指的元素，将i+1往后的元素添加到i之前的切片中
			// 相当于删除当前i所指的元素
			h.headerOrder = append(h.headerOrder[:i], h.headerOrder[i+1:]...)
			break
		}
	}
}

func (h *headers) Headers() []Header {
	result := make([]Header, 0)
	for _, headerName := range h.headerOrder {
		result = append(result, h.headers[headerName]...)
	}
	return result
}

func (h *headers) AppendHeader(header Header) {
	name := strings.ToLower(header.Name())
	if _, ok := h.headers[name]; ok {
		// 如果当前头在map中存在，则追加
		h.headers[name] = append(h.headers[name], header)
	} else {
		// 不存在则初始化切片，并且存放到顺序切片中
		h.headers[name] = []Header{header}
		h.headerOrder = append(h.headerOrder, name)
	}
}

func (h *headers) CloneHeaders() []Header {
	hdrs := make([]Header, 0)
	for _, header := range h.Headers() {
		hdrs = append(hdrs, header.Clone())
	}

	return hdrs
}

func (h headers) String() string {
	buffer := bytes.Buffer{}
	// 根据顺序切片去构造请求头
	for typeIdx, name := range h.headerOrder {
		headers := h.headers[name]
		for idx, header := range headers {
			buffer.WriteString(header.String())
			// 该请求头的最后添加 '\r\n'
			if typeIdx < len(h.headerOrder) || idx < len(headers) {
				buffer.WriteString("\r\n")
			}
		}
	}
	return buffer.String()
}

func (h *headers) Via() (ViaHeader, bool) {
	header := h.GetHeader("Via")
	if len(header) == 0 {
		return nil, false
	}

	viaHeader, ok := header[0].(ViaHeader)
	if !ok {
		return nil, false
	}

	return viaHeader, true
}

func (h *headers) ViaParams() (*viaParams, bool) {
	via, ok := h.Via()
	if !ok {
		return nil, false
	}
	hops := []*viaParams(via)
	if len(hops) == 0 {
		return nil, false
	}

	return hops[0], true
}

func (h *headers) From() (*From, bool) {
	header := h.GetHeader("From")
	if len(header) == 0 {
		return nil, false
	}

	from, ok := header[0].(*From)
	if !ok {
		return nil, false
	}
	return from, true
}

func (h *headers) To() (*To, bool) {
	header := h.GetHeader("To")
	if len(header) == 0 {
		return nil, false
	}

	to, ok := header[0].(*To)
	if !ok {
		return nil, false
	}
	return to, true
}

func (h *headers) CallID() (*CallId, bool) {
	header := h.GetHeader("Call-ID")
	if len(header) == 0 {
		return nil, false
	}

	callId, ok := header[0].(*CallId)
	if !ok {
		return nil, false
	}
	return callId, true
}

func (h *headers) CSeq() (*CSeq, bool) {
	header := h.GetHeader("CSeq")
	if len(header) == 0 {
		return nil, false
	}
	cSeq, ok := header[0].(*CSeq)
	if !ok {
		return nil, false
	}
	return cSeq, true
}

func (h *headers) ContentLength() (*ContentLength, bool) {
	header := h.GetHeader("Content-Length")
	if len(header) == 0 {
		return nil, false
	}
	contentLength, ok := header[0].(*ContentLength)
	if !ok {
		return nil, false
	}
	return contentLength, true
}

func (h *headers) ContentType() (*ContentType, bool) {
	header := h.GetHeader("Content-Type")
	if len(header) == 0 {
		return nil, false
	}
	contentType, ok := header[0].(*ContentType)
	if !ok {
		return nil, false
	}
	return contentType, true
}

func (h *headers) Contact() (*Contact, bool) {
	header := h.GetHeader("Contact")
	if len(header) == 0 {
		return nil, false
	}
	contact, ok := header[0].(*Contact)
	if !ok {
		return nil, false
	}
	return contact, true
}
