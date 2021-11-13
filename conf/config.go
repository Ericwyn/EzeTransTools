package conf

import (
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/spf13/viper"
)

const ConfigKeyBaiduTransAppId = "baiduTransAppId"
const ConfigKeyBaiduTransAppSecret = "baiduTransAppSecret"

const ConfigKeyFormatSpace = "formatSpace"
const ConfigKeyFormatCarriageReturn = "formatCarriageReturn"
const ConfigKeyFormatAnnotation = "formatAnnotation"

func InitConfig() {
	viper.SetDefault(ConfigKeyBaiduTransAppId, "baiduTransAppId-xxxxxxxxxxxxxxx")
	viper.SetDefault(ConfigKeyBaiduTransAppSecret, "baiduTransAppSecret-xxxxxxxxxxxxxxx")

	viper.SetDefault(ConfigKeyFormatSpace, false)
	viper.SetDefault(ConfigKeyFormatCarriageReturn, false)
	viper.SetDefault(ConfigKeyFormatAnnotation, false)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".conf")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.E("载入配置时候出错")
		panic(err)
	}
}

// 返回百度翻译 api 的 appId 和 appSecret
func GetBaiduTransApiMsg() (string, string) {
	return "", ""
}
