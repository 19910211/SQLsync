package logic

import (
	"SQLsync/common"
	"SQLsync/logic/command"
	_ "SQLsync/logic/command"
	"SQLsync/svc"
)

type Command interface {
	Execute(svcCtx *svc.ServiceContext)
}

var commandManager = map[common.Cmd]Command{
	common.NewVersionCmd:  &command.NewSqlVersionCommand{},
	common.SyncVersionCmd: &command.SyncVersionCommand{},
}
