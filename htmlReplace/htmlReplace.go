package htmlReplace

import (
	"io/ioutil"
	"nginx_config_builder/helper"
	"os"
	"strings"
)

const (
	SedIndex    = `https://lwy.yungengxin.com`
	SedCloudOld = `<div id=run-platform style="display: none;">lwy</div>`
	SedCloudNew = `<div id=run-platform style="display: none;">ylwp</div>`
)

func BuildStartShellFile(path string, dirs []string) error {
	shContext := `#!/bin/bash`
	shContext += "\n"

	for _, v := range dirs {
		shContext += `sed -i "s|https://lwy.yungengxin.com|$HOST|g" /usr/share/nginx/html/`
		shContext += v
		shContext += "/index.html"
		shContext += "\n"
	}

	shContext += "rm -rf /usr/local/lwyview/dist/*\n"
	shContext += "cp -rf /usr/share/nginx/html/latest/* /usr/local/lwyview/dist/\n"

	shContext += `exec nginx -g 'daemon off;'`

	_ = os.Remove(path + "/start.sh")
	err := ioutil.WriteFile(path+"/start.sh", []byte(shContext), 0777)
	if err != nil {
		return err
	}
	return nil
}

func HtmlReplace(dir, host string) error {

	// 替换主页
	// sed -i "s|https://lwy.yungengxin.com|$HOST|g" /usr/share/nginx/html/index.html
	htmlIndexPagePath := dir + "/index.html"
	bytesContext, err := ioutil.ReadFile(htmlIndexPagePath)
	if err != nil {
		return err
	}
	htmlIndexPageContext := strings.Replace(string(bytesContext), SedIndex, host, -1)

	// 云领无盘
	if strings.Contains(host, "ylwp") {
		_ = os.Remove(dir + "/favicon.ico")
		_ = helper.FileCopy(dir+"/favicon-cafe.ico", dir+"/favicon.ico")
		htmlIndexPageContext = strings.Replace(htmlIndexPageContext, SedCloudOld, SedCloudNew, -1)
	}

	_ = os.Remove(htmlIndexPagePath)
	err = ioutil.WriteFile(htmlIndexPagePath, []byte(htmlIndexPageContext), 0755)
	if err != nil {
		return err
	}
	return nil
}
