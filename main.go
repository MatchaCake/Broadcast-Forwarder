package main

import (
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/Broadcast-Forwarder/utils"
	"gopkg.in/yaml.v2"
)

var (
	config = flag.String("bilibiliConfig", "", "config file path")
)

type Config struct {
	CommonConfig      `yaml:"common"`
	BilibiliConfig    `yaml:"bilibili"`
	TwitCastingConfig `yaml:"twitcasting"`
	YoutubeConfig     `yaml:"youtube"`
	TranslateConfig   `yaml:"translate"`
}

type CommonConfig struct {
	ChannelSize int `yaml:"channel_size"`
}

type BilibiliConfig struct {
	Name    string        `yaml:"name"`
	RoomId  string        `yaml:"room_id"`
	Timeout time.Duration `yaml:"timeout"`
}

type TwitCastingConfig struct {
	Name     string        `yaml:"name"`
	ClientId string        `yaml:"client_id"`
	Timeout  time.Duration `yaml:"timeout"`
}

type YoutubeConfig struct {
	Name      string        `yaml:"name"`
	ChannelId string        `yaml:"channel_id"`
	Timeout   time.Duration `yaml:"timeout"`
}

type TranslateConfig struct {
	AppId  string `yaml:"app_id"`
	Secret string `yaml:"secret"`
}

var msg chan string

func main() {
	flag.Parse()

	if len(*config) == 0 {
		log.Fatal("please specify config file")
	}

	Config := &Config{}
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, Config)
	if err != nil {
		log.Fatal(err)
	}

	msg = make(chan string, Config.ChannelSize)
	defer func() {
		close(msg)
	}()

	yMonitor := utils.NewYoutubeMonitor(Config.TranslateConfig.AppId, Config.TranslateConfig.Secret)
	bMonitor := utils.NewBilibiliMonitor()
	tMonitor := utils.NewTwitCastingMonitor(Config.TranslateConfig.AppId, Config.TranslateConfig.Secret)

	go func() {
		time.Sleep(Config.YoutubeConfig.Timeout)
		yMonitor.Monitor(&msg, Config.YoutubeConfig.Name, Config.YoutubeConfig.ChannelId)
	}()

	go func() {
		time.Sleep(Config.BilibiliConfig.Timeout)
		bMonitor.Monitor(&msg, Config.BilibiliConfig.Name, Config.BilibiliConfig.RoomId)
	}()

	go func() {
		time.Sleep(Config.TwitCastingConfig.Timeout)
		tMonitor.Monitor(&msg, Config.TwitCastingConfig.Name, Config.TwitCastingConfig.ClientId)
	}()

}
