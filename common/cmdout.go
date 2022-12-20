package common

import (
	"fmt"
	"os/user"
	"strings"
)

func CmdOutput() func(data ...string) {
	var cmdPrefix string // 命令行前缀

	u, err := user.Current()
	if err != nil {
		fmt.Println(fmt.Sprintf("user.Current error:%+v", err))
		cmdPrefix = "$>"
	} else {
		cmdPrefix = fmt.Sprintf("%s>", u.Username)
	}

	return func(data ...string) {
		fmt.Print(fmt.Sprintf("%s\n%s", strings.Join(data, "\n"), cmdPrefix))
	}
}

func CmdInput(data string) {
	fmt.Print(data)
}
