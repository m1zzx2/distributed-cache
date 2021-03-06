package impl

import (
	"distributed-cache/log"
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"stathat.com/c/consistent"
	"time"
)

type Node struct{
	*consistent.Consistent
	addr string
}


func (n *Node) Addr() string {
	return n.addr
}


func (n *Node) ShouldProcess(key string)(string, bool){
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}



func NewNode(addr, cluster string)(*Node, error){
	conf := memberlist.DefaultLocalConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil{
		return nil, e
	}
	if cluster == ""{
		cluster = addr
	}
	clu := []string{cluster}
	_, e = l.Join(clu)
	if e != nil{
		return nil, e
	}
	circle := consistent.New()
	circle.NumberOfReplicas = 256
	go func() {
		for {
			m := l.Members()
			nodes := make([]string, len(m))
			for i , n := range m{
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	log.Infof("new node addr :%+v",addr)
	return &Node{circle, addr}, nil
}