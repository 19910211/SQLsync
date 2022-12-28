package logic

import (
	"SQLsync/common"
	"SQLsync/svc"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SyncLogic struct {
	svcCtx *svc.ServiceContext
}

func NewSyncLogic(svcCtx *svc.ServiceContext) *SyncLogic {
	return &SyncLogic{
		svcCtx: svcCtx,
	}
}

//
// NewSqlVersionCommand
// @Description: 创建新的sql版本
func (s SyncLogic) NewSqlVersionCommand() {

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
	// 获取目录
	dirPath := filepath.Join(s.svcCtx.Config.Path, strconv.Itoa(datetime.Year()), strconv.Itoa(int(datetime.Month())))
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	// 获取文件路径
	fileName = filepath.Join(dirPath, fmt.Sprintf("%s-%s.sql", datetime.Format(DatetimeLayout), fileName))

	// 创建文件
	if f1, err := os.Create(fileName); err != nil {
		fmt.Println(err)
	} else {
		defer f1.Close()
		fmt.Println(fmt.Sprintf("创建完成地址：%s", fileName))
	}

}

// 同步sql文件
func (s SyncLogic) SyncVersionCommand() {
	// 获取当前版本号
	currentVersion, err := s.GetCurrentVersion()
	if err != nil {
		fmt.Println("获取版本号错误", fmt.Sprintf("error:%+v", err))
		return
	}

	// 获取SQLVersion文件列表
	sqlFileList := GetFileList(s.svcCtx.Config.Path, dirFilter(currentVersion), fileFilter(currentVersion))
	if len(sqlFileList) == 0 {
		fmt.Println("已经是最新的了")
		return
	}

	// 同步数据
	if err := s.sync(sqlFileList); err != nil {
		fmt.Println(fmt.Sprintf("同步失败：error:%+v", err))
		return
	}
}

func (s SyncLogic) sync(sqlFileList []*SqlFile) error {
	for _, sqlFile := range sqlFileList {
		// 更新
		if err := s.doSync(sqlFile); err != nil {
			return err
		}
	}
	return nil
}

func (s SyncLogic) doSync(sqlFile *SqlFile) error {

	sqlCommands, err := ioutil.ReadFile(sqlFile.Path)
	if err != nil {
		return err
	}

	return common.Transaction(s.svcCtx.DB, func(tx *sqlx.Tx) error {

		// 执行需要更新的数据
		if _, err := tx.Exec(string(sqlCommands)); err != nil {
			fmt.Println(fmt.Sprintf("sqlCommand:%s error: %+v", string(sqlCommands), err))
			return err
		}

		// 记录当前的版本
		versionSql := fmt.Sprintf("insert into %s values (%s)", s.svcCtx.Config.DataSource.Table, sqlFile.VersionNo())
		if _, err := tx.Exec(versionSql); err != nil {
			fmt.Println(fmt.Sprintf("versionSql:%s error:%+v", versionSql, err))
			return err
		}

		fmt.Println(fmt.Sprintf("VersionNo:%s sqlCommand:%s", sqlFile.VersionNo(), string(sqlCommands)))
		return nil
	})
}

// GetCurrentVersion 获取当前的版本号
func (s SyncLogic) GetCurrentVersion() (Version, error) {
	var version *string
	if err := s.svcCtx.DB.Get(&version, fmt.Sprintf("select max(version_no) from %s", s.svcCtx.Config.DataSource.Table)); err != nil {
		fmt.Println(fmt.Sprintf("select max(version_no) from %s", s.svcCtx.Config.DataSource.Table))
		return EmptyVersion, err
	}

	if version == nil {
		return EmptyVersion, nil
	}

	return VersionOfVersionNo(*version)
}

//
// OnStandby
// @Description: 等待指令
// @param syncLogic
func OnStandby(svcCtx *svc.ServiceContext) {
	syncLogic := NewSyncLogic(svcCtx)

	common.CmdInput(`请输入命令:(new|sync)>`)
	//执行
	var cmd string
	if _, err := fmt.Scan(&cmd); err != nil {
		fmt.Println(err.Error())
		return
	}
	switch cmd {
	case NewVersion:
		syncLogic.NewSqlVersionCommand()
	case SyncVersion:
		syncLogic.SyncVersionCommand()
	default:
		fmt.Println(fmt.Sprintf("%s 命令不存在", cmd))
	}
}

// fileFilter 过滤掉旧版本的sql文件
func fileFilter(currentVersion Version) func(info fs.FileInfo) bool {
	return func(info fs.FileInfo) bool {
		fileName := info.Name()
		// 判断是否是sql文件
		if filepath.Ext(fileName) == "sql" {
			return false
		}

		// 文件为空
		if info.Size() == 0 {
			return false
		}

		version := MustVersionOfFileName(fileName)
		// 判断版本是否在当前版本之后
		return version.After(currentVersion)
	}
}

//
// dirFilter 过滤旧版本的目录下的所有sql文件
func dirFilter(currentVersion Version) func(info fs.FileInfo) bool {
	return func(info fs.FileInfo) bool {
		// 过滤掉 版本之前的年月
		switch len(info.Name()) {
		case 4: // 年
			year, err := strconv.Atoi(info.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("%s %s 解析报错", info.Name(), YearsLayout))
				return false
			}
			// 过滤掉 版本之前的年
			return year >= currentVersion.Year()
		case 2: // 月
			month, err := strconv.Atoi(info.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("%s %s 解析报错", info.Name(), MonthLayout))
				return false
			}
			// 过滤掉 版本之前的月
			return month >= int(currentVersion.Month())
		}

		return false
	}
}

//
// GetFileList 获取所有sql版本文件
// 可以通过 dirFilter 过滤掉旧版本的目录
// 可以通过 fileFilter 过滤掉旧版本的sql文件
func GetFileList(rootPath string, dirFilter, fileFilter func(info fs.FileInfo) bool) []*SqlFile {
	var sqlFileList []*SqlFile

	list, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println("获取文件列表错误")
		return nil
	}

	for _, f := range list {
		if f.IsDir() {
			if dirFilter(f) {
				sqlFileList = append(sqlFileList, GetFileList(filepath.Join(rootPath, f.Name()), dirFilter, fileFilter)...)
			}
		} else {
			if fileFilter(f) {
				if version, err := VersionOfFileName(f.Name()); err == nil {
					sqlFileList = append(sqlFileList, &SqlFile{
						Version: version,
						Path:    filepath.Join(rootPath, f.Name()),
					})
				} else {
					fmt.Println(fmt.Sprintf("error:%+v", err))
				}
			}
		}
	}

	// 升序
	sort.Slice(sqlFileList, func(i, j int) bool { return sqlFileList[i].Version.Before(sqlFileList[i].Version) })
	return sqlFileList
}
