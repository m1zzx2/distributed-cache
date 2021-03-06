package http

import (
	"bytes"
	"distributed-cache/cache"
	"distributed-cache/cluster"
	"distributed-cache/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Server struct {
	cache.Cache
	Node cluster.Node
}

func (s *Server) Listen() {
	//todo
	http.HandleFunc("/GET", s.GetKey)

	http.HandleFunc("/PUT", s.PutObject)

	http.HandleFunc("/INFO", s.Info)

	http.HandleFunc("/cluster", s.Cluster)

	http.HandleFunc("/rebalance",s.Rebalanced)
	log.Infof("start service port : %v", "8080")
	http.ListenAndServe(s.Node.Addr()+":8080", nil)
}

func NewServer(cache cache.Cache, node cluster.Node) *Server {
	return &Server{
		Cache: cache,
		Node:  node,
	}
}

func (s *Server) GetKey(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		key := r.URL.Query().Get("key")
		value, err := s.Get(key)
		if err != nil {
			Response(w, nil, 1, err.Error())
			return
		}
		Response(w, value, 0, "")
		return
	}
}

func (s *Server) PutObject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		byte, _ := ioutil.ReadAll(r.Body)
		var r Request
		err := json.Unmarshal(byte, &r)
		if err != nil {
			Response(w, nil, 1, err.Error())
			return
		}
		log.Infof("start put node k :%+v v :%+v",r.Key, r.Value)
		s.Cache.Set(r.Key, r.Value)
		Response(w, nil, 0, "successful")
	}
}

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stat := s.Cache.GetStat()
		Response(w, stat, 0, "")
	}
}

func (s *Server) Cluster(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		m := s.Node.Members()
		b, e := json.Marshal(m)
		if e != nil {
			log.Errorf("Cluster marshal error err :%+v", e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

func (h *Server) rebalance() {
	s := h.NewScanner()
	defer s.Close()
	c := &http.Client{}
	for s.Scan() {
		k := s.Key()
		log.Infof("scanner k :%+v v :%+v", k, string(s.Value()))
		n, ok := h.Node.ShouldProcess(k)
		if !ok {
			req := Request{
				Value: s.Value(),
				Key:   k,
			}
			byteReq, _ := json.Marshal(req)
			r, _ := http.NewRequest(http.MethodPut, "http://"+n+":8080/PUT", bytes.NewBuffer(byteReq))
			resp, err := c.Do(r)
			if err != nil {
				log.Errorf("put cache node :%+v req :%+v rebalance err :%+v", n, req, err)
				continue
			}
			respByte, _ := ioutil.ReadAll(resp.Body)
			log.Infof("put cache node :%+v req :%+v rebalance resp :%s", n, req, string(respByte))
			h.Del(k)
		}
	}
}

func (s *Server) Rebalanced(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		go s.rebalance()
		Response(w, "",0,"success")
	}
}

func (s *Server) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stat := s.Cache.GetStat()
		Response(w, stat, 0, "")
	}
}
