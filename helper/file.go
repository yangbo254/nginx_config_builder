package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// FileCopy copies a single file from src to dst
func FileCopy(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func(srcFd *os.File) {
		_ = srcFd.Close()
	}(srcFd)

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func(dstFd *os.File) {
		_ = dstFd.Close()
	}(dstFd)

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// DirCopy copies a whole directory recursively
func DirCopy(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = DirCopy(srcFp, dstFp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = FileCopy(srcFp, dstFp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// PathExists 判断文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
