package filter

import (
	"SQLsync/common"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
)

// fileFilter 过滤掉旧版本的sql文件
func FileFilter(currentVersion common.Version) func(info fs.FileInfo, dir string) bool {
	return func(info fs.FileInfo, dir string) bool {
		fileName := info.Name()
		// 判断是否是sql文件
		if filepath.Ext(fileName) != common.SqlCommandExt {
			return false
		}

		// 文件为空
		if info.Size() == 0 {
			return false
		}

		// 判断版本是否在当前版本之后
		return common.MustVersionOfFileName(fileName).After(currentVersion)
	}
}

//
// dirFilter 过滤旧版本的目录下的所有sql文件
func DirFilter(currentVersion common.Version) func(info fs.FileInfo, dir string) bool {
	return func(info fs.FileInfo, dir string) bool {
		// 过滤掉 版本之前的年月
		switch len(info.Name()) {
		case 4: // 年
			return checkYear(info.Name(), currentVersion)
		case 2: // 月

			return checkMonth(info, dir, currentVersion)

		}

		return false
	}
}

func checkYear(name string, version common.Version) bool {
	return yearSub(name, version) > -1
}

func checkMonth(info fs.FileInfo, dir string, currentVersion common.Version) bool {
	s := yearSub(filepath.Base(dir), currentVersion)
	switch {
	case s > 0: // 年大于 直接返回true
		return true
	case s == 0: // 年相等 判断月
		month, err := strconv.Atoi(info.Name())
		if err != nil {
			fmt.Println(fmt.Sprintf("%s %s 解析报错", info.Name(), common.MonthLayout))
			return false
		}
		// 过滤掉 版本之前的月
		return month >= int(currentVersion.Month())
	}

	return false
}

func yearSub(yearName string, currentVersion common.Version) int {
	year, err := strconv.Atoi(yearName)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s %s 解析报错", yearName, common.YearsLayout))
		return -1
	}
	return year - currentVersion.Year()
}
