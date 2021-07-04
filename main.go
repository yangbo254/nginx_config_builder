package main

import (
	"flag"
	"io/ioutil"
	"log"
	"nginx_config_builder/template"
	"os"
	"regexp"
	"sort"
	"strings"
)

func scanDir(match, dirPath string, num int) []string {
	var dirs []string
	files, _ := ioutil.ReadDir(dirPath)
	log.Printf("load dir info and scanning...")
	for _, dir := range files {
		reg, _ := regexp.Compile(match)
		if "" != reg.FindString(dir.Name()) {
			dirs = append(dirs, dir.Name())
		}
	}
	log.Printf("sort dir...")
	sort.Strings(dirs)
	if num == -1 || len(dirs) <= num {
		return dirs
	} else {
		return dirs[(len(dirs) - num):]
	}

}

func buildNginxConfig(fileName string, dirs []string) error {
	log.Printf("build nginx config...")
	latestPath := ""
	replaceVersion := ""
	for _, version := range dirs {
		replaceContext := template.NginxVersionReplaceFormat
		replaceContext = strings.Replace(replaceContext, "{{version}}", version, -1)
		replaceContext = strings.Replace(replaceContext, "{{path}}", version, -1)
		replaceVersion += replaceContext
		latestPath = version
	}
	nginxContext := template.NginxDefaultConfigTemplate
	nginxContext = strings.Replace(nginxContext, "{{version_replace}}", replaceVersion, -1)
	nginxContext = strings.Replace(nginxContext, "{{default_path}}", latestPath, -1)

	_ = os.Remove(fileName)
	return ioutil.WriteFile(fileName, []byte(nginxContext), 0755)
}

func buildDockerfile(configFileName string, dirs []string) error {
	log.Printf("build dockerfile...")
	copyCommand := ""
	for _, version := range dirs {
		replaceContext := template.DockerfileCopyDirCommand
		replaceContext = strings.Replace(replaceContext, "{{source_path}}", version, -1)
		copyCommand += replaceContext
	}
	configCommand := strings.Replace(template.DockerfileCopyNginxConfigCommand,
		"{{source_path}}", configFileName, -1)
	dockerfileContext := template.DockerfileTemplate
	dockerfileContext = strings.Replace(dockerfileContext, "{{copy_dst}}", copyCommand, -1)
	dockerfileContext = strings.Replace(dockerfileContext, "{{copy_config}}", configCommand, -1)

	const dockerFileName = "./Dockerfile"
	_ = os.Remove(dockerFileName)
	return ioutil.WriteFile(dockerFileName, []byte(dockerfileContext), 0755)
}

var parameterMaxNum = flag.Int("maxnum", -1, "the maximum retention versions")
var parameterMatch = flag.String("match", `^(\d+)\.(\d+)$`, "the regular rules(default:version{#.#})")
var parameterScanDir = flag.String("path", ".", "scan directory(default:work dir)")

func main() {
	flag.Parse()
	dirs := scanDir(*parameterMatch, *parameterScanDir, *parameterMaxNum)
	log.Printf("load dir num: %d", len(dirs))
	if len(dirs) == 0 {
		log.Fatal("the required directory was not found.")
	}
	const nginxConfigFileName = "nginx_default.conf"
	err := buildNginxConfig(nginxConfigFileName, dirs)
	if err != nil {
		log.Fatal(err)
	}
	err = buildDockerfile(nginxConfigFileName, dirs)
	if err != nil {
		log.Fatal(err)
	}
	return
}
