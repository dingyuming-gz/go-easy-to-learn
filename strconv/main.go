package main

import (
	"fmt"
	"strconv"

	"github.com/Tinyming-GO/go-easy-to-learn/tools"
)

func demoAppendBool() {
	//AppendBool 将bool值追加到[]byte  并返回扩展缓冲区。
	b := []byte("bool:")
	b = strconv.AppendBool(b, true)
	fmt.Println(tools.RunFuncName() + " Output: " + string(b))
	//等价于append(dst, FormatBool(b)...)
}

func main() {
	demoAppendBool()
}
