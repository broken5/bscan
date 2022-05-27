package config

type Sylas struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Db     string `yaml:"db"`
}
