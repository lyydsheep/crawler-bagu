package internal

import (
	"crawler-bagu/internal/crawler"
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
)

type CrawlerService struct {
	// 待爬取的题目类型URL
	Destinations []string
	c            crawler.Crawler
	rl           ratelimit.Limiter
}

func (svc *CrawlerService) Run() {
	// 爬取所有的题目类型
	for _, destination := range svc.Destinations {
		if err := svc.c.Consume(destination); err != nil {
			log.Info().Msgf("svc.c.Visit(%s) error: %v", destination, err)
		}
	}
}

func NewCrawlerService(destinations []string, c crawler.Crawler, rl ratelimit.Limiter) *CrawlerService {
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
