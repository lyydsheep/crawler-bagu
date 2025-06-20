package crawler

type Response struct {
	Success bool `json:"success"`
	Data    struct {
		EmptyReason string `json:"emptyReason"`
	} `json:"data"`
}
