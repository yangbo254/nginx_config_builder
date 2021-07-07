package template

const (
	NginxCoreConfig = `
user 				nginx;
worker_processes 	auto;

error_log 			/var/log/nginx/error.log notice;
pid 				/var/run/nginx.pid

events {
	worker_connections 1024;
}
http {
	include /etc/nginx/mime.types;
	default_type application/octet-stream;
	
	log_format main '$remote_addr - $remote_user [$time_local] "$request" '
					'$status $body_bytes_sent "http_referer" '
					'"$http_user_agent" "$http_x_forwarded_for"';

	access_log /var/log/nginx/access.log main;
	
	sendfile on;
	keepalive_timeout 65;
	gzip on;
	
	include /etc/nginx/conf.d/*.conf;
}
`
	NginxDefaultConfigTemplate = `
server {
	listen 80;
	listen [::]:80;
	server_name localhost;

	gzip on;
	gzip_min_length 1k;
	gzip_comp_level 3;
	gzip_types application/javascript text/css text/javascript image/jpeg image/gif image/png image/x-icon;
	gzip_disable "MSIE [1-6]\.";
	add_header Cache-Control no-cache;
	underscores_in_headers on;
	
	location / {
		set $is_matched 0;
		{{version_replace}}
		if ($cookie_api_version = "debug") {
			set $is_matched 1;
			root /usr/share/nginx/html/debug;
			add_header X-MatchVersion debug;
		}
		if ($is_matched = 0) {
			root /usr/share/nginx/html/{{default_path}};
			add_header X-MatchVersion {{default_path}};
		}
		index index.html index.htm;
	}
}
`
	NginxVersionReplaceFormat = `
		if ($http_api_version = "{{version}}") {
			set $is_matched 1;
			root /usr/share/nginx/html/{{path}};
			add_header X-MatchVersion {{version}};
		}
		if ($cookie_api_version = "{{version}}") {
			set $is_matched 1;
			root /usr/share/nginx/html/{{path}};
			add_header X-MatchVersion {{version}};
		}
`
)
