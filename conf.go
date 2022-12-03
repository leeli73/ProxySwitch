package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type ProxyPoolInfo struct {
	ProxyPoolType string
	Address       string
	IPCount       int
}

var (
	OriginConfigFileData = `#设置代理池,格式: ProxyPool 代理池类型 代理池URL
#目前支持: https://github.com/jhao104/proxy_pool; https://github.com/Python3WebSpider/ProxyPool
#Python3WebSpider/ProxyPool 因为没有批量获取接口，故随机获取不同的IP,样例: ProxyPool Python3WebSpider/ProxyPool https://proxypool.scrape.center 30 ,表示随机获取30个不通的父级代理地址
ProxyPool jhao104/proxy_pool http://180.76.169.220:5010
ProxyPool Python3WebSpider/ProxyPool https://proxypool.scrape.center 30

#Proxy Server 基础设置

#HTTP Proxy监听的地址
ProxyServerAddress 127.0.0.1:8080

#可变父级代理本地流量转发的工作端口
EndProxyServerAddress 127.0.0.1:48901

#父级代理切换规则
#Time: 按时间进行切换,样例: ProxyRule Time 300 ,表示每300秒切换一个父级代理
#Romdom: 每个请求都是随机的父级代理,样例: ProxyRule Romdom
#Next: 每当程序检测当前父级代理无法正常访问验证目标时,自动随机切换一个父级代理,样例: ProxyRule Next 10 ,表示每10秒检测一次是否可以正常访问目标
#Once: 程序启动时随机获取一个父级代理并保持不变,重启后自动更换,样例: ProxyRule Once
#Static: 指定某个特定的父级代理服务器,样例: ProxyRule Static 127.0.0.1:8080
ProxyRule Next 10

#父级代理验证目标地址
CheckTarget https://ip.hao86.com/`
	ProxyPoolAddress      []ProxyPoolInfo
	ProxyServerAddress    string
	EndProxyServerAddress string
	ProxyRule             string
	ProxyRuleTime         int
	ProxyRuleStatic       string
	CheckTarget           string
)

func ReadConfig() {
	data, err := ioutil.ReadFile("conf.txt")
	if err != nil {
		WirteConfig()
		log.Fatalf("读取配置文件错误,已自动写出新的配置文件,请检查并配置后重新启动本程序: %v\n", err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if len(line) > 1 && line[0:1] != "#" {
			line = strings.ReplaceAll(line, "\r", "")
			keys := strings.Split(line, " ")
			if len(keys) > 0 {
				if keys[0] == "ProxyPool" && len(keys) > 2 {
					if keys[1] == "Python3WebSpider/ProxyPool" && len(keys) > 3 {
						ipcount, err := strconv.Atoi(keys[3])
						if err != nil {
							log.Fatalf("配置文件中的ProxyPool存在错误")
						}
						ProxyPoolAddress = append(ProxyPoolAddress, ProxyPoolInfo{
							ProxyPoolType: keys[1],
							Address:       keys[2],
							IPCount:       ipcount,
						})
					} else {
						ProxyPoolAddress = append(ProxyPoolAddress, ProxyPoolInfo{
							ProxyPoolType: keys[1],
							Address:       keys[2],
						})
					}
				} else if keys[0] == "ProxyServerAddress" && len(keys) > 1 {
					ProxyServerAddress = keys[1]
				} else if keys[0] == "EndProxyServerAddress" && len(keys) > 1 {
					EndProxyServerAddress = keys[1]
				} else if keys[0] == "ProxyRule" && len(keys) > 1 {
					ProxyRule = keys[1]
					if ProxyRule == "Time" && len(keys) > 2 {
						ProxyRuleTime, err = strconv.Atoi(keys[2])
						if err != nil {
							log.Fatalf("配置文件中的ProxyRule存在错误")
						}
					} else if ProxyRule == "Static" && len(keys) > 2 {
						ProxyRuleStatic = keys[2]
					} else if ProxyRule == "Next" && len(keys) > 2 {
						ProxyRuleTime, err = strconv.Atoi(keys[2])
						if err != nil {
							log.Fatalf("配置文件中的ProxyRule存在错误")
						}
					}
				} else if keys[0] == "CheckTarget" && len(keys) > 1 {
					CheckTarget = keys[1]
				}
			}
		}
	}
	if len(ProxyPoolAddress) == 0 {
		log.Fatalf("配置文件中的ProxyPool为空或未找到")
	} else if ProxyServerAddress == "" {
		log.Fatalf("配置文件中的ProxyServerAddress为空或未找到")
	} else if EndProxyServerAddress == "" {
		log.Fatalf("配置文件中的EndProxyServerAddress为空或未找到")
	} else if ProxyRule == "" {
		log.Fatalf("配置文件中的ProxyRule为空或未找到")
	} else if CheckTarget == "" {
		log.Fatalf("配置文件中的CheckTarget为空或未找到")
	} else {
		log.Printf("当前代理池 %d个,父级代理切换模式为 %s,父级代理验证目标地址为 %s\n", len(ProxyPoolAddress), ProxyRule, CheckTarget)
	}
}

func WirteConfig() {
	err := ioutil.WriteFile("conf.txt", []byte(OriginConfigFileData), 0644)
	if err != nil {
		log.Fatalf("写出配置文件错误: %v\n", err)
	}
}
