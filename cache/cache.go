package cache

type Cache interface {
	Set(string, []byte)
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
	NewScanner() Scanner
}

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) Add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) Del(k string, v []byte){
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
