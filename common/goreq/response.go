package goreq

import (
	"bscan/config"
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	reflect "reflect"
	"regexp"
	"strings"

	"github.com/google/cel-go/common/types"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	Charsets = []string{"utf-8", "gbk", "gb2312"}
)

func (r *AliveWeb) SetAWTitle(body []byte) {
	text := string(body)

	Decodegbk := func(s []byte) ([]byte, error) { // GBK解码
		I := bytes.NewReader(s)
		O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
		d, e := ioutil.ReadAll(O)
		if e != nil {
			return nil, e
		}
		return d, nil
	}

	getEncoding := func() string { // 判断Content-Type
		r1, err := regexp.Compile(`(?im)charset=\s*?([\w-]+)`)
		r2, err := regexp.Compile(`(?im)<meta.*?charset=['"]?([\w-]+)["']?.*?>`)
		if err != nil {
			return ""
		}
		headerCharset := r1.FindString(r.Headers["Content-Type"])
		htmlCharset := r2.FindString(text)
		for _, v := range Charsets { // headers 编码优先，所以放在前面
			if headerCharset != "" && strings.Contains(strings.ToLower(headerCharset), v) == true {
				return v
			}
		}
		for _, v := range Charsets {
			if htmlCharset != "" && strings.Contains(strings.ToLower(htmlCharset), v) == true {
				return v
			}
		}
		return ""
	}
	r.Title = func() (title string) { // 设置Title
		var re = regexp.MustCompile(`(?im)<title.*?>([\s\S]*?)</title>`)
		matchs := re.FindStringSubmatch(text)
		if len(matchs) == 2 {
			title = strings.Trim(html.UnescapeString(matchs[1]), "\r\n\t ")
		} else {
			title = ""
		}
		enc := getEncoding()
		if enc == "utf-8" || enc == "" {
			return title
		}

		if enc == "gbk" || enc == "gb2312" {
			titleGBK, err := Decodegbk([]byte(title))
			if err != nil {
				return
			}
			return string(titleGBK)
		}
		return
	}()
	r.Title = strings.Trim(r.Title, "\r\n \t")
	r.Title = strings.Replace(r.Title, "\n", "", -1)
	r.Title = strings.Replace(r.Title, "\r", "", -1)
}

func FormatHeaders(respHeader http.Header) map[string]string {
	headers := make(map[string]string)
	for k, v := range respHeader {
		headers[k] = strings.Join(v, ";")
	}
	return headers
}

func GetBoolResult(v interface{}) bool {
	switch v.(type) {
	case types.Bool:
		return reflect.ValueOf(v).Bool()
	default:
		return false
	}
}

func (a *AliveWeb) GetFingerPrint(resp *Response) {
	exec := ExecuteExpr(map[string]interface{}{"response": resp})
	if exec == nil {
		return
	}
	if config.FingerPrintData == nil {
		return
	}
	fps := make(map[string][]string)
	for _, fp := range *config.FingerPrintData {
		expr := fp.Expression
		out, err := exec(expr)
		if err != nil {
			continue
		}
		ret := GetBoolResult(out)
		if ret == false {
			continue
		}
		if fp.WebServer != "" {
			fps["WebServer"] = append(fps["WebServer"], fp.WebServer)
		}
		if fp.Application != "" {
			fps["Application"] = append(fps["Application"], fp.Application)
		}
		if fp.Framework != "" {
			fps["Framework"] = append(fps["Framework"], fp.Framework)
		}
		if fp.Os != "" {
			fps["Os"] = append(fps["Os"], fp.Os)
		}
		if fp.Desc != "" {
			fps["Desc"] = append(fps["Desc"], fp.Desc)
		}
	}
	a.Webserver = strings.Join(fps["WebServer"], ",")
	a.Application = strings.Join(fps["Application"], ",")
	a.Os = strings.Join(fps["Os"], ",")
	a.Framework = strings.Join(fps["Framework"], ",")
	a.Desc = strings.Join(fps["Desc"], ",")
}

func FilterResp(resp *Response) bool {
	if resp == nil {
		return true
	}
	exec := ExecuteExpr(map[string]interface{}{"response": resp})
	if exec == nil {
		return false
	}
	if config.BlackListData == nil {
		return false
	}
	for _, bl := range *config.BlackListData {
		expr := bl.Expression
		out, err := exec(expr)
		if err != nil {
			continue
		}
		ret := GetBoolResult(out)
		if ret == true {
			return true
		}
	}
	return false
}
