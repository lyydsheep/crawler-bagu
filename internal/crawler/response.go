package crawler

type Response struct {
	Success bool `json:"success"`
	Data    struct {
		Id          int         `json:"id"`
		Name        string      `json:"name"`
		Options     interface{} `json:"options"`
		Answer      string      `json:"answer"`
		Type        int         `json:"type"`
		Level       int         `json:"level"`
		Freq        float64     `json:"freq"`
		Analysis    string      `json:"analysis"`
		MoreAsk     string      `json:"more_ask"`
		Mindmap     string      `json:"mindmap"`
		Keynote     string      `json:"keynote"`
		GroupId     int         `json:"group_id"`
		Kps         []string    `json:"kps"`
		Years       []int       `json:"years"`
		Corps       []string    `json:"corps"`
		EmptyReason string      `json:"emptyReason"`
	} `json:"data"`
}
