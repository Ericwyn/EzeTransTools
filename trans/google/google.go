package google

import (
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"github.com/Ericwyn/EzeTranslate/trans"
	translator "github.com/Ericwyn/go-googletrans"
	"github.com/spf13/viper"
)

var translateApi *translator.TranslateApi

func generalTransApi() *translator.TranslateApi {
	var translatorConfig = translator.TranslateConfig{}

	if viper.GetString(conf.ConfigKeyGoogleTranslateUrl) != "" {
		var url = viper.GetString(conf.ConfigKeyGoogleTranslateUrl)
		translatorConfig.ServiceUrls = []string{url}
		log.I("为 google 翻译设置 URL:" + url)
	}

	if viper.GetString(conf.ConfigKeyGoogleTranslateProxy) != "" {
		var proxy = viper.GetString(conf.ConfigKeyGoogleTranslateProxy)
		translatorConfig.Proxy = proxy
		log.I("为 google 翻译设置代理:" + proxy)
	}

	return translator.New(translatorConfig)
}

func Translate(str string, toLang strutils.Lang, transCallback trans.TransResCallback) {
	if translateApi == nil {
		translateApi = generalTransApi()
	}

	log.D("Google 翻译文字:", str)
	fromLang := strutils.DetectLanguage(str)
	if toLang == "" {
		if fromLang == strutils.Chinese {
			toLang = strutils.English
		} else {
			toLang = strutils.Chinese
		}
	}
	if fromLang == toLang {
		transCallback(str, string(fromLang+"->"+toLang))
		return
	}
	var note string
	var transRes string

	//var err error
	if strutils.English == toLang {
		// 中文较多的时候，都会翻译成英文句子
		log.D("翻译中文句子为英文")
		note = "zh -> en"
		result, err := translateApi.Translate(str, "auto", "en")
		if err != nil {
			//err.Error()
			log.E(err.Error())
			transCallback("翻译错误, 请查看日志", note)
			return
		}
		transRes = result.Text
		note = result.Src + "->" + result.Dest
	} else {
		note = "en -> zh"
		result, err := translateApi.Translate(str, "auto", "zh-cn")
		if err != nil {
			log.E(err.Error())
			transCallback("翻译错误, 请查看日志", note)
			return
		}
		transRes = result.Text
		note = result.Src + "->" + result.Dest
	}
	transCallback(transRes, note)
}
