package main

import (
	limit_util "current-limit-demo/limit-util"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	sentinelUrl = "/sentinel"
	chanUrl     = "/chan"
	counter     = "/counter"
)

var (
	chan_limiter    *limit_util.ChannelLimiter
	counter_limiter *limit_util.CountLimiter
)

func main() {
	initAll()
}

func initAll() {
	// sentinel限流器
	sentinelInit()
	createFlowRule("limitWithSentinel", 1, 1000)
	// chan限流器
	chan_limiter = limit_util.NewChannelLimiter(1)
	counter_limiter = limit_util.NewCountLimiter(1*time.Second, 1)
	httpInit()
}

// 初始化Sentinel
func sentinelInit() {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化http，并注册接口
func httpInit() {
	server := http.Server{
		Addr: "127.0.0.1:8003",
	}
	http.HandleFunc(sentinelUrl, limitWithSentinel)
	http.HandleFunc(chanUrl, limitWithChan)
	http.HandleFunc(counter, limitWithCounter)
	server.ListenAndServe()
}

// chan单机限流器
func limitWithChan(w http.ResponseWriter, r *http.Request) {
	if chan_limiter.Allow() {
		fmt.Fprintf(w, "pass～")
		time.Sleep(1 * time.Second)
		chan_limiter.Release()
	} else {
		fmt.Fprintf(w, "限流！")
	}
}

// sentinel限流器
func limitWithSentinel(w http.ResponseWriter, r *http.Request) {
	resourceName := runFuncName()
	// 埋点（流控规则方式）
	e, b := sentinel.Entry(resourceName, sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		fmt.Fprintf(w, "限流！！！")
	} else {
		fmt.Fprintf(w, "测试流控规则~~~")
		e.Exit()
	}
}

// counter限流器
func limitWithCounter(w http.ResponseWriter, r *http.Request) {
	if counter_limiter.Allow() {
		fmt.Fprintf(w, "pass～")
	} else {
		fmt.Fprintf(w, "限流！")
	}
}

// 创建流控规则（默认基于QPS）
// threshold 阈值
// interval 统计间隔（毫秒）
func createFlowRule(resourceName string, threshold float64, interval uint32) {
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              threshold,
			StatIntervalInMs:       interval,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	totalFuncName := runtime.FuncForPC(pc[0]).Name()
	names := strings.Split(totalFuncName, ".")
	return names[len(names)-1]
}
