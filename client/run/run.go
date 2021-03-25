package main


import (
	"distributed-cache/log"
	"distributed-cache/client"
	"distributed-cache/client/impl"
	"flag"
)

func main(){
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "","key")
	value := flag.String("v", "","value")
	flag.Parse()
	log.Infof("server :%+v op :%+v key :%+v value :%+v ",server,op,key,value)
	tcpclient := impl.NewTcpClient(*server)
	cmd := &client.Cmd{*op, *key, *value, nil}
	tcpclient.Run(cmd)
	log.Infof("cmd res :%+v", cmd)
}
