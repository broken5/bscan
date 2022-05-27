package alive

import (
	"bscan/common/goreq"
	"bscan/common/output"
	"bscan/common/targetloader"
	"bscan/config"
	connect "bscan/database"
	"fmt"
	"html"
	"os"
	"strconv"

	"github.com/remeh/sizedwaitgroup"
)

// var wg = &sync.WaitGroup{}
var (
	index        = 0
	total        = 0
	AliveWebList = make([]*goreq.AliveWeb, 0, 0)
)

func HandleResp(req goreq.Request) *goreq.AliveWeb {
	index += 1
	output.Progress(index, total)
	resp := req.Get()
	if resp == nil {
		return nil
	}
	aliveweb := &goreq.AliveWeb{
		Status:  resp.Status,
		Headers: resp.Headers,
		Length:  int32(len(resp.Body)),
		Url:     req.Url,
	}
	aliveweb.GetFingerPrint(resp)
	aliveweb.SetAWTitle(resp.Body)
	AliveWebList = append(AliveWebList, aliveweb)
	return aliveweb
}

func Write(a *goreq.AliveWeb, filename string) {
	fd, _ := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0644)
	defer fd.Close()
	if filename != "" {
		s := "</td><td>"
		url := html.EscapeString(a.GetUrl())
		url = `<a href="` + url + `" target="_blank">` + url + "</a>"
		title := html.EscapeString(a.GetTitle())
		line := "<tr><td>" + url + s + title + s + strconv.Itoa(int(a.GetStatus())) + s + strconv.Itoa(int(a.GetLength())) + s + a.GetWebserver() + s + a.GetApplication() + s + a.GetFramework() + s + a.GetOs() + s + a.GetDesc() + "</td></tr>\n"
		fd.WriteString(line)
	}
	fd.Close()
}

func SwitchMode(options *config.Config) {
	var lines []string
	const COMMON = 1 // 常规扫描模式
	const Sylas = 2  // 全自动利用Sylas中的域名扫描
	var targetMgr = targetloader.GenerateTarget(options.Ports)
	switch options.TargetMode {
	case COMMON:
		if options.Target == "" {
			output.Error("You don't enter your target file.")
			os.Exit(0)
		}
		lines = targetloader.LoadTargetFile(options.Target)
		for _, line := range lines {
			targetMgr.AddTarget(line, options.Path)
		}
		aliveScan(targetMgr, options)
	case Sylas:
		// 按照rootdomain来读取，然后分批次写入。解决问题
		var subDomainTables = []string{"SubDomain", "SimilarSubDomain"}
		connect.LoadConfig(options)
		if !connect.DbConnectStatus {
			output.Error("You don't have mysql connection set up")
			os.Exit(0)
		}
		for _, table := range subDomainTables {
			var rootDomains = connect.GetRootDomain(table)
			var subDomains []string
			for _, rootDomain := range rootDomains {
				subDomains = connect.GetSubDomain(rootDomain, table)
				if subDomains != nil {
					for _, line := range subDomains {
						targetMgr.AddTarget(line, options.Path)
					}
					aliveScan(targetMgr, options, rootDomain, table)
					for _, line := range subDomains {
						connect.UpdateScanned(line, table)
					}
				}
				targetMgr.ClearTarget()
			}
		}
	}
}

func aliveScan(targetMgr *targetloader.TargetMgr, options *config.Config, args ...string) {
	total = len(targetMgr.TargetList)

	if total == 0 {
		output.Warning("Not found vaild target.")
		os.Exit(0)
	}
	//output.PrintAliveConfig(options.Threads, len(options.Ports), options.Request.Timeout, total)
	//if options.Outfile != "" {
	//	utils.WriteHTML(options.Outfile)
	//	fd, _ := os.OpenFile(options.Outfile, os.O_APPEND|os.O_RDWR, 0644)
	//	fd.WriteString(`<table class="hovertable"  width="100%"><tr><th>URL</th><th>Title</th><th>Status</th><th>Length</th><th>WebServer</th><th>Application</th><th>Framework</th><th>OS</th><th>Desc</th></tr>`)
	//	defer fd.Close()
	//	defer fd.WriteString("</table>")
	//}
	// return
	swg := sizedwaitgroup.New(options.Threads)
	for _, target := range targetMgr.TargetList {
		swg.Add()
		req := goreq.Request{
			Method:         options.Request.Method,
			Url:            target.GetUrl(),
			Query:          options.Request.Query,
			Headers:        options.Request.Headers,
			Timeout:        options.Request.Timeout,
			AllowRedirects: options.Request.AllowRedirects,
		}
		go func() {
			defer swg.Done()
			ret := HandleResp(req)
			if ret != nil {
				output.PrintAlive(
					ret.Url,
					ret.Status,
					ret.Length,
					ret.Title,
					ret.Application,
					ret.Webserver,
					ret.Desc,
					ret.Os,
					ret.Framework,
				)
				if options.TargetMode == 2 {
					connect.InsertAliveDomainInfo(ret.Url, ret.Status, ret.Title, args[0], args[1])
					//connect.UpdateScanned(target.Domain, args[1])
				}
				//Write(ret, options.Outfile)
			}
		}()
	}
	swg.Wait()
	fmt.Println()
	//connect.Sylas.Close() //关闭数据库连接
}
