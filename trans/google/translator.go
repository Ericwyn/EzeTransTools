package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ericwyn/EzeTranslate/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Translator struct {
	client             *http.Client
	googleTransApiPath string
	httpProxyUrl       string
}

func NewTranslatorWithConfig(googleTransApiPath, httpProxyUrl string) Translator {
	var client http.Client
	if httpProxyUrl != "" {
		log.D("init google translator with proxy: " + httpProxyUrl)
		// 设置代理地址
		proxyURL, err := url.Parse(httpProxyUrl)
		if err != nil {
			log.E("set google translator proxy URL error: ", err.Error())
		}

		// 创建自定义的 Transport
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		// 创建带有自定义 Transport 的 http.Client
		client = http.Client{
			Transport: transport,
		}
	} else {
		client = http.Client{}
	}

	return Translator{
		client:             &client,
		googleTransApiPath: googleTransApiPath,
		httpProxyUrl:       httpProxyUrl,
	}
}

func (t Translator) Translate(text, sourceLang, targetLang string) (string, error) {
	var result []interface{}
	var translated []string

	urlStr := fmt.Sprintf(
		t.googleTransApiPath+"/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		sourceLang,
		targetLang,
		url.QueryEscape(text),
	)

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	res, err := t.client.Do(req)

	if err != nil {
		log.E("google translate api error get translate : ", err.Error())
		return "err", errors.New("error getting translate.googleapis.com")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.E("google translate api error reading response body : ", err.Error())
		return "err", errors.New("error reading response body")
	}

	if res.StatusCode != 200 {
		return "err", errors.New("google translate api, api request failed" + string(rune(res.StatusCode)))
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", errors.New("google translate api error unmarshaling body" + string(body))
	}

	if len(result) > 0 {
		data := result[0]
		for _, slice := range data.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				translated = append(translated, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		return strings.Join(translated, ""), nil
	}
	return "err", errors.New("translation not found")
}
