package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/log"
	"io"
	"net"
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

	readState struct {
		//  isHeader isBody two field，at the same time，only one field will be true
		isHeader bool
		isBody   bool

		isHaveBody bool
		bodyLen    int

		t int8
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
		packet packet
	)

	for !p.isStop {
		packet = <-p.in
		line, err := packet.getLine()
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
				continue
			}

			log.Log.Debugf("method: %v, URI: %v, sipVersion: %v", method, uri, sipVersion)

		} else {
			log.Log.Info("that line is not SIP request line")
		}
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
func (p *packet) getLine() (string, error) {
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

// parser the request line of SIP protocol,
// return method and URI and sip version in request line
func parseRequestLine(startLine string) (method requestMethod, uri URI, sipVersion string, err error) {
	// SIP Request example :
	// REGISTER sip:44010200492000000001@4401020049 SIP/2.0
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("SIP request line should have 2 spaces,but that line not have: [%s]", startLine)
		return
	}
	method = requestMethod(strings.ToUpper(parts[0]))
	uri, err = parseURI(parts[1])

	sipVersion = parts[2]
	return
}

// parse request line of sip message, return URI object
func parseURI(startLine string) (u URI, err error) {

	index := strings.Index(startLine, ":")
	if index == -1 {
		err = fmt.Errorf("no contain ':' in %s", startLine)
		return
	}

	switch strings.ToLower(startLine[:index]) {
	case "sip", "sips":
		u, err = parseSipURI(startLine)
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
	u.FHost = String{Str: startLine[index+1:]}
	return
}
