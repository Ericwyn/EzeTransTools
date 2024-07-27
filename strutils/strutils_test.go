package strutils

import (
	"fmt"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	fmt.Println(DetectLanguage("Hello, World!"))
	fmt.Println(DetectLanguage("你好啊"))
	fmt.Println(DetectLanguage("我的 api 接口 response 是错误的"))
}
