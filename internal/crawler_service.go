package internal

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
	"math/rand/v2"
	"os"
	"time"
)

type CrawlerService struct {
	// 待爬取的题目类型URL
	Destinations []string
	c            *colly.Collector
	rl           ratelimit.Limiter
}

func (svc *CrawlerService) registerHook(c *colly.Collector) {
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
		if err := svc.processContent(data); err != nil {
			log.Error().Msgf("processContent error: %v.", err)
			return
		}
		log.Info().Msgf("processed %d questions successfully.", len(data.ProblemList))
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Error().Msgf("resp.Request.URL: %s, error: %v.", resp.Request.URL.String(), err)
	})
}

func (svc *CrawlerService) processContent(data KeyContent) (err error) {
	// 创建题目类型目录
	var dir string
	if dir, err = Touch(data.PageTitle); err != nil {
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

func (svc *CrawlerService) Run() {
	svc.registerHook(svc.c)
	// 爬取所有的题目类型
	for _, destination := range svc.Destinations {
		if err := svc.c.Visit(destination); err != nil {
			log.Info().Msgf("svc.c.Visit(%s) error: %v", destination, err)
		}
	}
}

func NewCrawlerService(destinations []string, c *colly.Collector, rl ratelimit.Limiter) *CrawlerService {
	return &CrawlerService{
		Destinations: destinations,
		c:            c,
		rl:           rl,
	}
}

func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
		colly.Async(false),
	)
	return c
}

func NewRL(count int) ratelimit.Limiter {
	return ratelimit.New(count)
}
