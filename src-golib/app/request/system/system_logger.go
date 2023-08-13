package request

type SystemLoggerRequest struct {
	Level   string `json:"level" bind:"query" process:"trim,cut(20),tolower"`
	Content string `json:"content" process:"trim"`
}

type SystemLoggerResponse struct{}

func (r *SystemLoggerRequest) Verify() error {

	return nil
}

func (r *SystemLoggerRequest) Clear() {

}
