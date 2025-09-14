package sysproxy

type ProxyConfig struct {
	HTTP   string
	HTTPS  string
	SOCKS5 string
	Enable bool
}

type SysProxy interface {
	GetProxy() (*ProxyConfig, error)
	SetProxy(cfg *ProxyConfig) error
	ClearProxy() error
}

func NewSysProxy() SysProxy {
	return newPlatformProxy()
}
