package template

const (
	DockerfileTemplate = `FROM nginx:1.21.0
{{copy_dst}}{{copy_config}}
COPY ./debug /usr/share/nginx/html/debug
COPY ./start.sh /start.sh
CMD ["/start.sh"]
`
	DockerfileCopyDirCommand = `
COPY ./{{source_path}}/ /usr/share/nginx/html/{{source_path}}/
`
	DockerfileCopyNginxConfigCommand = `
COPY ./{{source_path}} /etc/nginx/conf.d/default.conf
`
)
