package main

import (
	"SQLsync/common"
	"SQLsync/config"
	"SQLsync/logic"
	"SQLsync/svc"
	"flag"
	"fmt"
)

var configFile = flag.String("f", "conf.yaml", "the config file")

func main() {
	//加载配置文件
	conf := config.MustLoad(*configFile)

	// 初始化
	svcCtx := svc.NewServiceContext(conf)

	// 退出
	defer exit(svcCtx)

	// 等待指令
	logic.ScanCommand(svcCtx)
}

func exit(svcCtx *svc.ServiceContext) {
	svcCtx.DB.Close()
	common.CmdInput("请输入任意字符结束>")
	var quit string
	fmt.Scan(&quit)

}
