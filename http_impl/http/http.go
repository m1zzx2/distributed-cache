package http

import (
	"distributed-cache/cache"
	cacheImpl "distributed-cache/cache/impl"
	"distributed-cache/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	//todo
	http.HandleFunc("/GET", s.GetKey)

	http.HandleFunc("/PUT", s.PutObject)

	http.HandleFunc("/INFO", s.Info)
	log.Infof("start service port : %v", "8080")
	http.ListenAndServe(":8080", nil)
}

func NewServer() *Server {
	return &Server{
		Cache: cacheImpl.NewMemoryCache(),
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
	if r.Method == "PUT"{
		byte, _ := ioutil.ReadAll(r.Body)
		var r Request
		err := json.Unmarshal(byte, &r)
		if err != nil{
			Response(w, nil, 1, err.Error())
			return	
		}
		s.Cache.Set(r.Key, r.Value)
		Response(w, nil, 0, "")
	}
}

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		stat := s.Cache.GetStat()
		Response(w, stat, 0, "")
	}
}
