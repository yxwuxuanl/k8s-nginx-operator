package nginx

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

type Config struct {
	CommonConfig

	*ReverseProxyConfig

	*TCPReverseProxyConfig
}

// +kubebuilder:object:generate=true

type CommonConfig struct {
	Resolvers []string `json:"resolver,omitempty"`

	WorkerProcesses *int `json:"workerProcesses,omitempty"`

	// +kubebuilder:default:=1024
	WorkerConnections int `json:"workerConnections,omitempty"`
}

// +kubebuilder:object:generate=true

type TCPReverseProxyConfig struct {
	// +kubebuilder:validation:MinItems:=1
	Servers []Upstream `json:"servers"`

	Hash string `json:"hash,omitempty"`

	// +kubebuilder:validation:Enum=TCP;UDP
	// +kubebuilder:default:="TCP"
	Protocol string `json:"protocol,omitempty"`

	// +kubebuilder:default:="10m"
	// +kubebuilder:validation:Pattern:=^\d+(s|m)$
	ProxyTimeout string `json:"proxyTimeout,omitempty"`
}

// +kubebuilder:object:generate=true

type Upstream struct {
	Server   string `json:"server"`
	Weight   *int32 `json:"weight,omitempty"`
	MaxConns *int32 `json:"maxConns,omitempty"`
}

func (u *Upstream) String() string {
	var dir []string

	dir = append(dir, u.Server)

	if u.Weight != nil {
		dir = append(dir, fmt.Sprintf("weight=%d", *u.Weight))
	}

	if u.MaxConns != nil {
		dir = append(dir, fmt.Sprintf("max_conns=%d", *u.MaxConns))
	}

	return strings.Join(dir, " ")
}

// +kubebuilder:object:generate=true

type ReverseProxyConfig struct {
	ProxyPass string `json:"proxyPass"`

	Logfmt *string `json:"logFormat,omitempty"`

	Rewrite []RewriteConfig `json:"rewrite,omitempty"`

	HideHeaders []string `json:"proxyHideHeader,omitempty"`

	ProxyHeaders map[string]string `json:"proxySetHeader,omitempty"`

	// +kubebuilder:default:="15s"
	// +kubebuilder:validation:Pattern:=^\d+s$
	ReadTimeout string `json:"proxyReadTimeout,omitempty"`

	// +kubebuilder:default:="15s"
	// +kubebuilder:validation:Pattern:=^\d+s$
	SendTimeout string `json:"proxySendTimeout,omitempty"`

	// +kubebuilder:default:="15s"
	// +kubebuilder:validation:Pattern:=^\d+s$
	ConnectTimeout string `json:"proxyConnectTimeout,omitempty"`
}

type RewriteConfig struct {
	Regex       string `json:"regex"`
	Replacement string `json:"replacement"`

	// +kubebuilder:validation:Enum=last;break
	// +kubebuilder:default=break
	Flag string `json:"flag,omitempty"`
}

func (i *TCPReverseProxyConfig) IsUDP() bool {
	return i.Protocol == string(corev1.ProtocolUDP)
}

func (r *ReverseProxyConfig) GetLogFmt() string {
	if r.Logfmt != nil {
		return *r.Logfmt
	}

	return ""
}

func (c *CommonConfig) GetWorkerProcesses() string {
	if v := c.WorkerProcesses; v != nil {
		return strconv.Itoa(*v)
	}

	return "auto"
}
