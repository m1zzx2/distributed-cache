package tcp

import (
	"bufio"
	"distributed-cache/cache"
	"distributed-cache/log"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	cache cache.Cache
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":12345")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {

		}
		go s.process(c)
	}
}

func NewCache(c cache.Cache) *Server {
	return &Server{
		c,
	}
}
func readLen(r *bufio.Reader) (int, error) {
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0, e
	}
	return l, nil
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}
	k := make([]byte, klen)
	if e != nil {
		return "", e
	}
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	return string(k), nil
}

func (s *Server) readKeyAndVal(r *bufio.Reader) (string, []byte, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}

	vlen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	k := make([]byte, klen)
	if e != nil {
		return "", nil, e
	}

	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, e
	}

	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}
	return string(k), v, nil
}

func sendRespose(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	v, e := s.cache.Get(k)
	return sendRespose(v, e, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := s.readKeyAndVal(r)
	if e != nil {
		return e
	}
	s.cache.Set(k, v)
	return sendRespose(nil, nil, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	return sendRespose(nil, s.cache.Del(k), conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Errorf("close connection due to error :%v", e)
			}
			return
		}
		if op == 'S' {
			s.set(conn, r)
		}
		if op == 'G' {
			s.get(conn, r)
		}
		if op == 'D' {
			s.del(conn, r)
		} else {
			log.Infof("close connection due to invalid operator: %v", op)
			return
		}
		if e != nil {
			log.Infof("close connection due to error :%v", e)
		}
	}
}
