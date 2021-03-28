package main


import (
	"distributed-cache/log"
	"distributed-cache/client"
	"distributed-cache/client/impl"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main(){
	seed := time.Now().Unix()%1e6
	rand.Seed(seed)
	server := flag.String("h", "localhost", "cache server address")
	num :=  flag.Int("n", 1000, "set num")
	flag.Parse()
	tcpclient := impl.NewTcpClient(*server)

	for i := 0; i < *num;i++ {
		randNum := rand.Int()%1e6
		key := fmt.Sprintf("key-%d",i)
		value := fmt.Sprintf("value-%d",randNum)
		cmd := &client.Cmd{"set", key, value, nil}
		tcpclient.Run(cmd)
		log.Infof("cmd res :%+v", cmd)
	}
}
