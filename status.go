package cuimeter

type Status map[string]int64

func NewStatus(contents []string) *Status {
	m := Status{}
	for _, v := range contents {
		m[v] = 0
	}
	return &m
}

func (s Status) at(key string) int64 {
	return map[string]int64(s)[key]
}

func (s *Status) set(key string, val int64) {
	map[string]int64(*s)[key] = val
}

func (s Status) Sub(right *Status) *Status {
	m := map[string]int64{}
	for k := range map[string]int64(s) {
		m[k] = s.at(k) - right.at(k)
	}
	st := Status(m)
	return &st
}
