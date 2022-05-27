package config

type Request struct {
	Method         string            `yaml:"method"`
	Path           []string          `yaml:"path"`
	Query          string            `yaml:"query"`
	Body           string            `yaml:"body"`
	Headers        map[string]string `yaml:"headers"`
	Timeout        int               `yaml:"timeout"`
	AllowRedirects bool              `yaml:"allow_redirects"`
}

type Poc struct {
	Name       string  `yaml:"name"`
	Request    Request `yaml:"request"`
	FilterExpr string  `yaml:"filter_expr"`
	VerifyExpr string  `yaml:"verify_expr"`
}
