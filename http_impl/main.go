package main

import (
	"distributed-cache/http"
)

func main(){
	httpInstance := http.NewServer()
	httpInstance.Listen()
}