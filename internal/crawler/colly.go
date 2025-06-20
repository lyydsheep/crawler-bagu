package crawler

import (
	"bytes"
	"crawler-bagu/pkg"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

type Colly struct {
	colly *colly.Collector
}

func (c *Colly) Consume(url string) error {
	registerHook(c.colly)
	return c.colly.Visit(url)
}

func registerHook(c *colly.Collector) {
	var data KeyContent
	c.OnHTML("div[data-widget-id='problem-exercise-widget']", func(e *colly.HTMLElement) {
		log.Info().Msgf("visiting %s, the name of html is %s.", e.Request.URL.String(), e.Name)
		jsonScript := e.DOM.Find("script[type='application/json']").First()
		if jsonScript.Length() == 0 {
			log.Error().Msgf("can not find json script")
			return
		}

		// 获取JSON数据
		jsonContent := jsonScript.Text()
		if jsonContent == "" {
			log.Error().Msgf("jsonContent == ''")
			return
		}
		log.Info().Msgf("obtained json data successfully.")

		// 解析JSON数据
		err := json.Unmarshal([]byte(jsonContent), &data)
		if err != nil {
			log.Error().Msgf("json.Unmarshal error: %v.", err)
			return
		}
		log.Info().Msgf("parsed json data successfully, the problemList length is %d.", len(data.ProblemList))
	})

	c.OnScraped(func(resp *colly.Response) {
		log.Info().Msgf("scraped %s.", resp.Request.URL.String())
		// 处理数据
		if err := processContent(data); err != nil {
			log.Error().Msgf("processContent error: %v.", err)
			return
		}
		log.Info().Msgf("processed %d questions successfully.", len(data.ProblemList))
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Error().Msgf("resp.Request.URL: %s, error: %v.", resp.Request.URL.String(), err)
	})
}

func processContent(data KeyContent) (err error) {
	// 创建题目类型目录
	var dir string
	if dir, err = pkg.Touch(data.PageTitle); err != nil {
		// 不是文件已存在错误
		if !os.IsExist(err) {
			log.Error().Msgf("Touch error: %v", err)
			return err
		}
	}

	// 爬取有效题目
	cnt := 1
	for _, problem := range data.ProblemList {
		if len(problem.Corps) > 0 {
			flag := true
			if flag, err = fetchAndSaveAnswer(dir, data.GroupId, problem.Id); err != nil {
				log.Error().Msgf("fetchAndSaveAnswer error: %v, groupId is %d, problemId is %d", err, data.GroupId, problem.Id)
				continue
			}
			cnt++
			// 成功写入文件
			if flag {
				// 两秒一题
				time.Sleep(10*time.Second + time.Duration(8000*rand.Float32())*time.Millisecond)
			}
		}
	}
	log.Info().Msgf("the number of valid questions is %d", cnt)
	return nil
}

var (
	lastId int
	gId    int
)

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
	gId = groupId
	body, err := get(url)
	if err != nil {
		return true, err
	}

	var resp Response
	if err = json.NewDecoder(bytes.NewReader(body)).Decode(&resp); err != nil {
		log.Error().Msgf("json.NewDecoder error: %v", err)
		return true, err
	}
	if resp.Data.EmptyReason != "" {
		log.Warn().Msgf("empty reason: %s", resp.Data.EmptyReason)
		return true, nil
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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", fmt.Sprintf("https://www.bagujing.com/problem-exercise/%d?pid=%d", gId, lastId))

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
