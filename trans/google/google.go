package google

import (
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"github.com/Ericwyn/EzeTranslate/trans"
	"github.com/spf13/viper"
)

var translator *Translator

func initTranslator() Translator {
	var googleApiUrl, translatorHttpProxy string
	if viper.GetString(conf.ConfigKeyGoogleTranslateUrl) != "" {
		googleApiUrl = viper.GetString(conf.ConfigKeyGoogleTranslateUrl)
		log.I("为 google 翻译设置 URL:" + googleApiUrl)
	} else {
		googleApiUrl = "https://translate.googleapis.com"
	}

	if viper.GetString(conf.ConfigKeyGoogleTranslateProxy) != "" {
		translatorHttpProxy = viper.GetString(conf.ConfigKeyGoogleTranslateProxy)
		log.I("为 google 翻译设置代理:" + translatorHttpProxy)
	}

	return NewTranslatorWithConfig(googleApiUrl, translatorHttpProxy)
}

func Translate(str string, toLang strutils.Lang, transCallback trans.TransResCallback) {
	if translator == nil {
		t := initTranslator()
		translator = &t
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
		result, err := translator.Translate(str, "auto", "en")
		if err != nil {
			//err.Error()
			log.E(err.Error())
			transCallback("翻译错误, 请查看日志", note)
			return
		}
		transRes = result
		note = "auto -> en"
	} else {
		note = "en -> zh"
		result, err := translator.Translate(str, "auto", "zh-cn")
		if err != nil {
			log.E(err.Error())
			transCallback("翻译错误, 请查看日志", note)
			return
		}
		transRes = result
		note = "auto -> zh-cn"
	}
	transCallback(transRes, note)
}
