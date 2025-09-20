package xray

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StreamField struct {
	Network        string
	StreamSecurity string
	Path           string
	Host           string
	TCPHeaderType  string
}

type ParserVmess struct {
	Address  string
	Port     int
	UUID     string
	Security string
	*StreamField
}

func (p *ParserVmess) Parse(rawUri string) {
	raw := strings.TrimPrefix(rawUri, "vmess://")
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return
	}
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return
	}
	p.Address = fmt.Sprintf("%v", j["add"])
	p.Port, _ = strconv.Atoi(fmt.Sprintf("%v", j["port"]))
	p.UUID = fmt.Sprintf("%v", j["id"])
	p.Security = fmt.Sprintf("%v", j["security"])
	if p.Security == "" {
		p.Security = "none"
	}
	p.StreamField = &StreamField{
		Network: fmt.Sprintf("%v", j["net"]),
		Path:    fmt.Sprintf("%v", j["path"]),
	}
}

type ParserVless struct {
	Address    string
	Port       int
	UUID       string
	Encryption string
	Flow       string
	*StreamField
}

func (p *ParserVless) Parse(rawUri string) {
	u, _ := url.Parse(rawUri)
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.UUID = u.User.Username()
	q := u.Query()
	p.Encryption = q.Get("encryption")
	if p.Encryption == "" {
		p.Encryption = "none"
	}
	p.Flow = q.Get("flow")
	p.StreamField = &StreamField{
		Network: q.Get("type"),
		Path:    q.Get("path"),
	}
}

type ParserTrojan struct {
	Address  string
	Port     int
	Password string
	*StreamField
}

func (p *ParserTrojan) Parse(rawUri string) {
	u, _ := url.Parse(rawUri)
	p.Address = u.Hostname()
	p.Port, _ = strconv.Atoi(u.Port())
	p.Password = u.User.Username()
	q := u.Query()
	p.StreamField = &StreamField{
		Network: q.Get("type"),
		Path:    q.Get("path"),
	}
}

type ParserWireGuard struct {
	PrivateKey string
	Address    string
	DNS        string
	PublicKey  string
	Endpoint   string
	AllowedIPs []string
}

func (p *ParserWireGuard) Parse(raw string) error {
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "PrivateKey":
			p.PrivateKey = value
		case "Address":
			p.Address = value
		case "DNS":
			p.DNS = value
		case "PublicKey":
			p.PublicKey = value
		case "Endpoint":
			p.Endpoint = value
		case "AllowedIPs":
			p.AllowedIPs = strings.Split(value, ",")
		}
	}
	return nil
}

var (
	AllNodes []interface{}
	nodesMu  sync.Mutex
)

func getOutputDir(subdir string) string {
	dir, _ := os.UserCacheDir()
	dir = filepath.Join(dir, "freedom-core", subdir)
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

func parseXrayLink(link string) (interface{}, error) {
	link = strings.TrimSpace(link)
	scheme := ""
	if idx := strings.Index(link, "://"); idx != -1 {
		scheme = link[:idx+3]
	}
	switch scheme {
	case "vmess://":
		node := &ParserVmess{}
		node.Parse(link)
		if node.Address == "" {
			return nil, fmt.Errorf("invalid vmess")
		}
		return node, nil
	case "vless://":
		node := &ParserVless{}
		node.Parse(link)
		if node.Address == "" {
			return nil, fmt.Errorf("invalid vless")
		}
		return node, nil
	case "trojan://":
		node := &ParserTrojan{}
		node.Parse(link)
		if node.Address == "" {
			return nil, fmt.Errorf("invalid trojan")
		}
		return node, nil
	default:
		var node interface{}
		if err := json.Unmarshal([]byte(link), &node); err != nil {
			return nil, fmt.Errorf("unsupported scheme and invalid JSON: %v", err)
		}
		return node, nil
	}
}

func SaveNodes() error {
	nodesMu.Lock()
	defer nodesMu.Unlock()
	dir := getOutputDir("xray")
	outFile := filepath.Join(dir, "config.json")
	outbounds := []map[string]interface{}{}
	for _, n := range AllNodes {
		var network string
		var path string
		if nodeField, ok := n.(interface{ GetStreamField() *StreamField }); ok && nodeField.GetStreamField() != nil {
			network = nodeField.GetStreamField().Network
			path = nodeField.GetStreamField().Path
		}
		if network == "" {
			network = "tcp"
		}
		switch node := n.(type) {
		case *ParserVmess:
			outbounds = append(outbounds, map[string]interface{}{
				"protocol": "vmess",
				"settings": map[string]interface{}{
					"vnext": []map[string]interface{}{
						{
							"address": node.Address,
							"port":    node.Port,
							"users": []map[string]interface{}{
								{
									"id":       node.UUID,
									"alterId":  0,
									"security": node.Security,
								},
							},
						},
					},
				},
				"streamSettings": map[string]interface{}{
					"network": network,
					"wsSettings": map[string]interface{}{
						"path": path,
					},
				},
				"tag": "proxy",
			})
		case *ParserVless:
			outbounds = append(outbounds, map[string]interface{}{
				"protocol": "vless",
				"settings": map[string]interface{}{
					"vnext": []map[string]interface{}{
						{
							"address": node.Address,
							"port":    node.Port,
							"users": []map[string]interface{}{
								{
									"id":         node.UUID,
									"encryption": node.Encryption,
									"flow":       node.Flow,
								},
							},
						},
					},
				},
				"streamSettings": map[string]interface{}{
					"network": network,
					"wsSettings": map[string]interface{}{
						"path": path,
					},
				},
				"tag": "proxy",
			})
		case *ParserTrojan:
			outbounds = append(outbounds, map[string]interface{}{
				"protocol": "trojan",
				"settings": map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"address":  node.Address,
							"port":     node.Port,
							"password": node.Password,
						},
					},
				},
				"streamSettings": map[string]interface{}{
					"network": network,
					"wsSettings": map[string]interface{}{
						"path": path,
					},
				},
				"tag": "proxy",
			})
		case *ParserWireGuard:
			wgFile := filepath.Join(dir, "wireguard.conf")
			data, _ := json.MarshalIndent(node, "", "  ")
			_ = os.WriteFile(wgFile, data, 0o644)
		}
	}
	config := map[string]interface{}{
		"inbounds": []map[string]interface{}{
			{
				"port":     1080,
				"protocol": "socks",
				"settings": map[string]interface{}{},
			},
		},
		"outbounds": outbounds,
		"routing": map[string]interface{}{
			"domainStrategy": "IPIfNonMatch",
			"rules": []map[string]interface{}{
				{
					"type":        "field",
					"ip":          []string{"0.0.0.0/0", "::/0"},
					"outboundTag": "proxy",
				},
			},
		},
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(outFile, data, 0o644)
}

func AddXrayLink(link string) error {
	node, err := parseXrayLink(link)
	if err != nil {
		return err
	}
	nodesMu.Lock()
	AllNodes = []interface{}{node}
	nodesMu.Unlock()
	return SaveNodes()
}

func ParseXrayStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, "Missing 'link' query parameter", http.StatusBadRequest)
		return
	}
	node, err := parseXrayLink(link)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	nodesMu.Lock()
	AllNodes = []interface{}{node}
	nodesMu.Unlock()
	if err := SaveNodes(); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "time": time.Now().Format(time.RFC3339)})
}
