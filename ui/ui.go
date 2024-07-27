package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/Ericwyn/EzeTranslate/conf"
	"github.com/Ericwyn/EzeTranslate/ipc"
	"github.com/Ericwyn/EzeTranslate/log"
	"github.com/Ericwyn/EzeTranslate/ocr"
	"github.com/Ericwyn/EzeTranslate/strutils"
	"github.com/Ericwyn/EzeTranslate/trans"
	"github.com/Ericwyn/EzeTranslate/trans/baidu"
	"github.com/Ericwyn/EzeTranslate/trans/google"
	"github.com/Ericwyn/EzeTranslate/trans/openai"
	"github.com/Ericwyn/EzeTranslate/trans/youdao"
	"github.com/Ericwyn/EzeTranslate/ui/resource"
	"github.com/Ericwyn/GoTools/shell"
	"github.com/spf13/viper"
	"runtime"
	"strings"
	"time"
)

var mainApp fyne.App

func StartApp(xclip bool, ocr bool) {

	shell.Debug(true)

	conf.InitConfig()

	if trySendMessage(ipc.IpcMessagePing) {
		// 如果已经有其他翻译进程的话, 就发送一下消息，然后退出就好了
		sendSocketMessage(xclip, ocr, false)
		return
	} else {
		// 开启 server 监听来自其他进程的翻译请求
		startUnixSocketServer()
		// 此处需要异步，需要等 app 界面起来之后再去做消息发送
		go sendSocketMessage(xclip, ocr, true)
		openNewApp()
	}
}

func sendSocketMessage(xclip bool, ocr bool, sleep bool) {
	// 重新发送一遍参数
	// 如果启动时候带有参数的话，那么就应该直接获取一遍选中的文字并进行翻译
	if xclip {
		if sleep {
			time.Sleep(time.Millisecond * 500)
		}
		trySendMessage(ipc.IpcMessageNewSelection)
	} else if ocr {
		if sleep {
			time.Sleep(time.Millisecond * 500)
		} // 只获取文字，但是不进行翻译
		trySendMessage(ipc.IpcMessageOcr)
	}

}

func openNewApp() {
	// 设置整个 app 的信息
	mainApp = app.New()
	mainApp.SetIcon(resource.ResourceIcon())
	mainApp.Settings().SetTheme(&resource.CustomerTheme{})

	// 启动主页面
	if viper.GetBool(conf.ConfigKeyMiniMode) {
		showMiniUi(true)
	} else {
		showHomeUi(true)
	}
}

func startTrans() {

	// 针对不同模式有不同的获取 formatText 的方法
	var formatText string
	if homeInputBox != nil {
		formatText = strutils.FormatInputBoxText(homeInputBox.Text)
	} else {
		// 如果以 mini 模式启动的话，就只获取选择的文字就可以了
		formatText = strutils.FormatInputBoxText(trans.GetSelection())
	}

	var resBox *widget.Entry
	var resNoteLabel *widget.Label
	if homeTransResBox != nil {
		resBox = homeTransResBox
		resNoteLabel = homeNoteLabel
	} else if miniTransResBox != nil {
		resBox = miniTransResBox
		resNoteLabel = miniNoteLabel
	} else {
		log.E("resBox 为空")
		return
	}

	if strings.Trim(formatText, " ") == "" {
		resBox.SetPlaceHolder("请输入需要翻译的内容")
		return
	}

	resBox.SetText("")
	resBox.SetPlaceHolder("正在翻译..........")

	handleTransResult := func(result string, note string) {
		fmt.Println("翻译结果:", result)
		resBox.SetText(result)
		resNoteLabel.SetText(note)
	}

	var toLang = strutils.Lang(conf.ToLang)

	if viper.GetString(conf.ConfigKeyTranslateSelect) == "google" {
		go google.Translate(formatText, toLang, handleTransResult)
	} else if viper.GetString(conf.ConfigKeyTranslateSelect) == "baidu" {
		go baidu.Translate(formatText, toLang, handleTransResult)
	} else if viper.GetString(conf.ConfigKeyTranslateSelect) == "youdao" {
		go youdao.Translate(formatText, toLang, handleTransResult)
	} else if viper.GetString(conf.ConfigKeyTranslateSelect) == "openai" {
		go openai.Translate(formatText, toLang, handleTransResult)
	}

}

func trySendMessage(message ipc.IpcMessage) bool {
	if runtime.GOOS != "linux" {
		log.D("not linux, don't send socket msg")
		return false
	}
	err := ipc.SendMessage(message)
	if err == nil {
		log.D("已发送给其他翻译进程 : " + string(message))
		return true
	} else {
		log.D("IPC 消息发送失败")
		return false
	}
}

// 开启一个 UnixSocketServer, 接收 IPC 消息
func startUnixSocketServer() {
	if runtime.GOOS != "linux" {
		log.D("not linux, don't start socket server")
		return
	}
	go ipc.StartUnixSocketListener(func(message ipc.IpcMessage) {
		log.D("接收到 IPC 消息 : " + string(message))
		switch message {
		case ipc.IpcMessagePing:
			break
		case ipc.IpcMessageOcr:
			setOcrTextToInputBox()
		case ipc.IpcMessageOcrAndTrans:
			if setOcrTextToInputBox() {
				startTrans()
			}
			break
		case ipc.IpcMessageNewSelection:
			if setSelectTextToInputBox() {
				startTrans()
				break
			}
		}
	})
}

// setWindowsFocus 请求焦点
func setWindowsFocus() (*EzeInputEntry, *widget.Entry) {
	var inputBox *EzeInputEntry
	var transResBox *widget.Entry

	if homeWindow != nil {
		// 请求焦点
		homeWindow.RequestFocus()
		inputBox = homeInputBox
		transResBox = homeTransResBox
	} else if miniWindow != nil {
		miniWindow.RequestFocus()
		transResBox = miniTransResBox
	}
	return inputBox, transResBox
}

func setSelectTextToInputBox() bool {
	inputBox, transResBox := setWindowsFocus()

	selectText := trans.GetSelection()
	selectText = strutils.FormatInputBoxText(selectText)
	fmt.Println("获取的划词:", selectText)

	miniSelectTextNow = selectText

	if inputBox != nil && strings.Trim(inputBox.Text, " ") ==
		strings.Trim(selectText, " ") {

		// 如果翻译框有数据的话，就不进行翻译
		if transResBox != nil && strings.Trim(transResBox.Text, " ") != "" {
			log.D("获取的划词与当前 homeInputBox 中文字一致，不进行翻译")
			return false
		}
	}

	if inputBox != nil {
		// 刷新当前数据
		inputBox.SetText(selectText)
	}
	return true
}

func setOcrTextToInputBox() bool {
	inputBox, transResBox := setWindowsFocus()

	ocrRes, successFlag := ocr.RunOcr()

	if successFlag {
		log.D("获取的 OCR: " + ocrRes)
		ocrRes = strutils.FormatInputBoxText(ocrRes)
		homeInputBox.SetText(ocrRes)
	} else {
		log.E("ocr 识别失败")
		return false
	}

	miniSelectTextNow = ocrRes

	if inputBox != nil && strings.Trim(inputBox.Text, " ") ==
		strings.Trim(ocrRes, " ") {

		// 如果翻译框有数据的话，就不进行翻译
		if transResBox != nil && strings.Trim(transResBox.Text, " ") != "" {
			log.D("获取的划词与当前 homeInputBox 中文字一致，不进行翻译")
			return false
		}
	}

	if inputBox != nil {
		// 刷新当前数据
		inputBox.SetText(ocrRes)
	}

	return true
}

func closeMiniWindow() {
	if miniWindow != nil {
		miniWindow.Close()
		miniWindow = nil
	}
}
