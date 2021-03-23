package main


import (
	"distributed-cache/log"
	"flag"
)

func main(){
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "","key")
	value := flag.String("v", "","value")
	flag.Parse()
	log.Infof("server :%+v op :%+v key :%+v value :%+v ",server,op,key,value)

}
