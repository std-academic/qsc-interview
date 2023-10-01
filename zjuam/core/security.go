package core

// 使用 JS 引擎进行 RSA 加密
// 原计划使用自带 Cipher 库
// 但认证系统采用的是过时的 zero-padding 方案
// 操作失败，故直接运行 js 来模拟

import (
	"os"
	"path"
	"runtime"

	goja "github.com/dop251/goja"
)

var ctx *goja.Runtime

func init() {
	// 初始化 JS 引擎
	ctx = goja.New()

	_, filename, _, _ := runtime.Caller(0)
	js, err := os.ReadFile(path.Join(path.Dir(filename), "security.js"))
	if err != nil {
		panic(err)
	}

	_, err = ctx.RunString(string(js))
	if err != nil {
		panic(err)
	}
}

func GetEncryptedString(modulusString string, exponentString string, body string) string {
	ctx.Set("public_exponent", exponentString)
	ctx.Set("Modulus", modulusString)
	ctx.Set("password", body)

	_, err := ctx.RunString(`var key = new RSAUtils.getKeyPair(public_exponent, '', Modulus);
	var rpwd = password.split('').reverse().join('')
	var epwd = RSAUtils.encryptedString(key,rpwd)`)
	if err != nil {
		panic(err)
	}

	return ctx.Get("epwd").String()
}
