package google

import (
	"fmt"
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := NewTranslatorWithConfig(
		"https://translate.googleapis.com",
		"http://127.0.0.1:10809",
	)
	text := "Hello"
	mdText := "你好"

	translated, _ := translator.Translate(
		text, "en", "zh")
	fmt.Println(translated)
	if translated != mdText {
		t.Logf("expected: %s", mdText)
		t.Logf("given: %s", translated)
		t.Fail()
	}
}

func TestGoogleTranslate(t *testing.T) {
	conf.InitConfig()

	Translate("你好啊", strutils.English, func(result string, note string) {
		fmt.Println("你好啊 -> " + result)
		fmt.Println(note)
	})

	Translate("hello", strutils.Chinese, func(result string, note string) {
		fmt.Println("hello -> " + result)
		fmt.Println(note)
	})

}
