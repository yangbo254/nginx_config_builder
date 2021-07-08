package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"log"
	"math/rand"
	"nginx_config_builder/helper"
	"nginx_config_builder/htmlReplace"
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
		reg, _ := regexp.Compile(`^` + match + `$`)
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

func gitScanVersionAndCopyDir(path, match, selectDir string, maxnum int) []string {
	// latest目录
	_ = helper.DirCopy(path+"/dist", path+"/latest")

	// 其他主版本
	repository, err := git.PlainOpen(path)
	if err != nil {
		return nil
	}
	tagRefs, err := repository.Tags()
	if err != nil {
		return nil
	}
	var tags []string
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		tags = append(tags, t.Name().String())
		return nil
	})
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	const regMatch = `refs/tags/v`
	var dirs []string
	for _, v := range tags {
		reg, _ := regexp.Compile(`^(` + regMatch + `)` + match + `$`)
		if "" != reg.FindString(v) && maxnum != 0 {
			dirs = append(dirs, strings.Replace(v, regMatch, "", -1))
			if maxnum != -1 {
				maxnum--
			}
		}
	}

	rootPath := fmt.Sprintf("%s/tmp_%d/", path, rand.Int())
	dirs = append(dirs, "debug")
	for _, v := range dirs {
		verPath := rootPath + v

		_ = os.RemoveAll(verPath)
		repositoryRemote, err := repository.Remotes()
		if err != nil {
			log.Print(err)
			continue
		}
		repositoryConfig := repositoryRemote[0].Config()
		if repositoryConfig == nil {
			continue
		}
		if 0 != len(repositoryConfig.URLs) {
			log.Printf("git checkout tag %s...", v)
			_, err := git.PlainClone(verPath, false, &git.CloneOptions{
				URL:               repositoryConfig.URLs[0],
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				ReferenceName:     plumbing.NewTagReferenceName(`v` + v),
				Auth: &http.BasicAuth{
					Username: *parameterGitUser,
					Password: *parameterGitPasswd,
				},
			})
			if err != nil {
				log.Print(err)
				continue
			}
			exists, err := helper.PathExists(verPath + "/" + selectDir)
			if err != nil {
				log.Print(err)
				continue
			}
			if exists {
				log.Printf("copy from tempdir...")
				_ = helper.DirCopy(verPath+"/"+selectDir, path+"/"+v)
			} else {
				log.Printf("not found select dir,continue...")
			}

			_ = os.RemoveAll(verPath)
			continue
		}
	}
	_ = os.RemoveAll(rootPath)
	return nil
}

func buildNginxConfig(fileName string, dirs []string) error {
	log.Printf("build nginx config...")
	replaceVersion := ""
	for _, version := range dirs {
		replaceContext := template.NginxVersionReplaceFormat
		replaceContext = strings.Replace(replaceContext, "{{version}}", version, -1)
		replaceContext = strings.Replace(replaceContext, "{{path}}", version, -1)
		replaceVersion += replaceContext
	}
	nginxContext := template.NginxDefaultConfigTemplate
	nginxContext = strings.Replace(nginxContext, "{{version_replace}}", replaceVersion, -1)
	nginxContext = strings.Replace(nginxContext, "{{default_path}}", "latest", -1)

	log.Printf("nginx config:%s", nginxContext)

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

	log.Printf("dockerfile:%s", dockerfileContext)

	const dockerFileName = "./Dockerfile"
	_ = os.Remove(dockerFileName)
	return ioutil.WriteFile(dockerFileName, []byte(dockerfileContext), 0755)
}

var parameterMaxNum = flag.Int("maxnum", -1, "the maximum retention versions")
var parameterMatch = flag.String("match", `(\d+)\.(\d+)`, "the regular rules(default:version{#.#})")
var parameterSelectDir = flag.String("select", "dist", "target dir name")
var parameterScanDir = flag.String("path", ".", "scan directory(default:work dir)")
var parameterGitUser = flag.String("git_user", "yangbo", "scan directory(default:work dir)")
var parameterGitPasswd = flag.String("git_passwd", "{C59ABA87-16EB-435E-9449-A87753CBEC28}", "scan directory(default:work dir)")

func main() {
	log.Printf("system will build nginx config and dockerfile...")
	flag.Parse()

	log.Printf("The system begins to prepare the file...")
	gitScanVersionAndCopyDir(*parameterScanDir, *parameterMatch, *parameterSelectDir, *parameterMaxNum)

	dirs := scanDir(*parameterMatch, *parameterScanDir, *parameterMaxNum)
	dirs = append(dirs, "latest")
	log.Printf("load dir num: %d", len(dirs))
	if len(dirs) == 0 {
		log.Fatal("the required directory was not found.")
	}
	_ = htmlReplace.BuildStartShellFile(*parameterScanDir, dirs)

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
