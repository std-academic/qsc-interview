package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

type AccountEx struct {
	Username string
	Password string
}

func TestLoginPassword(t *testing.T) {
	var account AccountEx

	_, filename, _, _ := runtime.Caller(0)
	buff, err := os.ReadFile(path.Join(path.Dir(filename), "../account.json"))
	if err != nil {
		fmt.Printf("找不到账号密码文件！\n请在 zjuam 目录下创建 account.json，填入 username 与 password 字段。\n%s", err.Error())
		panic(err)
	}

	err = json.Unmarshal(buff, &account)
	if err != nil {
		fmt.Printf("账密文件读取失败！\n请检查 zjuam 目录下 account.json，填入 username 与 password 字段。\n%s", err.Error())
		panic(err)
	}

	if !LoginPassword(account.Username, account.Password) {
		panic("账密 登录失败")
	}
}

func TestLoginOAuth(t *testing.T) {
	if !LoginOAuth("http://appservice.zju.edu.cn/index") {
		panic("OAuth 登录失败")
	}
}
