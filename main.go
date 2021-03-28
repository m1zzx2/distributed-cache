package main

import (
	cacheImpl "distributed-cache/cache/impl"
	clusterImpl "distributed-cache/cluster/impl"
	"distributed-cache/log"
	"distributed-cache/tcp"
	"distributed-cache/http"
	"flag"
)

func main() {
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Infof("node :%+v clus: %+v", *node, *clus)
	ca := cacheImpl.NewMemoryCache()
	cluster, err := clusterImpl.NewNode(*node, *clus)
	if err != nil{
		panic(err)
	}
	httpInstance := http.NewServer(ca, cluster)
	go tcp.NewCache(ca, cluster).Listen()
	httpInstance.Listen()
}