package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/pkg/errors"
)

type Stream struct {
	VideoId            string
	ActualStartTime    string
	ScheduledStartTime string
	ConcurrentViewers  int
	ActiveLiveChatId   string
	IsLiveStream       bool
	client             *http.Client
}

func NewStream(videoId string) *Stream {
	return &Stream{
		VideoId: videoId,
		client:  &http.Client{Timeout: 1 * time.Second},
	}
}

func (s *Stream) GetStreamDetail(key string) error {
	url := "https://www.googleapis.com/youtube/v3/videos"
	req := fmt.Sprintf("%s?part=%s&id=%s&key=%s", url, "liveStreamingDetails", s.VideoId, key)
	result, err := s.client.Get(req)
	if err != nil {
		log.Fatalf("failed to call youtube stream API, err:%v", err)
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
		log.Fatalf("failed to read stream result body, err:%v", err)
		return err
	}
	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		log.Fatalf("failed to parse stream result json, err:%v", err)
		return err
	}

	exists := jsonParsed.Exists("items")
	if !exists {
		return errors.New("error: cannot get items from stream result json")
	}

	exists = jsonParsed.Exists("items", "0", "liveStreamingDetails")
	if exists {
		s.IsLiveStream = true
	}
	s.ConcurrentViewers, _ = jsonParsed.Search("items", "0", "liveStreamingDetails", "concurrentViewers").Data().(int)
	scheduledStartTime, ok := jsonParsed.Search("items", "0", "liveStreamingDetails", "scheduledStartTime").Data().(string)
	if ok {
		s.ScheduledStartTime = ParseToBeijingTime(scheduledStartTime)
	}
	actualStartTime, ok := jsonParsed.Search("items", "0", "liveStreamingDetails", "actualStartTime").Data().(string)
	if ok {
		s.ActualStartTime = ParseToBeijingTime(actualStartTime)
	}
	s.ActiveLiveChatId, _ = jsonParsed.Search("items", "0", "liveStreamingDetails", "activeLiveChatId").Data().(string)
	return nil
}
