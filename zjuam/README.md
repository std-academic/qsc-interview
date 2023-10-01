# ZJUAM

## 问题回答

### 你觉得解决这个任务的过程有意思吗？
> 额，我觉得挺有意思的，虽然一开始研究 security.js 里面的 RSA 加密很痛苦，但是最后解决了这些难题还是挺舒服的。

### 你在网上找到了哪些资料供你学习？你觉得去哪里/用什么方式搜索可以比较有效的获得自己想要的资料？
> 我从来没有写过 Golang，所以在一个下午内写出这个项目还是非常有挑战性的，
> 不过 Golang 的语法非常简洁，我很快就上手了。
> 一般来说，我只要在 Google 上面搜索我的问题（比如 golang http get set parameters），
> 就能找到解答。其中 Stackoverflow 和 Golang-Example 这两个网站一般是最有帮助的。
> 不过 golang 自带的 http 库和 python urllib 一样难用.. 我怀疑 golang 有类似 python requests 库这样的
> 第三方 HTTP 库，用它估计能节省不少时间。可惜一开始没想到。

### 在过程中，你遇到最大的困难是什么？你是怎么解决的？
> 估计是 security.js 那个奇奇怪怪的 old-fashioned RSA 加密。
> 根据我的搜索结果，这种 RSA 加密被称为 non-padding，即在数据末尾填充 0
> 不过我对 RSA 不是很了解。我只知道用 golang 官方库应该做不到这种 RSA，
> 所以最后我选择了用 js 引擎跑 security.js
> 一开始选择了一个叫做 v8go 的库。写完了发现不支持 Windows
> 立马换成了 goja。之后一些开发就相对容易很多，因为 zju 登录的脚本很简单。

### 完成任务之后，再回去阅读你写下的代码和文档，有没有看不懂的地方？如果再过一年，你觉得那时你还可以看懂你的代码吗？
> 没有。应该可以。

### 其他想说的想法或者建议。
> ZJU 的 login.js 感觉是改过很多次，一开始我看到 sms qrcode 以为能短信和二维码登录，
> 结果后来想实现的时候 发现这两个都是被弃用的。
> 短信登录直接不行了，二维码登录换成了钉钉的 OAuth。
> 钉钉的二维码 OAuth 登录要显示它的 iframe，我觉得在这个项目里
> 再实现一个 http 服务器过于麻烦，就没做。
> 最后的遗憾是异常处理有点水，基本上都是 log.Panicf，
> 不过这么短的时间要是再考虑精细的异常处理我会死的。

## 文件说明
### core/
`security.go`: 实现了模拟运行 `security.js` 进行 RSA 加密。

`http.go`: 实现了许多 HTTP 请求相关的API。主要是供 `login.go` 使用

`login.go`: 实现了登录的主要逻辑。

### main/
`main.go`: 作为一个使用例，调用 `core` 进行登录示例。

## API 文档
你需要 import `core`

其实看过一遍 `main.go` 就差不多会用了吧（

`core.LoginPassword(username string, password string) bool`: 账密登录。返回值代表成功与否

`core.LoginOAuth(endpoint string) bool`：OAuth登录。需要账密登录后才能进行。参数表示登录的服务器，返回值代表成功与否

`core.SetUserAgent(userAgent string)`: 设置 User-Agent。没啥用

`core.SetHttpProxy(httpProxy string)`: 设置代理服务器

`core.ResetCookies()`: 重置 Cookies。用于二次登陆其他账号

`core.Get/Post(url string)`: 用登录的 Cookies 进行 GET/POST。不能指定 payload。返回值： `content string, status_code int, headers http.Header`

`core.GetCookies(site string)`: 获取指定站点的 Cookies。如果 site 为空则返回认证系统的
