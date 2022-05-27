package poc

import (
	"bscan/common/goreq"
	"bscan/common/output"
	"bscan/common/utils"
	"bscan/config"
	"fmt"
	"os"

	"github.com/remeh/sizedwaitgroup"
)

func checkPocResp(resp *goreq.Response, expr string) bool {
	exec := goreq.ExecuteExpr(map[string]interface{}{"response": resp})
	if exec == nil {
		return false
	}
	out, err := exec(expr)
	if err != nil {
		return false
	}
	ret := goreq.GetBoolResult(out)
	if ret == true {
		return true
	}
	return false
}

func PocScan(options *config.Config, AliveWebList []*goreq.AliveWeb) {
	pocs := config.LoadPocs(options.PocsPath)
	if len(pocs) == 0 {
		return
	}
	if options.Outfile != "" {
		utils.WriteHTML(options.Outfile)
		fd, _ := os.OpenFile(options.Outfile, os.O_APPEND|os.O_RDWR, 0644)
		fd.WriteString(`<table class="hovertable"  width="100%"><tr><th>URL</th><th>Title</th><th>Status</th><th>Length</th><th>WebServer</th><th>Application</th><th>Framework</th><th>OS</th><th>Desc</th></tr>`)
		defer fd.Close()
		defer fd.WriteString("</table>")
	}
	output.PrintPocConfig(options.Threads, len(pocs))
	swg := sizedwaitgroup.New(options.Threads)
	for _, poc := range pocs {
		for _, aliveweb := range AliveWebList {
			for _, pocPath := range poc.Request.Path {
				swg.Add()
				url := utils.UrlFormat(aliveweb.Url, pocPath)
				req := goreq.Request{
					Method:         poc.Request.Method,
					Url:            url,
					Query:          poc.Request.Query,
					Headers:        poc.Request.Headers,
					Body:           poc.Request.Body,
					Timeout:        poc.Request.Timeout,
					AllowRedirects: poc.Request.AllowRedirects,
				}
				expr := poc.VerifyExpr
				name := poc.Name
				go func() {
					defer swg.Done()
					resp := req.Req()
					if checkPocResp(resp, expr) == true {
						output.PrintFound(req.Url, name, resp.GetStatus(), len(resp.GetBody()))

					}
				}()
			}
		}
	}
	swg.Wait()
	fmt.Println()
}
