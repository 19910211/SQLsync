package logic

import (
	"SQLsync/common"
	"SQLsync/svc"
	"fmt"
)

//
// ScanCommand
// @Description: 等待指令
// @param ScanCommand
func ScanCommand(svcCtx *svc.ServiceContext) {

	common.CmdInput(`请输入命令:(new|sync)>`)
	//执行
	var cmd common.Cmd
	if _, err := fmt.Scan(&cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	if command, ok := commandManager[cmd]; ok {
		command.Execute(svcCtx)
	} else {
		fmt.Println(fmt.Sprintf("%s 命令不存在", cmd))
	}
}
