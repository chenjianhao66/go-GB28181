package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/log"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	BufferSizeMax = 1<<16 - 20 - 8

	// SIP协议规定的字符
	tab = " \t"
)

type (
	parser struct {
		out    chan any
		in     chan packet
		isStop bool
	}

	packet struct {
		reader     *bufio.Reader
		addr       net.Addr
		bodyLength int
	}

	HeaderParse func(headerName, headerData string) ([]Header, error)
)

var (
	// 每一个请求头的解析函数对象
	sipMessageHeaderParse = map[string]HeaderParse{
		"Via":            parseVia,
		"From":           parseFrom,
		"To":             parseTo,
		"Call-ID":        parseCallID,
		"CSeq":           parseCSeq,
		"Contact":        parseContact,
		"Expires":        parseExpires,
		"Content-Length": parseContentLength,
		"Content-Type":   parseContentType,
		"User-Agent":     parseUserAgent,
		"Max-Forwards":   parseMaxForwards,
	}
)

func newParser() *parser {
	p := &parser{
		out:    make(chan any),
		in:     make(chan packet),
		isStop: false,
	}
	go p.parserPacket()
	return p
}

func (p *parser) parserPacket() {
	var (
		packet    packet
		msg       Message
		parserErr error
	)

	for !p.isStop {
		packet = <-p.in
		line, err := packet.nextLine()
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Log.Errorf("get SIP message request Line fail,\n%s", err)
		}
		// parser that line
		if isRequest(line) {
			log.Log.Info("that line is SIP request line")
			method, uri, sipVersion, err := parseRequestLine(line)
			if err != nil {
				log.Log.Error(err)
				parserErr = err
			}
			log.Log.Debugf("method: %v, URI: %s, sipVersion: %v", method, uri.String(), sipVersion)
			msg = NewRequest("", method, uri, sipVersion, []Header{}, []byte{})

		} else if isResponse(line) {
			log.Log.Debugf("this line si SIP response line")
		} else {
			log.Log.Info("that line is not SIP request line")
			continue
		}

		// if find error,that continue
		if parserErr != nil {
			continue
		}

		// parse header of message
		var buf bytes.Buffer
		headers := make([]Header, 0)
		// this flushBuffer function be use for reset data of buf variable
		flushBuffer := func() {
			if buf.Len() > 0 {
				newHeader, err := parseHeader(buf.String())
				if err != nil {
					log.Log.Errorf("parse header of message fail in %s", buf.String())
				} else {
					headers = append(headers, newHeader...)
				}
				buf.Reset()
			}
		}

		// 读取请求行或者状态行之后就开始读取请求头
		for {
			data, err := packet.nextLine()
			if err != nil {
				log.Log.Errorf("call nextLine function of packet variable fail,%s:", err)
				break
			}
			if len(data) == 0 {
				// if length of data == 0,parse header part of message over
				// call flushBuffer function take append data of buf variable to headers slice
				// 如果data=0,那么就代表消息头部分已经读取完毕了
				// 因为sip协议规定消息头和消息体之间隔着 '\r\n'
				flushBuffer()
				break
			}
			// if data not contain '/t',that call flushBuffer function and take data write to the buf variable
			if !strings.Contains(tab, string(data[0])) {
				// This line starts a new header.
				// Parse anything currently in the buffer, then store the new header line in the buffer.
				flushBuffer()
				buf.WriteString(data)
			} else if buf.Len() > 0 {
				// This is a continuation line, so just add it to the buf.
				// 这是一个连续行
				buf.WriteString(" ")
				buf.WriteString(line)
			}
		}

		// take add headers to the msg variable
		// 将headers切片中的struct添加到消息对象中
		for _, header := range headers {
			msg.AppendHeader(header)
		}
		// if content-length length of msg object >= field bodyLength of packet, bodyLength of packet = content-length of msg
		// 如果msg对象的content-length消息头对象长度 >= packet对象的bodyLength，那么赋值
		if length, ok := msg.ContentLength(); ok {
			if int(*length) > packet.bodyLength {
				packet.bodyLength = int(*length)
			}
		}

		body, err := packet.getBody()
		if err != nil {
			continue
		}
		if len(body) != 0 {
			// 代表消息体内有值，
			msg.SetBody(body, false)
		}
		msg.SetSource(packet.addr)
	}
}

func (p *parser) stop() {
	p.isStop = true
}

func newPacket(buf []byte, remoteAddr net.Addr) packet {
	return packet{
		reader:     bufio.NewReader(bytes.NewReader(buf)),
		addr:       remoteAddr,
		bodyLength: getMessageBodyLength(buf),
	}
}

// get one line string of packet object
func (p *packet) nextLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// if size of line == 2, explain the line string is empty
	// because '\r\n' == two char, that line only contains '\r\n'
	if len(line) >= 2 {
		line = line[:len(line)-2]
	}
	return line, nil
}

// 根据packet对象的bodyLength字段获取消息体
func (p *packet) getBody() ([]byte, error) {
	if p.bodyLength < 1 {
		return []byte{}, nil
	}
	body := make([]byte, p.bodyLength)
	if p.bodyLength > 0 {
		n, err := io.ReadFull(p.reader, body)
		if err != nil && err != io.ErrUnexpectedEOF {
			return body, err
		}
		if n != p.bodyLength {
			logrus.Warningf("body length err,%d!=%d,body:%s", n, p.bodyLength, string(body))
			return body[:n], nil
		}
	}
	return body, nil
}

// get size of SIP message body
func getMessageBodyLength(data []byte) int {
	// the Header and body of SIP message are separated by '/r/n'
	// so use Index function of strings package, take data as Params get '\r\n' start index
	// again use start index subtract length of data , get the final result
	s := string(data)
	index := strings.Index(s, "\r\n\r\n")
	if index == -1 {
		return -1
	}
	bodyStartIndex := index + 4
	return len(s) - bodyStartIndex
}

// check whether that line is SIP request line
func isRequest(line string) bool {
	if strings.Count(line, " ") != 2 {
		return false
	}
	// Check that the version string starts with SIP.
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return false
	} else if len(parts[2]) < 3 {
		return false
	} else {
		return strings.ToUpper(parts[2][:3]) == "SIP"
	}
}

// check whether that line is SIP response line
func isResponse(line string) bool {
	return false
}

// parser the request line of SIP protocol,
// return method and URI and sip version in request line
func parseRequestLine(startLine string) (method RequestMethod, uri *URI, sipVersion string, err error) {
	// SIP Request example :
	// REGISTER sip:44010200492000000001@4401020049 SIP/2.0
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("SIP request line should have 2 spaces,but that line not have: [%s]", startLine)
		return
	}
	method = RequestMethod(strings.ToUpper(parts[0]))
	uri, err = parseURI(parts[1])

	sipVersion = parts[2]
	return
}

// parse request line of sip message, return URI object
func parseURI(startLine string) (u *URI, err error) {

	index := strings.Index(startLine, ":")
	if index == -1 {
		err = fmt.Errorf("no contain ':' in %s", startLine)
		return
	}

	switch strings.ToLower(startLine[:index]) {
	case "sip", "sips":
		var uri URI
		uri, err = parseSipURI(startLine)
		u = &uri
	default:
		log.Log.Debugf("startLine[:index]: %s", startLine[:index])
		err = fmt.Errorf("parse URI fail in %s", startLine)
	}

	return
}

func parseSipURI(startLine string) (u URI, err error) {

	if strings.ToLower(startLine[:3]) != "sip" {
		err = fmt.Errorf("invalid SIP URI protocol Name in '%s'", startLine)
		return
	}
	startLine = startLine[4:]

	index := strings.Index(startLine, "@")
	if index == -1 {
		err = fmt.Errorf("parse sip URI fail in %s", startLine)
		return
	}

	u.FUser = String{Str: startLine[:index]}
	u.FHost = startLine[index+1:]
	return
}

func parseHostPort(line string) (host string, port *Port) {
	hostAndPort := strings.Split(line, ":")

	if !strings.ContainsAny(hostAndPort[0], ".") || len(hostAndPort) != 2 {
		log.Log.Errorf("主机部分不是合法ip地址: %s", line)
		return "", nil
	}
	portUint16, _ := strconv.ParseUint(hostAndPort[1], 10, 16)
	var p uint16
	p = uint16(portUint16)
	host = hostAndPort[0]
	port = (*Port)(&p)
	return
}

func parseOtherParams(paramsPart []string) (Params, error) {
	//if len(paramsPart) == 0 {
	//	return nil,nil
	//}
	params := NewParams()
	for _, v := range paramsPart {
		if strings.Contains(v, "=") {
			keyValue := strings.Split(v, "=")
			params.Add(keyValue[0], String{Str: keyValue[1]})
		} else {
			params.Add(v, String{Str: ""})
		}
	}
	return params, nil
}

func parseHeader(headerData string) (headers []Header, err error) {
	headers = make([]Header, 0)

	index := strings.Index(headerData, ":")
	if index == -1 {
		log.Log.Debugf("not contains ':' in %s", headers)
		return headers, fmt.Errorf("not contains ':' in %s", headers)
	}
	// get headerName and value in headerData
	fieldName := strings.TrimSpace(headerData[:index])
	value := strings.TrimSpace(headerData[index+1:])

	if headerParse, ok := sipMessageHeaderParse[fieldName]; ok {
		headers, err = headerParse(fieldName, value)
	} else {

		log.Log.Infof("parse of temporary un support other header:%s", fieldName)
	}

	return headers, err
}

func parseVia(headerName, headerData string) ([]Header, error) {
	var viaHeader = ViaHeader{}
	viaParams := new(viaParams)
	spaceIndex := strings.Index(headerData, " ")
	if spaceIndex == -1 {
		return nil, fmt.Errorf("via Header colud contains one space,but not: %s", headerData)
	}

	split := strings.Split(headerData, " ")
	headerPart := strings.Split(split[0], "/")
	if len(headerPart) != 3 {
		return nil, fmt.Errorf("header parr of via header colud contains portocolName and version and transport,bu not :%s", headerPart)
	}
	viaParams.protocolName, viaParams.protocolVersion, viaParams.transport = headerPart[0], headerPart[1], headerPart[2]
	headerData = headerData[spaceIndex+1:]

	paramsPart := strings.Split(headerData, ";")
	viaParams.host, viaParams.port = parseHostPort(paramsPart[0])

	if len(paramsPart) == 1 {
		viaHeader = append(viaHeader, viaParams)
		return []Header{viaHeader}, nil
	}
	paramsPart = paramsPart[1:]
	params, err := parseOtherParams(paramsPart)
	if err != nil {
		return nil, err
	}
	viaParams.Params = params
	viaHeader = append(viaHeader, viaParams)
	//log.Log.Debugf("after parse via header: %s", viaHeader)
	return []Header{viaHeader}, nil
}

func parseFrom(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "From")

	from := new(From)
	part := strings.Split(headerData, ";")
	uriStr := part[0]
	uriStr = uriStr[1:]
	uriStr = uriStr[:len(uriStr)-1]
	uri, err := parseURI(uriStr)
	if err != nil {
		return nil, err
	}
	from.Address = uri
	part = part[1:]
	params, err := parseOtherParams(part)
	if err != nil {
		return nil, err
	}
	from.Params = params
	//log.Log.Debugf("after parse from header: %s", from)
	return []Header{from}, nil
}

func parseTo(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "To")

	to := new(To)
	part := strings.Split(headerData, ";")
	uriStr := part[0]
	uriStr = uriStr[1:]
	uriStr = uriStr[:len(uriStr)-1]
	uri, err := parseURI(uriStr)
	if err != nil {
		return nil, err
	}
	to.Address = uri
	part = part[1:]
	params, err := parseOtherParams(part)
	if err != nil {
		return nil, err
	}
	to.Params = params
	//log.Log.Debugf("after parse to header: %s", to)
	return []Header{to}, nil
}

func parseCallID(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "Call-ID")
	headerData = strings.TrimSpace(headerData)
	var callID = CallId(headerData)

	if strings.ContainsAny(string(callID), tab) {
		err := fmt.Errorf("unexpected whitespace in CallID header body '%s'", headerData)
		return nil, err
	}
	if strings.Contains(string(callID), ";") {
		err := fmt.Errorf("unexpected semicolon in CallID header body '%s'", headerData)
		return nil, err
	}
	if len(string(callID)) == 0 {
		err := fmt.Errorf("empty Call-ID body")
		return nil, err
	}

	headers := []Header{&callID}
	return headers, nil
}

func parseCSeq(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "CSeq")
	cSeq := new(CSeq)
	part := strings.Split(headerData, " ")
	if len(part) != 2 {
		return nil, fmt.Errorf("CSeq header should be contains one spaces,but not  :%s", headerData)
	}

	p, err := strconv.ParseUint(part[0], 10, 32)
	if err != nil {
		return nil, err
	}
	cSeq.SeqNumber, cSeq.MethodName = uint32(p), RequestMethod(strings.TrimSpace(part[1]))
	//log.Log.Debugf("after parse CSeq header: %s", cSeq)
	return []Header{cSeq}, nil
}

func parseContact(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "Contact")
	c := new(Contact)
	part := strings.Split(headerData, ";")
	uriStr := part[0]
	uriStr = uriStr[1:]
	uriStr = uriStr[:len(uriStr)-1]
	uri, err := parseURI(uriStr)
	if err != nil {
		return nil, err
	}
	c.Address = uri
	part = part[1:]
	params, err := parseOtherParams(part)
	if err != nil {
		return nil, err
	}
	c.Params = params
	//log.Log.Debugf("after parse contact header: %s", c)
	return []Header{c}, nil
}

func parseExpires(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "Expires")
	var e Expires
	var value uint64
	value, err := strconv.ParseUint(strings.TrimSpace(headerData), 10, 32)
	if err != nil {
		return nil, err
	}
	e = Expires(value)

	return []Header{&e}, nil
}

func parseContentLength(headerName, headerData string) ([]Header, error) {
	//log.Log.Debugf("this is %s paraser!!", "Content-Length")
	var contentLength ContentLength
	var value uint64
	value, err := strconv.ParseUint(strings.TrimSpace(headerData), 10, 32)
	if err != nil {
		return nil, err
	}
	contentLength = ContentLength(value)
	return []Header{&contentLength}, nil
}

func parseContentType(headerName string, headerData string) ([]Header, error) {
	var contentType ContentType
	headerData = strings.TrimSpace(headerData)
	contentType = ContentType(headerData)
	return []Header{&contentType}, nil
}

func parseUserAgent(headerName string, headerData string) ([]Header, error) {
	var userAgent UserAgent
	headerData = strings.TrimSpace(headerData)
	userAgent = UserAgent(headerData)
	return []Header{&userAgent}, nil
}

func parseMaxForwards(headerName string, headerText string) ([]Header, error) {
	var maxForwards MaxForwards
	var value uint64
	value, err := strconv.ParseUint(strings.TrimSpace(headerText), 10, 32)
	if err != nil {
		return nil, err
	}
	maxForwards = MaxForwards(value)

	return []Header{&maxForwards}, nil
}
