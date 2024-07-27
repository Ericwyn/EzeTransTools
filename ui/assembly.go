package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"github.com/Ericwyn/EzeTranslate/ui/resource/cusWidget"
	"github.com/spf13/viper"
)

func buildFormatCheckBox() *fyne.Container {
	return container.NewHBox(
		widget.NewLabel("输入优化    "),
		cusWidget.CreateCheckGroup(
			[]cusWidget.LabelAndInit{
				{"注释", viper.GetBool(conf.ConfigKeyFormatAnnotation)},
				{"空格", viper.GetBool(conf.ConfigKeyFormatSpace)},
				{"回车", viper.GetBool(conf.ConfigKeyFormatCarriageReturn)},
				{"驼峰", viper.GetBool(conf.ConfigKeyFormatCamelCase)},
			},
			true,  // 横向
			false, // 单选
			func(label string, checked bool) {
				if label == "注释" {
					viper.Set(conf.ConfigKeyFormatAnnotation, checked)
					log.I("输入优化: 注释: ", checked)
				} else if label == "空格" {
					viper.Set(conf.ConfigKeyFormatSpace, checked)
					log.I("输入优化: 空格: ", checked)
				} else if label == "回车" {
					viper.Set(conf.ConfigKeyFormatCarriageReturn, checked)
					log.I("输入优化: 回车: ", checked)
				} else if label == "驼峰" {
					viper.Set(conf.ConfigKeyFormatCamelCase, checked)
					log.I("输入优化: 驼峰: ", checked)
				}
				conf.SaveConfig()
			},
		),
	)
}

func buildTransApiCheckBox() *fyne.Container {
	return container.NewHBox(
		widget.NewLabel("翻译结果  "),
		cusWidget.CreateCheckGroup(
			[]cusWidget.LabelAndInit{
				{"Google", viper.GetString(conf.ConfigKeyTranslateSelect) == "google"},
				{"Baidu", viper.GetString(conf.ConfigKeyTranslateSelect) == "baidu"},
				{"Youdao", viper.GetString(conf.ConfigKeyTranslateSelect) == "youdao"},
				{"OpenAI", viper.GetString(conf.ConfigKeyTranslateSelect) == "openai"},
			},
			true, // 横向
			true, // 单选
			func(label string, checked bool) {
				if label == "Google" {
					viper.Set(conf.ConfigKeyTranslateSelect, "google")
				} else if label == "Baidu" {
					viper.Set(conf.ConfigKeyTranslateSelect, "baidu")
				} else if label == "Youdao" {
					viper.Set(conf.ConfigKeyTranslateSelect, "youdao")
				} else if label == "Youdao" {
					viper.Set(conf.ConfigKeyTranslateSelect, "youdao")
				} else if label == "OpenAI" {
					viper.Set(conf.ConfigKeyTranslateSelect, "openai")
				}
				e := viper.WriteConfig()
				if e != nil {
					log.E("配置文件保存失败")
					log.E(e)
				}
			},
		),
	)
}

func buildToLangDropDown() *fyne.Container {
	return container.NewHBox(
		//widget.NewLabel(""),
		cusWidget.CreateDropDown(
			[]cusWidget.LabelAndInit{
				{"自动", conf.ToLang == ""},
				{"中文", conf.ToLang == string(strutils.Chinese)},
				{"英文", conf.ToLang == string(strutils.English)},
			},
			func(label string, checked bool) {
				if checked {
					if label == "自动" {
						conf.ToLang = ""
					} else if label == "中文" {
						conf.ToLang = string(strutils.Chinese)
					} else if label == "英文" {
						conf.ToLang = string(strutils.English)
					}
				}
			},
		),
	)
}
