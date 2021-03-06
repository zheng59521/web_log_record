package main

/*
函数定义
*/

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgutz/str"
)

// 逐行读取access日志文件
func ReadLogTOLine(params *CmfParams, logChannel chan string) error {
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		log.Infof("readLogTOLine can not read")
	}
	bufferRead := bufio.NewReader(fd)
	count := 0
	for {
		line, err := bufferRead.ReadString('\n')
		if err != nil {
			if err == io.EOF { // 文件读取完毕
				log.Infof("文件读取完毕, raedline:%d", count)
				time.Sleep(10 * time.Second)
			} else {
				log.Warningln("readLogTOLine read line log have error")
			}
		}
		count++
		// log.Infoln("line", line)
		// redisPool.Cmd("set", "name", count)
		logChannel <- line

	}
	defer fd.Close()
	return nil
}

/*
分析日志
分析URL参数,返回参数结构体
*/
func fomatURLParams(logStr string) JsData {
	logStr = strings.TrimSpace(logStr)
	// fmt.Println("logStr", logStr)
	pos1 := strings.Index(logStr, "/log.")
	if pos1 == -1 {
		return JsData{}
	}
	pos1 += len(HANDLE_LOG) - 1
	pos2 := str.IndexOf(logStr, " 200 43", pos1)

	d := str.Substr(logStr, pos1, pos2-pos1)

	urlObj, err := url.Parse("http://localhost/" + d)
	if err != nil {
		return JsData{}
	}
	data := urlObj.Query()

	// log.Infof("url :" + data.Get("url") + "refer: " + data.Get("refer") + "ua:" + data.Get("ua"))
	// fmt.Println("url :" + data.Get("url") + "refer: " + data.Get("refer") + "ua:" + data.Get("ua"))
	return JsData{
		data.Get("time"),
		data.Get("url"),
		data.Get("refer"),
		data.Get("ua"),
		data.Get("ip"),
	}

}

/*
日志处理
*/
func LogHandle(logChannel chan string, pvChannel chan URLData, uvChannel chan URLData) {
	num := 0
	for logStr := range logChannel {
		num++
		// 切割逐行分析出的文本
		data := fomatURLParams(logStr)
		// fmt.Println("url :" + data.url + "refer: " + data.refer + "ua:" + data.ua)
		if data.url != "" {
			// 生成用户uid
			hasher := md5.New()
			hasher.Write([]byte(data.ip))
			uid := hex.EncodeToString(hasher.Sum(nil))
			urlObj := URLData{data, uid, formatURLData(data.url, data.time, data.ip)}
			// log.Infof("uid:" + uid + " url :" + data.url + " refer: " + data.refer + " ua:" + data.ua + " ip:" + data.ip)
			// log.Infof("time", data.time)
			// fmt.Println("node is: ", urlObj.unode)
			pvChannel <- urlObj
			uvChannel <- urlObj
		}
	}
}

/*
pv|uv
分析URL, 返回URL数据
*/
func formatURLData(u, t, ip string) URLNODE {
	str1 := "/id/"
	str2 := ".html"
	cutStr := func(start int, DEF string, types string) URLNODE {
		// end := len(DEF)
		// URLType := str.Substr(u, start+1, end-1)
		IDstart := str.IndexOf(u, str1, 0)
		IDend := str.IndexOf(u, str2, 0)
		ID := 1 // redis 0默认为无效数字
		if IDstart != -1 && IDend != -1 {
			IDstart = IDstart + len(str1)
			IDstr := strings.Split(str.Substr(u, IDstart, IDend), ".")[0]
			ID, _ = strconv.Atoi(IDstr)
		}
		// log.Infof("url is" + u + "type is " + URLType + "id is" + strconv.Itoa(ID))
		return URLNODE{types, ID, u, t, ip}
	}

	u = strings.TrimSpace(u)
	// fmt.Println("url is: ", u)
	if pos1 := str.IndexOf(u, HANDLE_INDEX, 0); pos1 != -1 { // 如果是首页
		return cutStr(pos1, HANDLE_INDEX, "index")
	} else if pos2 := str.IndexOf(u, HANDLE_LIST, 0); pos2 != -1 { // 如果是列表页
		return cutStr(pos2, HANDLE_LIST, "list")
	} else if pos3 := str.IndexOf(u, HANDLE_ARTICLE, 0); pos3 != -1 { // 如果是文章页
		return cutStr(pos3, HANDLE_ARTICLE, "article")
	} else { // 默认首页
		return URLNODE{"index", 1, u, t, ip}
	}
	return URLNODE{}

}

/*
处理pv
storage<-channel<-pv数据
*/
func PvCounter(pvChannel chan URLData, storageChannel chan StorageBlock) {
	for data := range pvChannel {
		if data.data.url == "" {
			continue
		}
		storageChannel <- StorageBlock{"uv", data.unode}
	}
	DataStorage(storageChannel)
}

/*
处理uv
*/
func UvCounter(uvChannel chan URLData, storageChannel chan StorageBlock) {

	for data := range uvChannel {
		if data.data.url == "" {
			continue
		}
		// today := formatData(data.data.time, "day")
		// log.Infof("log" + today + data.uid)
		// num, err := redisPool.Cmd("PFADD", "log"+today, data.uid).Int()
		// if err != nil {
		// 	log.Warningf("redis HyperLogLog add fail")
		// }
		// if num != 1 {
		// 	continue
		// }
		storageChannel <- StorageBlock{"uv", data.unode}
		DataStorage(storageChannel)
	}
}

/*
创建存储器
格式化数据存入redis
*/
func DataStorage(storageChannel chan StorageBlock) {
	redisPool.Cmd("AUTH", "zjhredis") // 认证

	for it := range storageChannel {
		// prefix := block.counterType + "_"
		// fmt.Printf("storageChannel %+v\n", it)
		dJson, err := json.Marshal(it)
		if err != nil {
			fmt.Println("error", err.Error())
		}
		fmt.Println(string(dJson))
		// switch it.pType {
		// case "pv":
		// 	redisPool.Cmd("AUTH", "zjhredis")
		// case "uv":
		// 	redisPool.Cmd("AUTH", "zjhredis")
		// }
	}
}

/*
格式化时间
时间戳转对应格式时间字符串
*/
func formatData(dataStr, tT string) string {
	var item string
	timestamp, _ := strconv.ParseInt(dataStr, 10, 64)
	tm := time.Unix(timestamp, 0)
	switch tT {
	case "day":
		item = "2006-01-02"
		break
	}
	newTime := tm.Format(item)
	return newTime
}

/*
开始函数
*/
func start() {
	// 初始化channel, 用于传递数据
	var (
		logChannel     = make(chan string, params.routineNum)
		pvChannel      = make(chan URLData, params.routineNum)
		uvChannel      = make(chan URLData, params.routineNum)
		storageChannel = make(chan StorageBlock, 2*params.routineNum)
	)
	/*
		创建日志消费者
		处理哪个日志, 用多少协程
	*/
	go ReadLogTOLine(params, logChannel)

	// 创建一组日志处理
	go LogHandle(logChannel, pvChannel, uvChannel)

	// // 创建PV UV统计器
	go PvCounter(pvChannel, storageChannel)
	go UvCounter(uvChannel, storageChannel)
}
