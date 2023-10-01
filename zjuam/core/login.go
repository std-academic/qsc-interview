package core

import (
	"encoding/json"
	"log"
	"net/url"
	"regexp"
)

// 与认证系统交互相关的 API

const URL_BASE = "https://zjuam.zju.edu.cn/cas/"
const URL_LOGIN = URL_BASE + "login" // ?service={SERVICE}

const URL_APIKEY = URL_BASE + "v2/getPubKey"                        // GET
const URL_APISMS = URL_BASE + "v1/services/sedsms"                  // GET, ?mobile=
const URL_APIQRCODE = "https://login.dingtalk.com/login/qrcode.htm" // GET, ?qrcode=(uuid)

var loginPage string
var loginToken map[string]string

type ApiKey struct {
	Modulus  string
	Exponent string
}

var loginKey ApiKey

func init() {
	loginToken = make(map[string]string)
}

func getLoginPage() {
	// 从登录页面获取 execution
	page, status, _ := get(URL_LOGIN, nil)

	if status != 200 {
		log.Panicf("无法获取登录页面: HTTP %d\n%s", status, page)
	}

	loginPage = page

	r := regexp.MustCompile(`name="execution" value="(.*?)"`)
	b := r.MatchString(page)

	if !b {
		log.Panicf("登录页面中找不到 execution:\n%s", page)
	}

	ex := r.FindStringSubmatch(page)[1]
	loginToken["execution"] = ex
}

func getLoginKey() {
	// 获取 Modulus 与 Exponent
	page, status, _ := get(URL_APIKEY, nil)

	if status != 200 {
		log.Panicf("无法获取公钥页面: HTTP %d\n%s", status, page)
	}

	err := json.Unmarshal([]byte(page), &loginKey)

	if err != nil {
		log.Panicf("无法获取公钥数据：%s\n%s", err.Error(), page)
	}
}

func LoginPassword(username string, password string) {
	// 密码登录
	getLoginPage()
	getLoginKey()

	payload := map[string]string{ // 构造 form 数据
		"username":  username,
		"password":  GetEncryptedString(loginKey.Modulus, loginKey.Exponent, password),
		"execution": loginToken["execution"],
		"_eventId":  "submit",
		"authcode":  "",
	}

	page, status, headers := post(URL_LOGIN, payload)

	if status == 200 {
		// 登录失败。尝试获取失败原因
		r := regexp.MustCompile(`<p class="error text-left" id="errormsg">(.*?)</p>`)
		if r.MatchString(page) {
			msg := r.FindStringSubmatch(page)[1]
			log.Panicf("密码登录失败：HTTP 200\n错误消息：%s\n", msg)
		} else {
			r = regexp.MustCompile(`exception.message=(.*?)[&"]`)
			if r.MatchString(page) {
				msg := r.FindStringSubmatch(page)[1]
				log.Panicf("密码登录失败：HTTP 200\n错误消息：%s\n", msg)
			} else {
				log.Panicf("密码登录失败：HTTP 200\n%s\n", page)
			}
		}
	} else if status == 403 {
		log.Panicf("密码登录失败：HTTP 403\n可能是访问次数过多，请稍后再试！\n")
	} else if status == 502 {
		log.Panicf("密码登录失败：HTTP 502\n可能是服务器维护！\n")
	} else if status != 302 {
		log.Panicf("密码登录失败：HTTP %d\n%s\n", status, page)
	} else {
		url, _ := url.Parse(URL_LOGIN)
		cookies := jar.Cookies(url)
		s := ""
		for _, cookie := range cookies {
			s += cookie.String() + "\n"
		}

		log.Printf("密码登录成功！\n重定向到：%s\nCookies：%s\n", headers["Location"], s)
	}
}
