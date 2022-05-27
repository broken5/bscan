package goreq

import (
	"bscan/common/utils"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method         string
	Url            string
	Headers        map[string]string
	Body           string
	Query          string
	Timeout        int
	AllowRedirects bool
}

// Req 请求接口
func (r Request) Req() (resp *Response) {
	switch r.Method {
	case "GET":
		return r.Get()
	case "POST":
		return r.Post()
	default:
		return r.Get()
	}
}

func (r *Request) getClient() *http.Client {
	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(r.Timeout) * time.Second,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		if r.AllowRedirects == true {
			return nil
		} else {
			return http.ErrUseLastResponse
		}
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Duration(r.Timeout) * time.Second,
	}
	return client
}

// SetReqHeaders 设置请求头
func (r *Request) SetReqHeaders(req *http.Request) { // 设置请求头
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
}

// Post POST请求
func (r *Request) Post() (newresp *Response) {
	var url string
	if r.Query != "" {
		url = utils.UrlQueryFormat(r.Url, r.Query)
	} else {
		url = r.Url
	}
	client := r.getClient()
	req, err := http.NewRequest("POST", url, strings.NewReader(r.Body))
	if err != nil {
		return
	}
	r.SetReqHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body := func() []byte { // 获取HTML
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			content = []byte{}
		}
		return content
	}()
	newresp = &Response{
		Status:  int32(resp.StatusCode),
		Headers: FormatHeaders(resp.Header),
		Body:    body,
	}
	if FilterResp(newresp) == true {
		return nil
	}
	return newresp
}

// Get GET请求
func (r *Request) Get() (newresp *Response) {
	var url string
	client := r.getClient()
	if r.Query != "" {
		url = utils.UrlQueryFormat(r.Url, r.Query)
	} else {
		url = r.Url
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	r.SetReqHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body := func() []byte { // 获取HTML
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			content = []byte{}
		}
		return content
	}()
	newresp = &Response{
		Status:  int32(resp.StatusCode),
		Headers: FormatHeaders(resp.Header),
		Body:    body,
	}
	if FilterResp(newresp) == true {
		return nil
	}
	return newresp
}
