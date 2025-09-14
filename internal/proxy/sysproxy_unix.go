//go:build linux || darwin

package sysproxy

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type unixProxy struct{}

func newPlatformProxy() SysProxy {
	return &unixProxy{}
}

func PrepareCore() (string, error) {
	if runtime.GOOS == "darwin" {
		return "proxy-core", nil
	}
	return "proxy-core", nil
}

func (u *unixProxy) GetProxy() (*ProxyConfig, error) {
	cfg := &ProxyConfig{
		HTTP:   os.Getenv("http_proxy"),
		HTTPS:  os.Getenv("https_proxy"),
		SOCKS5: os.Getenv("socks_proxy"),
		Enable: os.Getenv("http_proxy") != "" || os.Getenv("https_proxy") != "" || os.Getenv("socks_proxy") != "",
	}
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("networksetup", "-getwebproxy", "Wi-Fi").CombinedOutput()
		if err == nil {
			if bytes.Contains(out, []byte("Enabled: Yes")) {
				host := parseNetworksetupOutput(out, "Server:")
				port := parseNetworksetupOutput(out, "Port:")
				if host != "" && port != "" {
					cfg.HTTP = host + ":" + port
					cfg.Enable = true
				}
			}
		}
		out2, err := exec.Command("networksetup", "-getsocksfirewallproxy", "Wi-Fi").CombinedOutput()
		if err == nil {
			if bytes.Contains(out2, []byte("Enabled: Yes")) {
				host := parseNetworksetupOutput(out2, "Server:")
				port := parseNetworksetupOutput(out2, "Port:")
				if host != "" && port != "" {
					cfg.SOCKS5 = host + ":" + port
					cfg.Enable = true
				}
			}
		}
	}
	return cfg, nil
}

func (u *unixProxy) SetProxy(cfg *ProxyConfig) error {
	if runtime.GOOS == "darwin" {
		if cfg.HTTP != "" {
			host, port := splitHostPort(stripScheme(cfg.HTTP))
			if host == "" {
				return errors.New("invalid http address")
			}
			if err := exec.Command("networksetup", "-setwebproxy", "Wi-Fi", host, port).Run(); err != nil {
				return err
			}
			if err := exec.Command("networksetup", "-setwebproxystate", "Wi-Fi", "on").Run(); err != nil {
				return err
			}
		}
		if cfg.HTTPS != "" {
			host, port := splitHostPort(stripScheme(cfg.HTTPS))
			if host == "" {
				return errors.New("invalid https address")
			}
			if err := exec.Command("networksetup", "-setsecurewebproxy", "Wi-Fi", host, port).Run(); err != nil {
				return err
			}
			if err := exec.Command("networksetup", "-setsecurewebproxystate", "Wi-Fi", "on").Run(); err != nil {
				return err
			}
		}
		if cfg.SOCKS5 != "" {
			host, port := splitHostPort(stripScheme(cfg.SOCKS5))
			if host == "" {
				return errors.New("invalid socks address")
			}
			if err := exec.Command("networksetup", "-setsocksfirewallproxy", "Wi-Fi", host, port).Run(); err != nil {
				return err
			}
			if err := exec.Command("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "on").Run(); err != nil {
				return err
			}
		}
		return nil
	}
	if hasCommand("gsettings") {
		if cfg.HTTP != "" || cfg.HTTPS != "" || cfg.SOCKS5 != "" {
			if err := exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "manual").Run(); err != nil {
				return err
			}
		}
		if cfg.HTTP != "" {
			h, p := splitHostPort(stripScheme(cfg.HTTP))
			if h != "" {
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "host", h).Run(); err != nil {
					return err
				}
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "port", p).Run(); err != nil {
					return err
				}
			}
		}
		if cfg.HTTPS != "" {
			h, p := splitHostPort(stripScheme(cfg.HTTPS))
			if h != "" {
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "host", h).Run(); err != nil {
					return err
				}
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "port", p).Run(); err != nil {
					return err
				}
			}
		}
		if cfg.SOCKS5 != "" {
			h, p := splitHostPort(stripScheme(cfg.SOCKS5))
			if h != "" {
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "host", h).Run(); err != nil {
					return err
				}
				if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "port", p).Run(); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if cfg.HTTP != "" {
		os.Setenv("http_proxy", cfg.HTTP)
	}
	if cfg.HTTPS != "" {
		os.Setenv("https_proxy", cfg.HTTPS)
	}
	if cfg.SOCKS5 != "" {
		os.Setenv("socks_proxy", cfg.SOCKS5)
	}
	return nil
}

func (u *unixProxy) ClearProxy() error {
	if runtime.GOOS == "darwin" {
		_ = exec.Command("networksetup", "-setwebproxystate", "Wi-Fi", "off").Run()
		_ = exec.Command("networksetup", "-setsecurewebproxystate", "Wi-Fi", "off").Run()
		_ = exec.Command("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off").Run()
		return nil
	}
	if hasCommand("gsettings") {
		_ = exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
		return nil
	}
	_ = os.Unsetenv("http_proxy")
	_ = os.Unsetenv("https_proxy")
	_ = os.Unsetenv("socks_proxy")
	return nil
}

func splitHostPort(addr string) (string, string) {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i], addr[i+1:]
		}
	}
	return addr, "80"
}

func stripScheme(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "socks5://")
	s = strings.TrimPrefix(s, "socks://")
	return s
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func parseNetworksetupOutput(out []byte, key string) string {
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if strings.Contains(l, key) {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
