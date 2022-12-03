package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AddressInfo struct {
	Address string
	IsUsed  bool
	Status  string
}

var (
	AllProxyAddress []AddressInfo
)

func GetProxyAddr() {
	AllProxyAddress = make([]AddressInfo, 0)
	for _, pool := range ProxyPoolAddress {
		if pool.ProxyPoolType == "jhao104/proxy_pool" {
			getFromProxyPool(pool)
		} else if pool.ProxyPoolType == "Python3WebSpider/ProxyPool" {
			getFromproxy_pool(pool)
		}
	}
	log.Printf("当前共加载代理地址%d个", len(AllProxyAddress))
	//log.Printf("%v\n", AllProxyAddress)
}

func getFromProxyPool(pool ProxyPoolInfo) {
	data, err := getData(pool.Address + "/all")
	if err != nil {
		log.Fatalf("代理池%v获取失败,请检查代理池或配置文件, %v\n", pool.Address, err)
	}
	var v interface{}
	json.Unmarshal(data, &v)
	jsonData := v.([]interface{})
	count := 0
	for _, line := range jsonData {
		temp := line.(map[string]interface{})
		if v, ok := temp["proxy"]; ok {
			AllProxyAddress = append(AllProxyAddress, AddressInfo{
				Address: v.(string),
				IsUsed:  false,
				Status:  "uncheck",
			})
			count++
		}
	}
	log.Printf("从代理池 %v 成功加载代理IP地址%d个", pool.Address, count)
}

func getFromproxy_pool(pool ProxyPoolInfo) {
	result := make(map[string]string, 0)
	for len(result) < pool.IPCount {
		data, err := getData(pool.Address + "/random")
		if err != nil {
			log.Fatalf("代理池%v获取失败,请检查代理池或配置文件, %v\n", pool.Address, err)
		}
		if strings.Contains(string(data), ":") {
			result[string(data)] = string(data)
		}
		time.Sleep(1)
	}
	for _, v := range result {
		AllProxyAddress = append(AllProxyAddress, AddressInfo{
			Address: v,
			IsUsed:  false,
			Status:  "uncheck",
		})
	}
	log.Printf("从代理池 %v 成功加载代理IP地址%d个", pool.Address, len(result))
}

func getData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
