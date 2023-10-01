package main

import (
	"core"
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
)

type Account struct {
	Username string
	Password string
}

func main() {
	var account Account

	_, filename, _, _ := runtime.Caller(0)
	buff, err := os.ReadFile(path.Join(path.Dir(filename), "../account.json"))
	if err != nil {
		log.Panicf("找不到账号密码文件！\n请在 zjuam 目录下创建 account.json，填入 username 与 password 字段。\n%s", err.Error())
	}

	err = json.Unmarshal(buff, &account)
	if err != nil {
		log.Panicf("账密文件读取失败！\n请检查 zjuam 目录下 account.json，填入 username 与 password 字段。\n%s", err.Error())
	}

	core.LoginPassword(account.Username, account.Password)
}
