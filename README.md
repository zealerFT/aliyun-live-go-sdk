# aliyun-live-go-sdk
阿里云视频直播新版sdk，简单方便，无需使用ak/sk。
提供推流和拉流的基础接口，包括rtmp，rts，flv和m3u8，后续会不断更新其他功能。

## 官方文档
- 阿里云直播后台：https://live.console.aliyun.com/
- url鉴权规则文档：https://help.aliyun.com/document_detail/199349.html
- 阿里云视频直播后台配置url鉴权教程：https://help.aliyun.com/document_detail/85018.htm

## 使用
```shell
go get github.com/zealerFT/aliyun-live-go-sdk@latest
```
``` go
        package live
	import (
		sdk "github.com/zealerFT/aliyun-live-go-sdk"
        )
	live := sdk.NewServer(
		sdk.AppNameOption("live"), // 直播app名
		sdk.PushKeyOption("pushkey"), // 推流鉴权key，在阿里云视频直播-域名管理-点击对应"推流"域名-访问控制-鉴权URL设置-主key
		sdk.PlayKeyOption("playkey"), // 推流鉴权key，在阿里云视频直播-域名管理-点击对应"拉流"域名-访问控制-鉴权URL设置-主key
		sdk.PushDomainOption("push.aliyun.com"), // 推流域名，在阿里云视频直播-域名管理里配置
		sdk.PlayDomainOption("pull.aliyun.com"), // 拉流域名，在阿里云视频直播-域名管理里配置
		sdk.DelayOption(false),
		sdk.UidOption("1"), // 可选参数，用户的Uid，方便排查分析
		sdk.SnowflakeOption("127.0.0.1"), // 雪花算法，传入机器唯一标识，可以是机器ip
	)
	// 直播链接有效期
	expireTime := 20*3600 + 60
	liveDuration := sdk.ShanghaiTime(expireTime)
	// 区分每场直播
	streamName := live.Snowflake.Generate().String()
	// live.Delay = true 可选算法延迟直播方式
	pushurl = live.PushUrl(streamName, liveDuration)
	playurl = live.PlayUrl(streamName, liveDuration)
```
## 测试
可以使用ffmpeg来测试推流和拉流是否成功，注意推流是推自己的流到阿里云推流地址，可以是c++等语言控制推流，也可以直接推录播的视频
```shell
ffmpeg -i /Users/fei.teng/Desktop/test.mp4  -c:v libx264 -c:a aac -f flv "pushurl" #推流
ffplay -i "playurl" # 拉流
```
