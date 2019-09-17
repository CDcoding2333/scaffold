package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

var agent *gorequest.SuperAgent

func init() {
	agent = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	agent.Client.Timeout = time.Second * 5
}

// Dorequest ...
func Dorequest(method, url string, body interface{}, out interface{}) error {
	log.Debug(url, fmt.Sprintf("%+v", body))
	agent := agent.CustomMethod(method, url)
	if method != "GET" && body != nil {
		agent = agent.SendStruct(body)
	}
	resp, bs, errs := agent.EndStruct(out)
	if len(errs) > 0 {
		log.Error(errs)
		return errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("request error, code: ", resp.StatusCode)
		return errors.New(resp.Status)
	}
	log.Debug(string(bs))
	return nil
}
