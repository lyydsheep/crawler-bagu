package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
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

func fetchAndSaveAnswer(dir string, groupId int, problemId int) (bool, error) {
	// 检查一下是否已经写入了
	filename := fmt.Sprintf("%s/%d_%d.json", dir, groupId, problemId)
	if _, err := os.Stat(filename); err == nil {
		log.Info().Msgf("answer already exists, skip %s", filename)
		return false, nil
	} else if !os.IsNotExist(err) {
		// 如果错误不是“文件不存在”，则返回错误
		log.Error().Msgf("os.Stat error: %v", err)
		return false, err
	}

	// fetch answer
	url := fmt.Sprintf("https://www.bagujing.com/api/problems/%d?groupId=%d", problemId, groupId)
	body, err := get(url)
	if err != nil {
		return true, err
	}

	// 保存 answer
	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		log.Error().Msgf("os.WriteFile error: %v", err)
		return true, err
	}
	log.Info().Msgf("saved answer to %s", filename)
	return true, nil
}

func get(url string) ([]byte, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msgf("http.NewRequest error: %v", err)
		return nil, err
	}

	// 设置Cookie
	cookieStr := os.Getenv("Cookie")
	if len(cookieStr) == 0 {
		log.Error().Msg("Cookie is empty")
		panic(err)
	}
	req.Header.Set("Cookie", cookieStr)

	// 设置其他Header
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("client.Do error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("io.ReadAll error: %v", err)
		return nil, err
	}

	return body, err
}
