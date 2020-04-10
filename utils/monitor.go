package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
)

type YoutubeMonitor struct {
	isOnLive bool
	Title    string
	TitleZh  string
	client   *http.Client
	appId    string
	secret   string
}

type BilibiliMonitor struct {
	isOnLive bool
	Title    string
	client   *http.Client
}

type TwitCastingMonitor struct {
	isOnLive bool
	Title    string
	TitleZh  string
	client   *http.Client
	appId    string
	secret   string
}

func NewYoutubeMonitor(appId string, secret string) *YoutubeMonitor {
	return &YoutubeMonitor{
		client: &http.Client{Timeout: 1 * time.Second},
		appId:  appId,
		secret: secret,
	}
}

func NewBilibiliMonitor() *BilibiliMonitor {
	return &BilibiliMonitor{
		client: &http.Client{Timeout: 1 * time.Second},
	}
}

func NewTwitCastingMonitor(appId string, secret string) *TwitCastingMonitor {
	return &TwitCastingMonitor{
		client: &http.Client{Timeout: 1 * time.Second},
		appId:  appId,
		secret: secret,
	}
}

func (y *YoutubeMonitor) Monitor(msg *chan string, name string, channelId string) error {
	return nil
}

func (b *BilibiliMonitor) Monitor(msgChannel *chan string, name string, roomId string) error {
	req := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%s&from=room", roomId)
	result, err := b.client.Get(req)
	if err != nil {
		log.Fatalf("failed to call bilibili API, err:%v", err)
		return err
	}
	defer func() {
		err := result.Body.Close()
		if err != nil {
			log.Fatalf("failed to close result body, err:%v", err)
		}
	}()

	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatalf("failed to read bilibili result body, err:%v", err)
		return err
	}
	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		log.Fatalf("failed to parse bilibili result json, err:%v", err)
		return err
	}
	liveStatus, ok := jsonParsed.Search("data", "live_status").Data().(int)
	title, Ok := jsonParsed.Search("data", "title").Data().(string)
	if ok && Ok {
		if liveStatus == 1 && !b.isOnLive {
			addr := fmt.Sprintf("https://live.bilibili.com/%s", roomId)
			msg := fmt.Sprintf("%s正在B站直播：\"%s\"\n，地址：%s\n", name, title, addr)
			log.Println(msg)
			select {
			case *msgChannel <- msg:
				// do nothing
			default:
				log.Printf("msg input channel is full|msg:%v\n", msg)
				start := time.Now()
				*msgChannel <- msg
				log.Printf("msg input channel insert delay=%dμs|msg:%v\n",
					time.Since(start).Nanoseconds()/1e3, msg)
			}
			b.isOnLive = true
		} else if liveStatus == 0 && b.isOnLive {
			b.isOnLive = false
		}
		// do nothing if the live statuses are synced.
		return nil
	}
	return errors.New("failed to get live status from bilibili result")
}

func (t *TwitCastingMonitor) Monitor(msg *chan string, name string, clientId string) error {
	return nil
}
