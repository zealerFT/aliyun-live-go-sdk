package live

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	live := NewServer(
		AppNameOption("live"),
		PushKeyOption("pushkey"),
		PlayKeyOption("playkey"),
		PushDomainOption("push.aliyun.com"),
		PlayDomainOption("pull.aliyun.com"),
		DelayOption(false),
		UidOption("1"),
		SnowflakeOption("127.0.0.1"),
	)
	// 可自己选择时区
	expireTime := 20*3600 + 60
	liveDuration := ShanghaiTime(expireTime)
	streamName := live.Snowflake.Generate().String()
	pushurl := live.PushUrl(streamName, liveDuration)
	playurl := live.PlayUrl(streamName, liveDuration)
	rtsPushUrl := live.RtsPushUrl(streamName, liveDuration)
	rtsPlayUrl := live.RtsPlayUrl(streamName, liveDuration)
	flvPlayUrl := live.FlvPlayUrl(streamName, liveDuration)
	m3u8PlayUrl := live.M3u8PlayUrl(streamName, liveDuration)
	fmt.Println("pushurl:", pushurl)
	fmt.Println("playurl:", playurl)
	fmt.Println("rtsPushUrl:", rtsPushUrl)
	fmt.Println("rtsPlayUrl:", rtsPlayUrl)
	fmt.Println("flvPlayUrl:", flvPlayUrl)
	fmt.Println("m3u8PlayUrl:", m3u8PlayUrl)
	fmt.Println("-------------------------")
	live.Delay = true
	pushurl = live.PushUrl(streamName, liveDuration)
	playurl = live.PlayUrl(streamName, liveDuration)
	rtsPushUrl = live.RtsPushUrl(streamName, liveDuration)
	rtsPlayUrl = live.RtsPlayUrl(streamName, liveDuration)
	flvPlayUrl = live.FlvPlayUrl(streamName, liveDuration)
	m3u8PlayUrl = live.M3u8PlayUrl(streamName, liveDuration)
	fmt.Println("pushurl-Delay:", pushurl)
	fmt.Println("playurl-Delay:", playurl)
	fmt.Println("rtsPushUrl-Delay:", rtsPushUrl)
	fmt.Println("rtsPlayUrl-Delay:", rtsPlayUrl)
	fmt.Println("flvPlayUrl-Delay:", flvPlayUrl)
	fmt.Println("m3u8PlayUrl-Delay:", m3u8PlayUrl)
}
