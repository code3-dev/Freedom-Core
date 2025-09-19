//go:build windows

package sysproxy

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type winProxy struct{}

func newPlatformProxy() SysProxy {
	return &winProxy{}
}

func (w *winProxy) GetProxy() (*ProxyConfig, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue("ProxyEnable")
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return nil, err
	}
	server, _, err := k.GetStringValue("ProxyServer")
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return nil, err
	}
	cfg := &ProxyConfig{Enable: val == 1}
	if server != "" {
		parts := strings.Split(server, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "http=") {
				cfg.HTTP = strings.TrimPrefix(part, "http=")
			} else if strings.HasPrefix(part, "https=") {
				cfg.HTTPS = strings.TrimPrefix(part, "https=")
			} else if strings.HasPrefix(part, "socks=") {
				cfg.SOCKS5 = strings.TrimPrefix(part, "socks=")
			} else {
				if cfg.HTTP == "" {
					cfg.HTTP = part
				}
			}
		}
	}
	return cfg, nil
}

func (w *winProxy) SetProxy(cfg *ProxyConfig) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if cfg.Enable {
		server := []string{}
		if cfg.HTTP != "" {
			server = append(server, cfg.HTTP)
		}
		if cfg.HTTPS != "" {
			server = append(server, cfg.HTTPS)
		}
		if cfg.SOCKS5 != "" {
			server = append(server, cfg.SOCKS5)
		}
		if len(server) > 0 {
			if err := k.SetStringValue("ProxyServer", strings.Join(server, ";")); err != nil {
				return err
			}
		}
		if err := k.SetDWordValue("ProxyEnable", 1); err != nil {
			return err
		}
	} else {
		if err := k.SetDWordValue("ProxyEnable", 0); err != nil {
			return err
		}
		_ = k.DeleteValue("ProxyServer")
	}
	return nil
}

func (w *winProxy) ClearProxy() error {
	return w.SetProxy(&ProxyConfig{Enable: false})
}
