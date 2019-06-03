package main

import (
	"flag"
	"fmt"
	"github.com/XiaoXice/DHTNet"
)
const help = `
程序设计实验课作业.
`

func main(){
	flag.Usage = func() {
		fmt.Println(help)
		flag.PrintDefaults()
	}
	//configFile := flag.String("c","config.json","JSON配置文件路径")
	p2pport := flag.Int("l", 12000, "libp2p listen port")
	seed := flag.Int64("seed", 0, "给ID生成器配置随机种子")
	socksPort := flag.Int("p",8888,"代理的socks5监听端口")

	flag.Parse()
	DHTNet.New(p2pport, seed, socksPort)
}
