Proxy Switch
============

Proxy Switch是一个自动代理切换的工具，它在本地启动一个代理服务器，通过切换父级代理的形式，实现远程代理服务器的无缝切换。

本项目的代理均来自开源的代理池，与本项目及作者无关。

使用的代理池的项目地址：

proxy_pool： https://github.com/jhao104/proxy_pool

ProxyPool： https://github.com/Python3WebSpider/ProxyPool

![](https://github.com/leeli73/ProxySwitch/blob/main/img.png)

上面是Proxy Switch的预览效果，本地监听8080端口，启动了一个HTTP代理服务器，从代理池中获取了54个地址，程序通过对目标地址访问状态的验证，本地HTTP代理服务器自动将流量转发至对应的远程HTTP代理服务器。

支持的模式
====

Time模式

    每隔一段时间自动切换一个上级代理
    配置文件样例: ProxyRule Time 300 ,表示每300秒切换一个父级代理

Romdom模式

    每个请求都是随机的父级代理
    配置文件样例: ProxyRule Romdom

Next模式

    每当程序检测当前父级代理无法正常访问验证目标时,自动随机切换一个父级代理
    配置文件样例: ProxyRule Next 10 ,表示每10秒检测一次是否可以正常访问目标

Once模式

    程序启动时随机获取一个父级代理并保持不变,重启后自动更换
    配置文件样例: ProxyRule Once

Static模式

    指定某个特定的父级代理服务器
    配置文件样例: ProxyRule Static 127.0.0.1:8080 ,表示父级代理仅使用127.0.0.1:8080

目前已经完成
============

    从代理池中获取代理地址

    Time、Romdom、Next、Once、Static 5种父级代理切换模式

    对指定目标的访问状态检测

未来
====

    还没想好，欢迎issues
