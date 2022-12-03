package main

var (
	ParentProxyAddress = "23.234.252.229:8080"
)

func main() {
	ReadConfig()
	GetProxyAddr()
	go StartRuleProcess()
	StartPureMiddleProxy()
}
