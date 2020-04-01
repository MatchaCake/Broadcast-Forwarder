package utils

import "fmt"

type Notification struct {
	videoId      string
	videoUrl     string
	channelName  string
	channelId    string
	channelUrl   string
	title        string
	titleZh      string
	action       string
	publishTime  string
	deleteTime   string
	updateTime   string
	streamDetail Stream
}

func (n *Notification) String() string {
	if n.action == "delete" {
		return fmt.Sprintf("%s 删除了一个视频：https://www.youtube.com/watch?v=%s ，在%s\n", n.channelName, n.videoId, n.deleteTime)
	} else if n.action == "update" {
		return fmt.Sprintf("%s 上传了一个视频：https://www.youtube.com/watch?v=%s ，发布于%s ，"+
			" 更新于%s\n", n.channelName, n.videoId, n.publishTime, n.updateTime)
	}
	return ""
}

func (n *Notification) GetNotification() string {
	if n.streamDetail.IsLiveStream {
		return fmt.Sprintf("%s 添加了新直播。\n%s\n始于%s\n标题%s\n翻译%s\n",
			n.channelName, n.videoUrl, n.streamDetail.ScheduledStartTime, n.title, n.titleZh)
	}
	return ""
}
