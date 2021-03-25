package http

import (
	"distributed-cache/cache"
	"distributed-cache/cluster"
	"distributed-cache/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	//todo
	http.HandleFunc("/GET", s.GetKey)

	http.HandleFunc("/PUT", s.PutObject)

	http.HandleFunc("/INFO", s.Info)

	http.HandleFunc("/cluster", s.Cluster)
	log.Infof("start service port : %v", "8080")
	http.ListenAndServe(":8080", nil)
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
		s.Cache.Set(r.Key, r.Value)
		Response(w, nil, 0, "")
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
		m:= s.Members()
		b, e := json.Marshal(m)
		if e != nil{
			log.Errorf("Cluster marshal error err :%+v",e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}








func (s *Server) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stat := s.Cache.GetStat()
		Response(w, stat, 0, "")
	}
}
