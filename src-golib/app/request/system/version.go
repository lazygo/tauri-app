package request

type SystemVersionRequest struct {
}

type SystemVersionResponse struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	Debug   bool   `json:"debug"`
}

func (r *SystemVersionRequest) Verify() error {
	return nil
}

func (r *SystemVersionRequest) Clear() {
}
