package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Translator struct {
	client *http.Client
}

func (t *Translator) Translate(source string, appId string, secret string, target interface{}) error {
	salt := rand.Intn(time.Now().Second())
	merge := appId + source + strconv.Itoa(salt) + secret
	h := md5.New()
	sign := hex.EncodeToString(h.Sum([]byte(merge)))
	req := fmt.Sprintf("https://fanyi-api.baidu.com/api/trans/vip/translate?q=%s&from=%s&to=%s&appid=%s&salt=%v&sign=%s",
		source, "jp", "zh", appId, salt, sign)
	result, err := t.client.Get(req)
	if err != nil {
		log.Fatalf("failed to call translate API, err:%v", err)
		return err
	}
	defer func() {
		err := result.Body.Close()
		if err != nil {
			log.Fatalf("failed to close request body, err:%v", err)
		}
	}()

	return json.NewDecoder(result.Body).Decode(target)
}

func NewTranslator() *Translator {
	translator := &Translator{client: &http.Client{
		Timeout: 5 * time.Second,
	}}
	return translator
}
