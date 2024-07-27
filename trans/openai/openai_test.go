package openai

import (
	"fmt"
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/spf13/viper"
	"testing"
)

func TestDoTranslateReq(t *testing.T) {
	conf.InitConfig()
	log.D("openai url:", viper.GetString(conf.ConfigKeyOpenAIApiUrl))
	log.D("openai key:", viper.GetString(conf.ConfigKeyOpenAiKey))

	Translate("你好啊", "", func(result string, note string) {
		fmt.Println("你好啊 -> " + result)
		fmt.Println(note)
	})

	fmt.Println("========================================")

	Translate("You are a real, professional translation engine, please follow the following process "+
		"step by step to determine and translate the input content\\n1."+
		" Determine what language the input is in\\n2. if the input language is Chinese"+
		", then translate it into English", "", func(result string, note string) {
		fmt.Println("hello -> " + result)
		fmt.Println(note)
	})
}
