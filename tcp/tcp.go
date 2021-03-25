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
	log.Infof("tcp start listen :12346")
	l, e := net.Listen("tcp", ":12346")
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

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()+" "
		tmp := fmt.Sprintf("-%d", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	log.Infof("end sendRespose  msg :%v",string(append([]byte(vlen), value...)))
	return e
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	v, e := s.cache.Get(k)
	return sendResponse(v, e, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := s.readKeyAndVal(r)
	if e != nil {
		log.Errorf("set k :%+v v : %+v e :%+v", string(k), string(v), e)
	} else {
		log.Infof("set k :%+v v : %+v successful", string(k), string(v))
		s.cache.Set(k, v)
	}
	return sendResponse(v, e, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	return sendResponse(nil, s.cache.Del(k), conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		log.Infof("tcp server process receive op :%+v", string(op))
		if e != nil {
			if e != io.EOF {
				log.Errorf("close connection due to error :%v", e)
			}
			return
		}

		if string(op) == string('S') {
			s.set(conn, r)
		} else if string(op) == string('G') {
			s.get(conn, r)
		} else if string(op) == string('D') {
			s.del(conn, r)
		} else {
			log.Errorf("close connection due to invalid operator: %v", string(op))
			return
		}
		if e != nil {
			log.Infof("close connection due to error :%v", e)
		}
	}
}
