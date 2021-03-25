package impl

import (
	"bufio"
	"distributed-cache/client"
	"distributed-cache/log"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *tcpClient) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (c *tcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func (c *tcpClient) readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	log.Infof("tcpClient.readLen end readstring")
	if e != nil {
		log.Infof("error :%+v", e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Infof("error : %+V", e)
	}
	return l
}

func (c *tcpClient) recvResponse() (string, error) {
	log.Infof("recvResponse size : %+v",c.r.Size())
	bytes := make([]byte, 0)
	n, err := c.r.Read(bytes)
	if err != nil{
		log.Errorf("recvResponse r. read  error err :%+v",err)
	}
	log.Infof("recvResponse len n :%+v lenbyte :%d",n,len(bytes))
	vlen := c.readLen(c.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	log.Infof("tcpClient.recvResponse start read io stream vlen :%+v",vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (c *tcpClient) Run(cmd *client.Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
	}
	log.Errorf("cmd not found :%+v", cmd.Value)
}

func (c *tcpClient) PipelineRun(cmds []*client.Cmd) {
	for _, cmd := range cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value)
		}
		if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
	return
}

func NewTcpClient(server string) *tcpClient {
	conn, err := net.Dial("tcp", server+":12346")
	if err != nil {
		log.Errorf("NewTcpClient net.Dial error err :%+v", err)
		return nil
	}
	r := bufio.NewReader(conn)
	return &tcpClient{conn, r}
}
