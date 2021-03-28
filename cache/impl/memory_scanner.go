package impl

import "distributed-cache/log"

type InMemoryScanner struct {
	pair
	pairCh  chan *pair
	closeCh chan struct{}
}

func (i *InMemoryScanner) Scan() bool {
	p, ok := <-i.pairCh
	log.Infof("InMemoryScanner p:%+v ok :%+v",p,ok)
	if ok {
		i.k, i.v = p.k, p.v
	}
	return ok
}

func (i *InMemoryScanner) Key() string {
	return i.k
}

func (i *InMemoryScanner) Value() []byte {
	return i.v
}

func (i *InMemoryScanner) Close() {
	close(i.closeCh)
}
