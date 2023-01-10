package gb

import (
	"fmt"

	"github.com/ghettovoice/gosip/sip"
	"math/rand"
	"net/http"
	"time"
)

const (
	letterBytes      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultAlgorithm = "MD5"
	WWWHeader        = "WWW-Authenticate"
)

func RegisterHandler(req sip.Request, tx sip.ServerTransaction) {
	log.Infof("收到来自%s 的请求: \n%s\n", req.Source(), req.String())
	// 判断是否存在 Authorization 字段
	if headers := req.GetHeaders("Authorization"); len(headers) > 0 {
		// 存在 Authorization 头部字段
		authHeader := headers[0].(*sip.GenericHeader)
		log.Infof("存在Authorization头部信息，直接注册成功：\n%s\n", authHeader)

		// 发送OK信息
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "ok", ""))
	}

	// 没有存在 Authorization 头部字段
	response := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "")
	// 添加 WWW-Authenticate 头
	wwwHeader := &sip.GenericHeader{
		HeaderName: WWWHeader,
		Contents: fmt.Sprintf("Digest nonce=\"%s\", algorithm=%s, realm=\"%s\", qop=\"auth\"",
			"44010200491118000001",
			DefaultAlgorithm,
			RandString(32),
		),
	}
	response.AppendHeader(wwwHeader)
	log.Infof("没有Authorization头部信息，生成WWW-Authenticate头部返回：\n%s\n", response)
	_ = tx.Respond(response)
}

// RandString https://github.com/kpbird/golang_random_string
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	output := make([]byte, n)
	randomness := make([]byte, n)

	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)

	for pos := range output {
		random := randomness[pos]
		randomPos := random % uint8(l)
		output[pos] = letterBytes[randomPos]
	}

	return string(output)
}
