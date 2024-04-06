package nginx

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	ConfigDirPath = "/etc/nginx/nginx.conf"
	ListenPort    = 8080
	ProbeURL      = "/_/ping"
	ProbePort     = 3001
	GroupID       = 101
	UserID        = 101
)

const (
	LogFmtCombined = `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
)

var defaultResolvers []string

//go:embed config.tpl
var tplData string

var configTpl *template.Template

var funcMap = map[string]any{}

func init() {
	var err error
	configTpl, err = template.New("nginx-config").Funcs(funcMap).Parse(tplData)

	if err != nil {
		panic(err)
	}

	defaultResolvers, err = getNameservers()
	if err != nil {
		klog.ErrorS(err, "failed to get nameservers")
	}
}

type TemplateData struct {
	Config

	ListenPort int
	ProbePort  int
	ProbeURL   string
}

func BuildConfig(conf Config) ([]byte, error) {
	data := &TemplateData{
		Config:     conf,
		ListenPort: ListenPort,
		ProbeURL:   ProbeURL,
		ProbePort:  ProbePort,
	}

	var buf bytes.Buffer
	if err := configTpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestConfig(config []byte) error {
	tmpfile, err := os.CreateTemp("", "nginx-*.conf")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(config); err != nil {
		tmpfile.Close()
		return err
	}

	tmpfile.Close()

	cmd := exec.Command("nginx", "-c", tmpfile.Name(), "-t")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH")}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return errors.New(string(slurp))
	}

	return nil
}

func getNameservers() ([]string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var nameservers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				nameservers = append(nameservers, fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nameservers, nil
}

func (o *TemplateData) GetResolver() string {
	if len(o.Resolvers) > 0 {
		return strings.Join(o.Resolvers, " ")
	}

	return strings.Join(defaultResolvers, " ")
}
