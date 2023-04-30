package util

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

var client = gorequest.New()

func SendPost(url string, params map[string]interface{}) (b string, err error) {
	client.Post(url).SendMap(params).
		Timeout(3 * time.Second).
		End(func(response gorequest.Response, body string, errs []error) {

			if response.StatusCode != http.StatusOK || errs != nil {
				log.Error(errs)
				err = errs[0]
				return
			}
			b = body
		})
	return
}
