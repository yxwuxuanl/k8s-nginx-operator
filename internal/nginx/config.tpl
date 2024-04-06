daemon off;
pid /tmp/nginx.pid;

worker_processes {{ .GetWorkerProcesses }};
worker_rlimit_nofile 65535;

error_log /dev/stdout notice;
error_log /dev/stdout error;

events {
    multi_accept       on;
    worker_connections {{ .WorkerConnections }};
}

http {
    client_body_temp_path /tmp/client_temp;
    proxy_temp_path       /tmp/proxy_temp_path;
    fastcgi_temp_path     /tmp/fastcgi_temp;
    uwsgi_temp_path       /tmp/uwsgi_temp;
    scgi_temp_path        /tmp/scgi_temp;

    {{- if .ReverseProxyConfig }}
    {{- $logfmt := .ReverseProxyConfig.GetLogFmt }}
    {{- if $logfmt }}
    log_format operator '{{ $logfmt }}';
    access_log /dev/stdout operator;
    {{- else }}
    access_log off;
    {{- end }}

    server {
        listen {{ .ListenPort }};

        {{- $resolver := .GetResolver }}
        {{- if $resolver }}
        resolver {{ $resolver }} ipv6=off;
        {{- end }}

        location / {
            {{- range $_, $rewrite := .ReverseProxyConfig.Rewrite }}
            rewrite {{ $rewrite.Regex }} {{ $rewrite.Replacement }} {{ $rewrite.Flag }};
            {{- end }}

            proxy_read_timeout {{ .ReverseProxyConfig.ReadTimeout }};
            proxy_send_timeout {{ .ReverseProxyConfig.SendTimeout }};
            proxy_connect_timeout {{ .ReverseProxyConfig.ConnectTimeout }};

            proxy_ssl_server_name on;

            {{- range $_, $header := .ReverseProxyConfig.HideHeaders }}
            proxy_hide_header {{ $header }};
            {{- end }}

            proxy_set_header Upgrade $http_upgrade;

            {{- range $header, $value := .ReverseProxyConfig.ProxyHeaders }}
            proxy_set_header {{ $header }} {{ $value }};
            {{- end }}

            proxy_pass "{{ .ReverseProxyConfig.ProxyPass }}";
        }

        location = {{ .ProbeURL }} {
            access_log off;
            return 200 '.';
        }
    }
    {{- end }}

    server {
        listen {{ .ProbePort }};

         location = {{ .ProbeURL }} {
             access_log off;
             return 200 '.';
         }
    }
}

{{- if .TCPReverseProxyConfig }}
stream {
    upstream servers {
        {{- with .TCPReverseProxyConfig.Hash }}
        hash {{ . }};
        {{- end }}

        {{- range $_, $server := .TCPReverseProxyConfig.Servers }}
        server {{ $server.String }};
        {{- end }}
    }

    {{- $resolver := .GetResolver }}
    {{- if $resolver }}
    resolver {{ $resolver }} ipv6=off;
    {{- end }}

    server {
        listen {{ .ListenPort }} {{ if .TCPReverseProxyConfig.IsUDP }}udp{{ end }};

        proxy_timeout {{ .TCPReverseProxyConfig.ProxyTimeout }};

        proxy_pass servers;
    }
}
{{- end }}