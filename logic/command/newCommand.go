package command

import (
	"SQLsync/common"
	"SQLsync/svc"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NewSqlVersionCommand struct{}

//
// Execute
// @Description: 创建新的sql版本
func (c *NewSqlVersionCommand) Execute(svcCtx *svc.ServiceContext) {

	common.CmdInput(`新的sql文件版本名称：（不要有特殊字符）>`)
	var fileName string
	if _, err := fmt.Scan(&fileName); err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(strings.TrimSpace(fileName)) == 0 {
		fmt.Println("sql文件名称为空")
		return
	}

	datetime := time.Now().Local()
	// 获取目录 年/月
	dirPath := filepath.Join(svcCtx.Config.Path, datetime.Format("2006"), datetime.Format("01"))
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	// 获取文件路径 年/月/version_no-filename.sql
	fileName = filepath.Join(dirPath, fmt.Sprintf("%s-%s%s", datetime.Format(common.DatetimeLayout), fileName, common.SqlCommandExt))

	// 创建文件
	if f1, err := os.Create(fileName); err != nil {
		fmt.Println(err)
	} else {
		defer f1.Close()
		fmt.Println(fmt.Sprintf("创建完成地址：%s", fileName))
	}

}
