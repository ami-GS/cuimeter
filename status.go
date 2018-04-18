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

/*
func GetStatus(target string, status *Status) error {
	// TODO: will be changed in each file/network type
	// if starts with http, then http get
	// if start from '/' then it is raed line and save cursor?
	resp, _ := http.Get(target)
	data, _ := ioutil.ReadAll(resp.Body)

	// Hint should be written by user?
	sp := strings.Split(string(data), " ")
	ngxHint := map[string]int{
		"ActiveConnections": 2,
		"Accepts":           7,
		"Handled":           8,
		"Requests":          9,
		"Reading":           11,
		"Writing":           13,
		"Waiting":           15,
	}
	for k := range map[string]int64(*status) {
		dat, err := strconv.Atoi(sp[ngxHint[k]])
		if err != nil {
			return err
		}
		status.set(k, int64(dat))
	}

	return nil
}

type TrackStatus struct {
	StatusNow *Status
	StatusPre *Status
}

func NewTrackStatus(contents []string) *TrackStatus {
	s := &TrackStatus{
		StatusNow: NewStatus(contents),
		StatusPre: NewStatus(contents),
	}
	err := GetStatus(target, s.StatusNow)
	if err != nil {
		panic(err)
	}
	err = GetStatus(target, s.StatusPre)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *TrackStatus) Sub() *Status {
	return s.StatusNow.Sub(s.StatusPre)
}

func (s *TrackStatus) StoreNow() {
	*s.StatusPre = *s.StatusNow
}
*/
