package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
)

type Translator struct {
	client *http.Client
}

func (t *Translator) Translate(source string, appId string, secret string) (string, error) {
	salt := rand.Intn(time.Now().Second())
	merge := appId + source + strconv.Itoa(salt) + secret
	h := md5.New()
	sign := hex.EncodeToString(h.Sum([]byte(merge)))
	req := fmt.Sprintf("https://fanyi-api.baidu.com/api/trans/vip/translate?q=%s&from=%s&to=%s&appid=%s&salt=%v&sign=%s",
		source, "jp", "zh", appId, salt, sign)
	result, err := t.client.Get(req)
	if err != nil {
		log.Fatalf("failed to call translate API, err:%v", err)
		return "", err
	}
	defer func() {
		err := result.Body.Close()
		if err != nil {
			log.Fatalf("failed to close request body, err:%v", err)
		}
	}()

	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatalf("failed to read translate result body, err:%v", err)
		return "", err
	}
	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		log.Fatalf("failed to parse translate result json, err:%v", err)
		return "", err
	}
	translateResult, ok := jsonParsed.Search("trans_result", "0", "dst").Data().(string)
	if ok {
		return translateResult, nil
	}
	return "", errors.New("failed to find translate result dst")
}

func NewTranslator() *Translator {
	translator := &Translator{client: &http.Client{
		Timeout: 5 * time.Second,
	}}
	return translator
}
