package http

import (
	"distributed-cache/cache"
)

type Server struct {
	cache.Cache
}

func (s *Server)Listen(){
	//todo
}
