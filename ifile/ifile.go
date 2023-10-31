package ifile

import (
	"os"
)

// 调用os.MkdirAll递归创建文件夹
func MakeDir(filePath string, perm ...os.FileMode) error {
	var p os.FileMode
	if len(perm) == 0 {
		p = os.ModePerm
	} else {
		p = perm[0]
	}
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, p)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 返回文件最后修改时间
func UpdTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer func() {
		_ = f.Close()
	}()
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return 0
	}
	return fi.ModTime().Unix()
}

// FileStat 读取文件信息
func FileStat(path string) os.FileInfo {
	f, err := os.Stat(path)
	if err == nil {
		return f
	} else if os.IsExist(err) {
		return f
	} else {
		return nil
	}
}
