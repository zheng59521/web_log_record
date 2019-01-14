package main

import (
	"github.com/mediocregopher/radix.v2/pool"
)

var (
	params    *CmfParams
	redisPool *pool.Pool
)

const (
	ACCESS_LOG_FILE_PATH  = "./tmp/log.log"           // access日志存放路径
	CHANNEL_LOG_FILE_PATH = "./tmp/channel.log"       // 程序运行日志存放路径
	HANDLE_LOG            = "/log.gif?"               // 分隔符
	HANDLE_ARTICLE        = "/article/"               // 文章页分隔符
	HANDLE_LIST           = "/list/"                  // 列表页分隔符
	HANDLE_INDEX          = "/portal/list/index.html" // 首页分隔符
)

/*
js发送来的数据格式, 日志格式
*/
type JsData struct {
	time  string
	url   string
	refer string
	ua    string
}

/*
记录使用
*/
type CmfParams struct {
	logFilePath string // 日志文件路径
	routineNum  int    // 通道长度
}

/*
用作信息传输
存放url参数数据
*/
type URLData struct {
	data  JsData // 日志格式
	uid   string // 浏览用户id
	unode URLNODE
}

/*
用作存储pv|uv数据
*/
type URLNODE struct {
	unType string // 页面类型 首页|列表页|详情页
	unRid  int    // 资源ID
	unURL  string // 页面url
	unTime string // 浏览时间
}

/*
格式化数据,存储入数据库
*/
type StorageBlock struct {
	counterType  string  // 统计类型
	storageModel string  // redis相关
	uNode        URLNODE // url节点数据
}
