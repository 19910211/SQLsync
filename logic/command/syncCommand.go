package command

import (
	"SQLsync/common"
	"SQLsync/logic/filter"
	"SQLsync/svc"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var _emptyVersion = common.Version{} // 空版本

type SyncVersionCommand struct {
	svcCtx *svc.ServiceContext
}

// 同步sql文件
func (s *SyncVersionCommand) Execute(svcCtx *svc.ServiceContext) {
	s.svcCtx = svcCtx

	// 获取当前版本号
	currentVersion, err := s.GetCurrentVersion()
	if err != nil {
		fmt.Println("获取版本号错误", fmt.Sprintf("error:%+v", err))
		return
	}

	// 判断文件是否存在
	if _, err := os.Stat(s.svcCtx.Config.Path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(fmt.Sprintf("文件夹不存在 %s", s.svcCtx.Config.Path))
			return
		}
		fmt.Println(fmt.Sprintf("error:%+v", err))
		return
	}

	// 获取SQLVersion文件列表
	sqlFileList := GetFileList(svcCtx.Config.Path, filter.DirFilter(currentVersion), filter.FileFilter(currentVersion))
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

func (s *SyncVersionCommand) sync(sqlFileList []*common.SqlFile) error {
	for _, sqlFile := range sqlFileList {
		// 更新
		if err := s.doSync(sqlFile); err != nil {
			return err
		}
	}
	return nil
}

func (s *SyncVersionCommand) doSync(sqlFile *common.SqlFile) error {

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
func (s *SyncVersionCommand) GetCurrentVersion() (common.Version, error) {
	var version *string
	if err := s.svcCtx.DB.Get(&version, fmt.Sprintf("select max(version_no) from %s", s.svcCtx.Config.DataSource.Table)); err != nil {
		fmt.Println(fmt.Sprintf("select max(version_no) from %s", s.svcCtx.Config.DataSource.Table))
		return _emptyVersion, err
	}

	if version == nil {
		return _emptyVersion, nil
	}

	return common.VersionOfVersionNo(*version)
}

//
// GetFileList 获取所有sql版本文件
// 可以通过 dirFilter 过滤掉旧版本的目录
// 可以通过 fileFilter 过滤掉旧版本的sql文件
func GetFileList(rootPath string, dirFilter, fileFilter func(info fs.FileInfo, dir string) bool) []*common.SqlFile {
	var sqlFileList []*common.SqlFile

	list, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("获取文件列表错误:error:%+v", err))
		return nil
	}

	for _, f := range list {
		if f.IsDir() {
			if dirFilter(f, rootPath) {
				sqlFileList = append(sqlFileList, GetFileList(filepath.Join(rootPath, f.Name()), dirFilter, fileFilter)...)
			}
		} else {
			if fileFilter(f, rootPath) {
				if version, err := common.VersionOfFileName(f.Name()); err == nil {
					sqlFileList = append(sqlFileList, &common.SqlFile{
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
	sort.Slice(sqlFileList, func(i, j int) bool { return sqlFileList[i].Version.Before(sqlFileList[j].Version) })
	return sqlFileList
}
