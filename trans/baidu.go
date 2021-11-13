package trans

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/TransUtils/conf"
	"github.com/Ericwyn/TransUtils/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Md5(s string) string { //计算md5的值
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
func u2s(form string) (to string, err error) { //unicode转字符串
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

type BaiduTransResult struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
}
type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

var baiduAppId string
var baiduAppSecret string

func translate(word string) []byte { //调用api进行翻译
	if baiduAppId == "" || baiduAppSecret == "" {
		baiduAppId = viper.GetString(conf.ConfigKeyBaiduTransAppId)
		baiduAppSecret = viper.GetString(conf.ConfigKeyBaiduTransAppSecret)
	}

	data := make(url.Values)
	data["q"] = []string{word}
	data["from"] = []string{"auto"}
	data["to"] = []string{"auto"}
	data["appid"] = []string{baiduAppId}
	salt := "65"
	data["salt"] = []string{salt}
	s := baiduAppId + word + salt + baiduAppSecret //密匙
	sign := Md5(s)
	data["sign"] = []string{sign}
	res, err := http.PostForm("http://api.fanyi.baidu.com/api/trans/vip/translate", data)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//str := string(body)
	return body
}

type TransResCallback func(result string, note string)

func BaiduTrans(words string, callback TransResCallback) { //翻译函数
	body := translate(words)

	log.D("翻译结果", string(body))

	var transResult BaiduTransResult
	err := json.Unmarshal(body, &transResult)
	if err != nil {
		callback("error: "+err.Error(), "翻译错误")
	}

	res := ""
	for _, result := range transResult.TransResult {
		res += result.Dst + "\n"
	}

	callback(res, ""+transResult.From+" --> "+transResult.To)
}
