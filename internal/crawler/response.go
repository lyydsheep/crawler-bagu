package crawler

type Response struct {
	Success bool `json:"success"`
	Data    struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		GroupId     int    `json:"group_id"`
		EmptyReason string `json:"emptyReason"`
	} `json:"data"`
}
