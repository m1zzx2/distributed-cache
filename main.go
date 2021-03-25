package main

import (
	cacheImpl "distributed-cache/cache/impl"
	"distributed-cache/http"
	"distributed-cache/tcp"
)

func main(){
	httpInstance := http.NewServer()
	ca := cacheImpl.NewMemoryCache()
	go tcp.NewCache(ca).Listen()
	httpInstance.Listen()
}