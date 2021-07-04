package template

const (
	DockerfileTemplate = `
FROM nginx:1.21.0
{{copy_dst}}
{{copy_config}}
`
	DockerfileCopyDirCommand = `
COPY ./{{source_path}}/ /usr/share/nginx/html/{{source_path}}/
`
	DockerfileCopyNginxConfigCommand = `
COPY ./{{source_path}} /etc/nginx/conf.d/default.conf
`
)
