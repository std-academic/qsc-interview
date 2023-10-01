package core

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// HTTP 相关的 API

var USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
var HTTP_PROXY *url.URL = nil

var client *http.Client
var jar *cookiejar.Jar

func SetUserAgent(userAgent string) {
	USER_AGENT = userAgent
}

func SetHttpProxy(httpProxy string) {
	var err error
	HTTP_PROXY, err = url.Parse(httpProxy)
	if err != nil {
		log.Printf("错误：HTTP Proxy 设置失败: %s", httpProxy)
		HTTP_PROXY = nil
	}
	prepareClient()
}

func init() {
	prepareClient()
	ResetCookies()
}

func prepareClient() {
	if HTTP_PROXY == nil {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(HTTP_PROXY),
				// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
}

func ResetCookies() {
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Panicf("创建 Cookie Jar 失败: %s", err.Error())
	}
}

func get(url string, data map[string]string) (string, int, http.Header) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("创建 GET 请求失败: %s", err.Error())
	}

	req.Header.Add("User-Agent", USER_AGENT)
	if data != nil {
		q := req.URL.Query()
		for key, value := range data {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("发送 GET 请求失败: %s", err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("接收 GET 响应失败: %s", err.Error())
	}

	jar.SetCookies(req.URL, resp.Cookies())
	return string(body), resp.StatusCode, resp.Header
}

func post(endpoint string, data map[string]string) (string, int, http.Header) {
	// 发送 form 编码的 POST 请求
	var payload *strings.Reader = nil
	if data != nil {
		q := url.Values{}
		for key, value := range data {
			q.Set(key, value)
		}

		payload = strings.NewReader(q.Encode())
	}

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		log.Panicf("创建 POST 请求失败: %s", err.Error())
	}

	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("发送 POST 请求失败: %s", err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("接收 POST 响应失败: %s", err.Error())
	}

	jar.SetCookies(req.URL, resp.Cookies())
	return string(body), resp.StatusCode, resp.Header
}

func PostHelper(url string) (string, int) {
	// 由于我比较懒，不想在 main 里面再次实现一遍 http
	// 所以开放一个post接口供 example 函数使用

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Panicf("创建 POST 请求失败: %s", err.Error())
	}

	req.Header.Add("User-Agent", USER_AGENT)

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("发送 POST 请求失败: %s", err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("接收 POST 响应失败: %s", err.Error())
	}

	jar.SetCookies(req.URL, resp.Cookies())
	return string(body), resp.StatusCode
}

func GetCookies(site string) map[string]string {
	// 获取 cookie 信息并返回
	// 可用于其他程序登录后操作
	// 如果 site 为 "" 则获取认证系统

	if site == "" {
		site = URL_LOGIN
	}

	url, _ := url.Parse(site)
	cookies := jar.Cookies(url)
	s := map[string]string{}

	for _, cookie := range cookies {
		s[cookie.Name] = cookie.Value
	}

	return s
}
