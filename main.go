package main

/*
应用主体
*/

import (
	"flag"
	"fmt"
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
	// fmt.Println(params)
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
	redisPool, err = pool.New("tcp", "47.75.102.104:6379", 10)
	if err != nil {
		fmt.Println(err)
		log.Errorf("redis connect error!")
	} else {
		redisPool.Cmd("AUTH", "zjhredis") // 认证redis
		go func() {
			for {
				redisPool.Cmd("AUTH", "zjhredis")
				// redisPool.Cmd("ping")  .String()
				str := redisPool.Cmd("ping").String()
				fmt.Println("ping:", str)
				// 每隔10s ping一下redis服务,防止挂掉
				time.Sleep(time.Second * 10)
			}
		}()
	}
}
