package sysdns

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	helpers "github.com/Freedom-Guard/freedom-core/pkg/utils"
)

type DNSManager struct {
	cfg DNSConfig
}

func NewDNSManager() *DNSManager {
	return &DNSManager{}
}

func (d *DNSManager) SetDNS(cfg *DNSConfig) error {
	if cfg.Primary == "" {
		return errors.New("primary DNS is required")
	}
	d.cfg = *cfg
	logger.Log(logger.DEBUG, fmt.Sprintf("Setting DNS: Primary=%s Secondary=%s", cfg.Primary, cfg.Secondary))

	switch runtime.GOOS {
	case "windows":
		if !isAdmin() {
			return runAsAdmin()
		}
		return setDNSWindows(cfg)
	case "linux":
		return setDNSLinux(cfg)
	case "darwin":
		return setDNSMac(cfg)
	default:
		return errors.New("unsupported platform")
	}
}

func (d *DNSManager) GetDNS() (DNSConfig, error) {
	if d.cfg.Primary == "" {
		return DNSConfig{}, errors.New("no DNS configured")
	}
	return d.cfg, nil
}

func (d *DNSManager) ClearDNS() error {
	d.cfg = DNSConfig{}
	switch runtime.GOOS {
	case "windows":
		if !isAdmin() {
			return runAsAdmin()
		}
		return setDNSWindows(&DNSConfig{Primary: "", Secondary: ""})
	case "linux":
		return setDNSLinux(&DNSConfig{Primary: "", Secondary: ""})
	case "darwin":
		return setDNSMac(&DNSConfig{Primary: "", Secondary: ""})
	}
	return nil
}

// ---------------- Windows Helper ----------------
func isAdmin() bool {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func runAsAdmin() error {
	helpers.ShowInfo(
		"Running as Administrator",
		"The program is running with administrator privileges. Please retry the DNS registration request.",
	)

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	psCmd := fmt.Sprintf("Start-Process -FilePath '%s' -Verb RunAs -WindowStyle Normal", exe)
	cmd := exec.Command("powershell", "-Command", psCmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func setDNSWindows(cfg *DNSConfig) error {
	out, err := exec.Command("powershell", "-Command", `
Get-NetAdapter | ForEach-Object {
    $name = $_.Name
    if ('`+cfg.Primary+`' -ne '') {
        Set-DnsClientServerAddress -InterfaceAlias $name -ServerAddresses @('`+cfg.Primary+`','`+cfg.Secondary+`') -ErrorAction SilentlyContinue
    } else {
        Set-DnsClientServerAddress -InterfaceAlias $name -ResetServerAddresses -ErrorAction SilentlyContinue
    }
}
`).CombinedOutput()
	logger.Log(logger.DEBUG, "Windows DNS output: "+string(out))
	if err != nil {
		return errors.New("failed to set DNS on Windows: " + err.Error())
	}
	logger.Log(logger.DEBUG, fmt.Sprintf("System DNS successfully updated. Primary=%s Secondary=%s", cfg.Primary, cfg.Secondary))
	return nil
}

// ---------------- Linux Helper ----------------

func setDNSLinux(cfg *DNSConfig) error {
	var content string
	if cfg.Primary != "" {
		content = "nameserver " + cfg.Primary + "\n"
		if cfg.Secondary != "" {
			content += "nameserver " + cfg.Secondary + "\n"
		}
	}
	cmd := exec.Command("sh", "-c", "echo '"+content+"' > /etc/resolv.conf")
	out, err := cmd.CombinedOutput()
	logger.Log(logger.DEBUG, "Linux DNS output: "+string(out))
	if err != nil {
		return errors.New("failed to set DNS on Linux: " + err.Error())
	}
	logger.Log(logger.DEBUG, fmt.Sprintf("System DNS successfully updated. Primary=%s Secondary=%s", cfg.Primary, cfg.Secondary))
	return nil
}

// ---------------- macOS Helper ----------------

func setDNSMac(cfg *DNSConfig) error {
	ifacesCmd := exec.Command("networksetup", "-listallnetworkservices")
	ifaces, errList := ifacesCmd.CombinedOutput()
	if errList != nil {
		logger.Log(logger.DEBUG, "Mac list interfaces failed: "+string(ifaces))
		return errList
	}
	for _, iface := range strings.Split(string(ifaces), "\n") {
		if iface = strings.TrimSpace(iface); iface == "" || strings.HasPrefix(iface, "*") {
			continue
		}
		args := []string{"-setdnsservers", iface}
		if cfg.Primary != "" {
			args = append(args, cfg.Primary)
			if cfg.Secondary != "" {
				args = append(args, cfg.Secondary)
			}
		} else {
			args = append(args, "Empty")
		}
		c := exec.Command("networksetup", args...)
		out, err := c.CombinedOutput()
		logger.Log(logger.DEBUG, "Mac DNS output for "+iface+": "+string(out))
		if err != nil {
			return errors.New("failed to set DNS on Mac: " + err.Error())
		}
	}
	logger.Log(logger.DEBUG, fmt.Sprintf("System DNS successfully updated. Primary=%s Secondary=%s", cfg.Primary, cfg.Secondary))
	return nil
}
