package config

import (
	"bscan/common/customport"
	"bscan/common/output"
	"bscan/common/utils"
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type FingerPrint struct {
	WebServer   string `yaml:"webserver"`
	Application string `yaml:"application"`
	Framework   string `yaml:"framework"`
	Os          string `yaml:"os"`
	Desc        string `yaml:"desc"`
	Expression  string `yaml:"expression"`
}

type BlackList struct {
	name       string `yaml:"name"`
	Expression string `yaml:"expression"`
}

type Config struct {
	Threads         int              `yaml:"threads"`
	Ports           customport.Ports `yaml:"ports"`
	Path            string           `yaml:"path"`
	Request         Request          `yaml:"request"`
	FingerPrintPath string           `yaml:"fingerprint_path"`
	BlackListPath   string           `yaml:"blacklist_path"`
	PocsPath        string           `yaml:"pocs_path"`
	Target          string           `yaml:"target"`
	Outfile         string           `yaml:"outfile"`
	TargetMode      int              `yaml:"target_mode"`
	Sylas           Sylas            `yaml:"sylas"`
}

var (
	FingerPrintData = &[]FingerPrint{}
	BlackListData   = &[]BlackList{}
)

func init() { // 加载配置文件
	if !utils.FileExists("config.yml") { // 如果文件不存在，则新建默认配置
		output.Warning("config.yaml file does not exist and will be created automatically.")
		staticConfigExport()
	}
}

func loadConfigFromYaml() *Config {
	config := &Config{}
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
	// fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
	return config
}

func loadFingerPrint(path string) (data *[]FingerPrint) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		output.Warning(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		output.Warning(err.Error())
		return
	}
	return
}

func loadBlackList(path string) (data *[]BlackList) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		output.Warning(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		output.Warning(err.Error())
		return
	}
	return
}

func staticConfigExport() {
	request := Request{
		Timeout:        3,
		AllowRedirects: true,
		Headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
			"Accept":          "*/*",
			"Accept-Language": "en",
			"Connection":      "close",
		},
	}
	targetMode := 1
	Sylas := Sylas{
		User:   "root",
		Passwd: "123456",
		Host:   "127.0.0.1",
		Port:   "3306",
		Db:     "Sylas",
	}
	config := &Config{
		Threads:         50,
		Ports:           []int{80, 443},
		Request:         request,
		FingerPrintPath: "fingerprint.yml",
		BlackListPath:   "blacklist.yml",
		PocsPath:        "pocs",
		TargetMode:      targetMode,
		Sylas:           Sylas,
	}
	out, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	ioutil.WriteFile("config.yml", out, 0644)
}

func LoadPocs(path string) []Poc {
	data := make([]Poc, 0, 0)
	loadPoc := func(file string) (poc Poc) {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(yamlFile, &poc)
		if err != nil {
			return
		}
		return
	}
	fileList := utils.ListAllFileByName("yml", path)
	for _, v := range fileList {
		poc := loadPoc(v)
		if len(poc.Request.Path) > 0 {
			data = append(data, poc)
		}
	}
	return data
}

func ParseOptions() *Config {
	options := loadConfigFromYaml()
	FingerPrintData = loadFingerPrint(options.FingerPrintPath)
	BlackListData = loadBlackList(options.BlackListPath)
	flag.Var(&options.Ports, "ports", "ports range (nmap syntax: eg 1,2-10,11)")
	flag.IntVar(&options.Threads, "threads", options.Threads, "number of threads")
	flag.IntVar(&options.Request.Timeout, "timeout", options.Request.Timeout, "timeout of request")
	flag.StringVar(&options.Target, "target", "", "targets")
	flag.StringVar(&options.Outfile, "o", "", "save results")
	flag.StringVar(&options.Path, "path", "/", "targets")
	flag.BoolVar(&options.Request.AllowRedirects, "allow-redirects", options.Request.AllowRedirects, "http allow redirects")
	flag.Parse()
	return options
}
