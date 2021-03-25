package impl

import (
	"distributed-cache/client"
	"distributed-cache/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient)get(key string)string{
	resp, e := c.Get(c.server + key)
	if e != nil{
		log.Errorf("http client get error err :+v",e)
		panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if  resp.StatusCode != http.StatusOK{
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil{
		panic(e)
	}
	return string(b)
}


func (c *httpClient) set(key, value string){
	req, e := http.NewRequest(http.MethodPut, c.server + key, strings.NewReader(value))
	if e != nil{
		log.Errorf("NewRequest error err :%+v",e)
		panic(e)
	}
	resp, e := c.Do(req)
	if e != nil{
		log.Infof(key)
		panic(e)
	}
	if resp.StatusCode != http.StatusOK{
		panic(resp.Status)
	}
}

func (c *httpClient) Run(cmd *client.Cmd){
	if cmd.Name == "get"{
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set"{
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknow cmd name " +cmd.Name)
}

func (c *httpClient) PipelineRun([]*client.Cmd){
	return
}

func NewHttpClient(server string) *httpClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: 1,
		},
	}
	return &httpClient{
		client, "http://"+server+":12345/cache/"}
}
