package internal

import (
	"fmt"
	"google.golang.org/appengine/log"
	"io"
	"net/http"
	"os"
)

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

func fetchAndSaveAnswer(dir string, groupId int, problemId int) error {
	// fetch answer
	url := fmt.Sprintf("https://www.bagujing.com/api/problems/%d?groupId=%d", problemId, groupId)
	body, err := get(url)
	if err != nil {
		return err
	}
	// 保存 answer
	filename := fmt.Sprintf("%s/%d_%d.json", dir, groupId, problemId)
	err = os.WriteFile(filename, []byte(body), 0644)
	if err != nil {
		log.Errorf(nil, "os.WriteFile error: %v", err)
		return err
	}
	log.Infof(nil, "saved answer to %s", filename)
	return nil
}

func get(url string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf(nil, "http.NewRequest error: %v", err)
		return "", err
	}

	// 设置Cookie
	cookieStr := "session_id=abc123; user_token=xyz456"
	req.Header.Set("Cookie", cookieStr)

	// 设置其他Header
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(nil, "client.Do error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(nil, "io.ReadAll error: %v", err)
		return "", err
	}

	return string(body), err
}
