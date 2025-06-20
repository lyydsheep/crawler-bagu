package crawler

type Problem struct {
	Id        int      `json:"id"`
	Type      int      `json:"type"`
	BriefName string   `json:"brief_name"`
	Count     int      `json:"count"`
	Level     int      `json:"level"`
	Freq      float64  `json:"freq"`
	Kps       []string `json:"kps"`
	Corps     []string `json:"corps"`
}

type KeyContent struct {
	GroupId     int       `json:"groupId"`
	PageTitle   string    `json:"pageTitle"`
	ProblemList []Problem `json:"problemList"`
}
