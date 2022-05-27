package targetloader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Target struct {
	Scheme string
	Domain string
	Port   int
	Path   string
	Query  string
}

type TargetMgr struct {
	TargetList []*Target
	Ports      []int
}

func newTarget(Scheme string, Domain string, Port int, Path string, Query string) *Target {
	// if Scheme == "" && strings.HasSuffix(strconv.Itoa(Port), "443") {
	// 	return &Target{}
	// }
	return &Target{
		Scheme: Scheme,
		Domain: Domain,
		Port:   Port,
		Path:   Path,
		Query:  Query,
	}
}

func (t *Target) GetUrl() string {
	t_url := t.Scheme + "://" + t.Domain + ":" + strconv.Itoa(t.Port) + "/"
	if t.Path != "" {
		t_url = t_url + strings.TrimLeft(t.Path, "/")
	}
	if t.Query != "" {
		t_url = t_url + "?" + t.Query
	}
	return t_url
}

func CheckDomainVaild(domain string) bool {
	ret, err := regexp.MatchString(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`, domain)
	if err != nil {
		return false
	}
	return ret
}

func CheckIPVaild(ip string) bool {
	ret, err := regexp.MatchString(`^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`, ip)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// fmt.Println(ip, ret)
	return ret
}

func (mgr *TargetMgr) AddTarget(line string, path string) {
	// URL格式，直接添加
	if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
		fmt.Println(line)
		return
	}

	// 分割出Domain, Port
	domain, port := ParseHost(line)
	if domain == "" { // 过滤
		return
	}

	if port == 0 { // 没有指定端口的情况
		for _, p := range mgr.Ports {
			switch true {
			case strings.HasSuffix(strconv.Itoa(p), "443") || p == 443:
				mgr.TargetList = append(mgr.TargetList, newTarget("https", domain, p, path, ""))
			case p == 80:
				mgr.TargetList = append(mgr.TargetList, newTarget("http", domain, p, path, ""))
			default:
				mgr.TargetList = append(mgr.TargetList, newTarget("http", domain, p, path, ""))
				mgr.TargetList = append(mgr.TargetList, newTarget("https", domain, p, path, ""))
			}
		}
		return
	} else { // 指定了端口
		switch true {
		case strings.HasSuffix(strconv.Itoa(port), "443") || port == 443:
			mgr.TargetList = append(mgr.TargetList, newTarget("https", domain, port, path, ""))
		case port == 80:
			mgr.TargetList = append(mgr.TargetList, newTarget("http", domain, port, path, ""))
		default:
			mgr.TargetList = append(mgr.TargetList, newTarget("http", domain, port, path, ""))
			mgr.TargetList = append(mgr.TargetList, newTarget("https", domain, port, path, ""))
		}
		return
	}
}

func (mgr *TargetMgr) ClearTarget() {
	mgr.TargetList = make([]*Target, 0)
}

func ParseHost(host string) (string, int) {
	var (
		domain string
		port   int
	)
	hostSplit := strings.Split(host, ":")
	if len(hostSplit) > 2 {
		return "", 0
	}
	domain = hostSplit[0]
	if !CheckDomainVaild(domain) && !CheckIPVaild(domain) {
		return "", 0
	}
	if len(hostSplit) == 2 {
		portInt, err := strconv.Atoi(hostSplit[1])
		if err != nil {
			return "", 0
		}
		port = portInt
	}
	if port > 65535 || port < 0 {
		port = 0
	}
	return domain, port
}

func LoadTargetFile(filename string) []string {
	var lines = make([]string, 0)
	fileObj, err := os.Open(filename)
	if err != nil {
		return lines
	}
	reader := bufio.NewReader(fileObj)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return lines
		}
		lines = append(lines, strings.Trim(string(line), "\r\n "))
	}
}
func GenerateTarget(ports []int) *TargetMgr {
	var (
		//lines     = LoadTargetFile(filename)
		targetMgr = &TargetMgr{TargetList: make([]*Target, 0), Ports: ports}
	)
	return targetMgr
}

/*func GenerateTarget(lines []string, ports []int, path string) *TargetMgr {
	var (
		//lines     = LoadTargetFile(filename)
		targetMgr = &TargetMgr{TargetList: make([]*Target, 0), Ports: ports}
	)
	for _, line := range lines {
		targetMgr.AddTarget(line, path)
	}
	return targetMgr
}
func BscanGenerateTarget(line string, ports []int, path string) *TargetMgr {
	var (
		targetMgr = &TargetMgr{TargetList: make([]*Target, 0), Ports: ports}
	)
	targetMgr.AddTarget(line, path)
	return targetMgr
}
*/
