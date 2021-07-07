package htmlReplace

import (
	"io"
	"io/ioutil"
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
	shContext += `nginx -g 'daemon off;'`

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
		_, _ = copyFile(dir+"/favicon-cafe.ico", dir+"/favicon.ico")
		htmlIndexPageContext = strings.Replace(htmlIndexPageContext, SedCloudOld, SedCloudNew, -1)
	}

	_ = os.Remove(htmlIndexPagePath)
	err = ioutil.WriteFile(htmlIndexPagePath, []byte(htmlIndexPageContext), 0755)
	if err != nil {
		return err
	}
	return nil
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	return io.Copy(dst, src)
}
