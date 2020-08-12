package main

import (
	gopherjs "github.com/gopherjs/gopherjs/js"
	"github.com/mark07x/YanBiaoJiuMing/cmd"
	"github.com/mark07x/YanBiaoJiuMing/js"
)

func main() {
	if gopherjs.Global != nil && gopherjs.Global.Get("window") != gopherjs.Undefined { // 浏览器环境
		js.Main()
	} else { // Go原生环境或Node.js环境
		cmd.Main()
	}
}