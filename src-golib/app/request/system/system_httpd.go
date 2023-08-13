package request

type SystemHttpdRequest struct {
	Method string `json:"method" bind:"query" process:"trim,cut(20),tolower"`
	Addr   string `json:"addr" bind:"query" process:"trim,cut(20),tolower"`
}

type SystemHttpdResponse struct {
	BeforeState int `json:"before_state"`
	AfterState  int `json:"after_state"`
}

func (r *SystemHttpdRequest) Verify() error {
	if r.Addr == "" {
		r.Addr = "127.0.0.1:8087"
	}
	return nil
}

func (r *SystemHttpdRequest) Clear() {

}
