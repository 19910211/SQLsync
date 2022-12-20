package logic

import (
	"strings"
	"time"
)

type (
	Version = time.Time

	SqlFile struct {
		Version
		Path string
	}
)

func (f SqlFile) VersionNo() string {
	return f.Version.Format(_datetimeLayout)
}

func MustVersionOfFileName(fileName string) Version {
	version, err := time.Parse(_datetimeLayout, getVersionNo(fileName))
	if err != nil {
		panic(err)
	}
	return version
}

func VersionOfFileName(fileName string) (Version, error) {
	return time.Parse(_datetimeLayout, getVersionNo(fileName))
}

// 通过版本号获取当前版本
func VersionOfVersionNo(versionNo string) (Version, error) {
	return time.Parse(_datetimeLayout, versionNo)
}

// 通过文件名获取版本
func getVersionNo(fileName string) string {
	return strings.Split(fileName, "-")[0]
}
