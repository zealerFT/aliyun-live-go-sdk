package live

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

// Server
// @url鉴权：https://help.aliyun.com/document_detail/199349.html
type Server struct {
	PushKey    string          // 推流鉴权Key - 查看方式：在阿里云视频直播后台-域名管理-访问控制-URL鉴权-主KEY
	PushDomain string          // 推流域名
	PlayKey    string          // 播放（拉流）鉴权Key - 查看方式同PushKey
	PlayDomain string          // rtmp播放（拉流）域名 常规推流地址
	AppName    string          // 应用名称
	Delay      bool            // 是否延迟播放
	Uid        string          // 可以传递userid在url展示，方便排查问题，但会暴露数据，谨慎使用
	Snowflake  *snowflake.Node // 雪花算法的streamName，只做参考，需直接传参
}

func NewServer(options ...Option) *Server {
	s := &Server{}
	for _, option := range options {
		option(s)
	}
	return s
}

type Option func(*Server)

func AppNameOption(appName string) Option {
	return func(s *Server) {
		s.AppName = appName
	}
}

func PushKeyOption(pushKey string) Option {
	return func(s *Server) {
		s.PushKey = pushKey
	}
}

func PlayKeyOption(playKey string) Option {
	return func(s *Server) {
		s.PlayKey = playKey
	}
}

func PushDomainOption(pushDomain string) Option {
	return func(s *Server) {
		s.PushDomain = pushDomain
	}
}

func PlayDomainOption(playKey string) Option {
	return func(s *Server) {
		s.PlayDomain = playKey
	}
}

func DelayOption(delay bool) Option {
	return func(s *Server) {
		s.Delay = delay
	}
}

func UidOption(uid string) Option {
	return func(s *Server) {
		s.Uid = uid
	}
}

// SnowflakeOption streamName可以使用雪花算法
func SnowflakeOption(machineIp string) Option {
	return func(s *Server) {
		s.Snowflake = MustNewSnowflake(machineIp)
	}
}

// PushUrl 推流地址
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) PushUrl(streamName string, liveDuration int64) string {
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s-%s-0-%s-%s", s.AppName, streamName, liveDurationStr, s.Uid, s.PushKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("rtmp://%s/%s/%s?auth_key=%s-0-%s-%s", s.PushDomain, s.AppName, streamName, liveDurationStr, s.Uid, pathMd5Byte)
}

// RtsPushUrl RTS推流：超低延时直播（UDP）流量/带宽费用和标准直播不同
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) RtsPushUrl(streamName string, liveDuration int64) string {
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s-%s-0-%s-%s", s.AppName, streamName, liveDurationStr, s.Uid, s.PushKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("artc://%s/%s/%s?auth_key=%s-0-%s-%s", s.PushDomain, s.AppName, streamName, liveDurationStr, s.Uid, pathMd5Byte)
}

// PlayUrl 播流（拉流）地址
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) PlayUrl(streamName string, liveDuration int64) string {
	if s.Delay {
		streamName = fmt.Sprintf("%s%s", streamName, "-alidelay")
	}
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s-%s-0-%s-%s", s.AppName, streamName, liveDurationStr, s.Uid, s.PlayKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("rtmp://%s/%s/%s?auth_key=%s-0-%s-%s", s.PlayDomain, s.AppName, streamName, liveDurationStr, s.Uid, pathMd5Byte)
}

// RtsPlayUrl 超低延时直播 - 播流（拉流）地址
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) RtsPlayUrl(streamName string, liveDuration int64) string {
	if s.Delay {
		streamName = fmt.Sprintf("%s%s", streamName, "-alidelay")
	}
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s-%s-0-%s-%s", s.AppName, streamName, liveDurationStr, s.Uid, s.PlayKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("artc://%s/%s/%s?auth_key=%s-0-%s-%s", s.PlayDomain, s.AppName, streamName, liveDurationStr, s.Uid, pathMd5Byte)
}

// FlvPlayUrl 播流（拉流）地址 - https协议，适用于在web浏览器播放
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) FlvPlayUrl(streamName string, liveDuration int64) string {
	if s.Delay {
		streamName = fmt.Sprintf("%s%s", streamName, "-alidelay")
	}
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s%s-%s-0-%s-%s", s.AppName, streamName, ".flv", liveDurationStr, s.Uid, s.PlayKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("https://%s/%s/%s%s?auth_key=%s-0-%s-%s", s.PlayDomain, s.AppName, streamName, ".flv", liveDurationStr, s.Uid, pathMd5Byte)
}

// M3u8PlayUrl 播流（拉流）地址
// @streamName 区分同一AppName下不同直播
// @liveDuration 到期时间戳-秒 直播链接有效时间，可在阿里云后台设置，默认24小时。（鉴权URL生成后的有效时长，最小需设置为1分钟，无上限限制。）
func (s *Server) M3u8PlayUrl(streamName string, liveDuration int64) string {
	if s.Delay {
		streamName = fmt.Sprintf("%s%s", streamName, "-alidelay")
	}
	liveDurationStr := strconv.FormatInt(liveDuration, 10)
	path := fmt.Sprintf("/%s/%s%s-%s-0-%s-%s", s.AppName, streamName, ".m3u8", liveDurationStr, s.Uid, s.PlayKey)
	pathMd5Byte := Md5(path)
	return fmt.Sprintf("https://%s/%s/%s%s?auth_key=%s-0-%s-%s", s.PlayDomain, s.AppName, streamName, ".m3u8", liveDurationStr, s.Uid, pathMd5Byte)
}

func Md5(str string) string {
	if str == "" {
		return ""
	}
	h := md5.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	bigInt := new(big.Int).SetBytes(bs)
	// 16进制，不足32位高位补零
	result := bigInt.Text(16)
	for len(result) < 32 {
		result = "0" + result
	}
	return result
}

func MustNewSnowflake(machineIp string) *snowflake.Node {
	ipSegs := strings.Split(machineIp, ".")
	number := ipSegs[len(ipSegs)-1]
	no, _ := strconv.Atoi(number)
	node, err := snowflake.NewNode(int64(no))
	if err != nil {
		log.Fatalf("雪花算法初始化失败～ %v", err)
	}
	return node
}

func ShanghaiTime(expireTime int) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt := time.Now()
	unixTime := tt.In(loc).Unix()
	return unixTime + int64(expireTime)
}
