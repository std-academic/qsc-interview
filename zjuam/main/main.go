package main

import (
	"core"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
)

type Account struct {
	Username string
	Password string
}

func ServiceOauthExample() {
	// 使用本程序进行 OAuth 登录
	// 并自动获取数据
	// 以获取账号基本信息为例子

	// 注意：使用 core.LoginOAuth 前，必须已经登录认证系统。

	const URL_APIINFO = "http://appservice.zju.edu.cn/zju-smartcampus/zdyxxt/baseInfoCollection/getBaseInfo"
	const URL_APIAUTH = "http://appservice.zju.edu.cn/index"

	core.LoginOAuth(URL_APIAUTH) // OAuth 登录

	// 为方便起见，我们不是取出 Cookies 然后进行后续操作
	// 而是直接在 core 的 http 库中进行请求，方便一些
	// 取出 Cookies 给其他程序用也是可以的，在最后一个 example 里面
	content, status := core.PostHelper(URL_APIINFO)

	if status != 200 {
		log.Panicf("获取账号基本信息失败：%d\n%s\n", status, content)
	} else {
		get := func(name string) string {
			// 正则查找 JSON 中所需字段
			r := regexp.MustCompile(fmt.Sprintf(`"%s":"(.*?)"`, name))
			if r.MatchString(content) {
				return r.FindStringSubmatch(content)[1]
			} else {
				return "/"
			}
		}
		list := map[string]string{
			"xm":     "姓名",
			"jgmc":   "籍贯",
			"zymc":   "专业",
			"kzzd37": "家庭住址",
			"sfzh":   "身份证号",
		} // 我真的不是想开盒
		for key, value := range list {
			log.Printf("%s: %s\n", value, get(key))
		}
	}
}

func LoginPasswordExample() {
	// 使用账号密码登录例子
	// 请把 username password 填入 zjuam\account.json

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

func CookiesExample() {
	// 导出所有 Cookies

	log.Printf("认证系统 Cookies: \n%s\n\n\n", core.GetCookies(""))
	log.Printf("AppService Cookies: \n%s\n", core.GetCookies("http://appservice.zju.edu.cn/"))
}

func main() {
	LoginPasswordExample()
	ServiceOauthExample()
	CookiesExample()
}
