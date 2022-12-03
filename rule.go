package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var (
	LastRand = -1
)

func StartRuleProcess() {
	if ProxyRule == "Time" {
		LastAddress := ""
		NowAddress := AllProxyAddress[GetRand(0, len(AllProxyAddress))].Address
		for {
			if LastAddress == NowAddress {
				continue
			}
			log.Printf("Time模式,当前使用父级代理地址为 %v", ParentProxyAddress)
			ParentProxyAddress = NowAddress
			LastAddress = NowAddress
			NowAddress = AllProxyAddress[GetRand(0, len(AllProxyAddress))].Address
			time.Sleep(time.Duration(ProxyRuleTime) * time.Second)
		}
	} else if ProxyRule == "Romdom" {
		//
	} else if ProxyRule == "Next" {
		ParentProxyAddress := AllProxyAddress[GetRand(LastRand, len(AllProxyAddress))].Address
		log.Printf("Next模式,当前使用父级代理地址为 %v", ParentProxyAddress)
		for {
			code, err := getWithProxy(CheckTarget, ParentProxyAddress)
			if err == nil && code == 200 {
				time.Sleep(time.Duration(ProxyRuleTime) * time.Second)
			} else {
				log.Printf("Next模式,Check Status: %d, error: %v\n", code, err)
				ParentProxyAddress = AllProxyAddress[GetRand(LastRand, len(AllProxyAddress))].Address
				log.Printf("Next模式,当前使用父级代理地址为 %v", ParentProxyAddress)
				time.Sleep(10)
			}
		}
	} else if ProxyRule == "Once" {
		ParentProxyAddress = AllProxyAddress[GetRand(0, len(AllProxyAddress))].Address
		log.Printf("Once模式,当前使用父级代理地址为 %v", ParentProxyAddress)
	} else if ProxyRule == "Static" {
		ParentProxyAddress = ProxyRuleStatic
		log.Printf("Static模式,当前使用父级代理地址为 %v", ParentProxyAddress)
	}
}

func GetRand(s, e int) int {
	n := -1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n = r.Intn(e)
	time.Sleep(10)
	for s == n {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n = r.Intn(e)
		time.Sleep(10)
	}
	LastRand = n
	return n
}

func getWithProxy(req_url string, proxy_url string) (int, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxy_url)
	}
	transport := &http.Transport{Proxy: proxy}

	c := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		return -1, err
	}
	res, err := c.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}
