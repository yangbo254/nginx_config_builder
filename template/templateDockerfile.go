package template

const (
	// DockerfileTemplate docker build file
	DockerfileTemplate = `FROM nginx:1.21.0
{{copy_dst}}{{copy_config}}
COPY ./debug /usr/share/nginx/html/debug
COPY ./start.sh /start.sh
CMD ["/start.sh"]
`
	// DockerfileCopyDirCommand 拷贝文件命令行
	DockerfileCopyDirCommand = `
COPY ./{{source_path}}/ /usr/share/nginx/html/{{source_path}}/
`
	// DockerfileCopyNginxConfigCommand 拷贝nginx配置文件命令行
	DockerfileCopyNginxConfigCommand = `
COPY ./{{source_path}} /etc/nginx/conf.d/default.conf
`
)
