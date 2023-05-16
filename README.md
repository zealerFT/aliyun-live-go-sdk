# aliyun-live-go-sdk
阿里云视频直播新版sdk，简单方便，无需使用ak/sk。
提供推流和拉流的基础接口，包括rtmp，rts，flv和m3u8，后续会不断更新其他功能。

## 官方文档
url鉴权：https://help.aliyun.com/document_detail/199349.html
阿里云视频直播后台配置url鉴权：https://help.aliyun.com/document_detail/85018.htm

## 使用
``` shell
	live := NewServer(
		AppNameOption("live"), // 直播app名
		PushKeyOption("pushkey"),
		PlayKeyOption("playkey"),
		PushDomainOption("push.aliyun.com"),
		PlayDomainOption("pull.aliyun.com"),
		DelayOption(false),
		UidOption("1"), // 可选参数，用户的Uid，方便排查分析
		SnowflakeOption("127.0.0.1"), // 雪花算法，传入机器唯一标识，可以是机器ip
	)
	// 直播链接有效期
	expireTime := 20*3600 + 60
	liveDuration := ShanghaiTime(expireTime)
	// 区分每场直播
	streamName := live.Snowflake.Generate().String()
	// live.Delay = true 可选算法延迟直播方式
	pushurl = live.PushUrl(streamName, liveDuration)
	playurl = live.PlayUrl(streamName, liveDuration)
```
## 测试
可以使用ffmpeg来测试推流和拉流是否成功
```shell
ffmpeg -i /Users/fei.teng/Desktop/test.mp4  -c:v libx264 -c:a aac -f flv "pushurl" #推流
ffplay -i "playurl" # 拉流
```
