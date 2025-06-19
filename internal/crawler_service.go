package internal

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"go.uber.org/ratelimit"
	"google.golang.org/appengine/log"
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
		log.Infof(nil, "visiting %s, the name of html is %s.", e.Request.URL.String(), e.Name)
		jsonScript := e.DOM.Find("script[type='application/json']").First()
		if jsonScript.Length() == 0 {
			log.Errorf(nil, "jsonScript.Length() == 0")
			return
		}

		// 获取JSON数据
		jsonContent := jsonScript.Text()
		if jsonContent == "" {
			log.Errorf(nil, "jsonContent == ''")
			return
		}
		log.Debugf(nil, "jsonContent: %s.", jsonContent)
		log.Infof(nil, "obtained json data successfully.")

		// 解析JSON数据
		err := json.Unmarshal([]byte(jsonContent), &data)
		if err != nil {
			log.Errorf(nil, "json.Unmarshal error: %v", err)
			return
		}
		log.Debugf(nil, "data: %v", data)
		log.Infof(nil, "parsed json data successfully, the problemList length is %d.", len(data.ProblemList))
	})

	c.OnScraped(func(resp *colly.Response) {
		log.Infof(nil, "scraped %s, begin process data.", resp.Request.URL.String())
		// 处理数据
		if err := svc.processContent(data); err != nil {
			log.Errorf(nil, "processContent error: %v.", err)
			return
		}
		log.Infof(nil, "processed json data successfully.")
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Errorf(nil, "resp.Request.URL: %s, error: %v.", resp.Request.URL.String(), err)
	})
}

func (svc *CrawlerService) processContent(data KeyContent) (err error) {
	// 创建题目类型目录
	var dir string
	if dir, err = Touch(data.PageTitle); err != nil {
		log.Errorf(nil, "Touch error: %v", err)
		return err
	}

	// 爬取有效题目
	cnt := 1
	for _, problem := range data.ProblemList {
		if len(problem.Corps) > 0 {
			svc.rl.Take()
			if err = fetchAndSaveAnswer(dir, data.GroupId, problem.Id); err != nil {
				log.Errorf(nil, "fetchAndSaveAnswer error: %v, groupId is %d, problemId is %d", err, data.GroupId, problem.Id)
				continue
			}
			cnt++
		}
	}
	log.Infof(nil, "the number of valid questions is %d", cnt)
	return nil
}

func (svc *CrawlerService) Run() {
	// 爬取所有的题目类型
	for _, destination := range svc.Destinations {
		if err := svc.c.Visit(destination); err != nil {
			log.Errorf(nil, "svc.c.Visit(%s) error: %v", destination, err)
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
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36\n"),
		colly.IgnoreRobotsTxt(),
		colly.Async(false),
	)
	return c
}

func NewRL(count int) ratelimit.Limiter {
	return ratelimit.New(count)
}
