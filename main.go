package main

import (
	"flag"
	"os"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func main() {

	routineNum := flag.Int("channelLength", 200, "channel run num")
	flag.Parse()

	// 获取参数
	params = &CmfParams{logFilePath: ACCESS_LOG_FILE_PATH, routineNum: *routineNum}
	// 打日志(使用第三方库logrus)
	logFd, err := os.OpenFile(CHANNEL_LOG_FILE_PATH, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	} else {
		log.Errorln("日志记录文件打开失败")
		return
	}
	start()
	// 阻塞程序, 防止程序退出
	select {}
}
func init() {
	var err error
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
	redisPool, err = pool.New("tcp", "172.17.0.2:6379", 10)
	if err != nil {
		log.Errorf("redis connect error!")
	} else {
		go func() {
			for {
				// 每隔10s ping一下redis服务,防止挂掉
				time.Sleep(time.Second * 10)
			}
		}()
	}
}
