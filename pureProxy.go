package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

func StartPureMiddleProxy() {
	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse("http://" + EndProxyServerAddress)
	}
	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler("http://"+EndProxyServerAddress, nil)
	//middleProxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	/*middleProxy.OnRequest().DoFunc(
	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Printf("proxy rec %v\n", r.Host)
		return r, nil
	})*/
	log.Printf("Proxy Server 监听地址:  %v (请将本地代理设置为该地址)\n", ProxyServerAddress)
	go StartPureEndProxy()
	http.ListenAndServe(ProxyServerAddress, middleProxy)
}

func StartPureEndProxy() {
	ln, err := net.Listen("tcp", EndProxyServerAddress)
	if err != nil {
		log.Fatalf("End Proxy Server启动失败: %v\n", err)
		return
	}
	log.Printf("End Proxy Server 监听地址:  %v\n", EndProxyServerAddress)
	defer ln.Close()
	for {
		tcpConn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept:", err)
			return
		}
		go PureProxyForward(tcpConn)
	}
}

func PureProxyForward(tcpConn net.Conn) {
	if ProxyRule == "Romdom" {
		ParentProxyAddress = AllProxyAddress[GetRand(0, len(AllProxyAddress))].Address
	}
	RemoteProxy, err := net.Dial("tcp", ParentProxyAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	go io.Copy(RemoteProxy, tcpConn)
	go io.Copy(tcpConn, RemoteProxy)
}
